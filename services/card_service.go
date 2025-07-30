package services

import (
	"errors"
	"fmt"
	"kanban-app/api/auth"
	"kanban-app/api/database"
	"kanban-app/api/models"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardService struct {
	labelService *LabelService
}

func NewCardService() *CardService {
	return &CardService{
		labelService: NewLabelService(),
	}
}

func (s *CardService) CreateCard(listID, title, description string, dueDate *time.Time, userID string) (*models.Card, error) {
	var maxPosition int
	result := database.DB.Table("cards").Select("COALESCE(MAX(position), 0)").Where("list_id = ?", listID).Row().Scan(&maxPosition)
	if result != nil && !errors.Is(result, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get max card position: %w", result)
	}

	newCard := models.Card{
		ID:          uuid.New().String(),
		ListID:      listID,
		Title:       title,
		Description: description,
		Position:    maxPosition + 1,
		DueDate:     dueDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createResult := database.DB.Create(&newCard)
	if createResult.Error != nil {
		return nil, fmt.Errorf("failed to create card: %w", createResult.Error)
	}

	log.Printf("Card created: %s in list %s at position %d\n", newCard.Title, newCard.ListID, newCard.Position)

	_, err := auth.NewAuthorizationService().AddPolicy(userID, newCard.ID, "owner")
	if err != nil {
		return nil, fmt.Errorf("failed to add policy for new card: %w", err)
	}

	return &newCard, nil
}

func (s *CardService) GetCardsByListID(listID string) ([]models.Card, error) {
	var cards []models.Card
	result := database.DB.Where("list_id = ?", listID).Order("position ASC").Find(&cards)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve cards: %w", result.Error)
	}
	return cards, nil
}

func (s *CardService) GetCardByID(cardID string) (*models.Card, error) {
	var card models.Card
	result := database.DB.First(&card, "id = ?", cardID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("card not found")
		}
		return nil, fmt.Errorf("failed to retrieve card: %w", result.Error)
	}
	return &card, nil
}

func (s *CardService) UpdateCard(cardID string, updateReq models.UpdateCardRequest) (*models.Card, error) {
	card, err := s.GetCardByID(cardID)
	if err != nil {
		return nil, err
	}

	// Update text fields
	if updateReq.Title != "" {
		card.Title = updateReq.Title
	}
	if updateReq.Description != "" {
		card.Description = updateReq.Description
	}
	if updateReq.DueDate != nil {
		card.DueDate = updateReq.DueDate
	}

	card.UpdatedAt = time.Now()
	if err := database.DB.Save(&card).Error; err != nil {
		return nil, fmt.Errorf("failed to update card fields: %w", err)
	}

	// Handle moving the card
	if updateReq.Position != nil {
		// If ListID is not provided, use the card's current list ID
		newListID := card.ListID
		if updateReq.ListID != "" {
			newListID = updateReq.ListID
		}
		if err := s.MoveCard(cardID, newListID, *updateReq.Position); err != nil {
			return nil, err
		}
	}

	return s.GetCardByID(cardID)
}

func (s *CardService) DeleteCard(cardID string) error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	var cardToDelete models.Card
	if err := tx.First(&cardToDelete, "id = ?", cardID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("card not found")
		}
		return fmt.Errorf("failed to find card to delete: %w", err)
	}

	if err := tx.Delete(&models.Card{}, "id = ?", cardID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete card: %w", err)
	}

	if err := tx.Model(&models.Card{}).Where("list_id = ? AND position > ?", cardToDelete.ListID, cardToDelete.Position).Update("position", gorm.Expr("position - 1")).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update positions of subsequent cards: %w", err)
	}

	log.Printf("Card deleted: ID %s\n", cardID)
	return tx.Commit().Error
}

func (s *CardService) MoveCard(cardID string, newListID string, newPosition int) error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	var card models.Card
	if err := tx.First(&card, "id = ?", cardID).Error; err != nil {
		tx.Rollback()
		return errors.New("card not found")
	}

	oldListID := card.ListID
	oldPosition := card.Position

	if oldListID == newListID && oldPosition == newPosition {
		tx.Rollback()
		return nil // No change
	}

	// Decrement positions in the old list
	if err := tx.Model(&models.Card{}).Where("list_id = ? AND position > ?", oldListID, oldPosition).Update("position", gorm.Expr("position - 1")).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to shift cards in old list: %w", err)
	}

	// Increment positions in the new list
	if err := tx.Model(&models.Card{}).Where("list_id = ? AND position >= ?", newListID, newPosition).Update("position", gorm.Expr("position + 1")).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to shift cards in new list: %w", err)
	}

	card.ListID = newListID
	card.Position = newPosition
	card.UpdatedAt = time.Now()
	if err := tx.Save(&card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update card position: %w", err)
	}

	return tx.Commit().Error
}

func (s *CardService) AddLabelToCard(cardID, labelID string) (*models.Card, error) {
	card, err := s.GetCardByID(cardID)
	if err != nil {
		return nil, err
	}

	label, err := s.labelService.GetLabelByID(labelID)
	if err != nil {
		return nil, err
	}

	// Check if label is already associated
	for _, l := range card.Labels {
		if l.ID == label.ID {
			return nil, errors.New("label already associated with this card")
		}
	}

	if err := database.DB.Model(&card).Association("Labels").Append(label); err != nil {
		return nil, fmt.Errorf("failed to add label to card: %w", err)
	}

	// Reload the card to include the newly associated label
	updatedCard, err := s.GetCardByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload card after adding label: %w", err)
	}

	return updatedCard, nil
}

func (s *CardService) RemoveLabelFromCard(cardID, labelID string) (*models.Card, error) {
	card, err := s.GetCardByID(cardID)
	if err != nil {
		return nil, err
	}

	label, err := s.labelService.GetLabelByID(labelID)
	if err != nil {
		return nil, err
	}

	// Check if label is actually associated before attempting to remove
	found := false
	for _, l := range card.Labels {
		if l.ID == label.ID {
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("label not associated with this card")
	}

	if err := database.DB.Model(&card).Association("Labels").Delete(label); err != nil {
		return nil, fmt.Errorf("failed to remove label from card: %w", err)
	}

	// Reload the card to reflect the removal of the label
	updatedCard, err := s.GetCardByID(cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload card after removing label: %w", err)
	}

	return updatedCard, nil
}


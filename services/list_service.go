package services

import (
	"errors"
	"fmt"
	"kanban-app/api/auth"
	"kanban-app/api/database"
	"kanban-app/api/models"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ListService struct{}

func NewListService() *ListService {
	return &ListService{}
}

func (s *ListService) CreateList(boardID, name, userID string) (*models.List, error) {
	var maxPosition int
	result := database.DB.Table("lists").Select("COALESCE(MAX(position), 0)").Where("board_id = ?", boardID).Row().Scan(&maxPosition)
	if result != nil && !errors.Is(result, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get max list position: %w", result)
	}

	newList := models.List{
		ID:        uuid.New().String(),
		BoardID:   boardID,
		Name:      name,
		Position:  maxPosition + 1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createResult := database.DB.Create(&newList)
	if createResult.Error != nil {
		if strings.Contains(createResult.Error.Error(), "UNIQUE constraint failed") {
			return nil, errors.New("list name already exists within this board")
		}
		return nil, fmt.Errorf("failed to create list: %w", createResult.Error)
	}

	log.Printf("List created: %s in board %s at position %d\n", newList.Name, newList.BoardID, newList.Position)

	_, err := auth.NewAuthorizationService().AddPolicy(userID, newList.ID, "owner")
	if err != nil {
		return nil, fmt.Errorf("failed to add policy for new list: %w", err)
	}

	return &newList, nil
}

func (s *ListService) GetListsByBoardID(boardID string) ([]models.List, error) {
	var lists []models.List
	result := database.DB.Where("board_id = ?", boardID).Order("position ASC").Find(&lists)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve lists: %w", result.Error)
	}
	return lists, nil
}

func (s *ListService) GetListByID(listID string) (*models.List, error) {
	var list models.List
	result := database.DB.First(&list, "id = ?", listID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("list not found")
		}
		return nil, fmt.Errorf("failed to retrieve list: %w", result.Error)
	}
	return &list, nil
}

func (s *ListService) UpdateList(listID string, updateReq models.UpdateListRequest) (*models.List, error) {
	list, err := s.GetListByID(listID)
	if err != nil {
		return nil, err
	}

	if updateReq.Name != "" {
		list.Name = updateReq.Name
		list.UpdatedAt = time.Now()
		if err := database.DB.Save(&list).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				return nil, errors.New("list name already exists within this board")
			}
			return nil, fmt.Errorf("failed to update list name: %w", err)
		}
	}

	if updateReq.Position != nil {
		if err := s.MoveList(listID, *updateReq.Position); err != nil {
			return nil, err
		}
	}

	return s.GetListByID(listID)
}

func (s *ListService) DeleteList(listID string) error {
	// In a transaction, delete the list and re-order the remaining lists
	tx := database.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	var listToDelete models.List
	if err := tx.First(&listToDelete, "id = ?", listID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("list not found")
		}
		return fmt.Errorf("failed to find list to delete: %w", err)
	}

	if err := tx.Delete(&models.List{}, "id = ?", listID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete list: %w", err)
	}

	// Update positions of subsequent lists
	if err := tx.Model(&models.List{}).Where("board_id = ? AND position > ?", listToDelete.BoardID, listToDelete.Position).Update("position", gorm.Expr("position - 1")).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update positions of subsequent lists: %w", err)
	}

	log.Printf("List deleted: ID %s\n", listID)
	return tx.Commit().Error
}

func (s *ListService) MoveList(listID string, newPosition int) error {
	if newPosition < 1 {
		return errors.New("invalid position: must be greater than 0")
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	var list models.List
	if err := tx.First(&list, "id = ?", listID).Error; err != nil {
		tx.Rollback()
		return errors.New("list not found")
	}

	oldPosition := list.Position
	if oldPosition == newPosition {
		tx.Rollback()
		return nil // No change
	}

	// Shift items between old and new positions
	if oldPosition < newPosition {
		if err := tx.Model(&models.List{}).Where("board_id = ? AND position > ? AND position <= ?", list.BoardID, oldPosition, newPosition).Update("position", gorm.Expr("position - 1")).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to shift lists down: %w", err)
		}
	} else {
		if err := tx.Model(&models.List{}).Where("board_id = ? AND position >= ? AND position < ?", list.BoardID, newPosition, oldPosition).Update("position", gorm.Expr("position + 1")).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to shift lists up: %w", err)
		}
	}

	list.Position = newPosition
	list.UpdatedAt = time.Now()
	if err := tx.Save(&list).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update list position: %w", err)
	}

	return tx.Commit().Error
}


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

type BoardService struct{}

func NewBoardService() *BoardService {
	return &BoardService{}
}

func (s *BoardService) CreateBoard(projectID, name, description, userID string) (*models.Board, error) {
	board := models.Board{
		ID:          uuid.New().String(),
		ProjectID:   projectID,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result := database.DB.Create(&board)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") && strings.Contains(result.Error.Error(), "boards.name") {
			return nil, errors.New("board name already exists within this project")
		}
		return nil, fmt.Errorf("failed to create board: %w", result.Error)
	}

	log.Printf("Board created: %s in project %s\n", board.Name, board.ProjectID)

	// Add policy to Casbin
	_, err := auth.NewAuthorizationService().AddPolicy(userID, board.ID, "owner")
	if err != nil {
		return nil, fmt.Errorf("failed to add policy for new board: %w", err)
	}

	return &board, nil
}

func (s *BoardService) GetBoardsByProjectID(projectID string) ([]models.Board, error) {
	var boards []models.Board
	result := database.DB.Where("project_id = ?", projectID).Find(&boards)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve boards: %w", result.Error)
	}
	return boards, nil
}

func (s *BoardService) GetBoardByID(boardID string) (*models.Board, error) {
	var board models.Board
	result := database.DB.First(&board, "id = ?", boardID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("board not found")
		}
		return nil, fmt.Errorf("failed to retrieve board: %w", result.Error)
	}
	return &board, nil
}

func (s *BoardService) UpdateBoard(boardID string, updateReq models.UpdateBoardRequest) (*models.Board, error) {
	board, err := s.GetBoardByID(boardID)
	if err != nil {
		return nil, err
	}

	if updateReq.Name != "" {
		board.Name = updateReq.Name
	}
	if updateReq.Description != "" {
		board.Description = updateReq.Description
	}
	board.UpdatedAt = time.Now()

	result := database.DB.Save(&board)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") && strings.Contains(result.Error.Error(), "boards.name") {
			return nil, errors.New("board name already exists within this project")
		}
		return nil, fmt.Errorf("failed to update board: %w", result.Error)
	}
	log.Printf("Board updated: %s (ID: %s)\n", board.Name, board.ID)
	return board, nil
}

func (s *BoardService) DeleteBoard(boardID string) error {
	result := database.DB.Delete(&models.Board{}, "id = ?", boardID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete board: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("board not found or already deleted")
	}
	log.Printf("Board deleted: ID %s\n", boardID)
	return nil
}

func (s *BoardService) GetBoardDetails(boardID string) (*models.Board, error) {
	var board models.Board
	result := database.DB.Preload("Lists", func(db *gorm.DB) *gorm.DB {
		return db.Order("lists.position ASC")
	}).Preload("Lists.Cards", func(db *gorm.DB) *gorm.DB {
		return db.Order("cards.position ASC")
	}).Preload("Lists.Cards.Labels").Preload("Lists.Cards.Comments").Preload("Lists.Cards.Comments.User").Preload("Lists.Cards.Attachments").First(&board, "id = ?", boardID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("board not found")
		}
		return nil, fmt.Errorf("failed to retrieve board details: %w", result.Error)
	}
	return &board, nil
}


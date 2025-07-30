package services

import (
	"fmt"
	"kanban-app/api/database"
	"kanban-app/api/models"
	"time"

	"github.com/google/uuid"
)

type CommentService struct{}

func NewCommentService() *CommentService {
	return &CommentService{}
}

func (s *CommentService) CreateComment(cardID, userID, content string) (*models.Comment, error) {
	comment := models.Comment{
		ID:        uuid.New().String(),
		CardID:    cardID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return &comment, nil
}

func (s *CommentService) DeleteComment(commentID string) error {
	if err := database.DB.Delete(&models.Comment{}, "id = ?", commentID).Error; err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}

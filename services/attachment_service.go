package services

import (
	"fmt"
	"kanban-app/api/database"
	"kanban-app/api/models"
	"time"

	"github.com/google/uuid"
)

type AttachmentService struct{}

func NewAttachmentService() *AttachmentService {
	return &AttachmentService{}
}

func (s *AttachmentService) CreateAttachment(cardID, fileName, fileURL, fileType string) (*models.Attachment, error) {
	attachment := models.Attachment{
		ID:        uuid.New().String(),
		CardID:    cardID,
		FileName:  fileName,
		FileURL:   fileURL,
		FileType:  fileType,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&attachment).Error; err != nil {
		return nil, fmt.Errorf("failed to create attachment: %w", err)
	}

	return &attachment, nil
}

func (s *AttachmentService) DeleteAttachment(attachmentID string) error {
	if err := database.DB.Delete(&models.Attachment{}, "id = ?", attachmentID).Error; err != nil {
		return fmt.Errorf("failed to delete attachment: %w", err)
	}
	return nil
}

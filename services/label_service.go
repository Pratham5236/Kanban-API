package services

import (
	"errors"
	"fmt"
	"kanban-app/api/database"
	"kanban-app/api/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LabelService struct{}

func NewLabelService() *LabelService {
	return &LabelService{}
}

func (s *LabelService) CreateLabel(name, color string) (*models.Label, error) {
	label := models.Label{
		ID:        uuid.New().String(),
		Name:      name,
		Color:     color,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := database.DB.Create(&label)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") && strings.Contains(result.Error.Error(), "labels.name") {
			return nil, errors.New("label name already exists")
		}
		return nil, fmt.Errorf("failed to create label: %w", result.Error)
	}

	return &label, nil
}

func (s *LabelService) GetAllLabels() ([]models.Label, error) {
	var labels []models.Label
	result := database.DB.Find(&labels)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve labels: %w", result.Error)
	}
	return labels, nil
}

func (s *LabelService) GetLabelByID(labelID string) (*models.Label, error) {
	var label models.Label
	result := database.DB.First(&label, "id = ?", labelID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("label not found")
		}
		return nil, fmt.Errorf("failed to retrieve label: %w", result.Error)
	}
	return &label, nil
}

func (s *LabelService) UpdateLabel(labelID string, updateReq models.UpdateLabelRequest) (*models.Label, error) {
	label, err := s.GetLabelByID(labelID)
	if err != nil {
		return nil, err
	}

	if updateReq.Name != nil {
		label.Name = *updateReq.Name
	}
	if updateReq.Color != nil {
		label.Color = *updateReq.Color
	}
	label.UpdatedAt = time.Now()

	result := database.DB.Save(&label)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") && strings.Contains(result.Error.Error(), "labels.name") {
			return nil, errors.New("label name already exists")
		}
		return nil, fmt.Errorf("failed to update label: %w", result.Error)
	}
	return label, nil
}

func (s *LabelService) DeleteLabel(labelID string) error {
	result := database.DB.Delete(&models.Label{}, "id = ?", labelID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete label: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("label not found or already deleted")
	}
	return nil
}

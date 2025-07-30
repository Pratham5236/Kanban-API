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

type OrganizationService struct{}

func NewOrganizationService() *OrganizationService {
	return &OrganizationService{}
}

func (s *OrganizationService) CreateOrganization(name, ownerID string) (*models.Organization, error) {
	org := models.Organization{
		ID:        uuid.New().String(),
		Name:      name,
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := database.DB.Create(&org)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") && strings.Contains(result.Error.Error(), "organizations.name") {
			return nil, errors.New("organization name already exists")
		}
		return nil, fmt.Errorf("failed to create organization: %w", result.Error)
	}

	log.Printf("Organization created: %s by user %s\n", org.Name, org.OwnerID)

	// Add policy to Casbin
	_, err := auth.NewAuthorizationService().AddPolicy(ownerID, org.ID, "owner")
	if err != nil {
		return nil, fmt.Errorf("failed to add policy for new organization: %w", err)
	}

	return &org, nil
}

func (s *OrganizationService) GetOrganizationsByUser(ownerID string) ([]models.Organization, error) {
	var organizations []models.Organization
	result := database.DB.Where("owner_id = ?", ownerID).Find(&organizations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve organizations: %w", result.Error)
	}
	return organizations, nil
}

func (s *OrganizationService) GetOrganizationByID(orgID string) (*models.Organization, error) {
	var organization models.Organization
	result := database.DB.First(&organization, "id = ?", orgID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("organization not found")
		}
		return nil, fmt.Errorf("failed to retrieve organization: %w", result.Error)
	}
	return &organization, nil
}

func (s *OrganizationService) UpdateOrganization(orgID string, updateReq models.UpdateOrganizationRequest) (*models.Organization, error) {
	org, err := s.GetOrganizationByID(orgID)
	if err != nil {
		return nil, err
	}

	if updateReq.Name != "" {
		org.Name = updateReq.Name
	}
	org.UpdatedAt = time.Now()

	result := database.DB.Save(&org)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") && strings.Contains(result.Error.Error(), "organizations.name") {
			return nil, errors.New("organization name already exists")
		}
		return nil, fmt.Errorf("failed to update organization: %w", result.Error)
	}
	log.Printf("Organization updated: %s (ID: %s)\n", org.Name, org.ID)
	return org, nil
}

func (s *OrganizationService) DeleteOrganization(orgID string) error {
	result := database.DB.Delete(&models.Organization{}, "id = ?", orgID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete organization: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("organization not found or already deleted")
	}
	log.Printf("Organization deleted: ID %s\n", orgID)
	return nil
}

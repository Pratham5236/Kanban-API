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

type ProjectService struct{}

func NewProjectService() *ProjectService {
	return &ProjectService{}
}

func (s *ProjectService) CreateProject(organizationID, name, description, userID string) (*models.Project, error) {
	project := models.Project{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		Name:           name,
		Description:    description,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	result := database.DB.Create(&project)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "Unique constraint failed") && strings.Contains(result.Error.Error(), "projects.name") {
			return nil, errors.New("project name already exists within this organization")
		}
		return nil, fmt.Errorf("failed to create project: %w", result.Error)
	}

	log.Printf("Project created: %s in organization %s\n", project.Name, project.OrganizationID)

	// Add policy to Casbin
	_, err := auth.NewAuthorizationService().AddPolicy(userID, project.ID, "owner")
	if err != nil {
		return nil, fmt.Errorf("failed to add policy for new project: %w", err)
	}

	return &project, nil
}

func (s *ProjectService) GetProjectsByOrganizationID(organizationID string) ([]models.Project, error) {
	var projects []models.Project
	result := database.DB.Where("organization_id = ?", organizationID).Find(&projects)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve projects: %w", result.Error)
	}
	return projects, nil
}

func (s *ProjectService) GetProjectByID(projectID string) (*models.Project, error) {
	var project models.Project
	result := database.DB.First(&project, "id = ?", projectID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, fmt.Errorf("failed to retrieve project: %w", result.Error)
	}
	return &project, nil
}

func (s *ProjectService) UpdateProject(projectID string, updateReq models.UpdateProjectRequest) (*models.Project, error) {
	project, err := s.GetProjectByID(projectID)
	if err != nil {
		return nil, err
	}

	if updateReq.Name != "" {
		project.Name = updateReq.Name
	}
	if updateReq.Description != "" {
		project.Description = updateReq.Description
	}
	project.UpdatedAt = time.Now()

	result := database.DB.Save(&project)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") && strings.Contains(result.Error.Error(), "projects.name") {
			return nil, errors.New("project name already exists within this organization")
		}
		return nil, fmt.Errorf("failed to update project: %w", result.Error)
	}
	log.Printf("Project updated: %s (ID: %s)\n", project.Name, project.ID)
	return project, nil
}

func (s *ProjectService) DeleteProject(projectID string) error {
	result := database.DB.Delete(&models.Project{}, "id = ?", projectID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete project: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("project not found or already deleted")
	}
	log.Printf("Project deleted: ID %s\n", projectID)
	return nil
}

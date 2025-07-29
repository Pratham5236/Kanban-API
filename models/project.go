package models

import "time"

type Project struct {
	ID             string    `json:"id" gorm:"primaryKey"`
	OrganizationID string    `json:"organization_id" gorm:"not null"`
	Name           string    `json:"name" gorm:"not null"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"not null"`

	Organization Organization `json:"-" gorm:"foreignKey:OrganizationID"` //NOTE this acts as a foreignKey in gorm syntax
}

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"omitempty,min=3,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

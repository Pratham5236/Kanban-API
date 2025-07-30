package models

import "time"

type Board struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	ProjectID   string    `json:"project_id" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null"`

	Project Project `json:"-" gorm:"foreignKey:ProjectID"`
}

type CreateBoardRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

type UpdateBoardRequest struct {
	Name        string `json:"name" binding:"omitempty,min=3,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

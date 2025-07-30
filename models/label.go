package models

import "time"

type Label struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null;unique"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

type CreateLabelRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=50"`
	Color string `json:"color" binding:"required,hexcolor"`
}

type UpdateLabelRequest struct {
	Name  *string `json:"name" binding:"omitempty,min=1,max=50"`
	Color *string `json:"color" binding:"omitempty,hexcolor"`
}
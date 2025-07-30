package models

import (
	"time"
)

type Card struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	ListID      string    `json:"list_id" gorm:"not null"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Notes       string    `json:"notes"`
	Position    int       `json:"position" gorm:"not null"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null"`

	List        List         `json:"-" gorm:"foreignKey:ListID"`
	Labels      []*Label     `json:"labels" gorm:"many2many:card_labels;"`
	Comments    []*Comment   `json:"comments" gorm:"foreignKey:CardID"`
	Attachments []*Attachment `json:"attachments" gorm:"foreignKey:CardID"`
}

type CreateCardRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=200"`
	Description string     `json:"description" binding:"omitempty,max=1000"`
	DueDate     *time.Time `json:"due_date" binding:"omitempty"`
}

type UpdateCardRequest struct {
	Title       string     `json:"title" binding:"omitempty,min=1,max=200"`
	Description string     `json:"description" binding:"omitempty,max=1000"`
	Position    *int       `json:"position" binding:"omitempty"`
	DueDate     *time.Time `json:"due_date" binding:"omitempty"`
	ListID      string     `json:"list_id" binding:"omitempty,uuid"`
}

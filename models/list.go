package models

import "time"

type List struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	BoardID   string    `json:"board_id" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	Position  int       `json:"position" gorm:"not null;default:0"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`

	Board Board `json:"board" gorm:"foreignKey:BoardID"`
}

type CreateListRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

type UpdateListRequest struct {
	Name     string `json:"name" binding:"omitempty,min=1,max=100"`
	Position *int   `json:"position" binding:"omitempty"`
}

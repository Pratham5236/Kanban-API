package models

import "time"

type Comment struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	CardID    string    `json:"card_id" gorm:"not null"`
	UserID    string    `json:"user_id" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

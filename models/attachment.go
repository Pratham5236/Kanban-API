package models

import "time"

type Attachment struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	CardID    string    `json:"card_id" gorm:"not null"`
	FileName  string    `json:"file_name" gorm:"not null"`
	FileURL   string    `json:"file_url" gorm:"not null"`
	FileType  string    `json:"file_type"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
}

type CreateAttachmentRequest struct {
	FileName string `json:"file_name" binding:"required"`
	FileURL  string `json:"file_url" binding:"required,url"`
	FileType string `json:"file_type"`
}
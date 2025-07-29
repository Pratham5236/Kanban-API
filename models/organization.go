package models

import "time"

// TODO add other properties to an 'organization' such as Github prefix, accounts etc

type Organization struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"unique;not null"`
	OwnerID   string    `json:"owner_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
}

type UpdateOrganizationRequest struct {
	Name string `json:"name" binding:"omitempty,min=3,max=100"`
}

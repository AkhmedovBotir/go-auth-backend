package models

import "time"

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	FullName     string    `json:"full_name" gorm:"size:120;not null"`
	Email        string    `json:"email" gorm:"size:120;uniqueIndex;not null"`
	Phone        string    `json:"phone" gorm:"size:30"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

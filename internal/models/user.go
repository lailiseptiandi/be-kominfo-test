package models

import "time"

type User struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	Name         string    `gorm:"not null" json:"name"`
	Email        string    `gorm:"not null, unique" json:"email"`
	PasswordHash string    `gorm:"not null" json:"password_hash"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

package models

import "time"

type User struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	Username           string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	PasswordHash       string    `json:"-" gorm:"size:255;not null"`
	DisplayName        string    `json:"display_name" gorm:"size:100;not null"`
	IsAdmin            bool      `json:"is_admin" gorm:"default:false"`
	MustChangePassword bool      `json:"must_change_password" gorm:"default:true"`
	CreatedAt          time.Time `json:"created_at"`
}

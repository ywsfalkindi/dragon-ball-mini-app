package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TelegramID   int64     `gorm:"uniqueIndex;not null" json:"telegram_id"`
	Username     string    `gorm:"size:255" json:"username"`
	FirstName    string    `gorm:"size:255" json:"first_name"`
	PhotoURL     string    `json:"photo_url"`
	
	// Game Stats
	Energy       int       `gorm:"default:100" json:"energy"`
	MaxEnergy    int       `gorm:"default:100" json:"max_energy"`
	TotalScore   int       `gorm:"default:0" json:"total_score"` // تم توحيد الاسم لـ TotalScore
	Rank         string    `gorm:"default:'Low Class Warrior'" json:"rank"`

	// Security (الدرع الفولاذي)
	RefreshToken string    `json:"-"` // لا يرسل للعميل أبداً (JSON Ignore)

	// Timestamps
	LastLoginAt  time.Time `json:"last_login_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
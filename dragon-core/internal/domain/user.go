package domain

import (
	"time"
)

// User يمثل اللاعب في النظام
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TelegramID   int64     `gorm:"uniqueIndex;not null" json:"telegram_id"`
	Username     string    `gorm:"size:255" json:"username"`
	FirstName    string    `gorm:"size:255" json:"first_name"`
	PhotoURL     string    `json:"photo_url"`
	
	// Game Stats
	Energy       int       `gorm:"default:100" json:"energy"`
	MaxEnergy    int       `gorm:"default:100" json:"max_energy"`
	Score        int64     `gorm:"default:0" json:"score"`
	Coins        int64     `gorm:"default:0" json:"coins"`
	
	// Referral System
	ReferralCode string    `gorm:"uniqueIndex;size:12" json:"referral_code"`
	ReferredBy   *int64    `json:"referred_by"`

	// Timestamps
	LastLoginAt  time.Time `json:"last_login_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserRepo هذا هو الاسم الذي يبحث عنه الكود (تأكد أنه UserRepo وليس UserRepository)
type UserRepo interface {
	Create(user *User) error
	GetByTelegramID(id int64) (*User, error)
	GetByReferralCode(code string) (*User, error)
	Update(user *User) error
	IncrementScore(userID uint, points int64) error
}
package models

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey"`
	TelegramID int64     `gorm:"unique;not null"`
	Username   string    `gorm:"size:255"`
	Energy     int       `gorm:"default:10"`     // الطاقة الحالية
	TotalScore int       `gorm:"default:0"`      // مجموع النقاط التراكمي
	Rank       string    `gorm:"default:'Low Class Warrior'"` // اللقب الحالي
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func (User) TableName() string {
	return "users"
}
package models
import "time"

type Score struct {
	ID        uint      `gorm:"primaryKey"`
	UserID     int       `gorm:"index;not null"`
	Points    int       `gorm:"not null"`
	ObtainedAt time.Time `gorm:"autoCreateTime"`
}
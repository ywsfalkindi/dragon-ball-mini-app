package models

import "time"

// Match: يمثل سجلاً لمعركة واحدة (مباراة)
type Match struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"` // ربط المباراة باللاعب
	Points    int64     `json:"points"`      // النقاط المكتسبة
	Correct   int       `json:"correct"`     // عدد الإجابات الصحيحة
	Total     int       `json:"total"`       // عدد الأسئلة الكلية
	PlayedAt  time.Time `gorm:"autoCreateTime" json:"played_at"` // متى لعبت؟
}

func (Match) TableName() string {
	return "matches"
}
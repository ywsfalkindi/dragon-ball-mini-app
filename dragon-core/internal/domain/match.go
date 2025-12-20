package domain

import "time"

// Match يمثل سجلاً لمعركة واحدة
type Match struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"` // ربط المباراة باللاعب
	Points    int64     `json:"points"`      // النقاط المكتسبة
	Correct   int       `json:"correct"`     // عدد الإجابات الصحيحة
	Total     int       `json:"total"`       // عدد الأسئلة الكلية
	PlayedAt  time.Time `json:"played_at"`   // متى لعبت؟
}

// MatchRepo تعاقد وظائف المباريات
type MatchRepo interface {
	Create(match *Match) error
	GetHistory(userID uint, limit int) (Match, error)
}
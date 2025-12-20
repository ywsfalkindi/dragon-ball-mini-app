package repository

import (
	"errors"
	"dragon-core/internal/domain"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

// NewUserRepo ينشئ نسخة جديدة من المستودع
func NewUserRepo(db *gorm.DB) domain.UserRepo {
	return &userRepo{db: db}
}

// Create ينشئ مستخدماً جديداً
func (r *userRepo) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// GetByTelegramID يبحث عن مستخدم بمعرف تليجرام
func (r *userRepo) GetByTelegramID(id int64) (*domain.User, error) {
	var user domain.User
	// SELECT * FROM users WHERE telegram_id =? LIMIT 1
	err := r.db.Where("telegram_id =?", id).First(&user).Error
	if err!= nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // لم يتم العثور عليه، ليس خطأ تقنياً
		}
		return nil, err
	}
	return &user, nil
}

// GetByReferralCode يبحث عن مستخدم بكود الدعوة
func (r *userRepo) GetByReferralCode(code string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("referral_code =?", code).First(&user).Error
	if err!= nil {
		return nil, err
	}
	return &user, nil
}

// Update يحفظ التعديلات
func (r *userRepo) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// IncrementScore يزيد نقاط المستخدم بشكل ذري (Atomic)
// هذه الطريقة تمنع مشاكل التزامن (Race Conditions)
func (r *userRepo) IncrementScore(userID uint, points int64) error {
	// UPDATE users SET score = score + points WHERE id = userID
	return r.db.Model(&domain.User{}).Where("id =?", userID).
		UpdateColumn("score", gorm.Expr("score +?", points)).Error
}
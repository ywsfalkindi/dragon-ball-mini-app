package repository

import (
	"dragon-core/internal/models"
	"errors"
	"gorm.io/gorm"
)

// هيكل المستودع
type userRepo struct {
	db *gorm.DB
}

// دالة الإنشاء
func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{db: db}
}

// إنشاء مستخدم جديد
func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// البحث عن مستخدم بـ Telegram ID
func (r *userRepo) GetByTelegramID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.Where("telegram_id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// تحديث بيانات المستخدم
func (r *userRepo) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// زيادة النقاط
func (r *userRepo) IncrementScore(userID uint, points int64) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("total_score", gorm.Expr("total_score + ?", points)).Error
}

func (r *userRepo) DecreaseEnergy(userID uint, amount int) (bool, error) {
	// هذا الاستعلام يقول:
	// "يا داتابيز، حاولي تحديث الطاقة بإنقاص قيمتها، 
	// لكن بشرط: أن تكون الطاقة الحالية أكبر من أو تساوي المبلغ المطلوب"
	result := r.db.Model(&models.User{}).
		Where("id = ? AND energy >= ?", userID, amount).
		UpdateColumn("energy", gorm.Expr("energy - ?", amount))

	if result.Error != nil {
		return false, result.Error
	}

	// RowsAffected: يخبرنا كم صف تأثر بالعملية؟
	// إذا كان 0، فهذا يعني أن الشرط (energy >= amount) لم يتحقق، أي أن الطاقة غير كافية
	if result.RowsAffected == 0 {
		return false, nil // فشلت العملية (طاقة غير كافية)
	}

	return true, nil // نجحت العملية
}
package repositories

import (
	"dragon-core/internal/database"
	"dragon-core/internal/models"
	"errors"
	"gorm.io/gorm"
)

// 1. مهارة الإنشاء: Create User
// تأخذ بيانات مستخدم وتخزنه في الداتابيز
func CreateUser(user *models.User) error {
	// النتيجة: db.Create
	// تمرر له مؤشر المستخدم (&user)
	result := database.DB.Create(user)
	
	if result.Error != nil {
		return result.Error // أعد الخطأ إذا فشل (مثلاً التليجرام ID مكرر)
	}
	return nil
}

// 2. مهارة البحث: Get User by Telegram ID
// نبحث عن المستخدم برقم تليجرام لنعرف هل هو مسجل أم لا
func GetUserByTelegramID(tgID int64) (*models.User, error) {
	var user models.User
	
	// الترجمة: ابحث في جدول المستخدمين، حيث telegram_id يساوي tgID، وأعطني أول نتيجة
	result := database.DB.Where("telegram_id = ?", tgID).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // لم يتم العثور عليه (ليس خطأ تقنياً، بل مجرد عدم وجود)
		}
		return nil, result.Error // خطأ حقيقي في الداتابيز
	}

	return &user, nil // وجدنا المقاتل!
}

// 3. مهارة التطوير: Update Score / Energy
// نستخدمها لتنقيص الطاقة أو زيادة النقاط
func UpdateUserEnergy(userID uint, newEnergy int) error {
	// Model(&models.User{}): نحدد الجدول
	// Where("id = ?", userID): نحدد الصف
	// Update("energy", newEnergy): نحدد العمود والقيمة الجديدة
	result := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("energy", newEnergy)
	
	return result.Error
}

// دالة لزيادة النقاط وتحديث الطاقة
func AddPoints(userID uint, points int) (int, error) {
	var user models.User
	
	// 1. البحث عن المستخدم
	if err := database.DB.First(&user, userID).Error; err != nil {
		return 0, err
	}

	// 2. تحديث البيانات
	user.Energy = user.Energy - 1 // خصم طاقة واحدة لكل محاولة
	if points > 0 {
		// لو أجاب صح، نعطيه نقاطاً ولكن لا نعيد الطاقة (قرار تصميمي)
		// يمكنك تعديل المنطق هنا كما تشاء
	}
	
	// سنفترض هنا نظاماً بسيطاً: الإجابة لا تزيد "عموداً" للنقاط في جدول المستخدمين حالياً
	// لكننا سنقوم بإنشاء سجل في جدول Scores (النتائج)
	score := models.Score{
		UserID: int(user.ID), // تحويل uint إلى int
		Points: points,
	}
	database.DB.Create(&score)

	// تحديث الطاقة في جدول المستخدم
	database.DB.Save(&user)

	return points, nil // نرجع النقاط الحالية (للبساطة)
}
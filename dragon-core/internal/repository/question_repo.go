package repository

import (
	"dragon-core/internal/database"
	"dragon-core/internal/models"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

func CacheAllQuestionIDs() error {
	var questions []models.Question
	// نجلب فقط الـ ID لتوفير الذاكرة
	if err := database.DB.Select("id").Find(&questions).Error; err != nil {
		return err
	}

	if len(questions) == 0 {
		return nil
	}

	// اسم المجموعة في Redis
	key := "questions:all_ids"
	
	// نضيف كل الـ IDs للمجموعة
	pipeline := database.RDB.Pipeline()
	pipeline.Del(database.Ctx, key) // تنظيف القديم
	for _, q := range questions {
		pipeline.SAdd(database.Ctx, key, q.ID)
	}
	_, err := pipeline.Exec(database.Ctx)
	return err
}

// 2. دالة ذكية: GetRandomQuestionForUser
// تختار سؤالاً لم يجب عليه المستخدم من قبل
func GetRandomQuestionForUser(userID uint) (*models.Question, error) {
	allQuestionsKey := "questions:all_ids"
	userAnsweredKey := fmt.Sprintf("user:%d:answered", userID)

	// SDiff: عملية سحرية تعطينا الفرق بين المجموعتين (الكل - المجاب)
	availableIDs, err := database.RDB.SDiff(database.Ctx, allQuestionsKey, userAnsweredKey).Result()
	if err != nil {
		return nil, err
	}

	if len(availableIDs) == 0 {
		return nil, fmt.Errorf("no more questions available! you beat the game")
	}

	// نختار ID عشوائي من المتبقي
	randomIndex := rand.Intn(len(availableIDs))
	randomID := availableIDs[randomIndex]

	// نجلب تفاصيل السؤال المختار (نستخدم الدالة القديمة للبحث في الكاش أو الداتابيز)
	// ملاحظة: نحتاج تحويل string ID إلى uint
	var qID uint
	fmt.Sscanf(randomID, "%d", &qID)
	
	return GetQuestionCached(qID)
}

// 3. دالة مساعدة: MarkQuestionAsAnswered
// نضيف السؤال لقائمة "المجاب" الخاصة بالمستخدم
func MarkQuestionAsAnswered(userID uint, questionID uint) {
	key := fmt.Sprintf("user:%d:answered", userID)
	database.RDB.SAdd(database.Ctx, key, questionID)
}

// جلب السؤال من الكاش أو الداتابيز
func GetQuestionCached(questionID uint) (*models.Question, error) {
	cacheKey := fmt.Sprintf("question:%d", questionID)

	// 1. محاولة القراءة من Redis
	val, err := database.RDB.Get(database.Ctx, cacheKey).Result()
	
    // إذا لم يكن هناك خطأ (يعني وجدنا البيانات)
	if err == nil {
		var question models.Question
		json.Unmarshal([]byte(val), &question)
		return &question, nil
	}

	// 2. القراءة من الداتابيز (في حال لم نجدها في الكاش أو حدث خطأ)
	var question models.Question
	result := database.DB.First(&question, questionID)
	if result.Error != nil {
		return nil, result.Error
	}

	// 3. الحفظ في Redis للمستقبل
	jsonData, _ := json.Marshal(question)
	database.RDB.Set(database.Ctx, cacheKey, jsonData, time.Hour)

	return &question, nil
}
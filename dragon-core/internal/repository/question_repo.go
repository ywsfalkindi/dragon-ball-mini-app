package repository

import (
	"dragon-core/internal/database"
	"dragon-core/internal/models"
	"encoding/json"
	"fmt"
	"time"
    // تم حذف مكتبة redis من هنا لأننا لم نستخدمها مباشرة
)

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
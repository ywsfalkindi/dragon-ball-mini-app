package repositories

import (
	"dragon-core/internal/database"
	"dragon-core/internal/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// دالة تجلب السؤال من الكاش، وإن لم تجده تجلبه من الداتابيز
func GetQuestionCached(questionID uint) (*models.Question, error) {
	// 1. تحديد مفتاح البحث (المفتاح يجب أن يكون مميزاً)
	// مثلاً: question:1, question:55
	cacheKey := fmt.Sprintf("question:%d", questionID)

	// 2. محاولة القراءة من Redis (The Fast Way)
	val, err := database.RDB.Get(database.Ctx, cacheKey).Result()
	
	if err == nil {
		// --- سيناريو: وجدنا البيانات في الكاش (HIT) ⚡ ---
		fmt.Println("⚡ CACHE HIT: Getting question from RAM")
		
		var question models.Question
		// تحويل النص المحفوظ في ريديس (JSON) ليرجع كائن Go
		json.Unmarshal([]byte(val), &question)
		return &question, nil
	} else if err != redis.Nil {
		// حدث خطأ تقني في ريديس نفسه
		return nil, err
	}

	// --- سيناريو: لم نجد البيانات (MISS) 🐢 ---
	fmt.Println("🐢 CACHE MISS: Going to PostgreSQL...")

	// 3. الذهاب لقاعدة البيانات (The Slow Way)
	var question models.Question
	result := database.DB.First(&question, questionID)
	if result.Error != nil {
		return nil, result.Error
	}

	// 4. الحفظ في Redis للمرة القادمة (Set)
	// نحول الكائن لنص JSON
	jsonData, _ := json.Marshal(question)
	
	// الحفظ لمدة ساعة واحدة (time.Hour)
	// بعد ساعة سيحذف ريديس المعلومة ليتم تجديدها
	database.RDB.Set(database.Ctx, cacheKey, jsonData, time.Hour)

	return &question, nil
}
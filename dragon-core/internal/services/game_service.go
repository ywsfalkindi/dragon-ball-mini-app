package services

import (
	"dragon-core/internal/database"
	"dragon-core/internal/models"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	CostPerGame    = 1
	BaseScore      = 100
	MaxTimeSeconds = 30 // الحد الأقصى المسموح به للوقت
)

// ProcessAnswer: لم نعد نأخذ timeTaken من العميل!
func ProcessAnswer(userID uint, questionID uint, selectedOption string) (*models.AnswerResponse, error) {
	var user models.User
	var question models.Question

	// 1. جلب المستخدم
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("fighter not found")
	}

	// 2. التحقق من الطاقة
	if user.Energy < CostPerGame {
		return nil, errors.New("out of stamina! need senzu bean")
	}

	// 3. 🛡️ Security: حساب الوقت الحقيقي من السيرفر
	timerKey := fmt.Sprintf("game:timer:%d:%d", userID, questionID)
	startTimeStr, err := database.RDB.Get(database.Ctx, timerKey).Result()
	
	var timeTakenSeconds float64
	
	if err == redis.Nil {
		return nil, errors.New("session expired or invalid cheat attempt")
	} else if err != nil {
		return nil, err
	} else {
		// --- التصحيح هنا ---
		var startTime int64 // تم إضافة var
		fmt.Sscanf(startTimeStr, "%d", &startTime)
		
		now := time.Now().UnixMilli()
		diffMillis := now - startTime
		
		// تحويل لثواني
		timeTakenSeconds = float64(diffMillis) / 1000.0
		
		// حذف المفتاح لمنع الإجابة مرتين على نفس السؤال
		database.RDB.Del(database.Ctx, timerKey)
	}

	// 4. جلب السؤال للتصحيح
	if err := database.DB.First(&question, questionID).Error; err != nil {
		return nil, errors.New("question not found")
	}

	// 5. خصم الطاقة
	user.Energy -= CostPerGame

	isCorrect := (selectedOption == question.CorrectOption)
	pointsEarned := 0
	message := "You missed! 💥"

	if isCorrect {
		// أ) الصعوبة
		difficultyMultiplier := 1.0
		if question.Difficulty == 2 { difficultyMultiplier = 1.5 }
		if question.Difficulty == 3 { difficultyMultiplier = 2.0 }

		// ب) السرعة
		if timeTakenSeconds > MaxTimeSeconds {
			timeTakenSeconds = MaxTimeSeconds
		}
		if timeTakenSeconds < 0 {
			timeTakenSeconds = 0
		}

		timeSaved := float64(MaxTimeSeconds) - timeTakenSeconds
		speedBonus := timeSaved * 10 

		calcScore := (float64(BaseScore) * difficultyMultiplier) + speedBonus
		pointsEarned = int(math.Ceil(calcScore))

		user.TotalScore += pointsEarned
		user.Rank = calculateRank(user.TotalScore)
		message = fmt.Sprintf("Perfect! Time: %.1fs ⚡", timeTakenSeconds)
	}

	database.DB.Save(&user)

	// تسجيل المحاولة
	history := models.Score{
		UserID: int(user.ID),
		Points: pointsEarned,
	}
	database.DB.Create(&history)

	return &models.AnswerResponse{
		Correct:   isCorrect,
		Message:   message,
		NewScore:  user.TotalScore,
		NewEnergy: user.Energy,
	}, nil
}

func calculateRank(score int) string {
	if score >= 100000 { return "Angel 😇" }
	if score >= 20000 { return "God of Destruction 🟣" }
	if score >= 5000 { return "Super Saiyan 👱" }
	if score >= 1000 { return "Elite Warrior 👮" }
	return "Low Class Warrior 👨‍🌾"
}
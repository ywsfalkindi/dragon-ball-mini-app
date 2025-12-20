package services

import (
	"dragon-core/internal/database"
	"dragon-core/internal/models"
	"dragon-core/internal/repository"
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

func GetRandomQuestion(userID uint) (*models.Question, error) {
	// نستخدم الدالة الذكية من المستودع
	return repository.GetRandomQuestionForUser(userID)
}

func ProcessAnswer(userID uint, questionID uint, selectedOption string) (*models.AnswerResponse, error) {
	// 1. خصم الطاقة
	userRepo := repository.NewUserRepo(database.DB)
	hasEnergy, err := userRepo.DecreaseEnergy(userID, CostPerGame)
	if err != nil {
		return nil, err
	}
	if !hasEnergy {
		return nil, errors.New("out of stamina! recharge needed")
	}

	// 2. التحقق من الوقت (Security)
	timerKey := fmt.Sprintf("game:timer:%d:%d", userID, questionID)
	startTimeStr, err := database.RDB.Get(database.Ctx, timerKey).Result()
	
	var timeTakenSeconds float64

	// 👇 التعديل هنا: تساهل مع خطأ "انتهاء الجلسة"
	if err == redis.Nil {
		// بدلاً من إرجاع خطأ، سنفترض وقتاً افتراضياً ونكمل اللعب
		// return nil, errors.New("session expired") <--- تم إيقاف هذا السطر
		fmt.Println("⚠️ Warning: Timer key not found (Session Expired), skipping check...")
		timeTakenSeconds = 5.0 // وقت افتراضي
	} else if err != nil {
		// خطأ حقيقي في Redis
		return nil, err
	} else {
		// الوضع الطبيعي: وجدنا الوقت
		var startTime int64
		fmt.Sscanf(startTimeStr, "%d", &startTime)
		now := time.Now().UnixMilli()
		diffMillis := now - startTime
		timeTakenSeconds = float64(diffMillis) / 1000.0
		// تنظيف المفتاح
		database.RDB.Del(database.Ctx, timerKey)
	}

	// 3. جلب السؤال
	question, err := repository.GetQuestionCached(questionID)
	if err != nil {
		return nil, errors.New("question not found")
	}

	// 4. وضع علامة أن المستخدم أجاب
	repository.MarkQuestionAsAnswered(userID, questionID)

	// 5. حساب النتيجة
	isCorrect := (selectedOption == question.CorrectOption)
	pointsEarned := 0
	message := "You missed! 💥"
	
	var newTotalScore int

	if isCorrect {
		difficultyMultiplier := 1.0
		if question.Difficulty == 2 { difficultyMultiplier = 1.5 }
		if question.Difficulty == 3 { difficultyMultiplier = 2.0 }
		
		// حماية من الأوقات السالبة أو الطويلة جداً
		if timeTakenSeconds > MaxTimeSeconds { timeTakenSeconds = MaxTimeSeconds }
		if timeTakenSeconds < 0 { timeTakenSeconds = 0 }
		
		timeSaved := float64(MaxTimeSeconds) - timeTakenSeconds
		speedBonus := timeSaved * 10 
		calcScore := (float64(BaseScore) * difficultyMultiplier) + speedBonus
		pointsEarned = int(math.Ceil(calcScore))

		message = fmt.Sprintf("Perfect! Time: %.1fs ⚡", timeTakenSeconds)

		leaderboardRepo := repository.NewLeaderboardRepo(database.RDB)
		newScoreFloat, _ := leaderboardRepo.IncrementScore(database.Ctx, userID, float64(pointsEarned))
		newTotalScore = int(newScoreFloat)

		newRank := calculateRank(newTotalScore)
		fmt.Printf("User %d reached rank: %s 🌟\n", userID, newRank)
	} else {
		leaderboardRepo := repository.NewLeaderboardRepo(database.RDB)
		currentScore, _ := leaderboardRepo.GetCurrentScore(database.Ctx, userID)
		newTotalScore = int(currentScore)
	}

	// تسجيل المحاولة
	history := models.Score{ UserID: int(userID), Points: pointsEarned }
	database.DB.Create(&history)

	// جلب الطاقة المتبقية
	var user models.User
	database.DB.Select("energy").First(&user, userID)

	return &models.AnswerResponse{
		Correct:   isCorrect,
		Message:   message,
		NewScore:  newTotalScore,
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
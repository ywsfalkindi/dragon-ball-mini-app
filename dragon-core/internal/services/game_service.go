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
	// 1. خصم الطاقة (كما فعلنا في الفصل 2 - ممتاز)
	userRepo := repository.NewUserRepo(database.DB)
	hasEnergy, err := userRepo.DecreaseEnergy(userID, CostPerGame)
	if err != nil { return nil, err }
	if !hasEnergy { return nil, errors.New("out of stamina!") }

	// 2. التحقق من الوقت (Security)
	timerKey := fmt.Sprintf("game:timer:%d:%d", userID, questionID)
	startTimeStr, err := database.RDB.Get(database.Ctx, timerKey).Result()
	
	var timeTakenSeconds float64
	// ... (نفس منطق الوقت السابق ) ...
	if err == redis.Nil {
		return nil, errors.New("session expired")
	} else if err != nil {
		return nil, err
	} else {
		var startTime int64
		fmt.Sscanf(startTimeStr, "%d", &startTime)
		now := time.Now().UnixMilli()
		diffMillis := now - startTime
		timeTakenSeconds = float64(diffMillis) / 1000.0
		database.RDB.Del(database.Ctx, timerKey)
	}

	// 3. جلب السؤال
	// هنا نستخدم GetQuestionCached لأننا نعرف الـ ID مسبقاً من الطلب
	question, err := repository.GetQuestionCached(questionID)
	if err != nil { return nil, errors.New("question not found") }

	// 4. وضع علامة أن المستخدم أجاب على هذا السؤال
	// لكي لا يظهر له مرة أخرى في GetRandomQuestion
	repository.MarkQuestionAsAnswered(userID, questionID)

	// 5. حساب النتيجة
	isCorrect := (selectedOption == question.CorrectOption)
	pointsEarned := 0
	message := "You missed! 💥"
	
	var newTotalScore int

	if isCorrect {
		// ... (حساب النقاط والسرعة كما هو) ...
		difficultyMultiplier := 1.0
		if question.Difficulty == 2 { difficultyMultiplier = 1.5 }
		if question.Difficulty == 3 { difficultyMultiplier = 2.0 }
		if timeTakenSeconds > MaxTimeSeconds { timeTakenSeconds = MaxTimeSeconds }
		if timeTakenSeconds < 0 { timeTakenSeconds = 0 }
		timeSaved := float64(MaxTimeSeconds) - timeTakenSeconds
		speedBonus := timeSaved * 10 
		calcScore := (float64(BaseScore) * difficultyMultiplier) + speedBonus
		pointsEarned = int(math.Ceil(calcScore))

		message = fmt.Sprintf("Perfect! Time: %.1fs ⚡", timeTakenSeconds)

		// تحديث Redis (Write-Behind)
		leaderboardRepo := repository.NewLeaderboardRepo(database.RDB)
		newScoreFloat, _ := leaderboardRepo.IncrementScore(database.Ctx, userID, float64(pointsEarned))
		newTotalScore = int(newScoreFloat)

		// --- التصحيح هنا: استخدام الدالة المنسية ---
		// نحسب الرتبة الجديدة ونضعها في رسالة اللوج (أو يمكن إعادتها للمستخدم لاحقاً)
		newRank := calculateRank(newTotalScore)
		fmt.Printf("User %d reached rank: %s 🌟\n", userID, newRank)

		// ملاحظة: لا نحتاج لعمل db.Save(&user) للنقاط هنا! العامل سيقوم بذلك.
		// لكن، إذا أردنا تحديث الـ Rank في الواجهة، نستخدم المجموع الجديد من Redis.
	} else {
		// في حال الخسارة، نجلب السكور الحالي من Redis للعرض فقط
		leaderboardRepo := repository.NewLeaderboardRepo(database.RDB)
		currentScore, _ := leaderboardRepo.GetCurrentScore(database.Ctx, userID)
		newTotalScore = int(currentScore)
	}

	// تسجيل المحاولة في الأرشيف (هذه يمكن أن تبقى مباشرة لأنها Log)
	// أو يمكن أيضاً وضعها في طابور (Queue) لتحسين الأداء أكثر، لكن سنكتفي بهذا القدر حالياً
	history := models.Score{ UserID: int(userID), Points: pointsEarned }
	database.DB.Create(&history)

	// جلب الطاقة المتبقية للعرض
	var user models.User
	database.DB.Select("energy").First(&user, userID)

	return &models.AnswerResponse{
		Correct:   isCorrect,
		Message:   message,
		NewScore:  newTotalScore, // السكور القادم من Redis
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
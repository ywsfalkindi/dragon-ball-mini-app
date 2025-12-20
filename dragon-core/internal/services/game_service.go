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

func ProcessAnswer(userID uint, questionID uint, selectedOption string) (*models.AnswerResponse, error) {
	// 1. 🛡️ Security (Atomic Check): الخندق الدفاعي الأول
	// بدلاً من جلب المستخدم وقراءة طاقته، نحاول خصم الطاقة مباشرة
	userRepo := repository.NewUserRepo(database.DB)
	
	// نحاول خصم 1 طاقة. الدالة سترجع false إذا لم يكن لديه طاقة كافية
	hasEnergy, err := userRepo.DecreaseEnergy(userID, CostPerGame)
	if err != nil {
		return nil, err // خطأ في الداتابيز
	}
	if !hasEnergy {
		return nil, errors.New("out of stamina! need senzu bean") // لا توجد طاقة
	}

	// 2. 🛡️ Security (Time Check): الخندق الدفاعي الثاني
	timerKey := fmt.Sprintf("game:timer:%d:%d", userID, questionID)
	startTimeStr, err := database.RDB.Get(database.Ctx, timerKey).Result()
	
	var timeTakenSeconds float64
	if err == redis.Nil {
		// المستخدم خسر الطاقة التي خصمناها للتو لأنه حاول الغش!
		// (يمكننا إعادتها له إذا كنا لطفاء، لكن لنجعله درساً له حالياً)
		return nil, errors.New("session expired or invalid cheat attempt")
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

	// 3. جلب السؤال للتصحيح
	var question models.Question
	if err := database.DB.First(&question, questionID).Error; err != nil {
		return nil, errors.New("question not found")
	}

	// منطق اللعبة (Game Logic)
	isCorrect := (selectedOption == question.CorrectOption)
	pointsEarned := 0
	message := "You missed! 💥"
	
	// نحتاج جلب بيانات المستخدم الآن فقط لتحديث الـ Score والـ Rank ولنعرض له طاقته المتبقية
	// (لاحظ: نحن نجلب المستخدم *بعد* خصم الطاقة بنجاح)
	var user models.User
	database.DB.First(&user, userID)

	if isCorrect {
		// حساب النقاط (كما هو في كودك السابق)
		difficultyMultiplier := 1.0
		if question.Difficulty == 2 { difficultyMultiplier = 1.5 }
		if question.Difficulty == 3 { difficultyMultiplier = 2.0 }

		if timeTakenSeconds > MaxTimeSeconds { timeTakenSeconds = MaxTimeSeconds }
		if timeTakenSeconds < 0 { timeTakenSeconds = 0 }

		timeSaved := float64(MaxTimeSeconds) - timeTakenSeconds
		speedBonus := timeSaved * 10 
		calcScore := (float64(BaseScore) * difficultyMultiplier) + speedBonus
		pointsEarned = int(math.Ceil(calcScore))

		// تحديث النقاط والرتبة
		user.TotalScore += pointsEarned
		user.Rank = calculateRank(user.TotalScore)
		message = fmt.Sprintf("Perfect! Time: %.1fs ⚡", timeTakenSeconds)
		
		database.DB.Save(&user) // حفظ النقاط الجديدة
	}

	// تسجيل المحاولة في الأرشيف
	history := models.Score{
		UserID: int(user.ID),
		Points: pointsEarned,
	}
	database.DB.Create(&history)

	return &models.AnswerResponse{
		Correct:   isCorrect,
		Message:   message,
		NewScore:  user.TotalScore,
		NewEnergy: user.Energy, // هذه القيمة تم تحديثها ذرياً في الخطوة 1
	}, nil
}

func calculateRank(score int) string {
	if score >= 100000 { return "Angel 😇" }
	if score >= 20000 { return "God of Destruction 🟣" }
	if score >= 5000 { return "Super Saiyan 👱" }
	if score >= 1000 { return "Elite Warrior 👮" }
	return "Low Class Warrior 👨‍🌾"
}
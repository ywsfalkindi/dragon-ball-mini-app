package services

import (
	"dragon-core/internal/database"
	"dragon-core/internal/models"
	"errors"
	"math"
)

// ثوابت اللعبة
const (
	CostPerGame    = 1   // تكلفة الطاقة لكل سؤال
	BaseScore      = 100 // النقاط الأساسية
	MaxTimeSeconds = 30  // الوقت الأقصى للسؤال
)

// ProcessAnswer: الدالة الكبرى التي تدير عملية الإجابة كاملة
func ProcessAnswer(userID uint, questionID uint, selectedOption string, timeTakenSeconds int) (*models.AnswerResponse, error) {
	var user models.User
	var question models.Question

	// 1. جلب بيانات المستخدم
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("fighter not found")
	}

	// 2. التحقق من الطاقة (هل لديه كي كافٍ؟)
	if user.Energy < CostPerGame {
		return nil, errors.New("out of stamina! eat a senzu bean")
	}

	// 3. جلب السؤال (من الداتابيز مباشرة للدقة أو الكاش)
	if err := database.DB.First(&question, questionID).Error; err != nil {
		return nil, errors.New("question scroll missing")
	}

	// 4. خصم الطاقة فوراً (سواء أجاب صح أم خطأ)
	user.Energy -= CostPerGame

	// 5. التحقق من صحة الإجابة
	isCorrect := (selectedOption == question.CorrectOption)
	pointsEarned := 0
	message := "You missed! Frieza destroys the planet. 💥"

	if isCorrect {
		// --- هنا تبدأ الرياضيات ---
		
		// أ) حساب مضاعف الصعوبة
		difficultyMultiplier := 1.0
		switch question.Difficulty {
		case 2:
			difficultyMultiplier = 1.5 // متوسط
		case 3:
			difficultyMultiplier = 2.0 // صعب
		}

		// ب) حساب نقاط السرعة
		// نضمن أن الوقت لا يتجاوز 30 ولا يقل عن 0
		if timeTakenSeconds > MaxTimeSeconds {
			timeTakenSeconds = MaxTimeSeconds
		}
		timeSaved := MaxTimeSeconds - timeTakenSeconds
		speedBonus := timeSaved * 10 // 10 نقاط لكل ثانية موفرة

		// ج) المعادلة النهائية
		// Score = (100 * Diff) + SpeedBonus
		calcScore := (float64(BaseScore) * difficultyMultiplier) + float64(speedBonus)
		pointsEarned = int(math.Ceil(calcScore)) // تقريب الرقم لأعلى

		// تحديث مجموع نقاط اللاعب
		user.TotalScore += pointsEarned
		
		// تحديث الرتبة (الترقية)
		user.Rank = calculateRank(user.TotalScore)
		
		message = "Perfect Hit! 🎯"
	}

	// 6. حفظ التغييرات في قاعدة البيانات
	database.DB.Save(&user)

	// 7. تسجيل المحاولة في سجل النتائج (History)
	history := models.Score{
		UserID: int(user.ID),
		Points: pointsEarned,
	}
	database.DB.Create(&history)

	// 8. تجهيز الرد
	response := &models.AnswerResponse{
		Correct:   isCorrect,
		Message:   message,
		NewScore:  user.TotalScore, // نرجع المجموع الكلي
		NewEnergy: user.Energy,
	}

	return response, nil
}

// دالة داخلية مساعدة لحساب الرتبة
func calculateRank(score int) string {
	if score >= 100000 {
		return "Angel 😇"
	} else if score >= 20000 {
		return "God of Destruction 🟣"
	} else if score >= 5000 {
		return "Super Saiyan 👱"
	} else if score >= 1000 {
		return "Elite Warrior 👮"
	}
	return "Low Class Warrior 👨‍🌾"
}
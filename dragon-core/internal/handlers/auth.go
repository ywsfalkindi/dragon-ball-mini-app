package handlers

import (
	"dragon-core/internal/auth"
	"dragon-core/internal/database"
	"dragon-core/internal/models"
	"encoding/json"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	InitData string `json:"init_data"`
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// 1. الأمان: التحقق من التوقيع ومنع التكرار
	botToken := os.Getenv("BOT_TOKEN")
	isValid, err := auth.ValidateWebAppData(req.InitData, botToken)
	if err != nil || !isValid {
		return c.Status(401).JSON(fiber.Map{"error": "Security Alert: " + err.Error()})
	}

	// 2. استخراج البيانات
	parsedData, _ := url.ParseQuery(req.InitData)
	userJSON := parsedData.Get("user")
	type TelegramUser struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
		PhotoURL  string `json:"photo_url"`
	}
	var tgUser TelegramUser
	json.Unmarshal([]byte(userJSON), &tgUser)

	// 3. الداتابيز: إيجاد أو إنشاء المقاتل
	var user models.User
	result := database.DB.Where("telegram_id = ?", tgUser.ID).First(&user)

	if result.Error != nil {
		// مستخدم جديد
		user = models.User{
			TelegramID: tgUser.ID,
			Username:   tgUser.Username,
			FirstName:  tgUser.FirstName,
			PhotoURL:   tgUser.PhotoURL,
			Energy:     10, // طاقة البداية
			LastLoginAt: time.Now(),
		}
		database.DB.Create(&user)
	} else {
		// تحديث وقت الدخول والصورة
		user.LastLoginAt = time.Now()
		user.PhotoURL = tgUser.PhotoURL
		user.FirstName = tgUser.FirstName
		// لا نحفظ هنا لأننا سنحفظ التوكن في الخطوة التالية
	}

	// 4. توليد المفاتيح (Access & Refresh)
	accessToken, refreshToken, err := auth.GenerateTokens(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Token generation failed"})
	}

	// 5. حفظ الـ Refresh Token في الداتابيز (للطرد عند الحاجة)
	user.RefreshToken = refreshToken
	database.DB.Save(&user)

	// 6. الرد النهائي
	return c.JSON(fiber.Map{
		"access_token":  accessToken,  // يستخدم للطلبات (صلاحية 15 دقيقة)
		"refresh_token": refreshToken, // خبه في الجهاز لتجديد الجلسة
		"user":          user,
	})
}
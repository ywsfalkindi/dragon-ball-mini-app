package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// مفتاح التشفير السري الخاص بالسيرفر
// ناخذه من متغيرات البيئة، وإذا لم نجد، نستخدم مفتاح احتياطي (فقط للتطوير)
var jwtSecret = []byte(getSecretEnv())

func getSecretEnv() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "super-secret-dragon-ball-key-change-me" // ⚠️ يجب تغييره في الإنتاج
	}
	return secret
}

// هذه البيانات التي سنخبئها داخل التوكن
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken: دالة تصنع التوكن للمستخدم
func GenerateToken(userID uint) (string, error) {
	// تحديد صلاحية التوكن (مثلاً 7 أيام)
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	// تعبئة البيانات
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// التوقيع على التوكن باستخدام خوارزمية HS256 ومفتاحنا السري
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken: دالة تفحص التوكن وتستخرج منه هوية المستخدم
func ValidateToken(tokenString string) (uint, error) {
	claims := &Claims{}

	// محاولة فك التشفير
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	// نجاح! أرجع رقم المستخدم
	return claims.UserID, nil
}
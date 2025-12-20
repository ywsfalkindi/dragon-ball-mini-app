package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"dragon-core/internal/database" // نحتاج للوصول لـ Redis
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ValidateWebAppData: يتحقق من التوقيع + الوقت + عدم التكرار
func ValidateWebAppData(initData string, botToken string) (bool, error) {
	// 1. Parsing
	parsedData, err := url.ParseQuery(initData)
	if err != nil {
		return false, fmt.Errorf("error parsing data")
	}

	// 2. Time Check (صلاحية 24 ساعة)
	authDateStr := parsedData.Get("auth_date")
	if authDateStr == "" {
		return false, fmt.Errorf("auth_date is missing")
	}
	authDate, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		return false, fmt.Errorf("invalid auth_date format")
	}
	if time.Now().Unix()-authDate > 86400 {
		return false, fmt.Errorf("data is expired (older than 24h)")
	}

	// 3. Replay Attack Check (الجديد كلياً!) 🛡️
	receivedHash := parsedData.Get("hash")
	// مفتاح مميز في Redis لهذه العملية
	replayKey := fmt.Sprintf("auth:replay:%s", receivedHash)
	
	// نحاول الحفظ في Redis. إذا كان المفتاح موجوداً مسبقاً، فهذا هجوم إعادة!
	// SetNX يحفظ القيمة فقط إذا لم تكن موجودة (Not Exist)
	isUnique, err := database.RDB.SetNX(database.Ctx, replayKey, "used", 24*time.Hour).Result()
	if err != nil {
		// خطأ في اتصال Redis، نعتبره فشلاً أمنياً
		return false, fmt.Errorf("security check failed (redis error)")
	}
	if !isUnique {
		// تم استخدام هذا الهاش سابقاً! 🚫
		return false, fmt.Errorf("replay attack detected! this data was used before")
	}

	// 4. Hash Validation (كما كان سابقاً)
	parsedData.Del("hash")
	var keys []string
	for k := range parsedData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var dataCheckArr []string
	for _, k := range keys {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%s", k, parsedData.Get(k)))
	}
	dataCheckString := strings.Join(dataCheckArr, "\n")

	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(botToken))
	secret := secretKey.Sum(nil)

	h := hmac.New(sha256.New, secret)
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	if calculatedHash == receivedHash {
		return true, nil
	}
	return false, nil
}
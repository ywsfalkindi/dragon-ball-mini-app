package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv" // نحتاجه لتحويل الوقت من نص لرقم
	"strings"
	"time"    // نحتاجه لمعرفة الوقت الحالي
)

// ValidateWebAppData: الآن تتحقق من التوقيع ومن الزمن أيضاً!
func ValidateWebAppData(initData string, botToken string) (bool, error) {
	// 1. تحليل النص (Parsing)
	parsedData, err := url.ParseQuery(initData)
	if err != nil {
		return false, fmt.Errorf("error parsing data")
	}

	// 2. التحقق من تاريخ الصلاحية (The Time Check) 🕒
	// تليجرام يرسل الوقت بصيغة Unix Timestamp (رقم طويل بالثواني)
	authDateStr := parsedData.Get("auth_date")
	if authDateStr == "" {
		return false, fmt.Errorf("auth_date is missing")
	}

	authDate, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		return false, fmt.Errorf("invalid auth_date format")
	}

	// الوقت الحالي
	now := time.Now().Unix()
	
	// المعادلة: إذا كان الفرق بين الآن ووقت التليجرام أكثر من 24 ساعة (86400 ثانية)
	// فهذه البيانات قديمة ومتعفنة! 🧟‍♂️
	if now-authDate > 86400 {
		return false, fmt.Errorf("data is expired (older than 24h)")
	}

	// 3. استخراج الهاش والتحقق منه (كما كان سابقاً) [cite: 12]
	receivedHash := parsedData.Get("hash")
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

	// حساب التوقيع السري
	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(botToken))
	secret := secretKey.Sum(nil)

	h := hmac.New(sha256.New, secret)
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	if calculatedHash == receivedHash {
		return true, nil // ✅ متطابق وجديد!
	}
	return false, nil // ❌ مزور!
}
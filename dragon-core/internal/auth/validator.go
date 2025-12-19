package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// ValidateWebAppData: دالة تتحقق هل البيانات قادمة من تليجرام فعلاً أم مزورة
// initData: النص المشفر القادم من تليجرام
// botToken: مفتاح البوت الخاص بك
func ValidateWebAppData(initData string, botToken string) (bool, error) {
	// 1. تحليل النص (Parsing)
	// تحويل النص إلى أزواج (key=value)
	parsedData, err := url.ParseQuery(initData)
	if err != nil {
		return false, fmt.Errorf("error parsing data")
	}

	// 2. استخراج الهاش (التوقيع) وحذفه من القائمة
	// لأننا سنعيد حساب الهاش للباقي ونقارنه بهذا
	receivedHash := parsedData.Get("hash")
	parsedData.Del("hash")

	// 3. ترتيب البيانات أبجدياً (شرط أساسي من تليجرام)
	var keys []string
	for k := range parsedData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 4. بناء "سلسلة التحقق" (Data Check String)
	// نجمع البيانات المتبقية في نص واحد: key=value\nkey=value...
	var dataCheckArr []string
	for _, k := range keys {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%s", k, parsedData.Get(k)))
	}
	dataCheckString := strings.Join(dataCheckArr, "\n")

	// 5. حساب المفتاح السري (Secret Key)
	// المفتاح = HMAC_SHA256("WebAppData", BotToken)
	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(botToken))
	secret := secretKey.Sum(nil)

	// 6. حساب التوقيع النهائي (Signature)
	// التوقيع = HMAC_SHA256(dataCheckString, SecretKey)
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(dataCheckString))
	calculatedHash := hex.EncodeToString(h.Sum(nil))

	// 7. المقارنة النهائية (لحظة الحقيقة)
	if calculatedHash == receivedHash {
		return true, nil // ✅ متطابق!
	}
	return false, nil // ❌ مزور!
}
package models

// AnswerRequest: هذا ما سيرسله تطبيق React لنا
type AnswerRequest struct {
	UserID     uint   `json:"user_id"`     // من الذي يجيب؟
	QuestionID uint   `json:"question_id"` // أي سؤال يجيب عليه؟
	Selected   string `json:"selected"`    // ماذا اختار؟ (A, B, C, D)
}

// AnswerResponse: هذا ما سنرده عليه
type AnswerResponse struct {
	Correct   bool   `json:"correct"`    // صح أم خطأ؟
	Message   string `json:"message"`    // رسالة تشجيعية (مثل: أحسنت يا وحش!)
	NewScore  int    `json:"new_score"`  // نقاطه الجديدة
	NewEnergy int    `json:"new_energy"` // طاقته المتبقية
}
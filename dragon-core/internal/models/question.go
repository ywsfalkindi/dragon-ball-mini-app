package models

type Question struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	QuestionText  string `json:"question_text"`
	OptionA       string `json:"option_a"`
	OptionB       string `json:"option_b"`
	OptionC       string `json:"option_c"`
	OptionD       string `json:"option_d"`
	CorrectOption string `json:"-"` // العلامة - تعني لا ترسل الإجابة للمستخدم (سرية)
	Difficulty    int    `json:"difficulty" gorm:"default:1"` // 1: Easy, 2: Medium, 3: Hard
}
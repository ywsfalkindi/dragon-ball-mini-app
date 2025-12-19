package models

// هذا الهيكل يمثل شكل الرد الذي سيصل للمستخدم
// json:"message" تعني: عند تحويله لنص، سمّي الحقل message
type JSend struct {
	Status  string      `json:"status"` // success, fail, error
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"` // interface{} تعني أي شيء (مرنة)
}
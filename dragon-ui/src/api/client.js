import axios from "axios";

// إنشاء نسخة مخصصة من axios بإعداداتنا
export const apiClient = axios.create({
  baseURL: "http://localhost:3000/api", // عنوان السيرفر (Backend)
  headers: {
    "Content-Type": "application/json",
  },
});

// ملاحظة: لاحقاً سنضيف هنا الـ Interceptors لوضع توكن الحماية تلقائياً

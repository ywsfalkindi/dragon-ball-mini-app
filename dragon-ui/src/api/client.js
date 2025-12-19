import axios from "axios";
import WebApp from "@twa-dev/sdk"; // تأكد من استيراد هذا

export const apiClient = axios.create({
  baseURL: "http://localhost:3000/api",
  headers: {
    "Content-Type": "application/json",
    "ngrok-skip-browser-warning": "true",
  },
});

// --- هذا هو الجزء الناقص ---
apiClient.interceptors.request.use((config) => {
  // 1. محاولة جلب بيانات تليجرام الحقيقية
  const initData = WebApp.initData;

  if (initData) {
    // إذا كنا داخل تليجرام، أرسل التوكن الحقيقي
    config.headers.Authorization = initData;
  } else {
    // 2. وضع المطور (Localhost):
    // إذا كنا نختبر في المتصفح، نرسل توكن "وهمي" لكي لا يرفضنا السيرفر
    // ملاحظة: ستحتاج لتعطيل التحقق في السيرفر مؤقتاً لقبول هذا، أو استخدام بيانات حقيقية
    config.headers.Authorization = "test-token-for-goku";
  }

  return config;
});

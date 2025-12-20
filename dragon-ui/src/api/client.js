import axios from "axios";

export const apiClient = axios.create({
  baseURL: "http://localhost:3000/api", // تأكد أن هذا الرابط صحيح
  headers: {
    "Content-Type": "application/json",
    "ngrok-skip-browser-warning": "true",
  },
});

// Interceptor: يعترض كل طلب قبل خروجه ويضيف التوكن
apiClient.interceptors.request.use(
  (config) => {
    // 1. نحاول جلب التوكن المحفوظ في المتصفح
    const token = localStorage.getItem("dragon_token");

    // 2. إذا وجدنا توكن، والطلب ليس "تسجيل دخول"، نضيفه في الهيدر
    if (token && !config.url.includes("/auth/login")) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    // 3. (اختياري) وضع المطور: توكن وهمي إذا لم نجد توكن حقيقي وكنا في وضع التطوير
    else if (!token && import.meta.env.DEV) {
      // config.headers.Authorization = "test-token-for-goku"; // يمكنك تفعيل هذا للتجربة السريعة
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

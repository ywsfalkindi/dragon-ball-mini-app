import axios from "axios";

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "http://localhost:3000/api",
  headers: {
    "Content-Type": "application/json",
  },
});

// 1. Request Interceptor: إضافة التوكن لكل طلب
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("dragon_token");
    if (token && !config.url.includes("/auth/")) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// 2. Response Interceptor: التعامل مع انتهاء التوكن (401)
apiClient.interceptors.response.use(
  (response) => response, // إذا نجح الطلب، مرره كما هو
  async (error) => {
    const originalRequest = error.config;

    // إذا كان الخطأ 401 (غير مصرح) ولم نقم بإعادة المحاولة بعد
    if (
      error.response &&
      error.response.status === 401 &&
      !originalRequest._retry
    ) {
      originalRequest._retry = true; // علامة لمنع الدوران اللانهائي

      try {
        // محاولة جلب التوكن الطويل
        const refreshToken = localStorage.getItem("dragon_refresh_token");
        if (!refreshToken) throw new Error("No refresh token");

        // طلب تجديد التوكن من السيرفر
        const res = await axios.post("http://localhost:3000/api/auth/refresh", {
          refresh_token: refreshToken,
        });

        if (res.status === 200) {
          const { access_token, refresh_token } = res.data;

          // حفظ التوكنات الجديدة
          localStorage.setItem("dragon_token", access_token);
          localStorage.setItem("dragon_refresh_token", refresh_token);

          // تحديث هيدر الطلب الأصلي وإعادة إرساله
          apiClient.defaults.headers.common[
            "Authorization"
          ] = `Bearer ${access_token}`;
          originalRequest.headers["Authorization"] = `Bearer ${access_token}`;

          return apiClient(originalRequest);
        }
      } catch (refreshError) {
        console.error("Session expired completely:", refreshError);
        // إذا فشل التجديد، يجب طرد المستخدم
        localStorage.removeItem("dragon_token");
        localStorage.removeItem("dragon_refresh_token");
        // يمكن هنا توجيه المستخدم لصفحة البداية أو عمل Reload
        // window.location.reload();
      }
    }

    return Promise.reject(error);
  }
);

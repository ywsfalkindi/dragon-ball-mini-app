// src/api/client.js
import axios from "axios";

export const apiClient = axios.create({
  // إذا استخدمت نفقاً ثانياً للسيرفر (المنفذ 3000) ضع الرابط هنا
  // أو استخدم IP جهازك إذا كان الهاتف والكمبيوتر على نفس الواي فاي
  baseURL: "http://localhost:3000/api",
  headers: {
    "Content-Type": "application/json",
    // هذا السطر مهم جداً إذا قررت البقاء مع ngrok للسيرفر
    "ngrok-skip-browser-warning": "true",
  },
});

import { create } from "zustand";
import { persist } from "zustand/middleware";
import { apiClient } from "../api/client";
import WebApp from "@twa-dev/sdk";

const useGameStore = create(
  persist(
    (set, get) => ({
      user: null,
      energy: 10,
      score: 0,
      currentQuestion: null,
      isLoading: false,
      error: null,

      setUser: (userData) => set({ user: userData }),

      // 1. دالة تسجيل الدخول الجديدة (The Gate Key) 🔑
      login: async () => {
        const initData = WebApp.initData;
        const authData =
          initData ||
          "query_id=test&user=%7B%22id%22%3A1%2C%22first_name%22%3A%22Goku%22%7D&auth_date=1700000000&hash=test";

        try {
          const response = await apiClient.post("/auth/login", {
            init_data: authData,
          });

          // التصحيح: نستخرج refresh_token ونستخدمه
          const { access_token, refresh_token, user } = response.data;

          localStorage.setItem("dragon_token", access_token);
          localStorage.setItem("dragon_refresh_token", refresh_token); // <--- تم استخدامه الآن!

          set({ user: user, error: null });
          return true;
        } catch (err) {
          console.error("Login Failed:", err);
          set({ error: "فشل الدخول للسيرفر! هل السيرفر يعمل؟" });
          return false;
        }
      },

      fetchQuestion: async () => {
        set({ isLoading: true, error: null });
        try {
          // لاحظ: تم تحديث المسار ليكون protected
          const response = await apiClient.get("/protected/question");
          set({ currentQuestion: response.data.data, isLoading: false });
        } catch (err) {
          console.error("Fetch Error:", err);
          // إذا كان الخطأ 401 (غير مصرح)، ربما انتهى التوكن
          if (err.response && err.response.status === 401) {
            set({ error: "انتهت الجلسة، قم بإعادة تحميل التطبيق." });
          } else {
            set({ error: "لا يمكن استشعار طاقة الكي!", isLoading: false });
          }
        }
      },

      submitAnswer: async (selectedOptionKey) => {
        const { currentQuestion, energy } = get(); // نحتاج الطاقة الحالية
        if (!currentQuestion) return false;

        // 1. Snapshot: نحفظ نسخة من الطاقة الحالية (للعودة إليها في حال الخطأ)
        const previousEnergy = energy;

        // 2. Optimistic Update: نخصم الطاقة فوراً في الواجهة! 👊
        // اللاعب يرى الطاقة تنقص في جزء من الثانية
        set((state) => ({ energy: Math.max(0, state.energy - 1) }));

        try {
          const payload = {
            question_id: currentQuestion.id,
            selected: selectedOptionKey,
          };

          // 3. نرسل الطلب للسيرفر في الخلفية
          const response = await apiClient.post("/protected/answer", payload);
          const result = response.data;

          // 4. Sync: نحدث البيانات الحقيقية القادمة من السيرفر
          // (غالباً ستكون نفس الطاقة التي توقعناها، ولكن السكور سيزيد)
          set({
            score: result.new_score,
            energy: result.new_energy, // تأكيد الطاقة من السيرفر
          });
          return result.correct;
        } catch (err) {
          console.error("Answer Error:", err);

          // 5. Rollback: حدث خطأ! تراجع فوراً! ↩️
          // نعيد الطاقة للاعب وكأن شيئاً لم يحدث
          set({ energy: previousEnergy });

          // نظهر رسالة خطأ
          WebApp.showAlert("خطأ في الاتصال! تمت إعادة طاقة الكي الخاصة بك.");
          return false;
        }
      },

      decreaseEnergy: (amount) =>
        set((state) => ({ energy: Math.max(0, state.energy - amount) })),

      clearError: () => set({ error: null }),
    }),
    {
      name: "dragon-storage",
      partialize: (state) => ({
        // نحفظ فقط الطاقة والسكور، لا نحفظ اليوزر لأننا نجلبه مع اللوجن كل مرة
        score: state.score,
        energy: state.energy,
      }),
    }
  )
);

export default useGameStore;

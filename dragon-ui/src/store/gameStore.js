import { create } from "zustand";
import { persist } from "zustand/middleware";
import { apiClient } from "../api/client";

const useGameStore = create(
  persist(
    (set, get) => ({
      user: null,
      energy: 10,
      score: 0,
      currentQuestion: null,

      // حالات الواجهة (States)
      isLoading: false, // هل غوكو يجمع الطاقة؟
      error: null, // هل هزمنا فريزا؟

      setUser: (userData) => set({ user: userData }),

      // دالة جلب السؤال
      fetchQuestion: async () => {
        set({ isLoading: true, error: null }); // 1. ابدأ التحميل وامسح أي خطأ سابق
        try {
          const response = await apiClient.get("/question");
          // محاكاة تأخير بسيط لنرى شاشة التحميل (اختياري)
          // await new Promise(r => setTimeout(r, 1000));

          set({ currentQuestion: response.data.data, isLoading: false }); // 2. نجاح! أوقف التحميل
        } catch (err) {
          console.error("Fetch Error:", err);
          // 3. فشل! سجل الخطأ وأوقف التحميل
          set({
            error: "لا يمكن استشعار طاقة الكي! تأكد من تشغيل السيرفر.",
            isLoading: false,
          });
        }
      },

      // دالة إرسال الإجابة (تم تصحيح الرابط هنا)
      submitAnswer: async (selectedOptionKey) => {
        const { user, currentQuestion } = get();
        if (!user || !currentQuestion) return false;

        // لاحظ: هنا لن نشغل isLoading لأننا نريد تفاعلاً فورياً،
        // أو يمكننا تشغيله إذا أردنا منع اللاعب من الضغط مرتين.

        try {
          const payload = {
            user_id: user.id,
            question_id: currentQuestion.id,
            selected: selectedOptionKey,
            time_taken: 5,
          };

          // --- التصحيح الهام جداً في الرابط ---
          // يجب أن نرسل التوكن أيضاً (لاحقاً)، والمسار الصحيح هو /protected/answer
          // حالياً بما أننا لم نضبط التوكن في الفرونت إند، قد نحصل على 401.
          // لكن دعنا نضبط المسار الصحيح أولاً:
          const response = await apiClient.post("/protected/answer", payload);

          const result = response.data;
          set({
            score: result.new_score,
            energy: result.new_energy,
          });

          return result.correct;
        } catch (err) {
          console.error("Answer Error:", err);
          // هنا لا نوقف اللعبة كاملة، بل ربما نظهر تنبيهاً صغيراً
          alert("فشل إرسال الهجمة! تحقق من اتصالك.");
          return false;
        }
      },

      decreaseEnergy: (amount) =>
        set((state) => ({ energy: Math.max(0, state.energy - amount) })),

      // دالة لإعادة المحاولة (تصفير الخطأ)
      clearError: () => set({ error: null }),
    }),
    {
      name: "dragon-storage",
      partialize: (state) => ({
        user: state.user,
        score: state.score,
        energy: state.energy,
      }),
    }
  )
);

export default useGameStore;

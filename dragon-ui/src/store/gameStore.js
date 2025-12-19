import { create } from "zustand";
import { persist } from "zustand/middleware";
import { apiClient } from "../api/client"; // استدعاء عميل الاتصال

const useGameStore = create(
  persist(
    (set, get) => ({
      // --- البيانات (State) ---
      user: null,
      energy: 10,
      score: 0,

      // حالة السؤال الحالي
      currentQuestion: null, // السؤال الذي يظهر الآن
      isLoading: false, // هل نقوم بالتحميل الآن؟
      error: null, // هل حدث خطأ؟

      // --- الأفعال (Actions) ---

      setUser: (userData) => set({ user: userData }),

      // 1. دالة جلب سؤال جديد من السيرفر
      fetchQuestion: async () => {
        set({ isLoading: true, error: null }); // بدأ التحميل
        try {
          // الطلب: GET /question
          const response = await apiClient.get("/question");

          // النجاح: نضع السؤال في الحالة
          // response.data.data لأن الباك إند يرسل { status: success, data: { ... } }
          set({ currentQuestion: response.data.data, isLoading: false });
        } catch (err) {
          console.error("Failed to fetch question:", err);
          set({ error: "فشل الاتصال بكوكب ناميك!", isLoading: false });
        }
      },

      // 2. دالة إرسال الإجابة للسيرفر
      submitAnswer: async (selectedOption) => {
        const { user, currentQuestion } = get(); // نأخذ البيانات الحالية
        if (!user || !currentQuestion) return false;

        try {
          // الطلب: POST /answer
          // نرسل البيانات كما توقعناها في Go (Chapter 9)
          const payload = {
            user_id: user.id, // الآيدي القادم من تليجرام
            question_id: currentQuestion.id,
            selected: selectedOption,
          };

          const response = await apiClient.post("/answer", payload);
          const result = response.data; // { correct: true, new_score: 100, ... }

          // تحديث البيانات بناءً على رد السيرفر الحقيقي
          set({
            score: result.new_score, // تحديث النقاط من السيرفر
            energy: result.new_energy, // تحديث الطاقة من السيرفر
          });

          return result.correct; // نرجع هل الإجابة صحيحة أم لا للواجهة
        } catch (err) {
          console.error("Attack failed:", err);
          return false;
        }
      },

      // دالة مساعدة لتقليل الطاقة محلياً (لو احتجنا تحديثاً فورياً قبل رد السيرفر)
      decreaseEnergy: (amount) =>
        set((state) => ({ energy: Math.max(0, state.energy - amount) })),
    }),
    {
      name: "dragon-storage",
      partialize: (state) => ({
        user: state.user,
        score: state.score,
        energy: state.energy,
      }), // لا نحفظ السؤال والتحميل
    }
  )
);

export default useGameStore;

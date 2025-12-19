import React, { useState, useEffect } from "react";
import { AnimatePresence } from "framer-motion";
import HealthBar from "./components/HealthBar";
import QuestionCard from "./components/QuestionCard";
import AnswerButton from "./components/AnswerButton";
import useGameStore from "./store/gameStore";
import WebApp from "@twa-dev/sdk"; // الاستيراد موجود (ممتاز)
import UserProfile from "./components/UserProfile"; // الاستيراد موجود (ممتاز)

// قائمة أسئلة للتجربة
const mockQuestions = [
  {
    id: 1,
    question_text: "ما هي التقنية التي تعلمها غوكو من الملك كاي؟",
    options: ["Kamehameha", "Kaio-ken", "Final Flash", "Instant Transmission"],
    correct: "Kaio-ken",
  },
  {
    id: 2,
    question_text: "من هو أول سوبر ساياجين ظهر في الأنمي؟",
    options: ["Goku", "Vegeta", "Broly", "Gohan"],
    correct: "Goku",
  },
  {
    id: 3,
    question_text: "ما اسم كوكب بيكولو الأصلي؟",
    options: ["Vegeta", "Earth", "Namek", "Sadala"],
    correct: "Namek",
  },
];

function App() {
  const decreaseEnergy = useGameStore((state) => state.decreaseEnergy);

  // 1. نحتاج لهذه الدالة لحفظ بيانات المستخدم القادم من تليجرام
  const setUser = useGameStore((state) => state.setUser);

  const [currentQIndex, setCurrentQIndex] = useState(0);
  const [isWrong, setIsWrong] = useState(false);

  const currentQuestion = mockQuestions[currentQIndex];

  // 2. useEffect: هذا هو الجزء المفقود لتشغيل تليجرام SDK
  useEffect(() => {
    // التحقق هل نحن داخل تليجرام؟
    if (WebApp.initDataUnsafe.user) {
      WebApp.ready();
      WebApp.expand();
      WebApp.setHeaderColor("#000000");

      // حفظ بيانات المستخدم الحقيقي
      setUser({
        id: WebApp.initDataUnsafe.user.id,
        firstName: WebApp.initDataUnsafe.user.first_name,
        username: WebApp.initDataUnsafe.user.username,
        photoUrl: WebApp.initDataUnsafe.user.photo_url,
      });
    } else {
      // وضع تجريبي للمتصفح
      console.log("Running in browser mode (Mock User)");
      setUser({
        id: 999,
        firstName: "Test Goku",
        username: "kakarot",
        photoUrl: null,
      });
    }
  }, [setUser]);

  const handleAnswer = (selectedOptionText) => {
    // 3. إضافة الاهتزاز عند الضغط (تجربة مستخدم أفضل)
    if (WebApp.initDataUnsafe.user) {
      WebApp.HapticFeedback.impactOccurred("light");
    }

    const isAnswerCorrect = currentQuestion.options.some(
      (opt) =>
        // تأكد أننا نتحقق من الخيار الذي يحوي نص الإجابة الصحيحة
        selectedOptionText === opt && opt.includes(currentQuestion.correct)
    );

    if (isAnswerCorrect) {
      if (currentQIndex < mockQuestions.length - 1) {
        setCurrentQIndex((prev) => prev + 1);
        setIsWrong(false);
      } else {
        alert("مبروك! انتهت الأسئلة التجريبية");
      }
    } else {
      setIsWrong(true);
      decreaseEnergy(1);
      setTimeout(() => setIsWrong(false), 500);
    }
  };

  return (
    <div
      className="app-container"
      style={{
        position: "relative",
        height: "100vh",
        overflow: "hidden",
        padding: "20px",
      }} // أضفنا padding
    >
      {/* 4. إظهار بطاقة المستخدم هنا */}
      <UserProfile />

      <HealthBar />

      <div style={{ position: "relative", width: "100%", height: "400px" }}>
        <AnimatePresence mode="wait">
          <QuestionCard
            key={currentQuestion.id}
            question={{
              id: currentQuestion.id,
              question_text: currentQuestion.question_text,
            }}
            isWrong={isWrong}
          >
            {currentQuestion.options.map((opt, index) => (
              <AnswerButton
                key={index}
                text={opt}
                onClick={() => handleAnswer(opt)}
                state={null}
              />
            ))}
          </QuestionCard>
        </AnimatePresence>
      </div>
    </div>
  );
}

export default App;

import React, { useState, useEffect } from "react";
/* eslint-disable no-unused-vars */
import { AnimatePresence, motion } from "framer-motion";
import HealthBar from "./components/HealthBar";
import QuestionCard from "./components/QuestionCard";
import AnswerButton from "./components/AnswerButton";
import useGameStore from "./store/gameStore";
import WebApp from "@twa-dev/sdk";
import UserProfile from "./components/UserProfile";

function App() {
  // تم إزالة user من هنا لأنه مستخدم فقط داخل مكون UserProfile المستقل
  const {
    setUser,
    fetchQuestion,
    currentQuestion,
    submitAnswer,
    isLoading,
    error,
  } = useGameStore();

  const [isWrong, setIsWrong] = useState(false);

  useEffect(() => {
    if (WebApp.initDataUnsafe.user) {
      WebApp.ready();
      WebApp.expand();
      WebApp.setHeaderColor("#000000");

      setUser({
        id: WebApp.initDataUnsafe.user.id,
        firstName: WebApp.initDataUnsafe.user.first_name,
        username: WebApp.initDataUnsafe.user.username,
        photoUrl: WebApp.initDataUnsafe.user.photo_url,
      });
    } else {
      setUser({
        id: 1,
        firstName: "Test Goku",
        username: "kakarot",
        photoUrl: null,
      });
    }
  }, [setUser]);

  useEffect(() => {
    if (!currentQuestion) {
      fetchQuestion();
    }
  }, [fetchQuestion, currentQuestion]);

  const handleAnswer = async (selectedKey) => {
    if (WebApp.initDataUnsafe.user) {
      WebApp.HapticFeedback.impactOccurred("light");
    }

    const isCorrect = await submitAnswer(selectedKey);

    if (isCorrect) {
      setTimeout(() => {
        setIsWrong(false);
        fetchQuestion();
      }, 500);
    } else {
      setIsWrong(true);
      setTimeout(() => setIsWrong(false), 500);
    }
  };

  const optionsList = currentQuestion
    ? [
        { key: "A", text: currentQuestion.option_a },
        { key: "B", text: currentQuestion.option_b },
        { key: "C", text: currentQuestion.option_c },
        { key: "D", text: currentQuestion.option_d },
      ]
    : [];

  return (
    <div
      className="app-container"
      style={{
        position: "relative",
        height: "100vh",
        overflow: "hidden",
        padding: "20px",
      }}
    >
      <UserProfile />
      <HealthBar />

      <div
        style={{
          position: "relative",
          width: "100%",
          height: "400px",
          marginTop: "20px",
        }}
      >
        {isLoading && !currentQuestion && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            style={{
              textAlign: "center",
              marginTop: "100px",
              color: "var(--db-orange)",
              fontSize: "18px",
              fontWeight: "bold",
            }}
          >
            جاري استدعاء التنين... 🐉
            <br />
            <span style={{ fontSize: "12px", color: "gray" }}>
              (Connecting to Namek...)
            </span>
          </motion.div>
        )}

        {error && (
          <div
            style={{
              textAlign: "center",
              color: "var(--danger-red)",
              marginTop: "50px",
              background: "rgba(0,0,0,0.7)",
              padding: "20px",
              borderRadius: "10px",
            }}
          >
            🛑 {error}
          </div>
        )}

        <AnimatePresence mode="wait">
          {currentQuestion && !isLoading && (
            <QuestionCard
              key={currentQuestion.id}
              question={currentQuestion}
              isWrong={isWrong}
            >
              {optionsList.map((opt) => (
                <AnswerButton
                  key={opt.key}
                  text={`${opt.key}) ${opt.text}`}
                  onClick={() => handleAnswer(opt.key)}
                  state={null}
                />
              ))}
            </QuestionCard>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}

export default App;

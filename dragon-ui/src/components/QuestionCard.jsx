import React, { useEffect, useState } from "react"; // <--- 1. تم إضافة الاستيرادات هنا
/* eslint-disable no-unused-vars */
import { motion } from "framer-motion";
import WebApp from "@twa-dev/sdk";

const cardVariants = {
  hidden: { opacity: 0, x: 100, scale: 0.8 },
  visible: {
    opacity: 1,
    x: 0,
    scale: 1,
    transition: { type: "spring", stiffness: 120 },
  },
  exit: { opacity: 0, x: -100, scale: 0.8, transition: { duration: 0.3 } },
  shake: { x: [0, -20, 20, -20, 20, 0], transition: { duration: 0.4 } },
};

const QuestionCard = ({ question, children, isWrong }) => {
  // ⏳ حالة للعداد الزمني
  // ملاحظة: لأننا نستخدم key في App.jsx، سيتم إعادة بناء هذا المكون مع كل سؤال
  // لذلك القيمة ستبدأ دائماً من 30 تلقائياً
  const [timeLeft, setTimeLeft] = useState(30);

  useEffect(() => {
    // تم حذف setTimeLeft(30) من هنا لحل مشكلة التزامن
    // العداد سيعمل تلقائياً عند بناء المكون

    const timer = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev <= 1) {
          clearInterval(timer);
          return 0;
        }

        // ⚡ اهتزازة خفيفة مع كل ثانية
        try {
          WebApp.HapticFeedback.selectionChanged();
        } catch (e) {
          // ignore error in development
        }

        return prev - 1;
      });
    }, 1000);

    // تنظيف العداد عند تدمير المكون
    return () => clearInterval(timer);
  }, []); // المصفوفة فارغة لأننا نعتمد على إعادة بناء المكون (Remounting)

  // تغيير لون العداد عند اقتراب النهاية
  const timerColor = timeLeft < 10 ? "var(--danger-red)" : "var(--db-orange)";

  return (
    <motion.div
      variants={cardVariants}
      initial="hidden"
      animate={isWrong ? "shake" : "visible"}
      exit="exit"
      style={{
        background: "rgba(20, 20, 30, 0.8)",
        backdropFilter: "blur(10px)",
        padding: "25px",
        borderRadius: "20px",
        border: "1px solid rgba(255,255,255,0.1)",
        boxShadow: "0 8px 32px 0 rgba(0, 0, 0, 0.37)",
        position: "absolute",
        width: "100%",
        left: 0,
        right: 0,
      }}
    >
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          marginBottom: "15px",
        }}
      >
        <span style={{ color: "gray", fontSize: "14px" }}>
          QUESTION #{question.id}
        </span>

        {/* ⏳ عرض العداد الزمني */}
        <span
          style={{ color: timerColor, fontWeight: "bold", fontSize: "16px" }}
        >
          ⏱ {timeLeft}s
        </span>
      </div>

      <h2
        style={{
          marginTop: 0,
          marginBottom: "25px",
          fontSize: "22px",
          lineHeight: "1.4",
        }}
      >
        {question.question_text}
      </h2>

      <div>{children}</div>

      {/* شريط وقت مرئي في الأسفل */}
      <div
        style={{
          width: "100%",
          height: "4px",
          background: "rgba(255,255,255,0.1)",
          marginTop: "20px",
          borderRadius: "2px",
          overflow: "hidden",
        }}
      >
        <motion.div
          initial={{ width: "100%" }}
          animate={{ width: "0%" }}
          transition={{ duration: 30, ease: "linear" }}
          style={{ height: "100%", background: timerColor }}
        />
      </div>
    </motion.div>
  );
};

export default QuestionCard;

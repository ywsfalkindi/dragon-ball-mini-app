import React from "react";
/* eslint-disable no-unused-vars */
import { motion } from "framer-motion";

// تعريف "وصفة" الحركات (Variants)
const cardVariants = {
  hidden: {
    opacity: 0,
    x: 100,
    scale: 0.8,
  },
  visible: {
    opacity: 1,
    x: 0,
    scale: 1,
    transition: { type: "spring", stiffness: 120 },
  },
  exit: {
    opacity: 0,
    x: -100,
    scale: 0.8,
    transition: { duration: 0.3 },
  },
  shake: {
    x: [0, -20, 20, -20, 20, 0],
    transition: { duration: 0.4 },
  },
};

const QuestionCard = ({ question, children, isWrong }) => {
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
      <div style={{ color: "gray", marginBottom: "10px", fontSize: "14px" }}>
        QUESTION #{question.id}
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
    </motion.div>
  );
};

export default QuestionCard;

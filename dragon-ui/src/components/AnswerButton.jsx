import React from "react";
/* eslint-disable no-unused-vars */
import { motion } from "framer-motion";

// المكون يستقبل: النص، وهل تم اختياره؟ وهل هو صحيح؟ ودالة الضغط
const AnswerButton = ({ text, onClick, state }) => {
  // تحديد لون الخلفية بناءً على الحالة
  let bg = "var(--glass-bg)"; // الوضع العادي (شفاف)
  let border = "1px solid rgba(255,255,255,0.2)";

  if (state === "selected") {
    bg = "rgba(255, 153, 0, 0.5)"; // برتقالي عند الاختيار
    border = "2px solid var(--db-orange)";
  } else if (state === "correct") {
    bg = "rgba(57, 255, 20, 0.5)"; // أخضر عند الفوز
    border = "2px solid var(--energy-green)";
  } else if (state === "wrong") {
    bg = "rgba(255, 0, 85, 0.5)"; // أحمر عند الخسارة
    border = "2px solid var(--danger-red)";
  }

  return (
    <motion.button
      whileHover={{ scale: 1.02 }} // تكبير بسيط عند تمرير الماوس
      whileTap={{ scale: 0.95 }} // تصغير بسيط عند الضغط (شعور الضغطة)
      onClick={onClick}
      style={{
        width: "100%",
        padding: "15px",
        margin: "8px 0",
        background: bg,
        border: border,
        borderRadius: "12px",
        color: "white",
        fontSize: "16px",
        cursor: "pointer",
        textAlign: "left",
        position: "relative",
        overflow: "hidden", // لإخفاء أي توهج يخرج عن الحدود
      }}
    >
      {text}
    </motion.button>
  );
};

export default AnswerButton;

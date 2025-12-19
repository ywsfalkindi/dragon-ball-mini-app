import React from "react";
/* eslint-disable no-unused-vars */
import { motion } from "framer-motion";

const ErrorScreen = ({ message, onRetry }) => {
  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
        height: "300px",
        textAlign: "center",
        padding: "20px",
        background: "rgba(255, 0, 0, 0.1)",
        borderRadius: "15px",
        border: "1px solid var(--danger-red)",
      }}
    >
      <div style={{ fontSize: "50px", marginBottom: "10px" }}>💀</div>

      <h3 style={{ color: "var(--danger-red)", margin: "0 0 10px 0" }}>
        مهمة فاشلة!
      </h3>

      <p style={{ color: "white", marginBottom: "20px", fontSize: "14px" }}>
        {message}
      </p>

      {/* زر إعادة المحاولة */}
      <motion.button
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        onClick={onRetry}
        style={{
          background: "var(--db-orange)",
          border: "none",
          padding: "10px 20px",
          borderRadius: "20px",
          color: "black",
          fontWeight: "bold",
          cursor: "pointer",
        }}
      >
        استخدام حبة سينزو (Retry) 💊
      </motion.button>
    </div>
  );
};

export default ErrorScreen;

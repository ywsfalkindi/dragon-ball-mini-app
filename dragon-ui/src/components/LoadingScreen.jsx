import React from "react";
/* eslint-disable no-unused-vars */
import { motion } from "framer-motion";

const LoadingScreen = () => {
  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        alignItems: "center",
        height: "300px", // ارتفاع مناسب داخل البطاقة
        color: "var(--db-orange)",
      }}
    >
      {/* كرة طاقة تدور */}
      <motion.div
        animate={{ rotate: 360 }}
        transition={{ repeat: Infinity, duration: 1, ease: "linear" }}
        style={{
          width: "50px",
          height: "50px",
          border: "5px solid rgba(255, 153, 0, 0.3)",
          borderTop: "5px solid var(--db-orange)",
          borderRadius: "50%",
          marginBottom: "20px",
        }}
      />

      <motion.div
        animate={{ opacity: [0.5, 1, 0.5] }}
        transition={{ repeat: Infinity, duration: 1.5 }}
        style={{ fontWeight: "bold", fontSize: "18px" }}
      >
        جاري شحن الطاقة...
      </motion.div>

      <div style={{ fontSize: "12px", color: "gray", marginTop: "5px" }}>
        (Goku is eating ramen 🍜)
      </div>
    </div>
  );
};

export default LoadingScreen;

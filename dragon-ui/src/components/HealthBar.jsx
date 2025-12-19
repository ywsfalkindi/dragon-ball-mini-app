import React from "react";
/* eslint-disable no-unused-vars */
import { motion } from "framer-motion";
import useGameStore from "../store/gameStore";

const HealthBar = () => {
  const energy = useGameStore((state) => state.energy);
  const max = 10;

  const percentage = (energy / max) * 100;

  let color = "var(--energy-green)";
  if (percentage < 50) color = "yellow";
  if (percentage < 20) color = "var(--danger-red)";

  return (
    <div style={{ width: "100%", marginBottom: "20px" }}>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          marginBottom: "5px",
        }}
      >
        <span style={{ fontWeight: "bold", color: "var(--db-orange)" }}>
          STAMINA (KI)
        </span>
        <span>
          {energy} / {max}
        </span>
      </div>

      <div
        style={{
          height: "15px",
          background: "rgba(255,255,255,0.2)",
          borderRadius: "10px",
          overflow: "hidden",
          border: "1px solid rgba(255,255,255,0.3)",
        }}
      >
        <motion.div
          initial={{ width: 0 }}
          animate={{ width: `${percentage}%` }}
          transition={{ duration: 0.5 }}
          style={{
            height: "100%",
            background: color,
            boxShadow: `0 0 10px ${color}`,
          }}
        />
      </div>
    </div>
  );
};

export default HealthBar;

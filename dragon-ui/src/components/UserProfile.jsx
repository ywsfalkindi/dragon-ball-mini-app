import React from "react";
import useGameStore from "../store/gameStore";

const UserProfile = () => {
  // جلب بيانات المستخدم من المخزن
  const user = useGameStore((state) => state.user);

  // إذا لم يكن هناك مستخدم (أو جاري التحميل)، لا تظهر شيئاً
  if (!user) return null;

  return (
    <div
      style={{
        display: "flex",
        alignItems: "center",
        gap: "10px",
        marginBottom: "15px",
        background: "rgba(0, 0, 0, 0.3)",
        padding: "10px",
        borderRadius: "50px", // شكل كبسولة
        width: "fit-content",
      }}
    >
      {/* صورة المستخدم */}
      {user.photo_url ? (
        <img
          src={user.photo_url}
          alt="User"
          style={{
            width: "40px",
            height: "40px",
            borderRadius: "50%",
            border: "2px solid var(--db-orange)",
          }}
        />
      ) : (
        // أيقونة بديلة لو لم يكن لديه صورة
        <div
          style={{
            width: "40px",
            height: "40px",
            borderRadius: "50%",
            background: "gray",
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
          }}
        >
          👤
        </div>
      )}

      {/* الاسم والرتبة */}
      <div>
        <div style={{ fontWeight: "bold", fontSize: "14px" }}>
          {user.first_name}
        </div>
        <div style={{ fontSize: "11px", color: "var(--db-orange)" }}>
          Warrior ID: {user.id}
        </div>
      </div>
    </div>
  );
};

export default UserProfile;

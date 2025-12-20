package worker

import (
	"context"
	"dragon-core/internal/database"
	"dragon-core/internal/models"
	"log"
	"strconv"
	"time"
)

// StartSyncWorker: يبدأ عملية المزامنة في الخلفية
func StartSyncWorker() {
	// Ticker: مثل المنبه، يرن كل دقيقة
	ticker := time.NewTicker(1 * time.Minute)

	go func() {
		// التصحيح هنا: استخدام for range مباشرة مع القناة (Channel)
		// هذا أنظف وأكثر كفاءة في Go
		for range ticker.C {
			SyncScoresToPostgres()
		}
	}()
}

func SyncScoresToPostgres() {
	log.Println("🔄 Worker: Syncing Redis scores to Postgres...")
	
	key := "leaderboard:global"
	ctx := context.Background()

	// 1. نجلب كل اللاعبين ونقاطهم من Redis
	// ZRangeWithScores يجلب القائمة كاملة
	results, err := database.RDB.ZRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		log.Println("❌ Worker Error reading Redis:", err)
		return
	}

	if len(results) == 0 {
		return
	}

	// 2. التحديث في Postgres
	// للسرعة، سنقوم بتحديث كل مستخدم على حدة (يمكن تحسينه ليكون Bulk Update لاحقاً)
	for _, z := range results {
		userIDStr := z.Member.(string)
		score := int(z.Score)
		
		userID, _ := strconv.Atoi(userIDStr)

		// تحديث عمود total_score في جدول users
		// نستخدم Model للتحديث المباشر
		err := database.DB.Model(&models.User{}).
			Where("id = ?", userID).
			Update("total_score", score).Error
		
		if err != nil {
			log.Printf("⚠️ Worker failed to update user %d: %v", userID, err)
		}
	}

	log.Printf("✅ Worker: Synced %d players to Database.", len(results))
}
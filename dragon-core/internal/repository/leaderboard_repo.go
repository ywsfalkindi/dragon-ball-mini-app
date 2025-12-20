package repository

import (
	"context"
	"fmt"
	// "time" <--- قمنا بحذف هذا السطر لأنه غير مستخدم

	"github.com/redis/go-redis/v9"
)

type LeaderboardRepo struct {
	rdb *redis.Client
}

func (r *LeaderboardRepo) IncrementScore(ctx context.Context, userID uint, points float64) (float64, error) {
	key := "leaderboard:global"
	// Member يجب أن يكون string لضمان التوافق
	return r.rdb.ZIncrBy(ctx, key, points, fmt.Sprintf("%d", userID)).Result()
}

func NewLeaderboardRepo(rdb *redis.Client) *LeaderboardRepo {
	return &LeaderboardRepo{rdb: rdb}
}

// UpdateScore يحدث نقاط اللاعب في Redis
func (r *LeaderboardRepo) UpdateScore(ctx context.Context, userID uint, score float64) error {
	key := "leaderboard:global"
	return r.rdb.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: userID,
	}).Err()
}

// GetTopPlayers يجلب أفضل N لاعبين
// التصحيح: جعلنا القيمة المرجعة []redis.Z (قائمة) بدلاً من redis.Z (عنصر واحد)
func (r *LeaderboardRepo) GetTopPlayers(ctx context.Context, limit int64) ([]redis.Z, error) {
	key := "leaderboard:global"
	return r.rdb.ZRevRangeWithScores(ctx, key, 0, limit-1).Result()
}

// GetUserRank يجلب ترتيب لاعب معين
func (r *LeaderboardRepo) GetUserRank(ctx context.Context, userID uint) (int64, error) {
	key := "leaderboard:global"
	return r.rdb.ZRevRank(ctx, key, fmt.Sprintf("%d", userID)).Result()
}

func (r *LeaderboardRepo) GetCurrentScore(ctx context.Context, userID uint) (float64, error) {
	key := "leaderboard:global"
	return r.rdb.ZScore(ctx, key, fmt.Sprintf("%d", userID)).Result()
}
package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math"
)

type PlayerRankInfo struct {
	PlayerId string `json:"playerId"`
	Rank     int64  `json:"rank"`
	Score    int64  `json:"score"`
}

// UpdatePlayerScore 更新玩家分数（使用ZADD命令）
func UpdatePlayerScore(playerId string, score float64, timestamp int64) error {
	//todo 入口打个log记录下，方便后续查日志
	fmt.Printf("UpdatePlayerScore services. playerId=%s,score=%f,time=%d\n", playerId, score, timestamp)

	// 为了让同成绩的人按时间戳小的排前边， 后边取的时候 还原真实成绩都是按score为整数做的
	compositeScore := score*ScoreShift - float64(timestamp)

	// 将玩家添加到有序集合中
	err := RedisClient.ZAdd(GlobalCtx, LeadBoardRedisKey, redis.Z{
		Score:  compositeScore,
		Member: playerId,
	}).Err()

	return err
}

// GetPlayerRank 获取玩家排名（使用ZREVRANK命令）
func GetPlayerRank(ctx context.Context, playerId string) (PlayerRankInfo, error) {
	rank, err := RedisClient.ZRevRank(ctx, LeadBoardRedisKey, playerId).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return PlayerRankInfo{}, ErrPlayerNotFound
		}
		return PlayerRankInfo{}, err
	}

	// 获取实际分数（需解析 compositeScore）
	score, err := RedisClient.ZScore(ctx, LeadBoardRedisKey, playerId).Result()
	if err != nil {
		return PlayerRankInfo{}, err
	}
	actualScore := int64(math.Round(score / ScoreShift))        // 还原真实成绩
	return PlayerRankInfo{playerId, rank + 1, actualScore}, nil // 排名从1开始
}

// GetTopN 获取前N名玩家（使用ZREVRANGE命令）
func GetTopN(ctx context.Context, n int64) ([]PlayerRankInfo, error) {
	zs, err := RedisClient.ZRevRangeWithScores(ctx, LeadBoardRedisKey, 0, n-1).Result()
	if err != nil {
		return nil, err
	}

	list := make([]PlayerRankInfo, 0, len(zs))
	for i, z := range zs {
		member, ok := z.Member.(string)
		if !ok {
			continue
		}
		list = append(list, PlayerRankInfo{
			PlayerId: member,
			Rank:     int64(i) + 1,
			Score:    int64(math.Round(z.Score / ScoreShift)),
		})
	}
	return list, nil
}

// GetPlayerRankRange 获取玩家前后N名
func GetPlayerRankRange(ctx context.Context, playerId string, before, after int64) ([]PlayerRankInfo, error) {
	// 1) 查出玩家自己的排名（0-based）
	rank, err := RedisClient.ZRevRank(ctx, LeadBoardRedisKey, playerId).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrPlayerNotFound
		}
		return nil, err
	}
	// 2) 计算切片区间
	start := rank - before
	if start < 0 {
		start = 0
	}
	end := rank + after

	// 3) 批量取
	zs, err := RedisClient.ZRevRangeWithScores(ctx, LeadBoardRedisKey, start, end).Result()
	if err != nil {
		return nil, err
	}

	// 4) 构造返回
	list := make([]PlayerRankInfo, 0, len(zs))
	for i, z := range zs {
		member, ok := z.Member.(string)
		if !ok {
			continue
		}
		list = append(list, PlayerRankInfo{
			PlayerId: member,
			Rank:     int64(i) + 1,
			Score:    int64(math.Round(z.Score / ScoreShift)),
		})
	}
	return list, nil
}

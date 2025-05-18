package services

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis服务器地址
		Password: "ZyxTest554",     // 密码
		DB:       0,                // 默认数据库
	})

	// 测试连接
	_, err := RedisClient.Ping(GlobalCtx).Result()
	if err != nil {
		panic(fmt.Sprintf("Redis连接失败: %v", err))
	}
	fmt.Println("Redis连接成功")
}

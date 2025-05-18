package main

import (
	"github.com/Eacie/leaderboard-service/handlers"
	"github.com/Eacie/leaderboard-service/services"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	services.InitRedis()
	// 定义路由
	router.POST("/update-score", handlers.UpdateScore)
	router.GET("/get-rank", handlers.GetPlayerRank)
	router.GET("/get-top", handlers.GetTopN)
	router.GET("/get-range", handlers.GetPlayerRankRange)

	// 启动服务
	router.Run(":8080")
}

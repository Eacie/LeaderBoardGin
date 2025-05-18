package handlers

import (
	"errors"
	"fmt"
	"github.com/Eacie/leaderboard-service/models"
	"github.com/Eacie/leaderboard-service/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//todo 根据业务传参类型在service层对应不同的redisKey，一个enum值是类型，名字是key那种应该最方便，在C#可以实现，go后边再看看

// UpdateScore 这只是给个方便自己测试接口，实际更新应该是由项目业务触发的，就不校验什么了
func UpdateScore(c *gin.Context) {
	var req models.UpdateScoreRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "param error: " + err.Error()})
		//一般来说是调日志接口打日志，不过不在项目里跑先就直接print一下
		fmt.Println("UpdateScore Handler. error: ", err)
		return
	}
	//score暂时按整数的方案处理的，如果是有小数应该不做放大，而是把时间戳缩小，取负相加，不用ZRevRank
	if err := services.UpdatePlayerScore(req.PlayerId, req.Score, time.Now().UnixMilli()); err != nil {
		// 异步错误处理（记录日志）
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update error: " + err.Error()})
		fmt.Printf("UpdateScore Handler error.player=%s: %v\n", req.PlayerId, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// GetPlayerRank 获取玩家排名
func GetPlayerRank(c *gin.Context) {
	var req models.GetPlayerRankRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "param error: " + err.Error()})
		fmt.Println("GetPlayerRank Handler. error: ", err)
		return
	}

	rankInfo, err := services.GetPlayerRank(c.Request.Context(), req.PlayerId)
	if err != nil {
		if errors.Is(err, services.ErrPlayerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		fmt.Println("GetPlayerRank Handler. error: ", err)
		return
	}

	c.JSON(http.StatusOK, rankInfo)
}

// GetTopN 获取前N名
func GetTopN(c *gin.Context) {
	var req models.GetTopNRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "param error: " + err.Error()})
		fmt.Println("GetTopN Handler. error: ", err)
		return
	}

	topPlayers, err := services.GetTopN(c.Request.Context(), int64(req.N))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		fmt.Println("GetTopN Handler. error: ", err)
		return
	}

	c.JSON(http.StatusOK, topPlayers)
}

// GetPlayerRankRange 获取玩家周边排名，参数改成自由传前面要多少个后面要多少个了，更适用一点
func GetPlayerRankRange(c *gin.Context) {
	var req models.GetPlayerRankRangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "param error: " + err.Error()})
		fmt.Println("GetPlayerRankRange Handler. error: ", err)
		return
	}
	// 数量限制一下
	if req.Before+req.After > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request number greater than limit"})
		return
	}
	rangePlayers, err := services.GetPlayerRankRange(c.Request.Context(), req.PlayerId, req.Before, req.After)
	if err != nil {
		if errors.Is(err, services.ErrPlayerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		fmt.Println("GetPlayerRankRange Handler. error: ", err)
		return
	}

	c.JSON(http.StatusOK, rangePlayers)
}

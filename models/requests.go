package models

type UpdateScoreRequest struct {
	PlayerId string  `form:"playerId" binding:"required"`
	Score    float64 `form:"score" binding:"required"`
}

type GetPlayerRankRequest struct {
	PlayerId string `form:"playerId" binding:"required"`
}

type GetTopNRequest struct {
	N int `form:"n" binding:"required,min=1"`
}

type GetPlayerRankRangeRequest struct {
	PlayerId string `form:"playerId" binding:"required"`
	Before   int64  `form:"before"   binding:"min=0"`
	After    int64  `form:"after"    binding:"min=0"`
}

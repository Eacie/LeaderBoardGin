package services

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient       *redis.Client
	GlobalCtx         = context.Background()
	ErrPlayerNotFound = errors.New("player not found")
)

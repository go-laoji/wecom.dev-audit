package queue

import (
	"context"
	"github.com/go-redis/redis/v9"
	"os"
	"time"
	"wecom.dev/audit/logger"
)

type Redis struct {
	client *redis.Client
	queue  string
	ctx    context.Context
}

func NewRedis() (rdb Redis, err error) {
	rdb.client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("QueueHost"),
		Username: os.Getenv("QueueUser"),
		Password: os.Getenv("QueuePassword"),
		DB:       0,
	})
	rdb.ctx = context.Background()
	_, err = rdb.client.Ping(rdb.ctx).Result()
	if err != nil {
		logger.Surgar.Error("ping redis failed", err)
		os.Exit(0)
	}
	rdb.queue = os.Getenv("QueueName")
	return
}

func (rdb Redis) Push(T any) (err error) {
	_, err = rdb.client.LPush(rdb.ctx, rdb.queue, T).Result()
	return
}

func (rdb Redis) Pop() (interface{}, error) {
	content, err := rdb.client.BRPop(rdb.ctx, time.Second*15, rdb.queue).Result()
	logger.Surgar.Info(content, err)
	return nil, nil
}

func (rdb Redis) Size() (size int64, err error) {
	size, err = rdb.client.LLen(rdb.ctx, rdb.queue).Result()
	return
}

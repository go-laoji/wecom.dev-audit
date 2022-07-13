package queue

import (
	"github.com/joho/godotenv"
	"os"
	"strings"
	"wecom.dev/audit/logger"
)

var Q Queue
var err error

type Queue interface {
	Push(T any) error
	Pop() (interface{}, error)
	Size() (int64, error)
}

func init() {
	godotenv.Load()
	switch strings.TrimSpace(os.Getenv("StoreType")) {
	case "2", "3":
		switch os.Getenv("QueueType") {
		case "redis":
			Q, err = NewRedis()
			break
		default:
			logger.Surgar.Error("un support queue type")
			os.Exit(1)
		}
		break
	default:
		logger.Surgar.Info("不需要配置队列服务")
	}

	if err != nil {
		panic(err)
	}
}

package redisclient

import (
	"context"
	"encoding/json"
	"os"

	"github.com/redis/go-redis/v9"

	"github.com/joaofilippe/go-queue/models"
)

var rCli *redis.Client

func LoadRedisClient(ctx context.Context) {
	uri := os.Getenv("REDIS_TLS_URL")
	opts, err := redis.ParseURL(uri)
	if err != nil {
		panic(err)
	}

	rCli := redis.NewClient(opts)

	qStd := models.UserQueue{}
	qPri := models.UserQueue{}
	qOne := models.UserQueue{}
	qTwo := models.UserQueue{}

	jsonQStd, err := json.Marshal(qStd)
	if err != nil {
		println(err)
	}
	jsonQPri, err := json.Marshal(qPri)
	if err != nil {
		println(err)
	}
	jsonQOne, err := json.Marshal(qOne)
	if err != nil {
		println(err)
	}
	jsonQTwo, err := json.Marshal(qTwo)
	if err != nil {
		println(err)
	}

	rCli.Set(ctx, "standart", string(jsonQStd), 0)
	rCli.Set(ctx, "priority", string(jsonQPri), 0)
	rCli.Set(ctx, "priority", string(jsonQOne), 0)
	rCli.Set(ctx, "priority", string(jsonQTwo), 0)
}

// GetQueueFromRedis gets the queue from redis-cli
func GetQueueFromRedis(ctx context.Context, key string) *models.UserQueue {
	result := rCli.Get(ctx, key)
	queueStr, err := result.Result()
	if err != nil {
		panic(err)
	}

	queue := new(models.UserQueue)

	json.Unmarshal([]byte(queueStr), queue)

	return queue
}

// SendQueueToRedis sends the new queue to Redis
func SendQueueToRedis(ctx context.Context, queue *models.UserQueue, key string) {
	res, _ := json.Marshal(queue)
	err := rCli.Set(ctx, key, string(res), 0)
	if err != nil {
		println(err)
	}
}

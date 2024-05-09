package redisclient

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"

	"github.com/joaofilippe/go-queue/models"
)

var rCli *redis.Client

func LoadRedisClient(ctx context.Context) {
	rCli = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	qStd := models.UserQueue{}
	qPri := models.UserQueue{}

	jsonQStd, err := json.Marshal(qStd)
	if err != nil {
		println(err)
	}
	jsonQPri, err := json.Marshal(qPri)
	if err != nil {
		println(err)
	}

	rCli.Set(ctx, "standart", string(jsonQStd), 0)
	rCli.Set(ctx, "priority", string(jsonQPri), 0)
}

// GetQueueFromRedis gets the queue from redis-cli
func GetQueueFromRedis(ctx context.Context, key string) *models.UserQueue {
	result := rCli.Get(ctx, key)
	queueStr, err := result.Result()
	if err != nil {
		println(err)
	}

	queue := new(models.UserQueue)

	json.Unmarshal([]byte(queueStr), queue)

	return queue
}

// SendQueueToRedis sends the new queue to Redis
func SendQueueToRedis(ctx context.Context, queue *models.UserQueue) {
	res, _ := json.Marshal(queue)
	rCli.Set(ctx, "standart", string(res), 0)
}

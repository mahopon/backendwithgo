package redis

import (
	"log"
	"os"
	"sync"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisClient interface {
	Set(key string, value any, ttl time.Duration) error
	Get(key string) (string, error)
	Del(keys ...string) error
}

type redisClient struct {
	client redis.Cmdable
	raw    *redis.Client
	mu     sync.RWMutex
}

var (
	clientInstance *redisClient
	once           sync.Once
)

func startClient(redisurl string) (*redisClient, error) {
	log.Println("Creating Redis client")
	opts, err := redis.ParseURL(redisurl)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	log.Println("Successfully created Redis client")
	clientInstance = &redisClient{client: client, raw: client}
	return clientInstance, nil
}

func GetClient() *redisClient {
	once.Do(func() {
		log.Println("Initialising connection to Redis")
		var err error
		clientInstance, err = startClient(os.Getenv("REDIS_URL"))
		if err != nil {
			panic(err)
		}
	})
	log.Println("Retrieved Redis Client")
	return clientInstance
}

func (instance *redisClient) CloseClient() {
	instance.raw.Close()
}

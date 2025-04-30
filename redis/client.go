package redis

import (
	"context"
	"errors"
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
	log.Println("Start connection to Redis")
	log.Printf("URL:%s", redisurl)
	opts, err := redis.ParseURL(redisurl)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	log.Println("Successfully connect to Redis")
	clientInstance = &redisClient{client: client, raw: client}
	clientInstance.startHealthMonitor()
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

func (instance *redisClient) ping() error {
	instance.mu.RLock()
	defer instance.mu.RUnlock()
	if instance.client == nil {
		return errors.New("redis client is nil")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return instance.client.Ping(ctx).Err()
}

func reconnect(redisurl string, maxRetries int, backoff time.Duration) (*redisClient, error) {
	var err error
	for i := range maxRetries {
		log.Printf("Attempting to connect to Redis (%d/%d)", i, maxRetries)
		client, err := startClient(redisurl)
		if err == nil && client != nil && client.ping() == nil {
			log.Println("Reconnected to Redis successfully")
			return client, nil
		}
		time.Sleep(backoff)
		backoff *= 2
	}
	return nil, err
}

func (instance *redisClient) startHealthMonitor() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			if err := instance.ping(); err != nil {
				log.Printf("Redis health check failed: %v", err)
				newClient, err := reconnect(os.Getenv("REDIS_URL"), 5, 1*time.Second)
				if err != nil {
					log.Printf("Redis reconnection failed: %v", err)
					continue
				}
				instance.mu.Lock()
				instance.client = newClient.client
				instance.raw = newClient.raw
				instance.mu.Unlock()
			}
		}
	}()
}

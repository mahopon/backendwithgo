package redis

import (
	"context"
	"log"
	"time"
)

var ctx context.Context = context.Background()

func (c *redisClient) Set(key string, value any, ttl time.Duration) error {
	log.Printf("Setting key %s with value %v with duration %v", key, value, ttl)
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		log.Printf("Error occurred during set: %+v", err)
		return err
	}
	log.Printf("Successfully set key %s with value %v with duration %v", key, value, ttl)
	return nil
}

func (c *redisClient) Get(key string) (string, error) {
	log.Printf("Getting key %s", key)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		log.Printf("Error occurred during get: %+v", err)
		return "", err
	}
	log.Printf("Successfully retrieved key %s", key)
	return val, nil
}

func (c *redisClient) Delete(key ...string) error {
	log.Printf("Deleting key %s", key)

	err := c.client.Del(ctx, key...).Err()
	if err != nil {
		log.Printf("Error occurred during delete: %+v", err)
		return err
	}
	log.Printf("Successfully deleted key %s", key)
	return nil
}

func (c *redisClient) Exists(key string) (bool, error) {
	log.Printf("Checking exists key %s", key)
	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		log.Printf("Error occurred during exists: %+v", err)
		return false, err
	}
	log.Printf("Successfully checked key %s, value %+v", key, result)
	return result > 0, nil
}

func (c *redisClient) HSet(key string, field string, value any, ttl time.Duration) (int64, error) {
	log.Printf("HSetting key %s, field %s with value %v, duration: %v", key, field, value, ttl)
	fieldCount, err := c.client.HSet(ctx, key, field, value, ttl).Result()
	if err != nil {
		log.Printf("Error occurred during HSet: %v", err)
		return 0, err
	}
	log.Printf("Successfully HSet key %s, field %s with value %v, duration: %v", key, field, value, ttl)
	return fieldCount, err
}

func (c *redisClient) HDelete(key string, field string) (int64, error) {
	log.Printf("HDeleting key %s, field %s", key, field)
	fieldCount, err := c.client.HDel(ctx, key, field).Result()
	if err != nil {
		log.Printf("Error occurred during HSet: %v", err)
		return 0, err
	}
	log.Printf("Successfully HDelete key %s, field %s", key, field)
	return fieldCount, err
}

func (c *redisClient) HGet(key string, field string) (string, error) {
	log.Printf("HGetting key %s, field %s", key, field)
	result, err := c.client.HGet(ctx, key, field).Result()
	if err != nil {
		log.Printf("Error occurred during HGet: %v", err)
		return "", err
	}
	log.Printf("Successfully HGet key %s, field %s, with value %s", key, field, result)
	return result, err
}

func (c *redisClient) HExists(key string, field string) (bool, error) {
	log.Printf("Checking HExists key %s, field %s", key, field)
	exists, err := c.client.HExists(ctx, key, field).Result()
	if err != nil {
		log.Printf("Error occurred during HGet: %v", err)
		return false, err
	}
	log.Printf("Successfully HExists key %s, field %s, %v", key, field, exists)
	return exists, err
}

func (c *redisClient) HGetAll(key string) (map[string]string, error) {
	log.Printf("HGet all fields of key %s", key)
	result, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		log.Printf("Error occurred during HGetAll: %v", err)
		return nil, err
	}
	log.Printf("Successfully HGet all fields of key %s", key)
	return result, nil
}

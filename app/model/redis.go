package model

import (
	"context"
	"errors"
	"github.com/go-redis/redis"
	"github.com/goccy/go-json"
	"log"
	"time"
)

func SetRedis(key, value string) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = Rdb.Set(context.Background(), key, jsonData, 1*time.Hour).Err()
	if err != nil {
		// 可以在这里添加适当的错误处理逻辑，比如记录日志
		log.Printf("Error setting key %s in Redis: %v", key, err)
	}
	return err
}

func Getsession(key string) (string, error) {
	val, err := Rdb.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		log.Printf("Error setting key %s in Redis: %v", key, err)
	}
	return val, err
}

func DelRedis(key string) (int64, error) {
	result, err := Rdb.Del(context.Background(), key).Result()
	return result, err
}

package storage

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type InMemoryStorageI interface {
	Set(key, value string, exp time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
}

type storageRedis struct {
	client *redis.Client
}

type RedisData struct {
	Password string `json:"password"`
	Code     string `json:"code"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func NewInMemoryStorage(rdb *redis.Client) InMemoryStorageI {
	return &storageRedis{
		client: rdb,
	}
}

func (rd *storageRedis) Set(key, value string, exp time.Duration) error {
	err := rd.client.Set(context.TODO(), key, value, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rd *storageRedis) Get(key string) (string, error) {
	val, err := rd.client.Get(context.TODO(), key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (rd *storageRedis) Del(key string) error {
	_, err := rd.client.Del(context.TODO(), key).Result()
	if err != nil {
		return err
	}
	return nil
}

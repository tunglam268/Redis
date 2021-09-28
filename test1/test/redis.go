package main

import (
	"encoding/json"
	"sync"
	"time"

	redis "github.com/go-redis/redis"
)

var (
	client = &redisClient{}
)

type redisClient struct {
	c *redis.Client
}

func initialize(wg *sync.WaitGroup) *redisClient {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
	})
	if err := c.Ping().Err(); err != nil {
		panic("Unbale to connect to redis" + err.Error())
	}
	client.c = c
	return client
}
func (client *redisClient) getKey(key string, src interface{}) error {
	val, err := client.c.Get(key).Result()
	if err == redis.Nil || err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), &src)
	if err != nil {
		return err
	}
	return nil
}
func (client *redisClient) setKey(key string, value interface{}, expiration time.Duration) error {
	cacheEntry, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = client.c.Set(key, cacheEntry, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

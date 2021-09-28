package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	redis "github.com/go-redis/redis"
)

type valueRedis struct {
	Key   string
	Value string
}

var (
	client = &redisClient{}
)

type redisClient struct {
	c *redis.Client
}

func initialize() *redisClient {
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

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	login := initialize()
	key := "1"
	value := &valueRedis{Key: "CROSSCHEK_TRANSACTION_DEBIT_27_09_2021", Value: "22222_0"}
	// value2 := &valueRedis{}
	go login.setKey(key, value, time.Minute*1)
	go login.getKey(key, value)
	log.Printf("Key : %s", value.Key)
	log.Printf("Value :%s", value.Value)

	wg.Wait()
}

package eago

import (
	"consistent"
	"fmt"
	"gopkg.in/redis.v3"
	"sync"
	"time"
)
const (
	KeyForCrawledUrls = "crawledurls"
)
var redisInit sync.Once
var DefaultRedisClient *RedisClient

func KeyForCrawlByDay() string {
	t := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s:%s", KeyForCrawledUrls, t)
}

type RedisClient struct {
	Clients        map[string]*redis.Client
	hash           *consistent.Consistent
}

func GetRedisClient() *RedisClient {
	redisInit.Do(func() {
		DefaultRedisClient = &RedisClient{
			Clients:        make(map[string]*redis.Client),
			hash:           consistent.New(),
		}
		// Init all the clients
		for k, v := range Configs.Redis {
			DefaultRedisClient.AddClient(k, v)
		}
	})
	return DefaultRedisClient
}

func (r *RedisClient) AddClient(name string, re *RedisInstance) {
	r.hash.Add(name)
	r.Clients[name] = redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     re.Host,
		DB:       re.DB,
		PoolSize: re.Pool,
	})
}

func (r *RedisClient) GetClient(key string) *redis.Client {
	res, err := r.hash.Get(key)
	if err != nil {
		Error.Println(err)
	}
	client, ok := r.Clients[res]
	if !ok {
		Error.Println("Get redis Instance Failed")
		return nil
	}
	return client
}

type RedisInstance struct {
	Host string
	DB   int64
	Pool int
}

package eago

import (
	"consistent"
	"fmt"
	"gopkg.in/redis.v3"
	"sync"
	"time"
)

var redisInit sync.Once

var DefaultRedisClient *RedisClient

const (
	KeyForCrawledUrls = "crawledurls"
)

func KeyForCrawlByDay() string {
	t := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s:%s", KeyForCrawledUrls, t)
}

type RedisClient struct {
	Clients        map[string]*redis.Client
	hash           *consistent.Consistent
	RedisInstances []*RedisInstance
}

func GetRedisClient() *RedisClient {
	redisInit.Do(func() {
		DefaultRedisClient = &RedisClient{
			Clients:        make(map[string]*redis.Client),
			hash:           consistent.New(),
			RedisInstances: Configs.RedisInstances,
		}
		// Init all the clients
		for _, v := range DefaultRedisClient.RedisInstances {
			DefaultRedisClient.AddClient(v)
		}
	})
	return DefaultRedisClient
}

func (r *RedisClient) AddClient(re *RedisInstance) {
	r.hash.Add(re.Name)
	r.Clients[re.Name] = redis.NewClient(&redis.Options{
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
	Name string
	Host string
	DB   int64
	Pool int
}

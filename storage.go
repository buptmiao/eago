package eago

import (
	"fmt"
	"time"
)

// You can customize the storage strategy in your application
// by implementing the interface Storage
type Storage interface {
	Store(resp *UrlResponse)
}

//By default, store the response into Redis.
type DefaultStore struct {
	*RedisClient
}

const (
	KeyForStore = "url:%s"
	Expiration  = time.Second * 3600 * 24 * 7
)

func KeyForUrlStore(url string) string {
	return fmt.Sprintf(KeyForStore, url)
}

func NewDefaultStore(r *RedisClient) *DefaultStore {
	res := &DefaultStore{
		RedisClient: r,
	}
	return res
}

func (d *DefaultStore) Store(resp *UrlResponse) {
	url := resp.Src.Url
	client := d.GetClient(url)
	key := KeyForUrlStore(url)
	client.Set(key, resp.Body, 0)
	client.Expire(key, Expiration)
	return
}

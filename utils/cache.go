package utils

import (
    "context"
    "os"
    "time"

    "github.com/go-redis/redis"
    "github.com/pkg/errors"
)

type Cache struct {
    client redis.UniversalClient
    ttl    time.Duration
}

var apqPrefix = os.Getenv("APQ_PREFIX")

func NewCache(redisAddress string, ttl time.Duration) (*Cache, error) {
    client := redis.NewClient(&redis.Options{
        Addr: redisAddress,
    })
    err := client.Ping().Err()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &Cache{
        client: client,
        ttl:    ttl,
    }, nil
}
func (c *Cache) Add(ctx context.Context, hash, query string) {
    c.client.Set(apqPrefix+hash, query, c.ttl)
}
func (c *Cache) Get(ctx context.Context, hash string) (string, bool) {
    s, err := c.client.Get(apqPrefix + hash).Result()
    if err != nil {
        return "", false
    }
    return s, true
}

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
func (c *Cache) Subscribe(channels ...string) *redis.PubSub {
    return c.client.Subscribe(channels...)
}
func (c *Cache) Publish(channel string, message interface{}) *redis.IntCmd {
    return c.client.Publish(channel, message)
}

func (c *Cache) SAdd(key, value string) error {
    return c.client.SAdd(key, value).Err()
}
func (c *Cache) LPush(key string, values ...interface{}) error {
    return c.client.LPush(key, values...).Err()
}
func (c *Cache) LRange(key string, start, stop int64) *redis.StringSliceCmd {
    return c.client.LRange(key, start, stop)
}
func (c *Cache) SMembers(key string) *redis.StringSliceCmd {
    return c.client.SMembers(key)
}
func (c *Cache) Add(ctx context.Context, hash, query string) {
    c.client.Set(hash, query, c.ttl)
}
func (c *Cache) Get(ctx context.Context, hash string) (string, bool) {
    s, err := c.client.Get(hash).Result()
    if err != nil {
        return "", false
    }
    return s, true
}

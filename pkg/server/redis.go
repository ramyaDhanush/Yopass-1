package server

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/z0x0z/yopass/pkg/yopass"
)

// NewRedis returns a new Redis database client
func NewRedis(url string) (Database, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(options)
	return &Redis{client, context.Background()}, nil
}

// Redis client
type Redis struct {
	client *redis.Client
	ctx context.Context
}

// Get key from Redis
func (r *Redis) Get(key string) (yopass.Secret, error) {
	var s yopass.Secret
	v, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return s, err
	}

	if err := json.Unmarshal([]byte(v), &s); err != nil {
		return s, err
	}

	return s, nil
}

// Put key to Redis
func (r *Redis) Put(key string, secret yopass.Secret) error {
	data, err := secret.ToJSON()
	if err != nil {
		return err
	}
	return r.client.Set(
		r.ctx,
		key,
		data,
		time.Duration(secret.Expiration)*time.Second,
	).Err()
}

// Delete key from Redis
func (r *Redis) Delete(key string) (bool, error) {
	err := r.client.Del(r.ctx, key).Err()

	if err == redis.Nil {
		return false, nil
	}

	return err == nil, err
}

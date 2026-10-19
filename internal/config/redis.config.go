package config

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Username string
	Password string
	Host     string
	Port     string
}

func NewRedisClient(username, password, host, port string) *RedisClient {
	return &RedisClient{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
	}
}

func (r *RedisClient) Connect() *redis.Client {
	return redis.NewClient(&redis.Options{
		Username: r.Username,
		Password: r.Password,
		Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
	})
}

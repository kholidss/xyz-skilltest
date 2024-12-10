package redisnative

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	Host               string
	Port               int
	DB                 int
	Username           string
	Password           string
	PoolSize           int
	ReadTimeout        int
	WriteTimeout       int
	PoolTimeout        int
	MinIdleConn        int
	IdleTimeout        int
	TLS                bool
	InsecureSkipVerify bool
}

func NewConnection(config Config) *Config {
	return &config
}

func (c *Config) Connect() *redis.Client {
	redisNativeCfg := redis.Options{
		Addr:         fmt.Sprintf("%s:%d", c.Host, c.Port),
		DB:           c.DB,
		PoolSize:     c.PoolSize,
		PoolTimeout:  time.Duration(c.PoolTimeout) * time.Second,
		MinIdleConns: c.MinIdleConn,
		ReadTimeout:  time.Duration(c.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.WriteTimeout) * time.Second,
		Password:     c.Password,
	}

	if c.TLS {
		redisNativeCfg.TLSConfig = &tls.Config{
			InsecureSkipVerify: c.InsecureSkipVerify,
		}
	}

	r := redis.NewClient(&redisNativeCfg)

	if r == nil {
		logger.Fatal("nil open connection redis")
	}

	ping := r.Ping(context.Background())

	if err := ping.Err(); err != nil {
		logger.Fatal("error ping connection redis")
	}

	return r
}

package bootstrap

import (
	"github.com/kholidss/xyz-skilltest/pkg/config"
	"github.com/kholidss/xyz-skilltest/pkg/database/redis_native"
	"github.com/redis/go-redis/v9"
)

// RegistryRedisNative initiate redis session
func RegistryRedisNative(cfg *config.Config) *redis.Client {
	manager := redisnative.NewConnection(redisnative.Config{
		Host:               cfg.RedisConfig.Host,
		Port:               cfg.RedisConfig.Port,
		DB:                 cfg.RedisConfig.DB,
		Username:           cfg.RedisConfig.Username,
		Password:           cfg.RedisConfig.Password,
		PoolSize:           cfg.RedisConfig.PoolSize,
		ReadTimeout:        cfg.RedisConfig.ReadTimeout,
		WriteTimeout:       cfg.RedisConfig.WriteTimeout,
		PoolTimeout:        cfg.RedisConfig.PoolTimeout,
		MinIdleConn:        cfg.RedisConfig.MinIdleConn,
		IdleTimeout:        cfg.RedisConfig.IdleTimeout,
		TLS:                cfg.RedisConfig.TLS,
		InsecureSkipVerify: cfg.RedisConfig.InsecureSkipVerify,
	})
	return manager.Connect()
}

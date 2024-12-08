package bootstrap

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kholidss/xyz-skilltest/pkg/config"
	"github.com/kholidss/xyz-skilltest/pkg/database/mongodb"
	"github.com/kholidss/xyz-skilltest/pkg/database/mysql"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
)

func RegistryMySQLDatabase(cfg *config.Config) *mysql.DB {
	manager := mysql.NewConnection(mysql.Config{
		Host:            cfg.DatabaseConfig.DBHost,
		Port:            cfg.DatabaseConfig.DBPort,
		Name:            cfg.DatabaseConfig.DBName,
		User:            cfg.DatabaseConfig.DBUser,
		Password:        cfg.DatabaseConfig.DBPassword,
		MaxOpenConn:     cfg.DatabaseConfig.MaxOpenConn,
		MaxIdleConn:     cfg.DatabaseConfig.MaxIdleConn,
		MaxConnLifetime: cfg.DatabaseConfig.MaxConnLifetime,
		MaxIdleTime:     cfg.DatabaseConfig.MaxIdleTime,
		CAPath:          cfg.DatabaseConfig.CAPath,
		ClientCertPath:  cfg.DatabaseConfig.ClientCertPath,
		ClientKeyPath:   cfg.DatabaseConfig.ClientKeyPath,
	})

	var (
		db  *sqlx.DB
		err error
	)

	if cfg.DatabaseConfig.TLS {
		db, err = manager.Connect(cfg, mysql.WithTLSTransport())
	} else {
		db, err = manager.Connect(cfg)
	}

	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to connect to mysql database: %v", err))
	}

	return mysql.New(db, false, cfg.DatabaseConfig.DBName)
}

func RegistryMongoDB(cfg *config.Config) *mongodb.DB {
	db, err := mongodb.Connect(cfg)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	return mongodb.New(db, false, cfg.DatabaseConfig.DBName)
}

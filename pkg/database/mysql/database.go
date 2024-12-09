package mysql

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/kholidss/xyz-skilltest/internal/consts"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
	"github.com/kholidss/xyz-skilltest/pkg/util"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kholidss/xyz-skilltest/pkg/config"
)

type Config struct {
	Host            string
	Port            int
	Name            string
	User            string
	Password        string
	MaxOpenConn     int
	MaxIdleConn     int
	MaxConnLifetime int
	MaxIdleTime     int
	TLS             bool
	CAPath          string
	ClientCertPath  string
	ClientKeyPath   string
}

// DBConfigOption is type of option
type DBConfigOption func(c *Config)

// WithTLSTransport is option to connect db use TLS transport
func WithTLSTransport() DBConfigOption {
	return func(c *Config) {
		c.TLS = true
	}
}

func connect(cnf *config.Config) (*sqlx.DB, error) {
	var (
		err      error
		dbConfig = cnf.DatabaseConfig
	)

	conf, err := NewMysqlConfig(cnf)
	if err != nil {
		logger.Fatal("Failed to create mysql config")
	}

	dsn := conf.FormatDSN()
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	db.SetMaxIdleConns(dbConfig.MaxIdleConn)
	db.SetMaxOpenConns(dbConfig.MaxOpenConn)
	db.SetConnMaxLifetime(time.Duration(dbConfig.MaxConnLifetime) * time.Hour)
	db.SetConnMaxIdleTime(time.Duration(dbConfig.MaxIdleTime) * time.Hour)

	return db, nil
}

func NewMysqlConfig(cnf *config.Config) (*mysql.Config, error) {
	dbConfig := cnf.DatabaseConfig
	conf := mysql.NewConfig()
	conf.Net = "tcp"
	conf.Addr = fmt.Sprintf("%v:%v", dbConfig.DBHost, dbConfig.DBPort)
	conf.User = dbConfig.DBUser
	conf.Passwd = dbConfig.DBPassword
	conf.DBName = dbConfig.DBName

	tlsConfig, err := dbConfig.TlsConfig(cnf.AppEnv)
	if err != nil {
		return nil, err
	}

	if tlsConfig != nil {
		if err = mysql.RegisterTLSConfig("custom", tlsConfig); err != nil {
			return nil, err
		}

		conf.TLSConfig = "custom"
	}

	return conf, nil
}

func ConnectDatabase(cnf *config.Config) (*sqlx.DB, error) {
	db, err := connect(cnf)

	if err != nil {
		return nil, err
	}
	return db, nil
}

func (c *Config) connecting(mysqlCfg *mysql.Config) (*sqlx.DB, error) {
	dsn := mysqlCfg.FormatDSN()
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	db.SetMaxIdleConns(c.MaxIdleConn)
	db.SetMaxOpenConns(c.MaxOpenConn)
	db.SetConnMaxLifetime(time.Duration(c.MaxConnLifetime) * time.Hour)
	db.SetConnMaxIdleTime(time.Duration(c.MaxIdleTime) * time.Hour)

	return db, nil
}

func NewConnection(cfgDb Config) *Config {
	return &cfgDb
}

func (c *Config) Connect(cfg *config.Config, opts ...DBConfigOption) (*sqlx.DB, error) {
	for _, opt := range opts {
		opt(c)
	}

	mysqlCfg := mysql.NewConfig()
	mysqlCfg.Net = "tcp"
	mysqlCfg.Addr = fmt.Sprintf("%v:%v", c.Host, c.Port)
	mysqlCfg.User = c.User
	mysqlCfg.Passwd = c.Password
	mysqlCfg.DBName = c.Name
	mysqlCfg.ParseTime = true

	if c.TLS {
		tlsConfig, err := c.setTLS(cfg.AppEnv)
		if err != nil {
			return nil, err
		}
		if tlsConfig != nil {
			if err = mysql.RegisterTLSConfig("custom", tlsConfig); err != nil {
				return nil, err
			}

			mysqlCfg.TLSConfig = "custom"
		}
	}

	return c.connecting(mysqlCfg)
}

func (c *Config) setTLS(env string) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: util.EnvironmentTransform(env) != consts.AppProduction,
	}

	pool := x509.NewCertPool()
	pem, err := os.ReadFile(c.CAPath)
	if err != nil {
		return nil, err
	}

	if ok := pool.AppendCertsFromPEM(pem); !ok {
		return nil, errors.New("unable to append root cert to pool")
	}

	cert, err := tls.LoadX509KeyPair(c.ClientCertPath, c.ClientKeyPath)
	if err != nil {
		return nil, err
	}

	tlsConfig.Certificates = []tls.Certificate{cert}

	return tlsConfig, nil
}

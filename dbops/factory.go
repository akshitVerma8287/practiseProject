package dbops

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type DatabaseType string

const (
	RedisType    = "redis"
	PostgresType = "postgres"
)

type DatabaseConfig struct {
	Type               DatabaseType
	Host               string
	Port               string
	Username           string
	Password           string
	Database           string
	PoolSize           int
	MaxConnection      int
	MaxIdleConnections int
}

type DatabaseConnector interface {
	Connect() error
	Close() error
	GetRepository() interface{}
}


type RedisConnector struct {
	client *redis.Client
	config DatabaseConfig
}

func (r *RedisConnector) Connect() error {

	opts := &redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			r.config.Host,
			r.config.Port,
		),

		Username: r.config.Username,
		Password: r.config.Password,
		PoolSize: r.config.PoolSize,
	}

	r.client = redis.NewClient(opts)

	_, err := r.client.Ping(
		context.Background(),
	).Result()

	if err != nil {
		return fmt.Errorf(
			"failed to connect redis: %w",
			err,
		)
	}

	logrus.Printf(
		"Connected Redis at %s:%s",
		r.config.Host,
		r.config.Port,
	)

	return nil
}

func (r *RedisConnector) Close() error {

	if r.client != nil {
		return r.client.Close()
	}

	return nil
}

func (r *RedisConnector) GetRepository() interface{} {

	return NewRedisRepository(r.client)
}

///////////////////////////////////////////////////////////
////////////////// POSTGRES CONNECTOR /////////////////////
///////////////////////////////////////////////////////////

type PostgresConnector struct {
	db     *sql.DB
	config DatabaseConfig
}

func (p *PostgresConnector) Connect() error {

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		p.config.Host,
		p.config.Port,
		p.config.Username,
		p.config.Password,
		p.config.Database,
	)

	db, err := sql.Open(
		"postgres",
		connectionString,
	)

	if err != nil {
		return fmt.Errorf(
			"failed postgres connection: %w",
			err,
		)
	}

	err = db.PingContext(
		context.Background(),
	)

	if err != nil {
		return fmt.Errorf(
			"failed postgres ping: %w",
			err,
		)
	}

	db.SetMaxOpenConns(
		p.config.MaxConnection,
	)

	db.SetMaxIdleConns(
		p.config.MaxIdleConnections,
	)

	db.SetConnMaxLifetime(
		24 * time.Hour,
	)

	p.db = db

	logrus.Printf(
		"Connected Postgres at %s:%s",
		p.config.Host,
		p.config.Port,
	)

	return nil
}

func (p *PostgresConnector) Close() error {

	if p.db != nil {
		return p.db.Close()
	}

	return nil
}

func (p *PostgresConnector) GetRepository() interface{} {

	return NewPostgresRepository(p.db)
}

///////////////////////////////////////////////////////////
//////////////////// FACTORY //////////////////////////////
///////////////////////////////////////////////////////////

type DatabaseConnectorFactory struct{}

func NewDatabaseConnectorFactory() *DatabaseConnectorFactory {
	return &DatabaseConnectorFactory{}
}

func (f *DatabaseConnectorFactory) GetConnector(config DatabaseConfig) (DatabaseConnector, error) {
	var connector DatabaseConnector

	switch config.Type {
	case RedisType:
		connector = &RedisConnector{config: config}
	case PostgresType:
		connector = &PostgresConnector{config: config}
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}

	err := connector.Connect()			// This is method connector 
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", config.Type, err)
	}

	return connector, nil
}
package dbops

import (
	"context"
	"database/sql"
	"time"
)


type RedisRepository interface {

	Set(key string,	value interface{}, expiration time.Duration,) error
	SetCtx(	ctx context.Context, key string, value interface{},	expiration time.Duration,) error
	Get(key string) (string, error)
	GetCtx(ctx context.Context, key string,) (string, error)
	Delete(key string) error
	Exists(key string) (int64, error)
	GetStatus() error
}

type PostgresRepository interface {

	GetStatus() error
	Insert(query string, args ...any) (*sql.Rows, error)
	Update(query string, args ...any) (*sql.Rows, error)
	Fetch(query string, args ...any) (*sql.Rows, error)
	Delete(query string, args ...any) (*sql.Rows, error)
	BeginTx() (*sql.Tx, error)
	BeginTxWithOptions(opts *sql.TxOptions) (*sql.Tx, error)
}

//////////////////// FACTORY ////////////////////

type DatabaseFactory interface {

	CreateRedisRepository()(RedisRepository,error)
	CreatePostgresRepository()(PostgresRepository,error)
}
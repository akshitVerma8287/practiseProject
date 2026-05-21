package dbops

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(
	client *redis.Client,
) RedisRepository {

	return &redisRepository{
		client: client,
	}
}

func (r *redisRepository) GetStatus() error {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5000*time.Millisecond,
	)

	defer cancel()

	_, err := r.client.Ping(ctx).Result()

	if err != nil {
		return err
	}

	return nil
}

///////////////////////////////////////////////////////////
//////////////////////// SET //////////////////////////////
///////////////////////////////////////////////////////////

func (r *redisRepository) Set(
	key string,
	value interface{},
	expiration time.Duration,
) error {

	return r.SetCtx(
		context.Background(),
		key,
		value,
		expiration,
	)
}

func (r *redisRepository) SetCtx(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) error {

	var strValue string

	switch v := value.(type) {

	case string:

		strValue = v

	case []byte:

		strValue = string(v)

	default:

		jsonBytes, err := json.Marshal(v)

		if err != nil {
			return err
		}

		strValue = string(jsonBytes)
	}

	return r.client.Set(
		ctx,
		key,
		strValue,
		expiration,
	).Err()
}

///////////////////////////////////////////////////////////
//////////////////////// SETNX ////////////////////////////
///////////////////////////////////////////////////////////

func (r *redisRepository) SetNX(
	key string,
	value interface{},
	expiration time.Duration,
) (bool, error) {

	return r.SetNXCtx(
		context.Background(),
		key,
		value,
		expiration,
	)
}

func (r *redisRepository) SetNXCtx(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) (bool, error) {

	var strValue string

	switch v := value.(type) {

	case string:

		strValue = v

	case []byte:

		strValue = string(v)

	default:

		jsonBytes, err := json.Marshal(v)

		if err != nil {
			return false, err
		}

		strValue = string(jsonBytes)
	}

	return r.client.SetNX(
		ctx,
		key,
		strValue,
		expiration,
	).Result()
}

///////////////////////////////////////////////////////////
//////////////////////// GET //////////////////////////////
///////////////////////////////////////////////////////////

func (r *redisRepository) Get(
	key string,
) (string, error) {

	return r.GetCtx(
		context.Background(),
		key,
	)
}

func (r *redisRepository) GetCtx(
	ctx context.Context,
	key string,
) (string, error) {

	result, err := r.client.Get(
		ctx,
		key,
	).Result()

	if err != nil {

		if err == redis.Nil {
			return "", nil
		}

		return "", err
	}

	return result, nil
}

///////////////////////////////////////////////////////////
/////////////////////// DELETE ////////////////////////////
///////////////////////////////////////////////////////////

func (r *redisRepository) Delete(
	key string,
) error {

	return r.client.Del(
		context.Background(),
		key,
	).Err()
}

///////////////////////////////////////////////////////////
//////////////// DEL IF VALUE MATCH ///////////////////////
///////////////////////////////////////////////////////////

func (r *redisRepository) DelIfValueMatch(
	ctx context.Context,
	key string,
	value string,
) (bool, error) {

	script := `
	if redis.call('get', KEYS[1]) == ARGV[1]
	then
		return redis.call('del', KEYS[1])
	else
		return 0
	end
	`

	result, err := r.client.Eval(
		ctx,
		script,
		[]string{key},
		value,
	).Result()

	if err != nil {
		return false, err
	}

	return result.(int64) == 1, nil
}

func (r *redisRepository) Exists(key string,) (int64, error) {
	return r.client.Exists(
		context.Background(),
		key,
	).Result()
}

func (r *redisRepository) Increment(key string,) int64 {
	return r.client.Incr(
		context.Background(),
		key,
	).Val()
}

func (r *redisRepository) HSet(key string,field string,	value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.HSet(
		context.Background(),
		key,
		field,
		jsonValue,
	).Err()
}

func (r *redisRepository) HGet(key string,field string,) (string, error) {

	return r.client.HGet(
		context.Background(),
		key,
		field,
	).Result()
}

func (r *redisRepository) HGetAll(key string) (map[string]string, error) {

	return r.client.HGetAll(
		context.Background(),
		key,
	).Result()
}

func (r *redisRepository) HDelete(key string,fields ...string) error {

	return r.client.HDel(
		context.Background(),
		key,
		fields...,
	).Err()
}

func (r *redisRepository) HExists(key string,field string) (bool, error) {

	return r.client.HExists(
		context.Background(),
		key,
		field,
	).Result()
}

func (r *redisRepository) HKeys(key string,) ([]string, error) {

	return r.client.HKeys(
		context.Background(),
		key,
	).Result()
}
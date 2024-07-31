package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	rdb		*redis.Client
}

func NewRedisStorage(addr string, password string, db int) (*RedisStorage, error) {
	newRdb :=redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password,
		DB: db,
	})

	ctx := context.TODO()
	if err := newRdb.Ping(ctx).Err() ; err != nil {
		return nil, err
	}

	return &RedisStorage{
		rdb: newRdb,
	}, nil
}
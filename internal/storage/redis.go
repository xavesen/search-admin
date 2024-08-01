package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/xavesen/search-admin/internal/models"
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

func (s *RedisStorage) GetNewUserId(ctx context.Context) (int, error) {
	id, err := s.rdb.Incr(ctx, "user_id_seq").Result()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func (s *RedisStorage) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	id, err := s.GetNewUserId(ctx)
	if err != nil {
		return nil, err
	}

	user.Id = id

	redisUserId := fmt.Sprintf("user:%d", id)

	err = s.rdb.JSONSet(ctx, redisUserId, ".", user).Err()
	if err != nil {
		return nil, err
	}

	return user, nil
}
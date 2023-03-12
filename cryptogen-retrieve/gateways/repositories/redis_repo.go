package repositories

import (
	"context"
	"cryptogen-retrieve/domain"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)

const (
	queueName      = "crypto_queue"
	redisUrl       = "localhost:6379"
	connectionType = "tcp"
	retryInterval  = time.Second * 5
)

type RedisRepo struct {
	redisPool *redis.Pool
}

var _ CryptoMetadataStorage = (*RedisRepo)(nil)

func NewRedisRepo() *RedisRepo {
	return &RedisRepo{
		redisPool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return redis.Dial(connectionType, redisUrl)
			},
		},
	}
}

func (r *RedisRepo) SaveDataList(metadata []*domain.CryptoMetadata) error {
	conn := r.redisPool.Get()
	defer conn.Close()

	// If queue is not empty, delete all data
	queueLength, err := redis.Int(conn.Do("LLEN", queueName))
	if err != nil {
		return fmt.Errorf("failed to get queue length: %v", err)
	}
	if queueLength > 0 {
		if _, err := conn.Do("DEL", queueName); err != nil {
			return fmt.Errorf("failed to delete queue: %v", err)
		}
	}

	// Add metadata to the queue
	var args []interface{}
	args = append(args, queueName)
	for _, data := range metadata {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %v", err)
		}
		args = append(args, jsonData)
	}

	err = r.retry(func() error {
		_, err = conn.Do("RPUSH", args...)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to add metadata to queue: %v", err)
	}

	return nil
}

func (r *RedisRepo) retry(f func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), retryInterval)
	defer cancel()

	for {
		err := f()
		if err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("operation timed out: %v", err)
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}

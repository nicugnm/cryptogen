package repositories

import (
	"context"
	"cryptogen-retrieve/domain"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

const (
	queueName      = "crypto_queue"
	dbKey          = "crypto_key"
	workers        = 5 // Number of worker goroutines
	redisUrl       = "localhost:6379"
	connectionType = "tcp"
)

type RedisRepo struct {
	redisPool *redis.Pool
}

var _ domain.CryptoMetadataStorage = (*RedisRepo)(nil)

func NewRedisRepo() *RedisRepo {
	return &RedisRepo{
		redisPool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return redis.Dial(connectionType, redisUrl)
			},
		},
	}
}

func (r RedisRepo) SaveData(metadata *domain.CryptoMetadata) error {
	pool := r.redisPool

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(pool, &wg)
	}

	// Push data to Redis queue
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("LPUSH", queueName, metadata)
	if err != nil {
		fmt.Println("Error pushing data to Redis queue:", err)
		return err
	}

	// Wait for all workers to finish
	wg.Wait()

	fmt.Println("All workers finished")

	return nil
}

func (r RedisRepo) SaveDataList(metadataList []*domain.CryptoMetadata) error {
	pool := r.redisPool

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(pool, &wg)
	}

	// Push data to Redis queue
	conn := pool.Get()
	defer conn.Close()

	for _, metadata := range metadataList {
		_, err := conn.Do("LPUSH", queueName, metadata)
		if err != nil {
			fmt.Println("Error pushing data to Redis queue:", err)
			return err
		}
	}

	// Wait for all workers to finish
	wg.Wait()

	fmt.Println("All workers finished")

	return nil
}

func worker(pool *redis.Pool, wg *sync.WaitGroup) {
	defer wg.Done()

	// Process data from Redis queue
	for {
		// Pop data from Redis queue with retries
		data, err := popFromQueueWithRetry(pool)
		if err != nil {
			fmt.Println("Error popping data from Redis queue:", err)
			return
		}

		// Save data to Redis database with retries
		err = saveToDbWithRetry(pool, data)
		if err != nil {
			fmt.Println("Error saving data to Redis:", err)
			return
		}

		fmt.Println("Saved data to Redis:", data)
		time.Sleep(time.Second)
	}
}

func popFromQueueWithRetry(pool *redis.Pool) (string, error) {
	var data string
	err := backoff.Retry(func() error {
		// Connect to Redis with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err := pool.GetContext(ctx)

		if err != nil {
			return err
		}

		defer conn.Close()

		// Pop data from Redis queue
		reply, err := conn.Do("RPOP", queueName)
		if err != nil {
			return err
		}

		// Check if there is data in the queue
		if reply == nil {
			return backoff.Permanent(fmt.Errorf("no more data in the Redis queue"))
		}

		// Cast reply to string
		strData, ok := reply.(string)
		if !ok {
			return fmt.Errorf("invalid data type in Redis queue")
		}

		// Assign the string value to data
		data = strData

		return nil
	}, backoff.NewExponentialBackOff())

	if err != nil {
		return "", err
	}

	return data, nil
}

func saveToDbWithRetry(pool *redis.Pool, data string) error {
	err := backoff.Retry(func() error {
		// Connect to Redis with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conn, err := pool.GetContext(ctx)
		if err != nil {
			return err
		}
		defer conn.Close()

		// Save data to Redis database
		_, err = conn.Do("SET", dbKey, data)

		return nil
	}, backoff.NewExponentialBackOff())

	return err
}

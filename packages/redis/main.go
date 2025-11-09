package redis

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/Grs2080w/worker-knoteq/packages/supa/token"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func New() (*Redis, error) {
	
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Redis{Client: client}, nil
}

func (r *Redis) Close() {
	r.Client.Close()
}

// Get: Get value from redis, if not found, call function and set value on redis
func (r *Redis) Get(ctx context.Context, key string, getValueFromDb func(key string) (string, error)) (string, error) {

	result, err := r.Client.Get(ctx, key).Result()

	if err == nil {
		return result, nil
	}

	if err != redis.Nil {
		return "", err
	}

	newValue, err := getValueFromDb(key)

	if err != nil {
		return "", err
	}

	if err := r.Client.Set(ctx, key, newValue, 1*time.Hour).Err(); err != nil {
		return "", err
	}

	return newValue, nil

}

// Get: Get value from redis, if not found, call function and set value on redis
func (r *Redis) GetToken(ctx context.Context, key string, getValueFromDb func() (token.Token, error)) (token.Token, error) {

	result, err := r.Client.Get(ctx, key).Result()

	if err == nil {

		var tok token.Token
		err := json.Unmarshal([]byte(result), &tok)
		if err != nil {
			return token.Token{}, err
		}

		return tok, nil
	}

	if err != redis.Nil {
		return token.Token{}, err
	}

	newValue, err := getValueFromDb()

	if err != nil {
		return token.Token{}, err
	}

	tok, err := json.Marshal(newValue)
	if err != nil {
		return token.Token{}, err
	}

	if err := r.Client.Set(ctx, key, string(tok), 1*time.Hour).Err(); err != nil {
		return token.Token{}, err
	}

	return newValue, nil

}


// InvalidateKey: Invalidate key on redis
func (r *Redis) InvalidateKey(ctx context.Context, key string) error {

	if err := r.Client.Del(ctx, key).Err(); err != nil {
		return err
	}

	return nil

}


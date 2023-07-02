package redisclient

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	client *redis.Client
	ctx    context.Context
}

func Default() *redisClient {
	return &redisClient{redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}), context.Background()}
}

func (r *redisClient) GetAnimal(id *string) (string, bool) {
	animal, err := r.client.Get(r.ctx, *id).Result()
	if err != nil {
		return "", false
	}
	return animal, true
}

func (r *redisClient) SetAnimal(key *string, json string) error {
	return r.client.Set(r.ctx, *key, json, 0).Err()
}

func (r *redisClient) IsAnimalPresent(key *string) bool {
	intcmd := r.client.Exists(r.ctx, *key)
	return intcmd.Val() != 1
}

func (r *redisClient) DeleteAnimal(key *string) bool {
	intcmd := r.client.Del(r.ctx, *key)
	return intcmd.Val() != 1
}

func (r *redisClient) SetPicture(key *string, img string) error {
	return r.client.Set(r.ctx, *key, img, 0).Err()
}

func (r *redisClient) GetPicture(key *string) (string, bool) {
	img, err := r.client.Get(r.ctx, *key).Result()
	fmt.Printf("Retrieving %s... \n", *key)
	if err != nil {
		return "", false
	}
	return img, true
}

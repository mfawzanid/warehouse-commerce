package repository

import (
	"context"
	"fmt"
	"mfawzanid/warehouse-commerce/core/entity"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepositoryInterface interface {
	LockOrderProduct(ctx context.Context, req *entity.LockOrderProductRequest) error
	InvalidateLockOrderProduct(ctx context.Context, req *entity.LockOrderProductRequest) error
	GetReservedProductQuantity(ctx context.Context, productId, warehouseId string) (int, error)
}

type redisRepository struct {
	redisClient *redis.Client
}

func NewRedisRepository(redisClient *redis.Client) RedisRepositoryInterface {
	return &redisRepository{
		redisClient: redisClient,
	}
}

func (r *redisRepository) LockOrderProduct(ctx context.Context, req *entity.LockOrderProductRequest) error {
	key := generateProductReservedKey(req)
	expiration := entity.OrderExpireTimeInMinute * time.Minute

	_, err := r.redisClient.Set(ctx, key, req.Quantity, expiration).Result()
	if err != nil {
		return fmt.Errorf("error cache repo lock order product: %v", err.Error())
	}

	return nil
}

func (r *redisRepository) InvalidateLockOrderProduct(ctx context.Context, req *entity.LockOrderProductRequest) error {
	key := generateProductReservedKey(req)

	if err := r.redisClient.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("error cache repo invalidate lock order product: %v", err.Error())
	}

	return nil
}

func (r *redisRepository) GetReservedProductQuantity(ctx context.Context, productId, warehouseId string) (int, error) {
	key := generateAllProductReservedKey(productId, warehouseId)

	keys, err := r.redisClient.Keys(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("error cache repo getting reserved product quantity: %v", err)
	}

	sum := 0
	for _, key := range keys {
		val, _ := r.redisClient.Get(ctx, key).Int()
		sum += val
	}

	return sum, nil
}

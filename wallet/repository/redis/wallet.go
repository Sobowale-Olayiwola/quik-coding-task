package redis

import (
	"context"
	"encoding/json"
	"quik/domain"

	"github.com/go-redis/redis/v8"
)

type redisWalletInMemoryDB struct {
	db *redis.Client
}

func NewRedisInMemoryDB(redisClient *redis.Client) domain.WalletInMemoryDB {
	return &redisWalletInMemoryDB{redisClient}
}

func (r *redisWalletInMemoryDB) Get(ctx context.Context, id string) (domain.Wallet, error) {
	var wallet domain.Wallet
	val, err := r.db.Get(ctx, id).Result()
	if err == redis.Nil {
		return domain.Wallet{}, domain.ErrKeyNotFound
	} else if err != nil {
		return domain.Wallet{}, err
	}
	json.Unmarshal([]byte(val), &wallet)
	return wallet, nil
}

func (r *redisWalletInMemoryDB) Set(ctx context.Context, id string, wallet *domain.Wallet) error {
	data, err := json.Marshal(wallet)
	if err != nil {
		return err
	}
	err = r.db.Set(ctx, id, string(data), 0).Err()
	return err
}

func (r *redisWalletInMemoryDB) Delete(ctx context.Context, id string) error {
	err := r.db.Del(ctx, id).Err()
	return err
}

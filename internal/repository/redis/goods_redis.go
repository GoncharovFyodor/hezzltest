package redis

import (
	"context"
	"encoding/json"
	"github.com/GoncharovFyodor/hezzltest/internal/models"
	"github.com/redis/go-redis/v9"
	"time"
)

type GoodsRepository struct {
	rdb *redis.Client
}

func NewGoodsRepository(rdb *redis.Client) *GoodsRepository {
	return &GoodsRepository{rdb: rdb}
}

func (repo *GoodsRepository) GetByIDCtx(ctx context.Context, key string) (*models.Good, error) {
	goodBytes, err := repo.rdb.Get(ctx, "good:"+key).Bytes()

	if err != nil {
		return nil, err
	}

	var good *models.Good

	if err = json.Unmarshal(goodBytes, &good); err != nil {
		return nil, err
	}

	return good, nil
}

func (repo *GoodsRepository) SetGoodCtx(ctx context.Context, key string, seconds int, good *models.Good) error {
	goodBytes, err := json.Marshal(good)

	if err != nil {
		return err
	}

	if err := repo.rdb.Set(ctx, "good:"+key, goodBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return err
	}

	return nil
}

func (repo *GoodsRepository) DeleteGoodCtx(ctx context.Context, key string) error {
	if err := repo.rdb.Del(ctx, "good:"+key).Err(); err != nil {
		return err
	}

	return nil
}

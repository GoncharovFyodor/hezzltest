package repository

import (
	"context"
	"database/sql"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/GoncharovFyodor/hezzltest/internal/domain"
	"github.com/GoncharovFyodor/hezzltest/internal/models"
	"github.com/GoncharovFyodor/hezzltest/internal/repository/clickhouse"
	"github.com/GoncharovFyodor/hezzltest/internal/repository/postgres"
	redisRepo "github.com/GoncharovFyodor/hezzltest/internal/repository/redis"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	GoodsPostgres
	GoodsRedis
	GoodsClickhouse
}

type GoodsPostgres interface {
	Get(ctx context.Context, limit, offset int) (domain.GoodList, error)
	Create(ctx context.Context, projectID int, input domain.CreateGoodRequest) (models.Good, error)
	Update(ctx context.Context, projectID, ID int, input domain.UpdateGoodRequest) (models.Good, error)
	Delete(ctx context.Context, projectID, ID int) (models.DeletedGood, error)
	Reprioritize(ctx context.Context, projectID, ID int, input domain.ReprioritizeRequest) (models.GoodPriorities, error)
}

type GoodsRedis interface {
	GetByIDCtx(ctx context.Context, key string) (*models.Good, error)
	SetGoodCtx(ctx context.Context, key string, seconds int, good *models.Good) error
	DeleteGoodCtx(ctx context.Context, key string) error
}

type GoodsClickhouse interface {
	InsertGoods(ctx context.Context, rows []domain.Good) error
}

func NewRepository(db *sql.DB, rdb *redis.Client, clickConn driver.Conn) *Repository {
	return &Repository{
		GoodsPostgres:   postgres.NewGoodsRepository(db),
		GoodsRedis:      redisRepo.NewGoodsRepository(rdb),
		GoodsClickhouse: clickhouse.NewGoodsRepository(clickConn),
	}
}

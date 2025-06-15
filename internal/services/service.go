package services

import (
	"context"
	"github.com/GoncharovFyodor/hezzltest/internal/config"
	"github.com/GoncharovFyodor/hezzltest/internal/domain"
	"github.com/GoncharovFyodor/hezzltest/internal/models"
	"github.com/GoncharovFyodor/hezzltest/internal/nats"
	"github.com/GoncharovFyodor/hezzltest/internal/repository"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	Goods
}

type Goods interface {
	GetGoods(ctx context.Context, limit, offset int) (domain.GoodList, error)
	CreateGood(ctx context.Context, projectID int, input domain.CreateGoodRequest) (models.Good, error)
	UpdateGood(ctx context.Context, projectID, id int, input domain.UpdateGoodRequest) (models.Good, error)
	ReprioritizeGood(ctx context.Context, projectID, id int, input domain.ReprioritizeRequest) (models.GoodPriorities, error)
	DeleteGood(ctx context.Context, projectID, id int) (models.DeletedGood, error)
}

func NewService(cfg *config.Config, log *log.Logger, repo *repository.Repository, producer *nats.Producer) *Service {
	return &Service{Goods: NewGoodsService(cfg, log, repo, producer)}
}

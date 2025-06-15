package services

import (
	"context"
	"github.com/GoncharovFyodor/hezzltest/internal/config"
	"github.com/GoncharovFyodor/hezzltest/internal/domain"
	"github.com/GoncharovFyodor/hezzltest/internal/models"
	"github.com/GoncharovFyodor/hezzltest/internal/nats"
	"github.com/GoncharovFyodor/hezzltest/internal/repository"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type GoodsService struct {
	cfg      *config.Config
	log      *log.Logger
	repo     *repository.Repository
	producer *nats.Producer
}

func NewGoodsService(cfg *config.Config, log *log.Logger, repo *repository.Repository, producer *nats.Producer) *GoodsService {
	return &GoodsService{cfg: cfg, log: log, repo: repo, producer: producer}
}

func (s *GoodsService) GetGoods(ctx context.Context, limit, offset int) (domain.GoodList, error) {
	return s.repo.GoodsPostgres.Get(ctx, limit, offset)
}

func (s *GoodsService) CreateGood(ctx context.Context, projectID int, input domain.CreateGoodRequest) (models.Good, error) {
	g, err := s.repo.GoodsPostgres.Create(ctx, projectID, input)

	err = s.producer.Publish(ctx, g)

	if err != nil {
		s.log.Infof("не удалось отправить товар в nats: %v", err)
	}

	return g, nil
}

func (s *GoodsService) UpdateGood(ctx context.Context, projectID, id int, input domain.UpdateGoodRequest) (models.Good, error) {
	g, err := s.repo.GoodsPostgres.Update(ctx, projectID, id, input)

	if err != nil {
		return models.Good{}, err
	}

	// Инвалидация кеша
	if err = s.repo.GoodsRedis.DeleteGoodCtx(ctx, strconv.Itoa(id)); err != nil {
		s.log.Infof("не удалось удалить товар из redis: %v", err)
	}

	err = s.producer.Publish(ctx, g)

	if err != nil {
		s.log.Infof("не удалось отправить товар в nats: %v", err)
	}

	return g, nil
}

func (s *GoodsService) ReprioritizeGood(ctx context.Context, projectID, id int, input domain.ReprioritizeRequest) (models.GoodPriorities, error) {
	goods, err := s.repo.GoodsPostgres.Reprioritize(ctx, projectID, id, input)

	if err != nil {
		return models.GoodPriorities{}, err
	}

	for _, good := range goods.Priorities {
		// Инвалидация кеша
		if err = s.repo.GoodsRedis.DeleteGoodCtx(ctx, strconv.Itoa(good.ID)); err != nil {
			s.log.Infof("не удалось удалить товар из redis: %v", err)
		}

		err = s.producer.Publish(ctx, good)

		if err != nil {
			s.log.Infof("не удалось отправить товар в nats: %v", err)
		}
	}

	return goods, nil
}

func (s *GoodsService) DeleteGood(ctx context.Context, projectID, id int) (models.DeletedGood, error) {
	good, err := s.repo.GoodsPostgres.Delete(ctx, projectID, id)

	if err != nil {
		return models.DeletedGood{}, err
	}

	// Инвалидация кеша
	if err = s.repo.GoodsRedis.DeleteGoodCtx(ctx, strconv.Itoa(id)); err != nil {
		s.log.Infof("не удалось удалить товар из redis: %v", err)
	}

	err = s.producer.Publish(ctx, good)

	if err != nil {
		s.log.Infof("не удалось отправить товар в nats: %v", err)
	}

	return good, nil
}

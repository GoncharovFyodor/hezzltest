package nats

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GoncharovFyodor/hezzltest/internal/config"
	"github.com/GoncharovFyodor/hezzltest/internal/domain"
	"github.com/GoncharovFyodor/hezzltest/internal/repository"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"time"
)

type Consumer struct {
	log    *log.Logger
	cfg    *config.Config
	stream jetstream.Stream
	repo   *repository.Repository
}

func NewConsumer(log *log.Logger, cfg *config.Config, stream jetstream.Stream, repo *repository.Repository) *Consumer {
	return &Consumer{log: log, cfg: cfg, stream: stream, repo: repo}
}

func (c *Consumer) Run(ctx context.Context) error {
	consumer, err := c.stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable: "processor",
	})

	if err != nil {
		c.log.Errorf("ошибка при создании консьюмера: %v", err)
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	for i := 0; i < c.cfg.Nats.WorkersCount; i++ {
		eg.Go(func() error {
			return c.RunConsumer(ctx, consumer)
		})
	}

	if err = eg.Wait(); err != nil {
		c.log.Errorf("ошибка при запуске консьюиера: %v", err)
	}

	return nil
}

func (c *Consumer) RunConsumer(ctx context.Context, consumer jetstream.Consumer) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			msgs, err := consumer.FetchNoWait(c.cfg.Nats.BatchSize)
			if err != nil {
				if errors.Is(err, nats.ErrTimeout) {
					continue
				}
				c.log.Errorf("ошибка при получении сообщения: %v", err)
				return err
			}

			var goods []domain.Good
			for msg := range msgs.Messages() {
				var good domain.Good
				if err = json.Unmarshal(msg.Data(), &good); err != nil {
					c.log.Errorf("ошибка десериализации сообщения: %v", err)
					msg.Nak()
					continue
				}
				goods = append(goods, good)
				msg.Ack()
			}

			if len(goods) > 0 {
				if err = c.repo.GoodsClickhouse.InsertGoods(ctx, goods); err != nil {
					c.log.Errorf("ошибка при добавлении товаров: %v", err)
				}
			}
		}
	}
}

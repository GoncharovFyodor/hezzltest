package app

import (
	"context"
	"github.com/GoncharovFyodor/hezzltest/internal/config"
	"github.com/GoncharovFyodor/hezzltest/internal/db/clickhouse"
	"github.com/GoncharovFyodor/hezzltest/internal/db/postgres"
	"github.com/GoncharovFyodor/hezzltest/internal/db/redis"
	"github.com/GoncharovFyodor/hezzltest/internal/handler"
	nats2 "github.com/GoncharovFyodor/hezzltest/internal/nats"
	"github.com/GoncharovFyodor/hezzltest/internal/queue/nats"
	"github.com/GoncharovFyodor/hezzltest/internal/repository"
	"github.com/GoncharovFyodor/hezzltest/internal/services"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func Run(log *log.Logger, cfg *config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := postgres.New(ctx, cfg)
	rdb := redis.New(ctx, cfg)
	nc, js, stream := nats.New(ctx, cfg)
	producer := nats2.NewProducer(js)

	clickConn := clickhouse.New(ctx, cfg)
	repos := repository.NewRepository(db, rdb, clickConn)
	service := services.NewService(cfg, log, repos, producer)
	srv := handler.NewServer(log, cfg, service, nc)
	consumer := nats2.NewConsumer(log, cfg, stream, repos)

	go func() {
		if err := srv.Run(srv.InitRoutes()); err != nil {
			log.Fatalf("Error while start server: %v", err)
		}
	}()

	go func() {
		if err := consumer.Run(ctx); err != nil {
			log.Fatalf("error while run consumer: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(ctx); err != nil {
		log.Infof("Error occured on server shutting down: %v", err)
	}

	db.Close()
	if err := rdb.Close(); err != nil {
		log.Infof("error while close redis conn: %v", err)
	}

	if err := nc.Drain(); err != nil {
		log.Infof("error while close redis conn: %v", err)
	}

	nc.Close()
}

package clickhouse

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/GoncharovFyodor/hezzltest/internal/config"
	"log"
	"time"
)

func New(ctx context.Context, cfg *config.Config) driver.Conn {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.Clickhouse.Host},
		Auth: clickhouse.Auth{
			Database: cfg.Clickhouse.Name,
			Username: cfg.Clickhouse.User,
			Password: cfg.Clickhouse.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 10 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})

	if err != nil {
		log.Fatalf("ошибка подключения к clickhouse: %v", err)
	}

	if err = conn.Ping(ctx); err != nil {
		log.Fatalf("ошибка ping clickhouse: %v", err)
	}

	return conn
}

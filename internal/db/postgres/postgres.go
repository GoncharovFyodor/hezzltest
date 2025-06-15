package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/GoncharovFyodor/hezzltest/internal/config"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func New(ctx context.Context, cfg *config.Config) *sql.DB {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка открытия соединения с БД: ", err)
	}

	if err = db.PingContext(ctx); err != nil {
		log.Fatal("Ошибка пингования БД: ", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db
}

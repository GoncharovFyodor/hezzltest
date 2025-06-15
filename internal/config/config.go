package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Env        string     `yaml:"env" env:"ENV"`
	Server     Server     `yaml:"server"`
	DB         DB         `yaml:"db"`
	Redis      Redis      `yaml:"redis"`
	Nats       Nats       `yaml:"nats"`
	Clickhouse Clickhouse `yaml:"clickhouse"`
}

type Server struct {
	Host string `yaml:"host" env:"HOST"`
	Port string `yaml:"port" env:"PORT"`
}

type DB struct {
	User     string `yaml:"user" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	Name     string `yaml:"name" env:"DB_NAME"`
	SSLMode  string `yaml:"ssl_mode" env:"DB_SSLMODE"`
}

type Redis struct {
	Host     string `yaml:"host" env:"REDIS_HOST"`
	Password string `yaml:"Password" env:"REDIS_PASSWORD"`
}

type Nats struct {
	Host         string `yaml:"host" env:"NATS_HOST"`
	WorkersCount int    `yaml:"workers_count" env:"NATS_WORKERS_COUNT"`
	BatchSize    int    `yaml:"batch_size" env:"NATS_BATCH_SIZE"`
}

type Clickhouse struct {
	Host     string `yaml:"host" env:"CLICKHOUSE_HOST"`
	Name     string `yaml:"name" env:"CLICKHOUSE_DBNAME"`
	User     string `yaml:"user" env:"CLICKHOUSE_USER"`
	Password string `yaml:"password" env:"CLICKHOUSE_PASSWORD"`
}

func Load() *Config {
	var cfg Config

	err := cleanenv.ReadConfig("config.yml", &cfg)

	if err != nil {
		log.Fatalf("error while read config: %v", err)
	}

	return &cfg
}

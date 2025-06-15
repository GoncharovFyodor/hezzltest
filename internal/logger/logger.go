package logger

import (
	"github.com/GoncharovFyodor/hezzltest/internal/config"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func NewLogger(cfg *config.Config) *log.Logger {
	l := log.New()

	l.SetFormatter(&log.JSONFormatter{})

	switch cfg.Env {
	case envLocal:
		l.SetLevel(log.DebugLevel)
		l.SetOutput(os.Stdout)
		l.SetFormatter(&log.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	case envDev:
		l.SetLevel(log.InfoLevel)
		l.SetOutput(os.Stdout)
	case envProd:
		l.SetLevel(log.WarnLevel)
	}

	return l
}

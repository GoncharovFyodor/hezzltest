package main

import (
	"github.com/GoncharovFyodor/hezzltest/internal/app"
	"github.com/GoncharovFyodor/hezzltest/internal/config"
	"github.com/GoncharovFyodor/hezzltest/internal/logger"
)

func main() {
	cfg := config.Load()

	log := logger.NewLogger(cfg)

	app.Run(log, cfg)
}

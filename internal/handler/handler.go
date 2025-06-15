package handler

import (
	"context"
	"github.com/GoncharovFyodor/hezzltest/internal/config"
	"github.com/GoncharovFyodor/hezzltest/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	log        *log.Logger
	cfg        *config.Config
	httpServer *http.Server
	services   *services.Service
	nc         *nats.Conn
}

func NewServer(log *log.Logger, cfg *config.Config, services *services.Service, nc *nats.Conn) *Server {
	return &Server{log: log, cfg: cfg, services: services, nc: nc}
}

func (s *Server) Run(handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         s.cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) InitRoutes() *gin.Engine {
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:    []string{"Content-Type", "Authorization"},
	}))

	api := r.Group("/api/v1")

	good := api.Group("/good")
	{
		good.POST("/create", s.CreateGood)
		good.PATCH("/update", s.UpdateGood)
		good.DELETE("/remove", s.DeleteGood)
		good.PATCH("/reprioritize", s.ReprioritizeGood)
	}

	api.GET("/goods/list", s.GetGoods)

	return r
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

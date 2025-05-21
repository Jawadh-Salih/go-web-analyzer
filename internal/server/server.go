package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Jawadh-Salih/go-web-analyzer/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Server struct {
	port          string
	logger        *slog.Logger
	svr           *http.Server
	router        *gin.Engine
	withTemplates bool
}

func New(port string, logger *slog.Logger, withTemplates bool) *Server {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.TimeoutMiddleware(3 * time.Second))

	s := &Server{
		port:   port,
		router: r,
		logger: logger,
	}

	s.setupMiddleware()
	s.registerRoutes(withTemplates)

	svr := &http.Server{
		Addr:    port,
		Handler: r,
	}

	s.svr = svr
	return s
}

func (s *Server) Start() error {
	return s.router.Run(s.port)
}

func (s *Server) Stop(ctx context.Context) error {
	// Implement graceful shutdown logic if needed
	s.logger.Info("Server is stopping...")
	return s.svr.Shutdown(ctx)
}

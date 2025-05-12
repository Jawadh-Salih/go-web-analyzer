package server

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port   string
	router *gin.Engine
	logger *slog.Logger
}

func New(port string, logger *slog.Logger) *Server {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())

	s := &Server{
		port:   port,
		router: r,
		logger: logger,
	}

	s.setupMiddleware()
	s.registerRoutes()

	return s
}

func (s *Server) Start() error {
	return s.router.Run(s.port)
}

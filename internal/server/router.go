package server

import (
	"log/slog"
	"net/http"

	"github.com/Jawadh-Salih/go-web-analyzer/internal/analyzer"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *Server) setupMiddleware() {
	s.router.Use(func(c *gin.Context) {
		s.logger.Info(
			"request received",
			"method",
			c.Request.Method,
			"path",
			c.Request.URL.Path)
		c.Next()
	})
}

func (s *Server) registerRoutes(withTemplates bool) {
	if withTemplates {
		s.router.LoadHTMLGlob("web/*.html")
	}

	s.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	s.router.POST("/analyze", s.analyzeHandler)

	// Observability
	s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func (s *Server) analyzeHandler(c *gin.Context) {
	var req analyzer.AnalyzerRequest
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil || req.Url == "" {
		s.logger.Error("Invalid request", slog.Any("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid request",
		})
		return
	}

	result, err := analyzer.Analyze(req)
	if err != nil {
		s.logger.Error("Internal Server Error", slog.Any("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

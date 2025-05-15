package server

import (
	"log/slog"
	"net/http"

	"github.com/Jawadh-Salih/go-web-analyzer/internal/analyzer"
	"github.com/Jawadh-Salih/go-web-analyzer/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/xid"
)

func (s *Server) setupMiddleware() {
	s.router.Use(func(c *gin.Context) {
		s.logger.Info(
			"request received",
			"method",
			c.Request.Method,
			"path",
			c.Request.URL.Path)

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = xid.New().String()
		}

		c.Set("request_id", requestID)
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

	// Should we pass the same logger to the Analyze function?
	log := s.logger.With(slog.String("request_id", getRequestID(c)))
	ctx := logger.SetLogger(c.Request.Context(), log)
	result, err := analyzer.Analyze(ctx, req)
	if err != nil {
		// cast the error and see if it's an HttpApiError
		// if not 500, if return the relevant code
		s.logger.Error("Internal Server Error", slog.Any("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func getRequestID(c *gin.Context) string {
	if val, ok := c.Get("request_id"); ok {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}

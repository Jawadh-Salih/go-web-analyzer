package server

import (
	"net/http"

	"github.com/Jawadh-Salih/go-web-analyzer/internal/analyzer"
	"github.com/gin-gonic/gin"
)

func (s *Server) setupMiddleware() {
	s.router.Use(func(c *gin.Context) {
		s.logger.Info("request received", "method", c.Request.Method, "path", c.Request.URL.Path)
		c.Next()
	})
}

func (s *Server) registerRoutes() {
	s.router.LoadHTMLGlob("web/*.html")

	s.router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	s.router.POST("/analyze", s.analyzeHandler)

	// set the response to the html file.
}

func (s *Server) analyzeHandler(c *gin.Context) {
	var req analyzer.AnalyzerRequest
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil || req.Url == "" {
		s.logger.Error("Invalid request %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid request",
		})
		return
	}

	result, err := analyzer.Analyze(c.Request.Context(), req)
	if err != nil {
		s.logger.Error("Internal Server Error %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

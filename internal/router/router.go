package router

import (
	"net/http"

	"github.com/Jawadh-Salih/go-web-analyzer/internal/analyzer"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Engine *gin.Engine
	Port   string
}

func Init(port string) Handler {
	r := gin.Default()
	route := Handler{
		Engine: r,
		Port:   port,
	}

	route.setupRoutes()

	return route
}

func (e Handler) Run() error {
	err := e.Engine.Run(e.Port)

	if err != nil {
		return err
	}

	return nil
}

func (e Handler) setupRoutes() {
	e.Engine.LoadHTMLGlob("web/*.html")

	e.Engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	e.Engine.POST("/analyze", analyzeHandler)

	// set the response to the html file.
}

func analyzeHandler(c *gin.Context) {
	var req analyzer.AnalyzerRequest
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil || req.Url == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid request",
		})
		return
	}

	result, err := analyzer.Analyze(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

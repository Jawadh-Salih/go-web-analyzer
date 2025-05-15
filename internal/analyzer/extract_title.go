package analyzer

import (
	"log/slog"
	"sync"
	"time"

	"github.com/Jawadh-Salih/go-web-analyzer/internal/observability"
	"golang.org/x/net/html"
)

func ExtractTitle(logger *slog.Logger, root *html.Node, wg *sync.WaitGroup, resultChan chan AnalyzerResponse) {
	start := time.Now()
	status := "Success"
	functionName := "ExtractTitle"
	defer wg.Done()

	title := getTitle(root)
	resultChan <- AnalyzerResponse{PageTitle: title}

	duration := time.Since(start).Nanoseconds()
	logger.Info("Function Executed",
		slog.String("function", functionName),
		slog.Int64("duration", duration),
	)

	observability.
		DurationMetrics.
		WithLabelValues(functionName, status).
		Observe(float64(duration))
}

func getTitle(node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "title" {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				return child.Data
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if title := getTitle(child); title != "" {
			return title
		}

	}

	return ""
}

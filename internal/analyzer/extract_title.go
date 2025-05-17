package analyzer

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Jawadh-Salih/go-web-analyzer/internal/logger"
	"github.com/Jawadh-Salih/go-web-analyzer/internal/observability"
	"golang.org/x/net/html"
)

func ExtractTitle(ctx context.Context, root *html.Node, wg *sync.WaitGroup, resultChan chan AnalyzerResponse) {
	logger := logger.FromContext(ctx)
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
	nodes := make([]html.Node, 0)
	getMatchingNodes(node, &nodes, "title")

	if len(nodes) > 0 {
		for child := nodes[0].FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				return child.Data
			}
		}
	}

	return ""
}

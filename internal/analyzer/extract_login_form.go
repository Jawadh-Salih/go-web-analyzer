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

func ExtractLoginForm(ctx context.Context, root *html.Node, wg *sync.WaitGroup, resultChan chan AnalyzerResponse) {
	logger := logger.FromContext(ctx)
	start := time.Now()
	status := "Success"
	functionName := "ExtractLoginForm"
	defer wg.Done()

	resultChan <- AnalyzerResponse{HasLoginForm: hasLoginForm(root)}

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

func hasLoginForm(node *html.Node) bool {
	// if the node data is input check if the input type is password and submit
	nodes := make([]html.Node, 0)
	getMatchingNodes(node, &nodes, "input", "button")
	var hasPasswordField, hasSubmitButton bool
	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == "type" && attr.Val == "password" {
				hasPasswordField = true
			}

			if attr.Key == "type" && (attr.Val == "submit" || attr.Val == "button") {
				hasSubmitButton = true
			}
		}
	}

	return hasPasswordField && hasSubmitButton
}

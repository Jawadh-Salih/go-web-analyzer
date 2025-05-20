package analyzer

import (
	"log/slog"
	"sync"
	"time"

	"github.com/Jawadh-Salih/go-web-analyzer/internal/observability"
	"golang.org/x/net/html"
)

func ExtractHeadings(root *html.Node, wg *sync.WaitGroup, resultChan chan AnalyzerResponse) {
	start := time.Now()
	status := "Success"
	functionName := "ExtractHeadings"
	defer wg.Done()

	headingCounts := headingsMap(root)

	// Send the result to the channel
	resultChan <- AnalyzerResponse{Headings: headingCounts}
	duration := time.Since(start).Nanoseconds()
	analyzerLogger.Info("Function Executed",
		slog.String("function", functionName),
		slog.Int64("duration", duration),
	)

	observability.
		DurationMetrics.
		WithLabelValues(functionName, status).
		Observe(float64(duration))
}

func headingsMap(node *html.Node) map[string]int {
	// Check if the node is an element node and is a heading tag (h1 to h6)
	headingCounts := make(map[string]int)
	nodes := make([]html.Node, 0)
	getMatchingNodes(node, &nodes, "h1", "h2", "h3", "h4", "h5", "h6")

	for _, value := range nodes {
		headingCounts[value.Data]++
	}

	return headingCounts
}

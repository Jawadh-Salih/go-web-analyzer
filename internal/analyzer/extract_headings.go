package analyzer

import (
	"log/slog"
	"sync"
	"time"

	"github.com/Jawadh-Salih/go-web-analyzer/internal/observability"
	"golang.org/x/net/html"
)

func ExtractHeadings(logger *slog.Logger, root *html.Node, wg *sync.WaitGroup, resultChan chan AnalyzerResponse) {
	start := time.Now()
	status := "Success"
	functionName := "ExtractHeadings"
	defer wg.Done()

	headingCounts := make(map[string]int)

	// Traverse the HTML node tree and count headings
	headingsMap(root, headingCounts)

	// Send the result to the channel
	resultChan <- AnalyzerResponse{Headings: headingCounts}
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

func headingsMap(node *html.Node, headingCounts map[string]int) {
	// Check if the node is an element node and is a heading tag (h1 to h6)
	if node.Type == html.ElementNode && node.Data[0] == 'h' && len(node.Data) == 2 {
		level := node.Data[1] - '0' // Convert 'h1' to level 1, 'h2' to level 2, etc.
		if level >= 1 && level <= 6 {
			headingCounts[node.Data]++
		}
	}

	// Recursively traverse the child nodes
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		headingsMap(child, headingCounts)
	}
}

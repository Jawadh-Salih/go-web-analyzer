package analyzer

import (
	"log/slog"
	"sync"
	"time"

	"github.com/Jawadh-Salih/go-web-analyzer/internal/observability"
	"golang.org/x/net/html"
)

func ExtractLoginForm(logger *slog.Logger, root *html.Node, wg *sync.WaitGroup, resultChan chan AnalyzerResponse) {
	start := time.Now()
	status := "Success"
	functionName := "ExtractLoginForm"
	defer wg.Done()

	var pwdField, submitButton bool
	hasLoginForm(root, &pwdField, &submitButton)
	// logger.Info("Login form okay", slog.Bool("Value", pwdField && submitButton))
	resultChan <- AnalyzerResponse{HasLoginForm: pwdField && submitButton}

	duration := time.Since(start).Nanoseconds()
	logger.Info("Function Executed",
		slog.String("function", functionName),
		slog.Int64("duration", duration),
		slog.Bool("Value", pwdField && submitButton),
	)

	observability.
		DurationMetrics.
		WithLabelValues(functionName, status).
		Observe(float64(duration))
}

func hasLoginForm(node *html.Node, hasPasswordField, hasSubmitButton *bool) {
	// if the node data is input check if the input type is password and submit

	// if we can find these 2 info then
	if node.Type == html.ElementNode && (node.Data == "input" || node.Data == "button") {
		for _, attr := range node.Attr {
			if attr.Key == "type" && attr.Val == "password" {
				*hasPasswordField = true
			}

			if attr.Key == "type" && (attr.Val == "submit" || attr.Val == "button") {
				*hasSubmitButton = true
			}
		}
	}

	// recursively check for child nodes
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		hasLoginForm(child, hasPasswordField, hasSubmitButton)
	}
}

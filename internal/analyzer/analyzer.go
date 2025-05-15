package analyzer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Jawadh-Salih/go-web-analyzer/internal/logger"
	"github.com/Jawadh-Salih/go-web-analyzer/internal/observability"
	"golang.org/x/net/html"
)

// This will analyze the request url.
type AnalyzerRequest struct {
	Url string
}
type AnalyzerResponse struct {
	HtmlVersion  string         // HTML version
	PageTitle    string         // Page title
	Headings     map[string]int // Headings count
	Links        []Link         // Links
	HasLoginForm bool           // true if the page has a login form
	Errors       []string       // Errors encountered during analysis
	err          error
}

type Link struct {
	LinkType   string // internal or external
	LinkUrl    string // url
	Accessible bool   // true if the link is inaccessible

}

func Analyze(ctx context.Context, request AnalyzerRequest) (*AnalyzerResponse, error) {
	logger := logger.FromContext(ctx)
	result := AnalyzerResponse{
		Errors: make([]string, 0),
	}

	pageUrl, err := validateURL(request.Url)
	if err != nil {
		logger.Error("Invalid URL", slog.Any("Error", err))
		return nil, err
	}

	// TODO Timeouts to be configured
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(request.Url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Error("Failed to reach the URL", slog.String("srl", request.Url), slog.Int("status", resp.StatusCode))
		return nil, errors.New("Failed to reach URL")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response body", slog.Any("error", err))
		return nil, err

	}

	// check for html content type
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		err := fmt.Errorf("Invalid response: %s", resp.Header.Get("Content-Type"))
		logger.Error(err.Error(), slog.String("content-type", resp.Header.Get("Content-Type")))
		return nil, err
	}

	rootNode, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		logger.Error("failed to parse HTML", slog.Any("error", err))
		return nil, err
	}

	// following can be done in parallel
	// number 5 can be taken into configs based on the different anaylysis we need from the analyzer
	resultChan := make(chan AnalyzerResponse, 5)
	var wg sync.WaitGroup

	// -   What HTML version has the document?
	wg.Add(5)
	go func() {
		startTime := time.Now()
		defer wg.Done()
		buffer := int(math.Min(float64(len(body)), 2048))
		htmlSnippet := string(body[:buffer])
		status := "Success"
		if htmlSnippet == "" {
			resultChan <- AnalyzerResponse{err: errors.New("empty HTML snippet")}
			status = "Fail"
			return
		}

		htmlV := detectHTMLVersion(htmlSnippet)
		resultChan <- AnalyzerResponse{HtmlVersion: htmlV}

		duration := time.Since(startTime).Seconds()
		functionName := "HtmlVersion Check"
		logger.Info("Function Executed",
			slog.String("function", functionName),
			slog.String("status", status),
			slog.Float64("duration", duration),
		)

		observability.
			DurationMetrics.
			WithLabelValues(functionName, status).
			Observe(duration)
	}()

	// -   What is the page title?
	go ExtractTitle(logger, rootNode, &wg, resultChan)

	// -   How many headings of what level are in the document?
	go ExtractHeadings(logger, rootNode, &wg, resultChan)

	// -   How many internal and external links are in the document? Are there any inaccessible links and how many?
	go ExtrackLinks(logger, rootNode, pageUrl, &wg, resultChan)

	// -   Does the page contain a login form?
	go ExtractLoginForm(logger, rootNode, &wg, resultChan)

	// Close the result channel after all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for res := range resultChan {
		if res.err != nil {
			result.Errors = append(result.Errors, res.err.Error())
			continue
		} else {
			if res.HtmlVersion != "" {
				result.HtmlVersion = res.HtmlVersion
			}
			if res.PageTitle != "" {
				result.PageTitle = res.PageTitle
			}
			if len(res.Headings) > 0 {
				result.Headings = res.Headings
			}
			if len(res.Links) > 0 {
				result.Links = append(result.Links, res.Links...)
			}
		}
	}

	return &result, nil
}

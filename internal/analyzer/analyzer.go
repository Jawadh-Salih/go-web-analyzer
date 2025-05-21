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
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Jawadh-Salih/go-web-analyzer/internal/logger"
	"github.com/Jawadh-Salih/go-web-analyzer/internal/observability"
	"golang.org/x/net/html"
)

// This will analyze the request url.
type AnalyzerRequest struct {
	Url string `json:"url" binding:"required,url"`
}
type AnalyzerResponse struct {
	HtmlVersion  string         // HTML version
	PageTitle    string         // Page title
	Headings     map[string]int // Headings count
	LinkSummary  *LinkSummary   // Links
	HasLoginForm bool           // true if the page has a login form
	Errors       []string       // Errors encountered during analysis
	err          error
}

type LinkSummary struct {
	Links             []Link
	InternalLinks     int
	ExternalLinks     int
	AccessibleLinks   int
	InaccessibleLinks int
}

type Link struct {
	LinkType   string // internal or external
	LinkUrl    string // url
	Accessible bool   // true if the link is inaccessible

}

var analyzerLogger *slog.Logger

func Analyze(ctx context.Context, request AnalyzerRequest) (*AnalyzerResponse, error) {
	analyzerLogger = logger.FromContext(ctx)
	result := AnalyzerResponse{
		Errors: make([]string, 0),
	}

	pageUrl, err := url.Parse(request.Url)
	if err != nil {
		return nil, fmt.Errorf("invalid URL syntax: %w", err)
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
		analyzerLogger.Error("Failed to reach the URL", slog.String("srl", request.Url), slog.Int("status", resp.StatusCode))
		return nil, errors.New("Failed to reach URL")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		analyzerLogger.Error("Failed to read response body", slog.Any("error", err))
		return nil, err

	}

	// check for html content type
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		err := fmt.Errorf("Invalid response: %s", resp.Header.Get("Content-Type"))
		analyzerLogger.Error(err.Error(), slog.String("content-type", resp.Header.Get("Content-Type")))
		return nil, err
	}

	rootNode, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		analyzerLogger.Error("failed to parse HTML", slog.Any("error", err))
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

		duration := time.Since(startTime).Nanoseconds()
		functionName := "HtmlVersion Check"
		analyzerLogger.Info("Function Executed",
			slog.String("function", functionName),
			slog.String("status", status),
			slog.Int64("duration", duration),
		)

		observability.
			DurationMetrics.
			WithLabelValues(functionName, status).
			Observe(float64(duration))
	}()

	// -   What is the page title?
	go ExtractTitle(rootNode, &wg, resultChan)

	// -   How many headings of what level are in the document?
	go ExtractHeadings(rootNode, &wg, resultChan)

	// -   How many internal and external links are in the document? Are there any inaccessible links and how many?
	go ExtrackLinks(rootNode, pageUrl, &wg, resultChan)

	// -   Does the page contain a login form?
	go ExtractLoginForm(rootNode, &wg, resultChan)

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
			if res.LinkSummary != nil && len(res.LinkSummary.Links) > 0 {
				result.LinkSummary = res.LinkSummary
			}
			if res.HasLoginForm {
				result.HasLoginForm = res.HasLoginForm
			}

		}
	}

	return &result, nil
}

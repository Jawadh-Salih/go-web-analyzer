package analyzer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

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
	result := AnalyzerResponse{
		Errors: make([]string, 0),
	}

	pageUrl, err := validateURL(request.Url)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %s", request.Url)
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
		return nil, fmt.Errorf("failed to reach the URL: %s with status code %d", request.Url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// check for html content type
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return nil, fmt.Errorf("Invalid response: %s", resp.Header.Get("Content-Type"))
	}

	rootNode, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// following can be done in parallel
	// number 5 can be taken into configs based on the different anaylysis we need from the analyzer
	resultChan := make(chan AnalyzerResponse, 5)
	var wg sync.WaitGroup

	// -   What HTML version has the document?
	wg.Add(1)
	go func() {
		defer wg.Done()
		htmlSnippet := string(body[:2048])
		if htmlSnippet == "" {
			resultChan <- AnalyzerResponse{err: errors.New("empty HTML snippet")}
			return
		}

		htmlV := detectHTMLVersion(htmlSnippet)
		resultChan <- AnalyzerResponse{HtmlVersion: htmlV}
	}()

	// -   What is the page title?
	wg.Add(1)
	go func() {
		defer wg.Done()
		title := getTitle(rootNode)
		resultChan <- AnalyzerResponse{PageTitle: title}
	}()

	// -   How many headings of what level are in the document?
	wg.Add(1)
	go func() {
		defer wg.Done()
		headings := make(map[string]int)
		headingsMap(rootNode, headings)

		resultChan <- AnalyzerResponse{Headings: headings}
	}()

	// -   How many internal and external links are in the document? Are there any inaccessible links and how many?
	wg.Add(1)
	go func() {
		defer wg.Done()
		links := make([]Link, 0)
		getLinks(rootNode, pageUrl, &links)

		if len(links) == 0 {
			resultChan <- AnalyzerResponse{err: errors.New("No links available")}
			return
		}

		linkChan := make(chan *Link, len(links))
		var linkWg sync.WaitGroup

		workers := int(math.Sqrt(float64(len(links))) * 3)
		for i := 0; i < workers; i++ {
			// now pass the link channel to
			linkWg.Add(1)
			// consider using context
			go checkLink(linkChan, &linkWg)

		}

		// feed the links to the linkChan
		for i := range links {
			linkChan <- &links[i]
		}

		close(linkChan)

		linkWg.Wait()
		resultChan <- AnalyzerResponse{Links: links}
	}()

	// -   Does the page contain a login form?
	wg.Add(1)
	go func() {
		defer wg.Done()
		var pwdField, submitButton bool
		loginForm := hasLoginForm(rootNode, &pwdField, &submitButton)
		resultChan <- AnalyzerResponse{HasLoginForm: loginForm}
	}()

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

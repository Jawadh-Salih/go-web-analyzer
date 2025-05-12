package analyzer

import (
	"fmt"
	"io"
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
}

type Link struct {
	LinkType   string // internal or external
	LinkUrl    string // url
	Accessible bool   // true if the link is inaccessible

}

func Analyze(request AnalyzerRequest) (*AnalyzerResponse, error) {
	// check if there is internet connection.
	// Here you would implement the logic to check for internet connection.
	// This is a placeholder implementation.

	result := AnalyzerResponse{}

	pageUrl, err := validateURL(request.Url)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %s", request.Url)
	}

	// all of these can be done using go concurrency parallely
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	req, _ := http.NewRequest("GET", request.Url, nil)

	resp, err := client.Do(req)
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

	htmlStr := string(body)

	rootNode, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// printTags(rootNode, 0)
	// following should be done in parallel
	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, 1)

	// -   What HTML version has the document?
	wg.Add(1)
	go func() {
		defer wg.Done()
		htmlV, err := getHtmlVersion(htmlStr)
		if err != nil {
			// put the error in a channel
			errCh <- err
		}
		mu.Lock()
		result.HtmlVersion = htmlV
		mu.Unlock()
	}()

	// -   What is the page title?
	wg.Add(1)
	go func() {
		defer wg.Done()
		title := getTitle(rootNode)

		mu.Lock()
		result.PageTitle = title
		mu.Unlock()
	}()

	// -   How many headings of what level are in the document?
	wg.Add(1)
	go func() {
		defer wg.Done()
		headings := make(map[string]int)
		headingsMap(rootNode, headings)

		mu.Lock()
		result.Headings = headings
		mu.Unlock()
	}()

	// -   How many internal and external links are in the document? Are there any inaccessible links and how many?
	wg.Add(1)
	go func() {
		defer wg.Done()
		links := make([]Link, 0)
		getLinks(rootNode, pageUrl, &links)

		mu.Lock()
		result.Links = links
		mu.Unlock()
	}()

	// -   Does the page contain a login form?
	wg.Add(1)
	go func() {
		defer wg.Done()

		var pwdField, submitButton bool
		loginForm := hasLoginForm(rootNode, &pwdField, &submitButton)

		fmt.Println("Login form: ", pwdField, submitButton)
		mu.Lock()
		result.HasLoginForm = loginForm
		mu.Unlock()
	}()

	wg.Wait()
	close(errCh)
	// Here you would implement the logic to analyze the URL.
	// check if you can access the url.
	// check if the url is reachable.

	return &result, nil
}

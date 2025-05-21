package analyzer

import (
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/Jawadh-Salih/go-web-analyzer/internal/observability"
	"golang.org/x/net/html"
)

func ExtrackLinks(root *html.Node, pageUrl *url.URL, wg *sync.WaitGroup, resultChan chan AnalyzerResponse) {
	start := time.Now()
	status := "Success"
	functionName := "ExtractLinks"
	defer wg.Done()

	links := make([]Link, 0)
	nodes := make([]html.Node, 0)
	getMatchingNodes(root, &nodes, "a")

	// can execute this parallely
	nodeChan := make(chan *html.Node, len(nodes))
	var linkWg sync.WaitGroup
	workers := int(math.Sqrt(float64(len(nodes))) * 3)

	for i := 0; i < workers; i++ {
		linkWg.Add(1)
		go setupLinks(nodeChan, pageUrl, &linkWg, &links)
	}

	// feed the nodes to the node channel
	for i := range nodes {
		nodeChan <- &nodes[i]
	}

	close(nodeChan)

	linkWg.Wait()
	var accessibles, internals int
	for _, link := range links {
		if link.LinkType == "internal" {
			internals++
		}

		if link.Accessible {
			accessibles++
		}
	}

	resultChan <- AnalyzerResponse{
		LinkSummary: &LinkSummaryResponse{
			Links:             links,
			InternalLinks:     internals,
			ExternalLinks:     len(links) - internals,
			AccessibleLinks:   accessibles,
			InaccessibleLinks: len(links) - accessibles,
		}}

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

func getLinkType(linkURL, baseURL *url.URL) string {
	// Check if the domain of the link matches the base URL
	if linkURL.Host == baseURL.Host {
		return "internal"
	}

	return "external"
}

func setupLinks(nodes <-chan *html.Node, baseUrl *url.URL, wg *sync.WaitGroup, links *[]Link) {
	defer wg.Done()
	for node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				// we have a link now.
				linkUrl, err := url.Parse(attr.Val)
				if err != nil {
					// validate the URL and ignore if invalid
					continue
				}

				// If the link is relative, resolve it to an absolute URL
				if !linkUrl.IsAbs() {
					linkUrl = baseUrl.ResolveReference(linkUrl)
				}

				// checking the accessibility of the link can be run in parallel
				// so it will be implemented once we have all the links extracted
				client := &http.Client{Timeout: 3 * time.Second}

				var accessible bool
				resp, _ := client.Head(linkUrl.String())
				if resp != nil && resp.StatusCode == http.StatusOK {
					accessible = true
					resp.Body.Close()
				}

				*links = append(*links, Link{
					LinkType:   getLinkType(linkUrl, baseUrl),
					LinkUrl:    linkUrl.String(),
					Accessible: accessible,
				})

			}
		}
	}
}

package analyzer

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/html"
)

func ExtrackLinks(root *html.Node, pageUrl *url.URL, wg *sync.WaitGroup, resultChan chan AnalyzerResponse) {
	defer wg.Done()
	links := make([]Link, 0)
	getLinks(root, pageUrl, &links)

	if len(links) == 0 {
		resultChan <- AnalyzerResponse{err: errors.New("No links available")}
		return
	}

	linkChan := make(chan *Link, len(links))
	var linkWg sync.WaitGroup

	// Following is a equation I recognize from what I saw on the internet
	workers := int(math.Sqrt(float64(len(links))) * 3)
	for i := 0; i < workers; i++ {
		linkWg.Add(1)
		go checkLink(linkChan, &linkWg)
	}

	// feed the links to the linkChan
	for i := range links {
		linkChan <- &links[i]
	}

	close(linkChan)

	linkWg.Wait()
	resultChan <- AnalyzerResponse{Links: links}
}

func getLinks(node *html.Node, baseUrl *url.URL, links *[]Link) {
	// should check for href attribute
	if node.Type == html.ElementNode && (node.Data == "a" || node.Data == "link") {
		// fmt.Println(node.Attr)
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				// we have a link now.
				linkUrl, err := url.Parse(attr.Val)
				if err != nil {
					fmt.Println("Error on parse ", err)
					continue
				}

				// If the link is relative, resolve it to an absolute URL
				if !linkUrl.IsAbs() {
					linkUrl = baseUrl.ResolveReference(linkUrl)
				}

				// checking the accessibility of the link can be run in parallel
				// so it will be implemented once we have all the links extracted

				*links = append(*links, Link{
					LinkType: getLinkType(linkUrl, baseUrl),
					LinkUrl:  linkUrl.String(),
				})
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		getLinks(child, baseUrl, links)
	}
}

func getLinkType(linkURL, baseURL *url.URL) string {
	// Check if the domain of the link matches the base URL
	if linkURL.Host == baseURL.Host {
		return "internal"
	}

	return "external"
}

func checkLink(links <-chan *Link, wg *sync.WaitGroup) {
	defer wg.Done()

	for link := range links {
		client := &http.Client{Timeout: 3 * time.Second}

		resp, err := client.Head(link.LinkUrl)
		if err != nil {
			// log the error
			link.Accessible = false
			// logger.Error("Error fetching URL", "worker_id", id, "url", link.LinkUrl, "error", err)
		} else {
			link.Accessible = resp.StatusCode == 200
			defer resp.Body.Close()
		}
	}
}

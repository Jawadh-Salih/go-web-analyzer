package analyzer

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

func validateURL(raw string) (*url.URL, error) {
	if raw == "" {
		return nil, fmt.Errorf("empty URL")
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid URL syntax: %w", err)
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("invalid URL: missing scheme or host")
	}

	return parsed, nil
}

func detectHTMLVersion(htmlStr string) string {
	lower := strings.ToLower(htmlStr)
	switch {
	case strings.Contains(lower, "<!doctype html>"):
		return "HTML5"
	case strings.Contains(lower, "html 4.01 transitional"):
		return "HTML 4.01 Transitional"
	case strings.Contains(lower, "html 4.01 strict"):
		return "HTML 4.01 Strict"
	case strings.Contains(lower, "xhtml 1.0 transitional"):
		return "XHTML 1.0 Transitional"
	case strings.Contains(lower, "xhtml 1.0 strict"):
		return "XHTML 1.0 Strict"
	case strings.Contains(lower, "xhtml 1.1"):
		return "XHTML 1.1"
	default:
		return "Unknown"
	}
}

func getTitle(node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "title" {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				return child.Data
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if title := getTitle(child); title != "" {
			return title
		}

	}

	return ""
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
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		headingsMap(c, headingCounts)
	}
}

func getLinks(node *html.Node, baseUrl *url.URL, links *[]Link) {
	// should check for href attribute
	// urlChan := make(chan string, 10)
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

				// this can be run in parallel
				// put the urls in a buffered channel and consume then with http.Get()
				// check if the link is accessible

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

func hasLoginForm(node *html.Node, hasPasswordField, hasSubmitButton *bool) bool {
	// if the node data is input. check if the input type is password
	// if the node data is input check if the input type is submit

	// if we can find these 2 info then
	if node.Type == html.ElementNode && node.Data == "input" {
		for _, attr := range node.Attr {
			if attr.Key == "type" && attr.Val == "password" {

				*hasPasswordField = true
			}

			if attr.Key == "type" && (attr.Val == "submit" || attr.Val == "button") {
				*hasSubmitButton = true
			}
		}

		if *hasPasswordField && *hasSubmitButton {
			return true
		}
	}

	// recursively check for child nodes
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if hasLoginForm(child, hasPasswordField, hasSubmitButton) {
			return true
		}
	}

	return *hasPasswordField && *hasSubmitButton
}

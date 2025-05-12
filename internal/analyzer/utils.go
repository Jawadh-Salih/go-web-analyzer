package analyzer

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

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

func getHtmlVersion(resp string) (string, error) {
	// Read first 2048 bytes (enough to catch the DOCTYPE) for Optimization
	buf := make([]byte, 2048)
	node, err := strings.NewReader(resp).Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}
	if node == 0 {
		return "", fmt.Errorf("empty response")
	}

	htmlSnippet := string(buf[:node])
	return detectHTMLVersion(htmlSnippet), nil
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

func printTags(node *html.Node, depth int) {
	if node.Type == html.ElementNode {
		fmt.Printf("%s<%s>\node", strings.Repeat("  ", depth), node.Data)
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		printTags(c, depth+1)
	}
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
	if node.Type == html.ElementNode && node.Data == "a" {
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

				resp, err := http.Get(linkUrl.Host)

				*links = append(*links, Link{
					LinkType:   getLinkType(linkUrl, baseUrl),
					LinkUrl:    linkUrl.Host,
					Accessible: (err == nil) && (resp.StatusCode >= 200 && resp.StatusCode < 300),
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
	if linkURL.Host == baseURL.Host || strings.HasPrefix(linkURL.Path, "/") {
		return "internal"
	}

	return "external"
}

func hasLoginForm(node *html.Node, hasPasswordField, hasSubmitButton *bool) bool {
	// if the node data is input. check if the input type is password
	// if the node data is input check if the input type is submit

	// if we can find these 2 info then
	if node.Type == html.ElementNode && node.Data == "input" {
		for _, attr := range node.Attr {
			if attr.Key == "type" && attr.Val == "password" {

				fmt.Println(attr.Key, " - ", attr.Val)
				*hasPasswordField = true
			}

			if attr.Key == "type" && (attr.Val == "submit" || attr.Val == "button") {
				fmt.Println(attr.Key, " - ", attr.Val)
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

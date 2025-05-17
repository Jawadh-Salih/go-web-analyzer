package analyzer

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func validateURL(raw string) (*url.URL, error) {
	// TODO a regex to validate the URL
	if raw == "" {
		return nil, fmt.Errorf("empty URL")
	}

	urlRegex := regexp.MustCompile(`^(https?:\/\/)?(www\.)?([a-zA-Z0-9_-]+(:[a-zA-Z0-9_-]+)?@)?((([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,5})|(\d{1,3}(\.\d{1,3}){3}))(:\d{1,5})?(\/.*)?$`)
	if !urlRegex.MatchString(raw) {
		return nil, fmt.Errorf("invalid URL: %s", raw)
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
	case strings.Contains(lower, "<!doctype html public \"-//w3c//dtd html 4.01 transitional//en\""):
		return "HTML 4.01 Transitional"
	case strings.Contains(lower, "<!doctype html public \"-//w3c//dtd html 4.01//en\""):
		return "HTML 4.01 Strict"
	case strings.Contains(lower, "<!doctype html public \"-//w3c//dtd xhtml 1.0 transitional//en\""):
		return "XHTML 1.0 Transitional"
	case strings.Contains(lower, "<!doctype html public \"-//w3c//dtd xhtml 1.0 strict//en\""):
		return "XHTML 1.0 Strict"
	default:
		return "Unknown"
	}
}

func getMatchingNodes(node *html.Node, nodes *[]html.Node, nodesData ...string) {
	// should check for href attribute

	// if we can find these 2 info then
	if node.Type == html.ElementNode {
		// for loop to filter only what is in the nodesData
		for _, data := range nodesData {
			if data == node.Data {
				*nodes = append(*nodes, *node)
			}
		}
	}

	// recursively check for child nodes
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		getMatchingNodes(child, nodes, nodesData...)
	}
}

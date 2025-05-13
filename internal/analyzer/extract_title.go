package analyzer

import (
	"sync"

	"golang.org/x/net/html"
)

func ExtractTitle(root *html.Node, wg *sync.WaitGroup, resultChan chan AnalyzerResponse) {
	defer wg.Done()

	title := getTitle(root)
	resultChan <- AnalyzerResponse{PageTitle: title}
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

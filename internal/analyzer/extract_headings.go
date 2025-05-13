package analyzer

import (
	"sync"

	"golang.org/x/net/html"
)

func ExtractHeadings(root *html.Node, wg *sync.WaitGroup, resultChan chan AnalyzerResponse) {
	defer wg.Done()

	headingCounts := make(map[string]int)

	// Traverse the HTML node tree and count headings
	headingsMap(root, headingCounts)

	// Send the result to the channel
	resultChan <- AnalyzerResponse{Headings: headingCounts}
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

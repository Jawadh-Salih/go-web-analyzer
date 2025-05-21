package analyzer

import (
	"strings"

	"golang.org/x/net/html"
)

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

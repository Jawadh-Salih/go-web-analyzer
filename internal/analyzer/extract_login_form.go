package analyzer

import (
	"sync"

	"golang.org/x/net/html"
)

func ExtractLoginForm(root *html.Node, wg *sync.WaitGroup, resultChan chan AnalyzerResponse) {
	defer wg.Done()

	var pwdField, submitButton bool
	loginForm := hasLoginForm(root, &pwdField, &submitButton)
	resultChan <- AnalyzerResponse{HasLoginForm: loginForm}
}

func hasLoginForm(node *html.Node, hasPasswordField, hasSubmitButton *bool) bool {
	// if the node data is input check if the input type is password and submit

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

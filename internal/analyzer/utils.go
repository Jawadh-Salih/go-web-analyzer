package analyzer

import (
	"fmt"
	"net/url"
	"strings"
)

func validateURL(raw string) (*url.URL, error) {
	// TODO a regex to validate the URL
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

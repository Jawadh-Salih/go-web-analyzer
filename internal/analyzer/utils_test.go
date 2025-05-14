package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlValidation(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"http://localhost:4200", true},
		{"http://example.com", true},
		{"https://example.com", true},
		{"https://example.com:8080", true},
		{"http://12.12.21.211:8080", true},
		{"http://12.12.21.211:8080/path?abc=123", true},
		{"ftp://example.com", false},
		{"invalid-url", false},
	}

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			result, err := validateURL(test.url)
			assert.Equal(t, test.expected, err == nil)
			if err == nil {
				assert.Equal(t, test.url, result.String())
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestHtmlVersions(t *testing.T) {
	// add tests for different content types
	htmls := []struct {
		content string
		version string
	}{
		{content: "<!DOCTYPE html><html><head><title>HTML5</title></head><body></body></html>", version: "HTML5"},
		{content: "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \"http://www.w3.org/TR/html4/strict.dtd\"><html><head><title>HTML 4.01</title></head><body></body></html>", version: "HTML 4.01 Strict"},
		{content: "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01 Transitional//EN\" \"http://www.w3.org/TR/html4/loose.dtd\"><html><head><title>HTML 4.01 Transitional</title></head><body></body></html>", version: "HTML 4.01 Transitional"},
		{content: "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\"><html><head><title>XHTML 1.0 Transitional</title></head><body></body></html>", version: "XHTML 1.0 Transitional"},
		{content: "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html><head><title>XHTML 1.0 Strict</title></head><body></body></html>", version: "XHTML 1.0 Strict"},
		{content: "<html></title></head><body></body></html>", version: "Unknown"},
	}

	for _, ct := range htmls {
		t.Run(ct.content, func(t *testing.T) {
			version := detectHTMLVersion(ct.content)
			assert.Equal(
				t,
				ct.version,
				version,
			)
		})
	}

}

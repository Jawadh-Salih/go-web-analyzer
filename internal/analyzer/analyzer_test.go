package analyzer

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyze_EmptyUrl(t *testing.T) {
	// Test with an empty URL
	request := AnalyzerRequest{Url: ""}
	ctx := context.Background()
	response, err := Analyze(ctx, request)

	assert.Error(t, err)
	assert.Equal(t, "invalid URL: ", err.Error())
	assert.Nil(t, response)
}

func TestAnalyze_InvalidUrl(t *testing.T) {
	// Test with an invalid URL
	request := AnalyzerRequest{Url: "invalid-url"}

	ctx := context.Background()
	response, err := Analyze(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, "invalid URL: invalid-url", err.Error())
	assert.Nil(t, response)
}

func TestAnalyze_UnreachableUrl(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable) // 503 - Service Unavailable
	}))

	defer testServer.Close()

	// Test with a URL that times out
	request := AnalyzerRequest{Url: testServer.URL} // Assuming this port is closed

	ctx := context.Background()
	response, err := Analyze(ctx, request)
	assert.Error(t, err)
	assert.Equal(
		t,
		fmt.Sprintf(`failed to reach the URL: %s with status code %d`, testServer.URL, http.StatusServiceUnavailable),
		err.Error(),
	)
	assert.Nil(t, response)
}

func TestAnalyze_ValidUrl_InvalidContent(t *testing.T) {
	// add tests for different content types
	contentTypes := []struct {
		contentType string
		content     []byte
		expectedErr string
	}{
		{"application/json", []byte(`{"message": "This is JSON content."}`), "application/json"},
		{"image/png", []byte{}, "image/png"},
		{"text/plain", []byte("This is plain text, not HTML."), "text/plain"},
	}

	ctx := context.Background()
	for _, ct := range contentTypes {
		t.Run(ct.contentType, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", ct.contentType)
				w.WriteHeader(http.StatusOK)
				w.Write(ct.content)
			}))
			defer testServer.Close()

			request := AnalyzerRequest{Url: testServer.URL}

			response, err := Analyze(ctx, request)
			assert.Error(t, err)
			assert.Equal(
				t,
				fmt.Sprintf(`Invalid response: %s`, ct.expectedErr),
				err.Error(),
			)
			assert.Nil(t, response)
		})
	}
}

func TestAnalyze_ValidUrl_UnknownHtmlVersion(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><head><title>Example Domain</title></head><body></body></html>"))
	}))

	// Test with a valid URL
	request := AnalyzerRequest{Url: testServer.URL}

	ctx := context.Background()
	response, err := Analyze(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Unknown", response.HtmlVersion)
	assert.Equal(t, "Example Domain", response.PageTitle)
	assert.Empty(t, response.Headings)
	assert.Empty(t, response.Links)
}

func TestAnalyze_ValidUrl_Html5Version(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<!DOCTYPE html><head><title>Example Domain</title></head><body></body></html>"))
	}))

	// Test with a valid URL
	request := AnalyzerRequest{Url: testServer.URL}

	ctx := context.Background()
	response, err := Analyze(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "HTML5", response.HtmlVersion)
	assert.Equal(t, "Example Domain", response.PageTitle)
	assert.Empty(t, response.Headings)
	assert.Empty(t, response.Links)
}

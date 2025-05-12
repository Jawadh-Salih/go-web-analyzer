package server

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(nil, nil))
	server := New(":8080", logger, false)

	assert.NotNil(t, server)
	assert.Equal(t, ":8080", server.port)
	assert.NotNil(t, server.router)
	assert.NotNil(t, server.logger)
}

func TestServerStart(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := New(":8080", logger, false)

	go func() {
		err := server.Start()
		assert.NoError(t, err)
	}()

	// Simulate a request to ensure the server is running
	req, err := http.NewRequest(http.MethodGet, "/not-found", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	err = server.Stop(context.Background()) // Stop the server after the test
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, w.Code) // Assuming no routes are defined

}

func TestServerMiddleware(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := New(":8080", logger, false)

	// Add a test route to verify middleware
	server.router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	err := server.Stop(context.Background()) // Stop the server after the test
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

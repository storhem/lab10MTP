package app_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go-service/app"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestPingEndpoint(t *testing.T) {
	r := app.SetupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "pong") {
		t.Errorf("expected body to contain 'pong', got %s", w.Body.String())
	}
}

func TestGetItemsEndpoint(t *testing.T) {
	r := app.SetupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Apple") {
		t.Errorf("expected body to contain items, got %s", w.Body.String())
	}
}

func TestGetItemByIDFound(t *testing.T) {
	r := app.SetupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/items/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGetItemByIDInvalidID(t *testing.T) {
	r := app.SetupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/items/abc", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestGetItemByIDNotFound(t *testing.T) {
	r := app.SetupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/items/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

// TestLoggerMiddlewareLogsFields проверяет, что middleware логирует
// метод, путь, статус, duration и ip для каждого запроса.
func TestLoggerMiddlewareLogsFields(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	r := app.SetupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)

	logOutput := buf.String()
	for _, expected := range []string{"GET", "/ping", "status=200", "duration=", "ip="} {
		if !strings.Contains(logOutput, expected) {
			t.Errorf("log should contain %q, got: %s", expected, logOutput)
		}
	}
}

// TestLoggerMiddlewareLogsNonOKStatus проверяет, что статус 404 тоже логируется корректно.
func TestLoggerMiddlewareLogsNonOKStatus(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	r := app.SetupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/items/999", nil)
	r.ServeHTTP(w, req)

	if !strings.Contains(buf.String(), "status=404") {
		t.Errorf("log should contain status=404, got: %s", buf.String())
	}
}

package gateway_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"api-gateway/gateway"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// mockBackend запускает тестовый HTTP-сервер, отвечающий фиксированным телом.
func mockBackend(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}))
}

func TestHealthEndpoint(t *testing.T) {
	ginSrv := mockBackend(`{"message":"pong"}`)
	defer ginSrv.Close()
	pythonSrv := mockBackend(`{"message":"pong"}`)
	defer pythonSrv.Close()

	cfg := gateway.Config{
		GinBackend:    ginSrv.URL,
		PythonBackend: pythonSrv.URL,
		ListenAddr:    ":8090",
	}
	r, err := gateway.SetupRouter(cfg)
	if err != nil {
		t.Fatalf("SetupRouter failed: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"status":"ok"`) {
		t.Errorf("expected status ok in body, got %s", w.Body.String())
	}
}

func TestGinRouteProxy(t *testing.T) {
	ginSrv := mockBackend(`{"message":"pong"}`)
	defer ginSrv.Close()
	pythonSrv := mockBackend(`{"message":"pong"}`)
	defer pythonSrv.Close()

	cfg := gateway.Config{GinBackend: ginSrv.URL, PythonBackend: pythonSrv.URL}
	r, _ := gateway.SetupRouter(cfg)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/gin/ping", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 from gin backend, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "pong") {
		t.Errorf("expected body from gin backend, got %s", w.Body.String())
	}
}

func TestPythonRouteProxy(t *testing.T) {
	ginSrv := mockBackend(`{"message":"pong"}`)
	defer ginSrv.Close()
	pythonSrv := mockBackend(`{"message":"pong"}`)
	defer pythonSrv.Close()

	cfg := gateway.Config{GinBackend: ginSrv.URL, PythonBackend: pythonSrv.URL}
	r, _ := gateway.SetupRouter(cfg)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/python/ping", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 from python backend, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "pong") {
		t.Errorf("expected body from python backend, got %s", w.Body.String())
	}
}

func TestGinAndPythonRoutedSeparately(t *testing.T) {
	ginSrv := mockBackend(`{"source":"gin"}`)
	defer ginSrv.Close()
	pythonSrv := mockBackend(`{"source":"python"}`)
	defer pythonSrv.Close()

	cfg := gateway.Config{GinBackend: ginSrv.URL, PythonBackend: pythonSrv.URL}
	r, _ := gateway.SetupRouter(cfg)

	wGin := httptest.NewRecorder()
	r.ServeHTTP(wGin, httptest.NewRequest(http.MethodGet, "/gin/items", nil))

	wPython := httptest.NewRecorder()
	r.ServeHTTP(wPython, httptest.NewRequest(http.MethodGet, "/python/items", nil))

	if !strings.Contains(wGin.Body.String(), "gin") {
		t.Errorf("/gin/* should route to gin backend, got %s", wGin.Body.String())
	}
	if !strings.Contains(wPython.Body.String(), "python") {
		t.Errorf("/python/* should route to python backend, got %s", wPython.Body.String())
	}
}

func TestBackendUnavailableReturns502(t *testing.T) {
	cfg := gateway.Config{
		GinBackend:    "http://localhost:19999",
		PythonBackend: "http://localhost:19998",
	}
	r, _ := gateway.SetupRouter(cfg)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/gin/ping", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("expected 502 when backend is unavailable, got %d", w.Code)
	}
}

func TestInvalidBackendURLReturnsError(t *testing.T) {
	cfg := gateway.Config{
		GinBackend:    "://bad-url",
		PythonBackend: "http://localhost:8000",
	}
	_, err := gateway.SetupRouter(cfg)
	if err == nil {
		t.Error("expected error for invalid backend URL, got nil")
	}
}

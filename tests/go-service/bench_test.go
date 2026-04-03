package app_test

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"go-service/app"
)

// silenceLogs подавляет вывод логгера на время бенчмарка и возвращает функцию restore.
func silenceLogs() func() {
	log.SetOutput(io.Discard)
	return func() { log.SetOutput(os.Stderr) }
}

func BenchmarkPingHandler(b *testing.B) {
	restore := silenceLogs()
	defer restore()

	r := app.SetupRouter()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkGetItemsHandler(b *testing.B) {
	restore := silenceLogs()
	defer restore()

	r := app.SetupRouter()
	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkGetItemByIDHandler(b *testing.B) {
	restore := silenceLogs()
	defer restore()

	r := app.SetupRouter()
	req := httptest.NewRequest(http.MethodGet, "/items/1", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

func BenchmarkMemoryHandler(b *testing.B) {
	restore := silenceLogs()
	defer restore()

	r := app.SetupRouter()
	req := httptest.NewRequest(http.MethodGet, "/memory", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

package gateway

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

// httpClient с явным таймаутом — предотвращает зависание при недоступном бэкенде.
var httpClient = &http.Client{Timeout: 10 * time.Second}

// NewProxy создаёт обратный прокси к указанному бэкенду.
// Возвращает ошибку, если target — невалидный URL.
func NewProxy(target string) (*httputil.ReverseProxy, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("invalid backend URL %q: %w", target, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.Transport = httpClient.Transport

	// Кастомный обработчик ошибок бэкенда.
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, `{"error":"backend unavailable","detail":%q}`, err.Error())
	}

	return proxy, nil
}

// noCloseNotifyWriter оборачивает http.ResponseWriter, скрывая http.CloseNotifier.
// Это необходимо, так как httputil.ReverseProxy пытается получить CloseNotifier
// из gin.ResponseWriter, который делегирует вызов к базовому writer'у,
// не реализующему этот интерфейс (например, в тестах — httptest.ResponseRecorder).
type noCloseNotifyWriter struct {
	http.ResponseWriter
}

// ProxyHandler возвращает gin-хэндлер, который стрипает prefix и проксирует запрос.
func ProxyHandler(proxy *httputil.ReverseProxy, prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Убираем prefix из пути, чтобы бэкенд получил оригинальный путь.
		// Например: /gin/items → /items
		c.Request.URL.Path = c.Param("path")
		if c.Request.URL.Path == "" {
			c.Request.URL.Path = "/"
		}
		proxy.ServeHTTP(noCloseNotifyWriter{c.Writer}, c.Request)
	}
}

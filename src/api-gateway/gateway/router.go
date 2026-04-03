package gateway

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Config хранит адреса бэкендов. Значения берутся из env или используются дефолтные.
type Config struct {
	GinBackend    string
	PythonBackend string
	ListenAddr    string
}

// ConfigFromEnv читает конфигурацию из переменных окружения.
func ConfigFromEnv() Config {
	cfg := Config{
		GinBackend:    "http://localhost:8080",
		PythonBackend: "http://localhost:8000",
		ListenAddr:    ":8090",
	}
	if v := os.Getenv("GIN_BACKEND"); v != "" {
		cfg.GinBackend = v
	}
	if v := os.Getenv("PYTHON_BACKEND"); v != "" {
		cfg.PythonBackend = v
	}
	if v := os.Getenv("GATEWAY_ADDR"); v != "" {
		cfg.ListenAddr = v
	}
	return cfg
}

// SetupRouter создаёт gin.Engine с маршрутами шлюза.
//
// Маршруты:
//
//	GET /health        — статус шлюза
//	ANY /gin/*path     — проксирование к Go-сервису (Gin)
//	ANY /python/*path  — проксирование к FastAPI-сервису
func SetupRouter(cfg Config) (*gin.Engine, error) {
	ginProxy, err := NewProxy(cfg.GinBackend)
	if err != nil {
		return nil, err
	}
	pythonProxy, err := NewProxy(cfg.PythonBackend)
	if err != nil {
		return nil, err
	}

	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":         "ok",
			"gin_backend":    cfg.GinBackend,
			"python_backend": cfg.PythonBackend,
		})
	})

	r.Any("/gin/*path", ProxyHandler(ginProxy, "/gin"))
	r.Any("/python/*path", ProxyHandler(pythonProxy, "/python"))

	return r, nil
}

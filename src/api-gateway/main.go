package main

import (
	"log"

	"api-gateway/gateway"
)

func main() {
	cfg := gateway.ConfigFromEnv()

	r, err := gateway.SetupRouter(cfg)
	if err != nil {
		log.Fatalf("failed to setup gateway: %v", err)
	}

	log.Printf("API Gateway listening on %s", cfg.ListenAddr)
	log.Printf("  /gin/*    → %s", cfg.GinBackend)
	log.Printf("  /python/* → %s", cfg.PythonBackend)

	if err := r.Run(cfg.ListenAddr); err != nil {
		log.Fatalf("gateway error: %v", err)
	}
}

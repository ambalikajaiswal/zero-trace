package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ambalikajaiswal/zero-trace/backend/handlers"
	"github.com/ambalikajaiswal/zero-trace/backend/middleware"
	"github.com/ambalikajaiswal/zero-trace/backend/store"
)

func main() {
	secretStore := store.NewSecretStore()

	// Start background goroutine for TTL-based secret cleanup
	go secretStore.StartCleanup()

	mux := http.NewServeMux()

	// REST API routes
	mux.HandleFunc("POST /api/secrets", handlers.CreateSecret(secretStore))
	mux.HandleFunc("GET /api/secrets/{id}", handlers.GetSecret(secretStore))
	mux.HandleFunc("GET /api/health", handlers.HealthCheck)

	// Apply middleware
	handler := middleware.CORS(middleware.RateLimit(middleware.Logger(mux)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("ZeroTrace server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

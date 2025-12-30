package main

import (
	"log"
	"net/http"
	"os"

	"github.com/teq0/musical-umbrella/pkg/api"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := api.NewRouter(logger)

	logger.Info("starting server",
		zap.String("port", port),
		zap.Strings("endpoints", []string{
			"GET /magic - Get a random XKCD comic quote",
		}),
	)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Fatal("server failed to start", zap.Error(err))
	}
}

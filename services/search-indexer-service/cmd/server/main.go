package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/b2b-platform/search-indexer-service/service"
	"github.com/b2b-platform/shared/events"
	"github.com/b2b-platform/shared/redis"
)

func main() {
	redisClient, err := redis.GetRedisClient()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	eventBus := events.NewRedisEventBus(redisClient)
	indexerService := service.NewIndexerService()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start health check HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8012"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"search-indexer-service"}`))
	})

	go func() {
		fmt.Printf("Health check server starting on port %s\n", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Printf("Health check server error: %v", err)
		}
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down search indexer...")
		cancel()
	}()

	// Subscribe to all events
	fmt.Println("Search indexer service starting...")
	fmt.Println("Subscribing to events...")

	handler := func(event *events.EventEnvelope) error {
		fmt.Printf("Received event: %s\n", event.Type)
		if err := indexerService.HandleEvent(event); err != nil {
			log.Printf("Error handling event %s: %v", event.Type, err)
			return err
		}
		fmt.Printf("Successfully indexed event: %s\n", event.Type)
		return nil
	}

	// Subscribe to all events
	if err := eventBus.SubscribeAll(ctx, handler); err != nil {
		if err == context.Canceled {
			fmt.Println("Search indexer stopped")
			return
		}
		log.Fatalf("Error in event subscription: %v", err)
	}
}

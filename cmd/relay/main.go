package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/filegate/filegate/internal/relay"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	domain := flag.String("domain", "filegate.app", "Base domain for subdomains")
	flag.Parse()

	// Allow environment variable override (PORT for Railway, RELAY_PORT as fallback)
	if envPort := os.Getenv("PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", port)
	} else if envPort := os.Getenv("RELAY_PORT"); envPort != "" {
		fmt.Sscanf(envPort, "%d", port)
	}
	if envDomain := os.Getenv("RELAY_DOMAIN"); envDomain != "" {
		*domain = envDomain
	}

	server := relay.NewServer(relay.Config{
		Domain: *domain,
		Port:   *port,
	})

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      server,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		httpServer.Shutdown(ctx)
	}()

	log.Printf("Relay server starting on :%d", *port)
	log.Printf("Domain: %s", *domain)
	log.Printf("Tunnel endpoint: ws://localhost:%d/tunnel", *port)
	log.Printf("Health check: http://localhost:%d/health", *port)

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

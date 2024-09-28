package main

import (
	"authService/internal/config"
	"authService/internal/handler"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var shutdownPassword = "password" // Will be secret - known to admins only and moved to encrypted config

// Channel to receive shutdown signals from the API call
var shutdownChan = make(chan struct{})

func main() {
	config.InitAppConfig()
	config.InitDB()
	defer config.CloseDB()

	handler.InitUserRoutes()
	handler.InitAuthRoutes()

	http.HandleFunc("/shutdown", shutdownHandler)

	port := fmt.Sprintf(":%s", config.ServerPort)
	server := &http.Server{
		Addr: port,
	}

	go func() {
		log.Printf("Auth service is running on port %s\n", config.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	shutdownServer(server)
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ShutdownPassword string `json:"shutdownPassword"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	// Will replace this check with config.shutdownPassword value
	if req.ShutdownPassword != shutdownPassword {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	go func() {
		shutdownChan <- struct{}{}
	}()
	w.Write([]byte("Server is shutting down..."))
}

func shutdownServer(server *http.Server) {
	<-shutdownChan
	log.Println("Shutdown signal received, initiating shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	} else {
		log.Println("Graceful shutdown complete.")
	}
}

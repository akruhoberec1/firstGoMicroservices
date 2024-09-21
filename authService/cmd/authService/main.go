package main

import (
	"authService/internal/config"
	"authService/internal/handler"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.InitAppConfig()
	config.InitDB()
	defer config.CloseDB()

	handler.InitUserRoutes()
	handler.InitAuthRoutes()

	port := fmt.Sprintf(":%s", config.ServerPort)
	log.Printf("Auth service is running on port %s\n", config.ServerPort)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

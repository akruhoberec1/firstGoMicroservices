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

	handler.InitRoutes()

	port := fmt.Sprintf(":%d", config.ServerPort)
	log.Printf("Auth service is running on port %s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

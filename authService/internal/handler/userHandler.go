package handler

import (
	"authService/internal/dto"
	"authService/internal/service"
	"encoding/json"
	"net/http"
)

func InitUserRoutes() {
	http.HandleFunc("/register", register)
}

func register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	err = service.RegisterUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("User registered successfully"))
}

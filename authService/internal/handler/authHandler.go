package handler

import (
	"authService/internal/dto"
	"authService/internal/helper"
	"authService/internal/service"
	"encoding/json"
	"net/http"
	"strings"
)

func InitAuthRoutes() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
}

func login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	tokens, err := service.Authenticate(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := helper.ValidateAccessToken(accessToken)
	if err != nil {
		http.Error(w, "Invalid or expired access token", http.StatusUnauthorized)
		return
	}

	userID, err := helper.ExtractUserIDFromToken(token)
	if err != nil {
		http.Error(w, "Failed to extract user ID from access token", http.StatusBadRequest)
		return
	}

	err = service.RevokeAllRefreshTokens(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Logged out successfully"))
}

package auth

import (
	"encoding/json"
	"net/http"
)

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string      `json:"token"`
	User  *UserPublic `json:"user"`
}

// UserPublic represents public user information (without password)
type UserPublic struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Get user from store
	userStore := GetUserStore()
	user, err := userStore.GetUserByUsernameOrEmail(req.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password
	if !userStore.VerifyPassword(user, req.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set cookie for HTML requests (24 hours expiration)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		MaxAge:   24 * 60 * 60, // 24 hours
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	// Return response
	response := LoginResponse{
		Token: token,
		User: &UserPublic{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// MeHandler returns the current authenticated user
func MeHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userStore := GetUserStore()
	user, err := userStore.GetUserByUsernameOrEmail(claims.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := &UserPublic{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// LogoutHandler handles user logout (mainly for frontend state management)
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Clear the auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// In a JWT-based system, logout is mainly client-side (token removal)
	// We could implement token blacklisting here if needed
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

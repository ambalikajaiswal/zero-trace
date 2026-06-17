package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ambalikajaiswal/zero-trace/backend/store"
)

// CreateSecretRequest represents the JSON body for creating a secret.
type CreateSecretRequest struct {
	Content    string `json:"content"`
	TTLSeconds int    `json:"ttl_seconds"`
	ViewOnce   bool   `json:"view_once"`
}

// CreateSecretResponse is returned after successful secret creation.
type CreateSecretResponse struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	ShareURL  string    `json:"share_url"`
}

// GetSecretResponse is returned when a secret is retrieved.
type GetSecretResponse struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ErrorResponse represents an API error.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// CreateSecret handles POST /api/secrets
func CreateSecret(s *store.SecretStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateSecretRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_request", "Invalid JSON body")
			return
		}

		// Validation
		if req.Content == "" {
			writeError(w, http.StatusBadRequest, "missing_content", "Secret content is required")
			return
		}

		if len(req.Content) > 10000 {
			writeError(w, http.StatusBadRequest, "content_too_large", "Secret content must be under 10,000 characters")
			return
		}

		if req.TTLSeconds <= 0 {
			req.TTLSeconds = 3600 // Default: 1 hour
		}

		if req.TTLSeconds > 604800 { // Max: 7 days
			req.TTLSeconds = 604800
		}

		id, token := s.Create(req.Content, req.TTLSeconds, req.ViewOnce)

		resp := CreateSecretResponse{
			ID:        id,
			Token:     token,
			ExpiresAt: time.Now().Add(time.Duration(req.TTLSeconds) * time.Second),
			ShareURL:  "/secret/" + id + "#" + token,
		}

		writeJSON(w, http.StatusCreated, resp)
	}
}

// GetSecret handles GET /api/secrets/{id}
func GetSecret(s *store.SecretStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			writeError(w, http.StatusBadRequest, "missing_id", "Secret ID is required")
			return
		}

		token := r.Header.Get("X-Access-Token")
		if token == "" {
			// Also check query param as fallback
			token = r.URL.Query().Get("token")
		}

		if token == "" {
			writeError(w, http.StatusUnauthorized, "missing_token", "Access token is required")
			return
		}

		secret, err := s.Get(id, token)
		if err != nil {
			switch err {
			case store.ErrSecretNotFound:
				writeError(w, http.StatusNotFound, "not_found", "Secret not found or has expired")
			case store.ErrInvalidToken:
				writeError(w, http.StatusForbidden, "forbidden", "Invalid access token")
			default:
				writeError(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred")
			}
			return
		}

		resp := GetSecretResponse{
			Content:   secret.Content,
			CreatedAt: secret.CreatedAt,
			ExpiresAt: secret.ExpiresAt,
		}

		writeJSON(w, http.StatusOK, resp)
	}
}

// HealthCheck handles GET /api/health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "zero-trace",
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, errCode, message string) {
	writeJSON(w, status, ErrorResponse{
		Error:   errCode,
		Message: message,
	})
}

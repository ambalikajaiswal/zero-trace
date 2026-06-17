package store

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Secret represents a stored secret with TTL metadata.
type Secret struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	ViewOnce  bool      `json:"view_once"`
	Token     string    `json:"-"` // Access token, never exposed in JSON
}

// SecretStore provides thread-safe in-memory storage with TTL-based expiration.
type SecretStore struct {
	mu      sync.RWMutex
	secrets map[string]*Secret
	stopCh  chan struct{}
}

var (
	ErrSecretNotFound = errors.New("secret not found or has expired")
	ErrInvalidToken   = errors.New("invalid access token")
)

// NewSecretStore creates a new thread-safe secret store.
func NewSecretStore() *SecretStore {
	return &SecretStore{
		secrets: make(map[string]*Secret),
		stopCh:  make(chan struct{}),
	}
}

// Create stores a new secret with the given TTL and returns its ID and access token.
func (s *SecretStore) Create(content string, ttlSeconds int, viewOnce bool) (string, string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	token := generateToken()

	secret := &Secret{
		ID:        id,
		Content:   content,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(ttlSeconds) * time.Second),
		ViewOnce:  viewOnce,
		Token:     token,
	}

	s.secrets[id] = secret
	return id, token
}

// Get retrieves a secret by ID and token. Automatically destroys view-once secrets.
func (s *SecretStore) Get(id, token string) (*Secret, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	secret, exists := s.secrets[id]
	if !exists {
		return nil, ErrSecretNotFound
	}

	// Check expiration
	if time.Now().After(secret.ExpiresAt) {
		delete(s.secrets, id)
		return nil, ErrSecretNotFound
	}

	// Validate access token (least-privilege access)
	if secret.Token != token {
		return nil, ErrInvalidToken
	}

	// If view-once, destroy after retrieval
	if secret.ViewOnce {
		delete(s.secrets, id)
	}

	return secret, nil
}

// StartCleanup runs a background goroutine that periodically removes expired secrets.
func (s *SecretStore) StartCleanup() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanup()
		case <-s.stopCh:
			return
		}
	}
}

// Stop signals the cleanup goroutine to stop.
func (s *SecretStore) Stop() {
	close(s.stopCh)
}

// cleanup removes all expired secrets from the store.
func (s *SecretStore) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for id, secret := range s.secrets {
		if now.After(secret.ExpiresAt) {
			delete(s.secrets, id)
		}
	}
}

// Count returns the number of active secrets (for health checks).
func (s *SecretStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.secrets)
}

// generateToken creates a cryptographically secure random token.
func generateToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to UUID if crypto/rand fails
		return uuid.New().String()
	}
	return hex.EncodeToString(bytes)
}

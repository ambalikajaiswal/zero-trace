package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"
)

// Logger logs incoming HTTP requests.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

// CORS adds Cross-Origin Resource Sharing headers.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Access-Token")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimit implements a simple in-memory rate limiter per IP.
func RateLimit(next http.Handler) http.Handler {
	type client struct {
		count    int
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Cleanup stale entries every minute
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		mu.Lock()
		c, exists := clients[ip]
		if !exists {
			clients[ip] = &client{count: 1, lastSeen: time.Now()}
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		// Reset counter if more than a minute has passed
		if time.Since(c.lastSeen) > time.Minute {
			c.count = 1
			c.lastSeen = time.Now()
			mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		c.count++
		c.lastSeen = time.Now()

		if c.count > 100 { // 100 requests per minute per IP
			mu.Unlock()
			http.Error(w, `{"error":"rate_limited","message":"Too many requests"}`, http.StatusTooManyRequests)
			return
		}

		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

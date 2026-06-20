package httpapi

import (
	"net/http"
	"time"
)

// CORS allows cross-origin requests from the dev frontend.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// InternalOnly returns middleware that rejects requests whose Authorization
// header does not match "Bearer <secret>". Use this to guard endpoints that
// must not be reachable from external clients.
func InternalOnly(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Authorization") != "Bearer "+secret {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimiter is a token-bucket middleware that caps the request rate on an
// endpoint. Excess requests receive 429 Too Many Requests.
type RateLimiter struct {
	tokens chan struct{}
}

// NewRateLimiter creates a RateLimiter with the given burst capacity that
// refills one token every refillInterval.
// Example: NewRateLimiter(10, 6*time.Second) → 10 req burst, 10 req/min sustained.
func NewRateLimiter(burst int, refillInterval time.Duration) *RateLimiter {
	rl := &RateLimiter{tokens: make(chan struct{}, burst)}
	for i := 0; i < burst; i++ {
		rl.tokens <- struct{}{}
	}
	go func() {
		ticker := time.NewTicker(refillInterval)
		defer ticker.Stop()
		for range ticker.C {
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		}
	}()
	return rl
}

// Wrap applies rate limiting to next.
func (rl *RateLimiter) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-rl.tokens:
			next.ServeHTTP(w, r)
		default:
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
		}
	})
}

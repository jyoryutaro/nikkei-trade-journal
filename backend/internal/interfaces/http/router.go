package httpapi

import (
	"net/http"
	"time"
)

// NewRouter builds the HTTP handler with routes and middleware applied.
// secret is used to protect the internal fetch endpoint.
func NewRouter(marketData *MarketDataHandler, journal *JournalHandler, secret string) http.Handler {
	mux := http.NewServeMux()

	// Public: reads stored candles from DB.
	mux.HandleFunc("GET /api/market-data", marketData.Get)
	mux.HandleFunc("GET /api/contracts", marketData.Contracts)

	// Internal: fetches from Yahoo Finance and persists to DB.
	// Rate-limited (burst 10, 10 req/min) and token-protected.
	limiter := NewRateLimiter(10, 6*time.Second)
	mux.Handle("POST /api/market-data/fetch",
		InternalOnly(secret)(limiter.Wrap(http.HandlerFunc(marketData.Fetch))))

	mux.HandleFunc("POST /api/journal-entries", journal.Create)
	mux.HandleFunc("GET /api/journal-entries", journal.List)

	return CORS(mux)
}

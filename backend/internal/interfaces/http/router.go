package httpapi

import "net/http"

// NewRouter builds the HTTP handler with routes and middleware applied.
func NewRouter(marketData *MarketDataHandler, journal *JournalHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/market-data", marketData.Get)
	mux.HandleFunc("POST /api/journal-entries", journal.Create)
	mux.HandleFunc("GET /api/journal-entries", journal.List)
	return CORS(mux)
}

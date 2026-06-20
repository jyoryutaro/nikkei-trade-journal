package httpapi

import "net/http"

// NewRouter builds the HTTP handler with routes and middleware applied.
func NewRouter(marketData *MarketDataHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/market-data", marketData.Get)
	return CORS(mux)
}

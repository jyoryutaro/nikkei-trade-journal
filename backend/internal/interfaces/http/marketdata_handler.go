package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/application"
	"github.com/min-legomain/nikkei-trade-journal/backend/internal/domain/marketdata"
)

// MarketDataHandler serves market-data HTTP endpoints.
type MarketDataHandler struct {
	svc *application.MarketDataService
}

// NewMarketDataHandler constructs the handler.
func NewMarketDataHandler(svc *application.MarketDataService) *MarketDataHandler {
	return &MarketDataHandler{svc: svc}
}

// Get handles GET /api/market-data?contract=&timeframe=.
func (h *MarketDataHandler) Get(w http.ResponseWriter, r *http.Request) {
	contract := r.URL.Query().Get("contract")

	tf, err := marketdata.ParseTimeframe(r.URL.Query().Get("timeframe"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	candles, err := h.svc.Candles(r.Context(), contract, tf)
	if err != nil {
		http.Error(w, "query error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toCandleDTOs(contract, tf, candles))
}

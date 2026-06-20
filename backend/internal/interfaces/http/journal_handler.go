package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/application"
	"github.com/min-legomain/nikkei-trade-journal/backend/internal/domain/journal"
)

// JournalHandler serves journal-entry HTTP endpoints.
type JournalHandler struct {
	svc *application.JournalService
}

// NewJournalHandler constructs the handler.
func NewJournalHandler(svc *application.JournalService) *JournalHandler {
	return &JournalHandler{svc: svc}
}

type createEntryRequest struct {
	Contract  string   `json:"contract"`
	Time      int64    `json:"time"` // Unix seconds (UTC)
	Side      string   `json:"side"` // "", "long", "short"
	TradeType string   `json:"tradeType"`
	Price     *float64 `json:"price"`
	Comment   string   `json:"comment"`
}

type entryDTO struct {
	ID        int64    `json:"id"`
	Contract  string   `json:"contract"`
	Time      int64    `json:"time"`
	Side      string   `json:"side"`
	TradeType string   `json:"tradeType"`
	Price     *float64 `json:"price"`
	Comment   string   `json:"comment"`
	CreatedAt int64    `json:"createdAt"`
}

func toEntryDTO(e journal.Entry) entryDTO {
	return entryDTO{
		ID:        e.ID,
		Contract:  e.Contract,
		Time:      e.Time.Unix(),
		Side:      string(e.Side),
		TradeType: string(e.TradeType),
		Price:     e.Price,
		Comment:   e.Comment,
		CreatedAt: e.CreatedAt.Unix(),
	}
}

// Create handles POST /api/journal-entries.
func (h *JournalHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	side, err := journal.ParseSide(req.Side)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tt, err := journal.ParseTradeType(req.TradeType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entry, err := journal.NewEntry(req.Contract, time.Unix(req.Time, 0).UTC(), side, tt, req.Price, req.Comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	saved, err := h.svc.Create(r.Context(), entry)
	if err != nil {
		http.Error(w, "could not save entry", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toEntryDTO(saved))
}

// List handles GET /api/journal-entries?contract=.
func (h *JournalHandler) List(w http.ResponseWriter, r *http.Request) {
	contract := r.URL.Query().Get("contract")
	entries, err := h.svc.List(r.Context(), contract)
	if err != nil {
		http.Error(w, "query error", http.StatusInternalServerError)
		return
	}
	out := make([]entryDTO, len(entries))
	for i, e := range entries {
		out[i] = toEntryDTO(e)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

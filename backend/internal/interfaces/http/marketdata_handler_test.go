package httpapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/application"
	"github.com/min-legomain/nikkei-trade-journal/backend/internal/domain/marketdata"
	httpapi "github.com/min-legomain/nikkei-trade-journal/backend/internal/interfaces/http"
)

type stubSource struct {
	contract string
	candles  []marketdata.Candle
}

func (s *stubSource) FetchCandles(_ context.Context, _ string) (string, []marketdata.Candle, error) {
	return s.contract, s.candles, nil
}

type stubRepo struct {
	candles []marketdata.Candle
}

func (r *stubRepo) FindBaseCandles(_ context.Context, _ string) ([]marketdata.Candle, error) {
	return r.candles, nil
}

func (r *stubRepo) BulkUpsert(_ context.Context, _ string, _ marketdata.Timeframe, candles []marketdata.Candle) (int, error) {
	return len(candles), nil
}

func (r *stubRepo) ListContracts(_ context.Context) ([]string, error) {
	return nil, nil
}

func newSvc(source *stubSource, repo *stubRepo) *application.MarketDataService {
	return application.NewMarketDataService(source, repo)
}

// TestGet_ValidRequestReturns200WithCandleJSON guarantees that a well-formed
// GET query returns 200 with a JSON array of candle objects from the DB.
func TestGet_ValidRequestReturns200WithCandleJSON(t *testing.T) {
	// Arrange
	candle := marketdata.NewCandle(time.Unix(0, 0).UTC(), 100, 110, 90, 105, 500)
	h := httpapi.NewMarketDataHandler(newSvc(&stubSource{}, &stubRepo{candles: []marketdata.Candle{candle}}))

	req := httptest.NewRequest(http.MethodGet, "/api/market-data?contract=%5EN225&timeframe=1m", nil)
	w := httptest.NewRecorder()

	// Act
	h.Get(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", w.Code)
	}
	var body []map[string]any
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if len(body) != 1 {
		t.Errorf("body length: got %d, want 1", len(body))
	}
}

// TestGet_InvalidTimeframeReturns400 guarantees that an unrecognised timeframe
// is rejected with Bad Request before reaching the service layer.
func TestGet_InvalidTimeframeReturns400(t *testing.T) {
	// Arrange
	h := httpapi.NewMarketDataHandler(newSvc(&stubSource{}, &stubRepo{}))

	req := httptest.NewRequest(http.MethodGet, "/api/market-data?timeframe=invalid", nil)
	w := httptest.NewRecorder()

	// Act
	h.Get(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", w.Code)
	}
}

// TestFetch_ValidRequestReturns200WithSavedCount guarantees that a POST to
// the fetch endpoint returns 200 with a JSON object containing the saved count.
func TestFetch_ValidRequestReturns200WithSavedCount(t *testing.T) {
	// Arrange
	candle := marketdata.NewCandle(time.Unix(0, 0).UTC(), 100, 110, 90, 105, 500)
	src := &stubSource{contract: "^N225", candles: []marketdata.Candle{candle}}
	h := httpapi.NewMarketDataHandler(newSvc(src, &stubRepo{}))

	req := httptest.NewRequest(http.MethodPost, "/api/market-data/fetch?symbol=%5EN225", nil)
	w := httptest.NewRecorder()

	// Act
	h.Fetch(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", w.Code)
	}
	var body map[string]int
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if body["saved"] != 1 {
		t.Errorf("saved: got %d, want 1", body["saved"])
	}
}

// TestFetch_MissingSymbolReturns400 guarantees that omitting the symbol
// parameter is rejected before calling the service.
func TestFetch_MissingSymbolReturns400(t *testing.T) {
	// Arrange
	h := httpapi.NewMarketDataHandler(newSvc(&stubSource{}, &stubRepo{}))

	req := httptest.NewRequest(http.MethodPost, "/api/market-data/fetch", nil)
	w := httptest.NewRecorder()

	// Act
	h.Fetch(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", w.Code)
	}
}

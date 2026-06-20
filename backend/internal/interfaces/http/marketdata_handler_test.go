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

type stubRepo struct {
	candles []marketdata.Candle
}

func (s *stubRepo) FindBaseCandles(_ context.Context, _ string) ([]marketdata.Candle, error) {
	return s.candles, nil
}

func (s *stubRepo) BulkUpsert(_ context.Context, _ string, _ marketdata.Timeframe, _ []marketdata.Candle) (int, error) {
	return 0, nil
}

// TestGet_ValidRequestReturns200WithCandleJSON guarantees that a well-formed
// query produces a 200 response whose body is a JSON array of candle objects.
func TestGet_ValidRequestReturns200WithCandleJSON(t *testing.T) {
	// Arrange
	candle := marketdata.NewCandle(time.Unix(0, 0).UTC(), 100, 110, 90, 105, 500)
	svc := application.NewMarketDataService(&stubRepo{candles: []marketdata.Candle{candle}})
	h := httpapi.NewMarketDataHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/market-data?contract=2609&timeframe=1m", nil)
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
// query parameter is rejected with Bad Request before reaching the service layer.
func TestGet_InvalidTimeframeReturns400(t *testing.T) {
	// Arrange
	svc := application.NewMarketDataService(&stubRepo{})
	h := httpapi.NewMarketDataHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/market-data?timeframe=invalid", nil)
	w := httptest.NewRecorder()

	// Act
	h.Get(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", w.Code)
	}
}

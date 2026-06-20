package yahoo_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/infrastructure/yahoo"
)

// TestFetchCandles_NonOKStatusReturnsError guarantees that a non-200 response
// from the upstream API surfaces as an error rather than returning empty candles.
func TestFetchCandles_NonOKStatusReturnsError(t *testing.T) {
	// Arrange
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()
	f := yahoo.NewFetcherAt(srv.Client(), srv.URL)

	// Act
	_, _, err := f.FetchCandles(context.Background(), "^N225")

	// Assert
	if err == nil {
		t.Fatal("expected error for non-200 status, got nil")
	}
}

// TestFetchCandles_EmptyContractFallsBackToSymbol guarantees that when the
// provider returns no YYMM contract code (e.g. a spot index), the symbol
// itself is used as the contract identifier so the DB row is never blank.
func TestFetchCandles_EmptyContractFallsBackToSymbol(t *testing.T) {
	// Arrange: valid JSON whose shortName has no parseable YYMM code
	body := `{"chart":{"result":[{"meta":{"symbol":"^N225","shortName":"Nikkei 225"},` +
		`"timestamp":[1000000],"indicators":{"quote":[{` +
		`"open":[100.0],"high":[110.0],"low":[90.0],"close":[105.0],"volume":[500]}]}}]}}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, body)
	}))
	defer srv.Close()
	f := yahoo.NewFetcherAt(srv.Client(), srv.URL)

	// Act
	contract, _, err := f.FetchCandles(context.Background(), "^N225")

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if contract != "^N225" {
		t.Errorf("contract: got %q, want %q (symbol fallback)", contract, "^N225")
	}
}

package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/application"
	"github.com/min-legomain/nikkei-trade-journal/backend/internal/domain/marketdata"
)

type stubRepo struct {
	candles []marketdata.Candle
	err     error
}

func (s *stubRepo) FindBaseCandles(_ context.Context, _ string) ([]marketdata.Candle, error) {
	return s.candles, s.err
}

func (s *stubRepo) BulkUpsert(_ context.Context, _ string, _ marketdata.Timeframe, _ []marketdata.Candle) (int, error) {
	return 0, nil
}

// TestCandles_AggregatesBaseCandlesToRequestedTimeframe guarantees that the
// service fetches base-timeframe candles from the repository and returns them
// collapsed into the caller's requested timeframe.
func TestCandles_AggregatesBaseCandlesToRequestedTimeframe(t *testing.T) {
	// Arrange: two 1m candles that fall inside the same 5m bucket
	base := []marketdata.Candle{
		marketdata.NewCandle(time.Unix(0, 0).UTC(), 100, 110, 95, 105, 10),
		marketdata.NewCandle(time.Unix(60, 0).UTC(), 105, 120, 90, 115, 20),
	}
	tf5m, _ := marketdata.ParseTimeframe("5m")
	svc := application.NewMarketDataService(&stubRepo{candles: base})

	// Act
	result, err := svc.Candles(context.Background(), "", tf5m)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 aggregated candle, got %d", len(result))
	}
	if result[0].Volume() != 30 {
		t.Errorf("Volume: got %d, want 30 (10+20 aggregated)", result[0].Volume())
	}
}

// TestCandles_PropagatesRepositoryError guarantees that a persistence failure
// surfaces unchanged to the caller and is never swallowed.
func TestCandles_PropagatesRepositoryError(t *testing.T) {
	// Arrange
	repoErr := errors.New("db unavailable")
	tf, _ := marketdata.ParseTimeframe("1m")
	svc := application.NewMarketDataService(&stubRepo{err: repoErr})

	// Act
	_, err := svc.Candles(context.Background(), "2609", tf)

	// Assert
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repo error to propagate, got: %v", err)
	}
}

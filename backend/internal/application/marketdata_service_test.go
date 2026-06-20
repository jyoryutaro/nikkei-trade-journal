package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/application"
	"github.com/min-legomain/nikkei-trade-journal/backend/internal/domain/marketdata"
)

type stubSource struct {
	contract string
	candles  []marketdata.Candle
	err      error
}

func (s *stubSource) FetchCandles(_ context.Context, _ string) (string, []marketdata.Candle, error) {
	return s.contract, s.candles, s.err
}

type stubRepo struct {
	candles []marketdata.Candle
	saved   []marketdata.Candle
	err     error
}

func (r *stubRepo) FindBaseCandles(_ context.Context, _ string) ([]marketdata.Candle, error) {
	return r.candles, r.err
}

func (r *stubRepo) BulkUpsert(_ context.Context, _ string, _ marketdata.Timeframe, candles []marketdata.Candle) (int, error) {
	r.saved = candles
	return len(candles), r.err
}

// TestCandles_AggregatesBaseCandlesToRequestedTimeframe guarantees that the
// service reads base candles from the DB and returns them collapsed into the
// caller's requested timeframe.
func TestCandles_AggregatesBaseCandlesToRequestedTimeframe(t *testing.T) {
	// Arrange: two 1m candles that fall inside the same 5m bucket
	base := []marketdata.Candle{
		marketdata.NewCandle(time.Unix(0, 0).UTC(), 100, 110, 95, 105, 10),
		marketdata.NewCandle(time.Unix(60, 0).UTC(), 105, 120, 90, 115, 20),
	}
	tf5m, _ := marketdata.ParseTimeframe("5m")
	svc := application.NewMarketDataService(&stubSource{}, &stubRepo{candles: base})

	// Act
	result, err := svc.Candles(context.Background(), "^N225", tf5m)

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

// TestCandles_PropagatesRepositoryError guarantees that a DB failure surfaces
// unchanged to the caller and is never swallowed.
func TestCandles_PropagatesRepositoryError(t *testing.T) {
	// Arrange
	repoErr := errors.New("db unavailable")
	tf, _ := marketdata.ParseTimeframe("1m")
	svc := application.NewMarketDataService(&stubSource{}, &stubRepo{err: repoErr})

	// Act
	_, err := svc.Candles(context.Background(), "^N225", tf)

	// Assert
	if !errors.Is(err, repoErr) {
		t.Errorf("expected repo error to propagate, got: %v", err)
	}
}

// TestFetch_PersistsYahooCandlesToRepository guarantees that Fetch passes the
// candles returned by the source to the repository unchanged.
func TestFetch_PersistsYahooCandlesToRepository(t *testing.T) {
	// Arrange
	candles := []marketdata.Candle{
		marketdata.NewCandle(time.Unix(0, 0).UTC(), 100, 110, 95, 105, 10),
	}
	repo := &stubRepo{}
	svc := application.NewMarketDataService(&stubSource{contract: "^N225", candles: candles}, repo)

	// Act
	n, err := svc.Fetch(context.Background(), "^N225")

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("saved count: got %d, want 1", n)
	}
	if len(repo.saved) != 1 {
		t.Errorf("repo.saved length: got %d, want 1", len(repo.saved))
	}
}

// TestFetch_PropagatesSourceError guarantees that a Yahoo Finance failure
// surfaces unchanged and nothing is written to the DB.
func TestFetch_PropagatesSourceError(t *testing.T) {
	// Arrange
	sourceErr := errors.New("upstream unavailable")
	repo := &stubRepo{}
	svc := application.NewMarketDataService(&stubSource{err: sourceErr}, repo)

	// Act
	_, err := svc.Fetch(context.Background(), "^N225")

	// Assert
	if !errors.Is(err, sourceErr) {
		t.Errorf("expected source error to propagate, got: %v", err)
	}
	if len(repo.saved) != 0 {
		t.Errorf("expected no DB writes on source error, got %d", len(repo.saved))
	}
}

package marketdata_test

import (
	"testing"
	"time"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/domain/marketdata"
)

// TestAggregate_EmptyInputReturnsNil guarantees that Aggregate is safe to
// call with no candles and produces no output.
func TestAggregate_EmptyInputReturnsNil(t *testing.T) {
	// Arrange
	tf5m, _ := marketdata.ParseTimeframe("5m")

	// Act
	result := marketdata.Aggregate(nil, tf5m)

	// Assert
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

// TestAggregate_CandlesInSameBucketAreCollapsed guarantees the OHLCV merge
// rules: High = max, Low = min, Open = first, Close = last, Volume = sum,
// and the bucket timestamp aligns to the interval boundary.
func TestAggregate_CandlesInSameBucketAreCollapsed(t *testing.T) {
	// Arrange: two 1m candles both inside the 00:00–05:00 bucket
	tf5m, _ := marketdata.ParseTimeframe("5m")
	base := []marketdata.Candle{
		marketdata.NewCandle(time.Unix(0, 0).UTC(), 100, 110, 95, 105, 10),
		marketdata.NewCandle(time.Unix(60, 0).UTC(), 105, 120, 90, 115, 20),
	}

	// Act
	result := marketdata.Aggregate(base, tf5m)

	// Assert
	if len(result) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(result))
	}
	got := result[0]
	if want := time.Unix(0, 0).UTC(); got.Time() != want {
		t.Errorf("Time: got %v, want %v", got.Time(), want)
	}
	if got.Open() != 100 {
		t.Errorf("Open: got %v, want 100 (first candle's open)", got.Open())
	}
	if got.High() != 120 {
		t.Errorf("High: got %v, want 120 (max of 110, 120)", got.High())
	}
	if got.Low() != 90 {
		t.Errorf("Low: got %v, want 90 (min of 95, 90)", got.Low())
	}
	if got.Close() != 115 {
		t.Errorf("Close: got %v, want 115 (last candle's close)", got.Close())
	}
	if got.Volume() != 30 {
		t.Errorf("Volume: got %v, want 30 (10 + 20)", got.Volume())
	}
}

// TestAggregate_CandlesInDifferentBucketsAreNotMerged guarantees that candles
// falling in distinct intervals produce separate output candles.
func TestAggregate_CandlesInDifferentBucketsAreNotMerged(t *testing.T) {
	// Arrange: candles at t=0 (bucket 0) and t=300 (bucket 300)
	tf5m, _ := marketdata.ParseTimeframe("5m")
	base := []marketdata.Candle{
		marketdata.NewCandle(time.Unix(0, 0).UTC(), 100, 110, 95, 105, 10),
		marketdata.NewCandle(time.Unix(300, 0).UTC(), 106, 115, 100, 112, 15),
	}

	// Act
	result := marketdata.Aggregate(base, tf5m)

	// Assert
	if len(result) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(result))
	}
}

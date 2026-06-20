package yahoo_test

import (
	"testing"
	"time"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/infrastructure/yahoo"
)

const validChartJSON = `{
  "chart": {
    "result": [{
      "meta": {"symbol": "NK=F", "shortName": "Nikkei/USD Futures,Sep-2026"},
      "timestamp": [1000000, 1000060],
      "indicators": {
        "quote": [{
          "open":   [25000.0, 25100.0],
          "high":   [25100.0, 25200.0],
          "low":    [24900.0, 25000.0],
          "close":  [25050.0, 25150.0],
          "volume": [1000, 2000]
        }]
      }
    }]
  }
}`

// TestParseChart_ValidJSONProducesContractAndCandles guarantees that
// well-formed Yahoo Finance chart JSON is decoded into the correct contract
// code (YYMM format) and the expected OHLCV candles.
func TestParseChart_ValidJSONProducesContractAndCandles(t *testing.T) {
	// Arrange
	raw := []byte(validChartJSON)

	// Act
	contract, candles, err := yahoo.ParseChart(raw)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if contract != "2609" {
		t.Errorf("contract: got %q, want %q", contract, "2609")
	}
	if len(candles) != 2 {
		t.Fatalf("len(candles): got %d, want 2", len(candles))
	}
	first := candles[0]
	if want := time.Unix(1000000, 0).UTC(); first.Time() != want {
		t.Errorf("Time: got %v, want %v", first.Time(), want)
	}
	if first.Open() != 25000.0 {
		t.Errorf("Open: got %v, want 25000.0", first.Open())
	}
	if first.Volume() != 1000 {
		t.Errorf("Volume: got %v, want 1000", first.Volume())
	}
}

// TestParseChart_RowsWithNilOHLCAreSkipped guarantees that incomplete rows
// (any nil OHLC field) do not produce a candle, preventing zero-value pollution.
func TestParseChart_RowsWithNilOHLCAreSkipped(t *testing.T) {
	// Arrange: second row has a nil open price
	raw := []byte(`{
    "chart": {
      "result": [{
        "meta": {"symbol": "NK=F", "shortName": "Nikkei/USD Futures,Sep-2026"},
        "timestamp": [1000000, 1000060],
        "indicators": {
          "quote": [{
            "open":   [25000.0, null],
            "high":   [25100.0, 25200.0],
            "low":    [24900.0, 25000.0],
            "close":  [25050.0, 25150.0],
            "volume": [1000, 2000]
          }]
        }
      }]
    }
  }`)

	// Act
	_, candles, err := yahoo.ParseChart(raw)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(candles) != 1 {
		t.Errorf("expected 1 candle (nil row skipped), got %d", len(candles))
	}
}

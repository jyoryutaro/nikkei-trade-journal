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
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	jstT := time.Unix(1000000, 0).In(jst)
	wantTime := time.Date(jstT.Year(), jstT.Month(), jstT.Day(), jstT.Hour(), jstT.Minute(), jstT.Second(), 0, time.UTC)
	if first.Time() != wantTime {
		t.Errorf("Time: got %v, want %v", first.Time(), wantTime)
	}
	if first.Open() != 25000.0 {
		t.Errorf("Open: got %v, want 25000.0", first.Open())
	}
	if first.Volume() != 1000 {
		t.Errorf("Volume: got %v, want 1000", first.Volume())
	}
}

// TestParseChart_NilOHLCForwardFilled guarantees that a row with a nil OHLC
// field is forward-filled with the previous close (a flat bar, volume 0) so
// the series stays continuous, rather than leaving a gap.
func TestParseChart_NilOHLCForwardFilled(t *testing.T) {
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
	if len(candles) != 2 {
		t.Fatalf("expected 2 candles (nil row forward-filled), got %d", len(candles))
	}
	filled := candles[1]
	jst2 := time.FixedZone("Asia/Tokyo", 9*60*60)
	jstT2 := time.Unix(1000060, 0).In(jst2)
	wantFilled := time.Date(jstT2.Year(), jstT2.Month(), jstT2.Day(), jstT2.Hour(), jstT2.Minute(), jstT2.Second(), 0, time.UTC)
	if filled.Time() != wantFilled {
		t.Errorf("Time: got %v, want %v", filled.Time(), wantFilled)
	}
	prevClose := 25050.0
	if filled.Open() != prevClose || filled.High() != prevClose || filled.Low() != prevClose || filled.Close() != prevClose {
		t.Errorf("filled OHLC: got O=%v H=%v L=%v C=%v, want all %v",
			filled.Open(), filled.High(), filled.Low(), filled.Close(), prevClose)
	}
	if filled.Volume() != 0 {
		t.Errorf("filled Volume: got %v, want 0", filled.Volume())
	}
}

// TestParseChart_LeadingGapDropped guarantees that nil rows before the first
// real value are dropped (nothing to forward-fill from).
func TestParseChart_LeadingGapDropped(t *testing.T) {
	// Arrange: first row is nil, second row is valid
	raw := []byte(`{
    "chart": {
      "result": [{
        "meta": {"symbol": "NK=F", "shortName": "Nikkei/USD Futures,Sep-2026"},
        "timestamp": [1000000, 1000060],
        "indicators": {
          "quote": [{
            "open":   [null, 25100.0],
            "high":   [null, 25200.0],
            "low":    [null, 25000.0],
            "close":  [null, 25150.0],
            "volume": [null, 2000]
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
		t.Fatalf("expected 1 candle (leading gap dropped), got %d", len(candles))
	}
	if candles[0].Open() != 25100.0 {
		t.Errorf("Open: got %v, want 25100.0", candles[0].Open())
	}
}

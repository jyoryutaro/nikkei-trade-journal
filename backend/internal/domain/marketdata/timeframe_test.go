package marketdata_test

import (
	"testing"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/domain/marketdata"
)

// TestParseTimeframe_KnownIdentifiersResolve guarantees that every supported
// timeframe string maps to the correct String() and Seconds() values.
func TestParseTimeframe_KnownIdentifiersResolve(t *testing.T) {
	cases := []struct {
		id      string
		seconds int64
	}{
		{"1m", 60},
		{"5m", 300},
		{"30m", 1800},
		{"1h", 3600},
		{"1d", 86400},
	}
	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			// Act
			tf, err := marketdata.ParseTimeframe(tc.id)

			// Assert
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tf.String() != tc.id {
				t.Errorf("String(): got %q, want %q", tf.String(), tc.id)
			}
			if tf.Seconds() != tc.seconds {
				t.Errorf("Seconds(): got %d, want %d", tf.Seconds(), tc.seconds)
			}
		})
	}
}

// TestParseTimeframe_EmptyStringDefaultsToBase guarantees that an absent
// timeframe parameter (empty string) falls back to the 1-minute base granularity.
func TestParseTimeframe_EmptyStringDefaultsToBase(t *testing.T) {
	// Act
	tf, err := marketdata.ParseTimeframe("")

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tf != marketdata.BaseTimeframe {
		t.Errorf("got %v, want BaseTimeframe", tf)
	}
}

// TestParseTimeframe_UnknownIdentifierReturnsError guarantees that an
// unrecognised string is rejected rather than silently defaulting to any value.
func TestParseTimeframe_UnknownIdentifierReturnsError(t *testing.T) {
	// Act
	_, err := marketdata.ParseTimeframe("2m")

	// Assert
	if err == nil {
		t.Fatal("expected an error for unsupported timeframe, got nil")
	}
}

package marketdata

import "fmt"

// Timeframe is a value object representing a candle interval (e.g. 1m, 1h).
type Timeframe struct {
	id      string
	seconds int64
}

// BaseTimeframe is the granularity persisted in the repository. All other
// timeframes are derived from it by aggregation.
var BaseTimeframe = Timeframe{id: "1m", seconds: 60}

// supported lists every timeframe the domain understands.
var supported = []Timeframe{
	BaseTimeframe,
	{id: "5m", seconds: 300},
	{id: "30m", seconds: 1800},
	{id: "1h", seconds: 3600},
	{id: "1d", seconds: 86400},
}

// ParseTimeframe converts an identifier into a Timeframe. An empty string
// defaults to the base timeframe.
func ParseTimeframe(s string) (Timeframe, error) {
	if s == "" {
		return BaseTimeframe, nil
	}
	for _, tf := range supported {
		if tf.id == s {
			return tf, nil
		}
	}
	return Timeframe{}, fmt.Errorf("unsupported timeframe: %q", s)
}

// String returns the timeframe identifier (e.g. "5m").
func (t Timeframe) String() string { return t.id }

// Seconds returns the timeframe duration in seconds.
func (t Timeframe) Seconds() int64 { return t.seconds }

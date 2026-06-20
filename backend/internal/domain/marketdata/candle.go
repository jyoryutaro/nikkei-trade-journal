package marketdata

import "time"

// Candle is an OHLCV bar. It is a value object: identified purely by its
// values, with no identity of its own.
type Candle struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

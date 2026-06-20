package marketdata

import "time"

// Candle is an OHLCV bar. It is a value object: identified purely by its
// values, with no identity of its own.
type Candle struct {
	t      time.Time
	open   float64
	high   float64
	low    float64
	close  float64
	volume int64
}

// NewCandle constructs a Candle value object.
func NewCandle(t time.Time, open, high, low, close float64, volume int64) Candle {
	return Candle{t: t, open: open, high: high, low: low, close: close, volume: volume}
}

func (c Candle) Time() time.Time { return c.t }
func (c Candle) Open() float64   { return c.open }
func (c Candle) High() float64   { return c.high }
func (c Candle) Low() float64    { return c.low }
func (c Candle) Close() float64  { return c.close }
func (c Candle) Volume() int64   { return c.volume }

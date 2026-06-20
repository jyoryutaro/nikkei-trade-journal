package marketdata

import (
	"sort"
	"time"
)

// Aggregate collapses base candles into buckets of the given timeframe.
// Buckets are aligned to the Unix epoch (e.g. 5m buckets start at :00, :05,
// :10 …). The input candles must be sorted by time ascending.
//
// This is a pure domain service: it has no dependency on persistence or
// transport and is deterministic given its inputs.
func Aggregate(candles []Candle, tf Timeframe) []Candle {
	interval := tf.Seconds()
	if len(candles) == 0 || interval <= 0 {
		return nil
	}

	type bucket struct {
		open   float64
		high   float64
		low    float64
		close  float64
		volume int64
	}

	buckets := map[int64]*bucket{}
	var keys []int64

	for _, c := range candles {
		k := (c.Time().Unix() / interval) * interval
		bkt, exists := buckets[k]
		if !exists {
			bkt = &bucket{open: c.Open(), high: c.High(), low: c.Low()}
			buckets[k] = bkt
			keys = append(keys, k)
		}
		if c.High() > bkt.high {
			bkt.high = c.High()
		}
		if c.Low() < bkt.low {
			bkt.low = c.Low()
		}
		bkt.close = c.Close()
		bkt.volume += c.Volume()
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	result := make([]Candle, len(keys))
	for i, k := range keys {
		bkt := buckets[k]
		result[i] = NewCandle(time.Unix(k, 0).UTC(), bkt.open, bkt.high, bkt.low, bkt.close, bkt.volume)
	}
	return result
}

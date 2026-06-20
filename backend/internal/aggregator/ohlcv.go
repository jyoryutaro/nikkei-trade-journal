package aggregator

import (
	"sort"
	"time"
)

// Bar is a single OHLCV candlestick.
type Bar struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

// IntervalSeconds maps timeframe identifiers to their duration in seconds.
var IntervalSeconds = map[string]int64{
	"1m":  60,
	"5m":  300,
	"30m": 1800,
	"1h":  3600,
	"1d":  86400,
}

// Aggregate collapses bars into buckets of intervalSec seconds.
// Buckets are aligned to Unix epoch (e.g. 5m buckets start at :00, :05, :10 …).
// The input bars must be sorted by time ascending.
func Aggregate(bars []Bar, intervalSec int64) []Bar {
	if len(bars) == 0 || intervalSec <= 0 {
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

	for _, b := range bars {
		k := (b.Time.Unix() / intervalSec) * intervalSec
		if _, exists := buckets[k]; !exists {
			buckets[k] = &bucket{open: b.Open, high: b.High, low: b.Low}
			keys = append(keys, k)
		}
		bkt := buckets[k]
		if b.High > bkt.high {
			bkt.high = b.High
		}
		if b.Low < bkt.low {
			bkt.low = b.Low
		}
		bkt.close = b.Close
		bkt.volume += b.Volume
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	result := make([]Bar, len(keys))
	for i, k := range keys {
		bkt := buckets[k]
		result[i] = Bar{
			Time:   time.Unix(k, 0).UTC(),
			Open:   bkt.open,
			High:   bkt.high,
			Low:    bkt.low,
			Close:  bkt.close,
			Volume: bkt.volume,
		}
	}
	return result
}

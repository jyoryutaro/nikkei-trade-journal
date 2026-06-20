package httpapi

import "github.com/min-legomain/nikkei-trade-journal/backend/internal/domain/marketdata"

// candleDTO is the JSON representation of a candle returned by the API.
type candleDTO struct {
	Contract  string  `json:"contract"`
	Timeframe string  `json:"timeframe"`
	Time      int64   `json:"time"` // Unix timestamp (seconds, UTC)
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    int64   `json:"volume"`
}

// toCandleDTOs maps domain candles to the transport representation.
func toCandleDTOs(contract string, tf marketdata.Timeframe, candles []marketdata.Candle) []candleDTO {
	out := make([]candleDTO, len(candles))
	for i, c := range candles {
		out[i] = candleDTO{
			Contract:  contract,
			Timeframe: tf.String(),
			Time:      c.Time().Unix(),
			Open:      c.Open(),
			High:      c.High(),
			Low:       c.Low(),
			Close:     c.Close(),
			Volume:    c.Volume(),
		}
	}
	return out
}

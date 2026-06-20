package marketdata

import "context"

// CandleSource fetches raw base-timeframe candles from an external data provider.
type CandleSource interface {
	// FetchCandles returns the contract code derived from the provider and
	// 1-minute candles for the given symbol.
	FetchCandles(ctx context.Context, symbol string) (contract string, candles []Candle, err error)
}

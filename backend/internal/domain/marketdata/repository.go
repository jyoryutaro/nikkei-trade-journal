package marketdata

import "context"

// Repository is the persistence boundary (port) for market data.
type Repository interface {
	// FindBaseCandles returns the stored base-timeframe candles for a contract,
	// sorted by time ascending. An empty contract returns all contracts.
	FindBaseCandles(ctx context.Context, contract string) ([]Candle, error)

	// BulkUpsert stores candles for a contract/timeframe and returns the number
	// of rows written.
	BulkUpsert(ctx context.Context, contract string, tf Timeframe, candles []Candle) (int, error)

	// ListContracts returns the distinct contract codes present in storage,
	// sorted ascending.
	ListContracts(ctx context.Context) ([]string, error)
}

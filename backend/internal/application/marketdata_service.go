package application

import (
	"context"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/domain/marketdata"
)

// MarketDataService orchestrates market-data use cases. It depends only on the
// domain repository port, keeping it free of persistence and transport details.
type MarketDataService struct {
	repo marketdata.Repository
}

// NewMarketDataService wires the service with a repository implementation.
func NewMarketDataService(repo marketdata.Repository) *MarketDataService {
	return &MarketDataService{repo: repo}
}

// Candles returns OHLCV for a contract aggregated to the requested timeframe.
// Only base candles are stored; higher timeframes are derived on the fly.
func (s *MarketDataService) Candles(ctx context.Context, contract string, tf marketdata.Timeframe) ([]marketdata.Candle, error) {
	base, err := s.repo.FindBaseCandles(ctx, contract)
	if err != nil {
		return nil, err
	}
	return marketdata.Aggregate(base, tf), nil
}

// Import stores imported candles for a contract/timeframe (used by the seeder).
func (s *MarketDataService) Import(ctx context.Context, contract string, tf marketdata.Timeframe, candles []marketdata.Candle) (int, error) {
	return s.repo.BulkUpsert(ctx, contract, tf, candles)
}

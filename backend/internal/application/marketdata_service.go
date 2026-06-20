package application

import (
	"context"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/domain/marketdata"
)

// MarketDataService orchestrates market-data use cases.
type MarketDataService struct {
	source marketdata.CandleSource
	repo   marketdata.Repository
}

// NewMarketDataService wires the service with a Yahoo Finance source and a DB repository.
func NewMarketDataService(source marketdata.CandleSource, repo marketdata.Repository) *MarketDataService {
	return &MarketDataService{source: source, repo: repo}
}

// Fetch pulls base candles for symbol from the external source and persists them.
// Returns the number of rows written.
func (s *MarketDataService) Fetch(ctx context.Context, symbol string) (int, error) {
	contract, candles, err := s.source.FetchCandles(ctx, symbol)
	if err != nil {
		return 0, err
	}
	return s.repo.BulkUpsert(ctx, contract, marketdata.BaseTimeframe, candles)
}

// Candles reads base candles for contract from the DB and aggregates to tf.
func (s *MarketDataService) Candles(ctx context.Context, contract string, tf marketdata.Timeframe) ([]marketdata.Candle, error) {
	base, err := s.repo.FindBaseCandles(ctx, contract)
	if err != nil {
		return nil, err
	}
	return marketdata.Aggregate(base, tf), nil
}

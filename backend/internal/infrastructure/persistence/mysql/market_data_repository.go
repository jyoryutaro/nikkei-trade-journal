package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/domain/marketdata"
)

// MarketDataRepository is a MySQL-backed marketdata.Repository.
type MarketDataRepository struct {
	db *sql.DB
}

var _ marketdata.Repository = (*MarketDataRepository)(nil)

// NewMarketDataRepository constructs the repository.
func NewMarketDataRepository(db *sql.DB) *MarketDataRepository {
	return &MarketDataRepository{db: db}
}

// FindBaseCandles returns the base-timeframe candles for a contract.
func (r *MarketDataRepository) FindBaseCandles(ctx context.Context, contract string) ([]marketdata.Candle, error) {
	query := `SELECT ts, open, high, low, close, volume
	          FROM market_data
	          WHERE timeframe = ?`
	args := []any{marketdata.BaseTimeframe.String()}
	if contract != "" {
		query += ` AND contract = ?`
		args = append(args, contract)
	}
	query += ` ORDER BY ts ASC LIMIT 10000`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candles []marketdata.Candle
	for rows.Next() {
		var t time.Time
		var open, high, low, close float64
		var volume int64
		if err := rows.Scan(&t, &open, &high, &low, &close, &volume); err != nil {
			return nil, err
		}
		candles = append(candles, marketdata.NewCandle(t, open, high, low, close, volume))
	}
	return candles, rows.Err()
}

// BulkUpsert inserts or updates candles for a contract/timeframe in one tx.
func (r *MarketDataRepository) BulkUpsert(ctx context.Context, contract string, tf marketdata.Timeframe, candles []marketdata.Candle) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO market_data (contract, timeframe, ts, open, high, low, close, volume)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE open=VALUES(open), high=VALUES(high),
		  low=VALUES(low), close=VALUES(close), volume=VALUES(volume)
	`)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	count := 0
	for _, c := range candles {
		if _, err := stmt.ExecContext(ctx, contract, tf.String(), c.Time().UTC(),
			c.Open(), c.High(), c.Low(), c.Close(), c.Volume()); err != nil {
			tx.Rollback()
			return 0, err
		}
		count++
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return count, nil
}

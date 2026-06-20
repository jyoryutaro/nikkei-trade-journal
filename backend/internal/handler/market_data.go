package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/aggregator"
)

type candlestick struct {
<<<<<<< Updated upstream
	Contract string    `json:"contract"`
	Time     time.Time `json:"time"`
	Open     float64   `json:"open"`
	High     float64   `json:"high"`
	Low      float64   `json:"low"`
	Close    float64   `json:"close"`
	Volume   int64     `json:"volume"`
=======
	Contract  string  `json:"contract"`
	Timeframe string  `json:"timeframe"`
	Time      int64   `json:"time"` // Unix timestamp (seconds UTC)
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    int64   `json:"volume"`
>>>>>>> Stashed changes
}

func MarketData(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contract := r.URL.Query().Get("contract")
<<<<<<< Updated upstream
		query := `SELECT contract, ts, open, high, low, close, volume FROM market_data`
=======
		tf := r.URL.Query().Get("timeframe")
		if tf == "" {
			tf = "1m"
		}

		intervalSec, ok := aggregator.IntervalSeconds[tf]
		if !ok {
			http.Error(w, "unsupported timeframe", http.StatusBadRequest)
			return
		}

		query := `SELECT ts, open, high, low, close, volume
		          FROM market_data
		          WHERE timeframe = '1m'`
>>>>>>> Stashed changes
		args := []any{}
		if contract != "" {
			query += ` WHERE contract = ?`
			args = append(args, contract)
		}
<<<<<<< Updated upstream
		query += ` ORDER BY ts ASC LIMIT 500`
=======
		query += ` ORDER BY ts ASC LIMIT 10000`
>>>>>>> Stashed changes

		rows, err := db.QueryContext(r.Context(), query, args...)
		if err != nil {
			http.Error(w, "query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var bars []aggregator.Bar
		for rows.Next() {
<<<<<<< Updated upstream
			var c candlestick
			if err := rows.Scan(&c.Contract, &c.Time, &c.Open, &c.High, &c.Low, &c.Close, &c.Volume); err != nil {
=======
			var b aggregator.Bar
			if err := rows.Scan(&b.Time, &b.Open, &b.High, &b.Low, &b.Close, &b.Volume); err != nil {
>>>>>>> Stashed changes
				http.Error(w, "scan error", http.StatusInternalServerError)
				return
			}
			bars = append(bars, b)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "rows error", http.StatusInternalServerError)
			return
		}

		agg := aggregator.Aggregate(bars, intervalSec)

		results := make([]candlestick, len(agg))
		for i, b := range agg {
			results[i] = candlestick{
				Contract:  contract,
				Timeframe: tf,
				Time:      b.Time.Unix(),
				Open:      b.Open,
				High:      b.High,
				Low:       b.Low,
				Close:     b.Close,
				Volume:    b.Volume,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type candlestick struct {
	Contract  string    `json:"contract"`
	Timeframe string    `json:"timeframe"`
	Time      time.Time `json:"time"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    int64     `json:"volume"`
}

func MarketData(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contract := r.URL.Query().Get("contract")
		timeframe := r.URL.Query().Get("timeframe")
		if timeframe == "" {
			timeframe = "1m"
		}

		query := `SELECT contract, timeframe, ts, open, high, low, close, volume FROM market_data WHERE timeframe = ?`
		args := []any{timeframe}
		if contract != "" {
			query += ` AND contract = ?`
			args = append(args, contract)
		}
		query += ` ORDER BY ts ASC LIMIT 2000`

		rows, err := db.QueryContext(r.Context(), query, args...)
		if err != nil {
			http.Error(w, "query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		results := []candlestick{}
		for rows.Next() {
			var c candlestick
			if err := rows.Scan(&c.Contract, &c.Timeframe, &c.Time, &c.Open, &c.High, &c.Low, &c.Close, &c.Volume); err != nil {
				http.Error(w, "scan error", http.StatusInternalServerError)
				return
			}
			results = append(results, c)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, "rows error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Yahoo Finance v8 chart API shape (minimal)
type yfResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Symbol    string `json:"symbol"`
				ShortName string `json:"shortName"`
			} `json:"meta"`
			Timestamps []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []*float64 `json:"open"`
					High   []*float64 `json:"high"`
					Low    []*float64 `json:"low"`
					Close  []*float64 `json:"close"`
					Volume []*int64   `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}

func main() {
	dataPath := flag.String("data", "", "path to Yahoo Finance v8 JSON file")
	contract := flag.String("contract", "", "contract code e.g. 2609 (auto-detected from shortName if empty)")
	timeframe := flag.String("timeframe", "1m", "timeframe e.g. 1m, 5m, 1d")
	flag.Parse()

	if *dataPath == "" {
		log.Fatal("usage: seed -data <path-to-json> [-contract 2609] [-timeframe 1m]")
	}

	raw, err := os.ReadFile(*dataPath)
	if err != nil {
		log.Fatalf("read file: %v", err)
	}

	var resp yfResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		log.Fatalf("parse json: %v", err)
	}

	if len(resp.Chart.Result) == 0 {
		log.Fatal("no results in JSON")
	}
	result := resp.Chart.Result[0]

	ct := *contract
	if ct == "" {
		ct = contractFromShortName(result.Meta.ShortName)
	}
	if ct == "" {
		log.Fatalf("could not determine contract from shortName %q — pass -contract explicitly", result.Meta.ShortName)
	}

	if len(result.Indicators.Quote) == 0 {
		log.Fatal("no quote data")
	}
	q := result.Indicators.Quote[0]

	dsn := getEnv("DB_DSN", "app:app@tcp(localhost:3306)/nikkei_trade?parseTime=true")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("begin: %v", err)
	}
	stmt, err := tx.Prepare(`
		INSERT INTO market_data (contract, timeframe, ts, open, high, low, close, volume)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE open=VALUES(open), high=VALUES(high),
		  low=VALUES(low), close=VALUES(close), volume=VALUES(volume)
	`)
	if err != nil {
		log.Fatalf("prepare: %v", err)
	}
	defer stmt.Close()

	inserted := 0
	for i, ts := range result.Timestamps {
		if i >= len(q.Open) || q.Open[i] == nil || q.High[i] == nil || q.Low[i] == nil || q.Close[i] == nil {
			continue
		}
		vol := int64(0)
		if i < len(q.Volume) && q.Volume[i] != nil {
			vol = *q.Volume[i]
		}
		t := time.Unix(ts, 0).UTC()
		if _, err := stmt.Exec(ct, *timeframe, t, *q.Open[i], *q.High[i], *q.Low[i], *q.Close[i], vol); err != nil {
			tx.Rollback()
			log.Fatalf("insert row %d: %v", i, err)
		}
		inserted++
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("commit: %v", err)
	}
	log.Printf("seeded %d bars → contract=%s timeframe=%s", inserted, ct, *timeframe)
}

// contractFromShortName extracts a YYMM code from strings like
// "Nikkei/USD Futures,Sep-2026" → "2609"
func contractFromShortName(name string) string {
	monthMap := map[string]string{
		"Jan": "01", "Feb": "02", "Mar": "03", "Apr": "04",
		"May": "05", "Jun": "06", "Jul": "07", "Aug": "08",
		"Sep": "09", "Oct": "10", "Nov": "11", "Dec": "12",
	}
	// find "Mon-YYYY" pattern
	for _, part := range strings.FieldsFunc(name, func(r rune) bool { return r == ' ' || r == ',' }) {
		if len(part) == 8 && part[3] == '-' {
			mon := part[:3]
			year := part[4:]
			if m, ok := monthMap[mon]; ok && len(year) == 4 {
				return fmt.Sprintf("%s%s", year[2:], m)
			}
		}
	}
	return ""
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

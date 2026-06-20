package main

import (
	"database/sql"
	"log"
	"math/rand/v2"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := getEnv("DB_DSN", "app:app@tcp(localhost:3306)/nikkei_trade?parseTime=true&loc=Asia%2FTokyo")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	const contract = "2506"
	jst := time.FixedZone("JST", 9*60*60)

	// 2025-04-01 08:45 JST から 1日足を60本生成
	base := time.Date(2025, 4, 1, 0, 0, 0, 0, jst)
	price := 35000.0

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("begin: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO market_data (contract, ts, open, high, low, close, volume)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE open=VALUES(open), high=VALUES(high),
		  low=VALUES(low), close=VALUES(close), volume=VALUES(volume)
	`)
	if err != nil {
		log.Fatalf("prepare: %v", err)
	}
	defer stmt.Close()

	for i := range 60 {
		ts := base.AddDate(0, 0, i)
		// skip weekends
		if ts.Weekday() == time.Saturday || ts.Weekday() == time.Sunday {
			price += (rand.Float64()-0.5)*200
			continue
		}
		chg := (rand.Float64() - 0.45) * 400
		open := round(price)
		close_ := round(price + chg)
		high := round(max(open, close_) + rand.Float64()*150)
		low := round(min(open, close_) - rand.Float64()*150)
		vol := int64(1000 + rand.IntN(5000))

		if _, err := stmt.Exec(contract, ts, open, high, low, close_, vol); err != nil {
			tx.Rollback()
			log.Fatalf("insert row %d: %v", i, err)
		}
		price = close_
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("commit: %v", err)
	}
	log.Printf("seeded %s market data", contract)
}

func round(f float64) float64 { return float64(int(f*10+0.5)) / 10 }

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

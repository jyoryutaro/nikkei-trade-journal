package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/application"
	"github.com/min-legomain/nikkei-trade-journal/backend/internal/domain/marketdata"
	"github.com/min-legomain/nikkei-trade-journal/backend/internal/infrastructure/db"
	"github.com/min-legomain/nikkei-trade-journal/backend/internal/infrastructure/persistence/mysql"
	"github.com/min-legomain/nikkei-trade-journal/backend/internal/infrastructure/yahoo"
)

func main() {
	dataPath := flag.String("data", "", "path to Yahoo Finance v8 JSON file")
	contractFlag := flag.String("contract", "", "contract code e.g. 2609 (auto-detected from shortName if empty)")
	timeframeFlag := flag.String("timeframe", "1m", "timeframe e.g. 1m, 5m, 1d")
	flag.Parse()

	if *dataPath == "" {
		log.Fatal("usage: seed -data <path-to-json> [-contract 2609] [-timeframe 1m]")
	}

	raw, err := os.ReadFile(*dataPath)
	if err != nil {
		log.Fatalf("read file: %v", err)
	}

	contract, candles, err := yahoo.ParseChart(raw)
	if err != nil {
		log.Fatalf("parse chart: %v", err)
	}
	if *contractFlag != "" {
		contract = *contractFlag
	}
	if contract == "" {
		log.Fatal("could not determine contract — pass -contract explicitly")
	}

	tf, err := marketdata.ParseTimeframe(*timeframeFlag)
	if err != nil {
		log.Fatalf("timeframe: %v", err)
	}

	dsn := getEnv("DB_DSN", "app:app@tcp(localhost:3306)/nikkei_trade?parseTime=true")
	database, err := db.Open(dsn)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer database.Close()

	repo := mysql.NewMarketDataRepository(database)
	svc := application.NewMarketDataService(repo)

	n, err := svc.Import(context.Background(), contract, tf, candles)
	if err != nil {
		log.Fatalf("import: %v", err)
	}
	log.Printf("seeded %d bars → contract=%s timeframe=%s", n, contract, tf)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

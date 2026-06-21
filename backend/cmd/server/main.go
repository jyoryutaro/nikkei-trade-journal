package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/application"
	"github.com/min-legomain/nikkei-trade-journal/backend/internal/infrastructure/db"
	"github.com/min-legomain/nikkei-trade-journal/backend/internal/infrastructure/persistence/mysql"
	"github.com/min-legomain/nikkei-trade-journal/backend/internal/infrastructure/yahoo"
	httpapi "github.com/min-legomain/nikkei-trade-journal/backend/internal/interfaces/http"
	"github.com/min-legomain/nikkei-trade-journal/backend/internal/scheduler"
)

func main() {
	dsn := getEnv("DB_DSN", "app:app@tcp(localhost:3306)/nikkei_trade?parseTime=true&loc=UTC")
	database, err := db.Open(dsn)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer database.Close()
	log.Println("DB connected")

	secret := getEnv("INTERNAL_SECRET", "")
	if secret == "" {
		log.Fatal("INTERNAL_SECRET env var is required")
	}

	marketRepo := mysql.NewMarketDataRepository(database)
	fetcher := yahoo.NewFetcher(&http.Client{})
	marketSvc := application.NewMarketDataService(fetcher, marketRepo)
	marketHandler := httpapi.NewMarketDataHandler(marketSvc)

	journalRepo := mysql.NewJournalRepository(database)
	journalSvc := application.NewJournalService(journalRepo)
	journalHandler := httpapi.NewJournalHandler(journalSvc)

	router := httpapi.NewRouter(marketHandler, journalHandler, secret)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	symbols := strings.Fields(getEnv("SYMBOLS", "^N225 NKD=F"))
	interval, err := time.ParseDuration(getEnv("SCHEDULER_INTERVAL", "15m"))
	if err != nil {
		log.Fatalf("invalid SCHEDULER_INTERVAL: %v", err)
	}
	sched := scheduler.New(marketSvc, symbols, interval)
	go sched.Run(ctx)

	addr := getEnv("ADDR", ":8080")
	srv := &http.Server{Addr: addr, Handler: router}
	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background()) //nolint:errcheck
	}()

	log.Printf("server listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

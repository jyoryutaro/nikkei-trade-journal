package main

import (
	"log"
	"net/http"
	"os"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/application"
	"github.com/min-legomain/nikkei-trade-journal/backend/internal/infrastructure/db"
	"github.com/min-legomain/nikkei-trade-journal/backend/internal/infrastructure/persistence/mysql"
	httpapi "github.com/min-legomain/nikkei-trade-journal/backend/internal/interfaces/http"
)

func main() {
	dsn := getEnv("DB_DSN", "app:app@tcp(localhost:3306)/nikkei_trade?parseTime=true&loc=Asia%2FTokyo")
	database, err := db.Open(dsn)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer database.Close()
	log.Println("DB connected")

	// Compose the layers: repository (infra) → service (application) → handler (interface).
	marketRepo := mysql.NewMarketDataRepository(database)
	marketSvc := application.NewMarketDataService(marketRepo)
	marketHandler := httpapi.NewMarketDataHandler(marketSvc)

	journalRepo := mysql.NewJournalRepository(database)
	journalSvc := application.NewJournalService(journalRepo)
	journalHandler := httpapi.NewJournalHandler(journalSvc)

	router := httpapi.NewRouter(marketHandler, journalHandler)

	addr := getEnv("ADDR", ":8080")
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

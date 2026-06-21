package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/application"
)

type Scheduler struct {
	svc      *application.MarketDataService
	symbols  []string
	interval time.Duration
}

func New(svc *application.MarketDataService, symbols []string, interval time.Duration) *Scheduler {
	return &Scheduler{svc: svc, symbols: symbols, interval: interval}
}

// Run fetches immediately on start, then repeats every interval until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	s.fetchAll(ctx)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.fetchAll(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Scheduler) fetchAll(ctx context.Context) {
	for _, sym := range s.symbols {
		n, err := s.svc.Fetch(ctx, sym)
		if err != nil {
			log.Printf("scheduler: fetch %s: %v", sym, err)
			continue
		}
		log.Printf("scheduler: fetch %s saved %d candles", sym, n)
	}
}

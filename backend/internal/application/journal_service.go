package application

import (
	"context"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/domain/journal"
)

// JournalService orchestrates journal-entry use cases.
type JournalService struct {
	repo journal.Repository
}

// NewJournalService wires the service with a repository implementation.
func NewJournalService(repo journal.Repository) *JournalService {
	return &JournalService{repo: repo}
}

// Create persists a validated entry. The entry is expected to have been built
// via journal.NewEntry (which enforces the domain invariants).
func (s *JournalService) Create(ctx context.Context, e journal.Entry) (journal.Entry, error) {
	return s.repo.Create(ctx, e)
}

// List returns the entries recorded for a contract.
func (s *JournalService) List(ctx context.Context, contract string) ([]journal.Entry, error) {
	return s.repo.ListByContract(ctx, contract)
}

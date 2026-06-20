package journal

import "context"

// Repository is the persistence boundary (port) for journal entries.
type Repository interface {
	// Create persists a new entry and returns it with ID and CreatedAt set.
	Create(ctx context.Context, e Entry) (Entry, error)

	// ListByContract returns entries for a contract, sorted by time ascending.
	ListByContract(ctx context.Context, contract string) ([]Entry, error)
}

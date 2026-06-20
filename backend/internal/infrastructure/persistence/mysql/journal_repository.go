package mysql

import (
	"context"
	"database/sql"
	"time"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/domain/journal"
)

// JournalRepository is a MySQL-backed journal.Repository.
type JournalRepository struct {
	db *sql.DB
}

var _ journal.Repository = (*JournalRepository)(nil)

// NewJournalRepository constructs the repository.
func NewJournalRepository(db *sql.DB) *JournalRepository {
	return &JournalRepository{db: db}
}

// Create inserts a new journal entry and returns it with ID/CreatedAt set.
func (r *JournalRepository) Create(ctx context.Context, e journal.Entry) (journal.Entry, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO journal_entries (contract, ts, side, trade_type, price, comment)
		VALUES (?, ?, ?, ?, ?, ?)
	`, e.Contract, e.Time.UTC(), nullString(string(e.Side)), nullString(string(e.TradeType)), nullFloat(e.Price), nullString(e.Comment))
	if err != nil {
		return journal.Entry{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return journal.Entry{}, err
	}
	e.ID = id
	e.CreatedAt = time.Now().UTC()
	return e, nil
}

// ListByContract returns entries for a contract, oldest first.
func (r *JournalRepository) ListByContract(ctx context.Context, contract string) ([]journal.Entry, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, contract, ts, side, trade_type, price, comment, created_at
		FROM journal_entries
		WHERE contract = ?
		ORDER BY ts ASC, id ASC
	`, contract)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []journal.Entry
	for rows.Next() {
		var (
			e         journal.Entry
			side      sql.NullString
			tradeType sql.NullString
			price     sql.NullFloat64
			comment   sql.NullString
		)
		if err := rows.Scan(&e.ID, &e.Contract, &e.Time, &side, &tradeType, &price, &comment, &e.CreatedAt); err != nil {
			return nil, err
		}
		e.Side = journal.Side(side.String)
		e.TradeType = journal.TradeType(tradeType.String)
		if price.Valid {
			p := price.Float64
			e.Price = &p
		}
		e.Comment = comment.String
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

// nullString returns a NULL for empty strings, otherwise the value.
func nullString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// nullFloat returns a NULL for nil pointers, otherwise the value.
func nullFloat(f *float64) any {
	if f == nil {
		return nil
	}
	return *f
}

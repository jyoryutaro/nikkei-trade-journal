// Package journal models trade-journal entries: a record attached to a point
// in time that is either a position record or a comment-only note.
package journal

import (
	"errors"
	"strings"
	"time"
)

// Side is the trade direction.
type Side string

const (
	SideNone  Side = "" // comment-only entry
	SideLong  Side = "long"
	SideShort Side = "short"
)

// TradeType distinguishes opening a position from closing one.
type TradeType string

const (
	TradeTypeNone  TradeType = ""
	TradeTypeOpen  TradeType = "open"  // 新規
	TradeTypeClose TradeType = "close" // 決済
)

// Domain validation errors.
var (
	ErrContractRequired  = errors.New("contract is required")
	ErrInvalidSide       = errors.New("invalid side")
	ErrInvalidTradeType  = errors.New("invalid trade type")
	ErrCommentRequired   = errors.New("comment is required for a comment-only entry")
	ErrTradeTypeRequired = errors.New("trade type is required when a side is selected")
	ErrPriceRequired     = errors.New("price is required when a side is selected")
	ErrPricePositive     = errors.New("price must be positive")
	ErrPositionFields    = errors.New("trade type and price are only allowed when a side is selected")
)

// Entry is a journal entry. It is either:
//   - a position record: Side is long/short, with a TradeType and Price; or
//   - a comment-only note: Side is none, with a non-empty Comment.
type Entry struct {
	ID        int64
	Contract  string
	Time      time.Time
	Side      Side
	TradeType TradeType
	Price     *float64
	Comment   string
	CreatedAt time.Time
}

// ParseSide validates and converts a raw side string.
func ParseSide(s string) (Side, error) {
	switch Side(s) {
	case SideNone, SideLong, SideShort:
		return Side(s), nil
	}
	return SideNone, ErrInvalidSide
}

// ParseTradeType validates and converts a raw trade-type string.
func ParseTradeType(s string) (TradeType, error) {
	switch TradeType(s) {
	case TradeTypeNone, TradeTypeOpen, TradeTypeClose:
		return TradeType(s), nil
	}
	return TradeTypeNone, ErrInvalidTradeType
}

// NewEntry constructs a validated Entry, enforcing the domain invariants:
//   - position entry (side long/short): trade type and a positive price are required
//   - comment-only entry (side none): a non-empty comment is required, and
//     position fields must be absent
func NewEntry(contract string, t time.Time, side Side, tt TradeType, price *float64, comment string) (Entry, error) {
	if strings.TrimSpace(contract) == "" {
		return Entry{}, ErrContractRequired
	}
	comment = strings.TrimSpace(comment)

	if side == SideNone {
		if tt != TradeTypeNone || price != nil {
			return Entry{}, ErrPositionFields
		}
		if comment == "" {
			return Entry{}, ErrCommentRequired
		}
	} else {
		if tt == TradeTypeNone {
			return Entry{}, ErrTradeTypeRequired
		}
		if price == nil {
			return Entry{}, ErrPriceRequired
		}
		if *price <= 0 {
			return Entry{}, ErrPricePositive
		}
	}

	return Entry{
		Contract:  contract,
		Time:      t.UTC(),
		Side:      side,
		TradeType: tt,
		Price:     price,
		Comment:   comment,
	}, nil
}

// IsValidationError reports whether err is a domain validation error (i.e. a
// client mistake, mapped to HTTP 400 by adapters).
func IsValidationError(err error) bool {
	switch {
	case errors.Is(err, ErrContractRequired),
		errors.Is(err, ErrInvalidSide),
		errors.Is(err, ErrInvalidTradeType),
		errors.Is(err, ErrCommentRequired),
		errors.Is(err, ErrTradeTypeRequired),
		errors.Is(err, ErrPriceRequired),
		errors.Is(err, ErrPricePositive),
		errors.Is(err, ErrPositionFields):
		return true
	}
	return false
}

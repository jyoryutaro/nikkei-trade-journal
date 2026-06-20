// Package yahoo parses Yahoo Finance v8 chart JSON into domain candles.
package yahoo

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/domain/marketdata"
)

// chartResponse is the minimal subset of the Yahoo Finance v8 chart API.
type chartResponse struct {
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

// ParseChart decodes Yahoo Finance v8 chart JSON, returning the detected
// contract code (from the meta shortName) and the candles. Timestamps with
// missing OHLC values (Yahoo returns null for minutes with no trade) are
// forward-filled with the previous close (a flat bar with volume 0) so the
// 1-minute series stays continuous. A leading gap before the first real value
// is dropped.
func ParseChart(raw []byte) (contract string, candles []marketdata.Candle, err error) {
	var resp chartResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "", nil, fmt.Errorf("parse json: %w", err)
	}
	if len(resp.Chart.Result) == 0 {
		return "", nil, fmt.Errorf("no results in JSON")
	}
	result := resp.Chart.Result[0]
	if len(result.Indicators.Quote) == 0 {
		return "", nil, fmt.Errorf("no quote data")
	}
	q := result.Indicators.Quote[0]

	contract = contractFromShortName(result.Meta.ShortName)

	var lastClose float64
	var havePrev bool
	for i, ts := range result.Timestamps {
		t := time.Unix(ts, 0).UTC()

		if i < len(q.Open) && q.Open[i] != nil && q.High[i] != nil && q.Low[i] != nil && q.Close[i] != nil {
			var vol int64
			if i < len(q.Volume) && q.Volume[i] != nil {
				vol = *q.Volume[i]
			}
			candles = append(candles, marketdata.NewCandle(t, *q.Open[i], *q.High[i], *q.Low[i], *q.Close[i], vol))
			lastClose = *q.Close[i]
			havePrev = true
			continue
		}

		// missing value → forward-fill with the previous close (flat bar, vol 0)
		if havePrev {
			candles = append(candles, marketdata.NewCandle(t, lastClose, lastClose, lastClose, lastClose, 0))
		}
	}
	return contract, candles, nil
}

// contractFromShortName extracts a YYMM code from strings like
// "Nikkei/USD Futures,Sep-2026" → "2609". Returns "" if not found.
func contractFromShortName(name string) string {
	monthMap := map[string]string{
		"Jan": "01", "Feb": "02", "Mar": "03", "Apr": "04",
		"May": "05", "Jun": "06", "Jul": "07", "Aug": "08",
		"Sep": "09", "Oct": "10", "Nov": "11", "Dec": "12",
	}
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

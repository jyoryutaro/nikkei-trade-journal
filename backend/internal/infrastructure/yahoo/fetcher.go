package yahoo

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/min-legomain/nikkei-trade-journal/backend/internal/domain/marketdata"
)

const yahooChartBaseURL = "https://query1.finance.yahoo.com/v8/finance/chart"

// Fetcher fetches base-timeframe candles from the Yahoo Finance v8 chart API.
// It implements marketdata.CandleSource.
type Fetcher struct {
	client  *http.Client
	baseURL string
}

// NewFetcher constructs a Fetcher pointing at the live Yahoo Finance API.
func NewFetcher(client *http.Client) *Fetcher {
	return &Fetcher{client: client, baseURL: yahooChartBaseURL}
}

// NewFetcherAt constructs a Fetcher pointing at baseURL.
// Intended for testing with a mock HTTP server.
func NewFetcherAt(client *http.Client, baseURL string) *Fetcher {
	return &Fetcher{client: client, baseURL: baseURL}
}

// FetchCandles retrieves 1-minute candles for symbol from Yahoo Finance and
// returns the contract code extracted from the response metadata.
// If no contract code is found in the response, symbol is used as fallback.
func (f *Fetcher) FetchCandles(ctx context.Context, symbol string) (string, []marketdata.Candle, error) {
	url := fmt.Sprintf("%s/%s?interval=1m&range=5d", f.baseURL, symbol)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := f.client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("fetch %s: %w", symbol, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("yahoo finance: status %d", resp.StatusCode)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("read body: %w", err)
	}

	contract, candles, err := ParseChart(raw)
	if err != nil {
		return "", nil, fmt.Errorf("parse chart: %w", err)
	}
	if contract == "" {
		contract = symbol
	}
	return contract, candles, nil
}

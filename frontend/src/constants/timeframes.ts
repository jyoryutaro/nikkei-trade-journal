export const TIMEFRAMES = [
  { value: '1m',  label: '1分' },
  { value: '5m',  label: '5分' },
  { value: '30m', label: '30分' },
  { value: '1h',  label: '1時間' },
  { value: '1d',  label: '日足' },
]

// Interval in seconds per timeframe. Mirrors the backend aggregator so entry
// times can be bucketed to the matching candle on each timeframe.
export const INTERVAL_SECONDS: Record<string, number> = {
  '1m': 60,
  '5m': 300,
  '30m': 1800,
  '1h': 3600,
  '1d': 86400,
}

export const WINDOW_OPTIONS: Record<string, { label: string; hours: number | null }[]> = {
  '1m':  [
    { label: '30分',  hours: 0.5 },
    { label: '1時間', hours: 1 },
    { label: '4時間', hours: 4 },
    { label: '全体',  hours: null },
  ],
  '5m':  [
    { label: '1時間', hours: 1 },
    { label: '4時間', hours: 4 },
    { label: '1日',   hours: 24 },
    { label: '全体',  hours: null },
  ],
  '30m': [
    { label: '4時間', hours: 4 },
    { label: '1日',   hours: 24 },
    { label: '全体',  hours: null },
  ],
  '1h':  [
    { label: '1日',   hours: 24 },
    { label: '全体',  hours: null },
  ],
  '1d':  [
    { label: '全体',  hours: null },
  ],
}

export function defaultWindow(tf: string): number | null {
  return WINDOW_OPTIONS[tf]?.[0]?.hours ?? null
}

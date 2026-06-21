import { API_BASE } from '../config'

export interface Candle {
  contract: string
  timeframe: string
  time: number  // Unix timestamp seconds (UTC)
  open: number
  high: number
  low: number
  close: number
  volume: number
}

/** Fetch the distinct contract codes available in the DB. */
export async function fetchContracts(): Promise<string[]> {
  const res = await fetch(`${API_BASE}/api/contracts`)
  if (!res.ok) throw new Error(`API error: ${res.status}`)
  const data: string[] | null = await res.json()
  return data ?? []
}

export async function fetchMarketData(contract: string, timeframe: string): Promise<Candle[]> {
  const params = new URLSearchParams({ contract, timeframe })
  const res = await fetch(`${API_BASE}/api/market-data?${params}`)
  if (!res.ok) throw new Error(`API error: ${res.status}`)
  const data: Candle[] | null = await res.json()
  return data ?? []
}

/** Format a Unix timestamp as "YYYY-MM-DD HH:mm JST" */
export function formatTimeJST(unix: number): string {
  const d = new Date((unix + 9 * 3600) * 1000)
  const s = d.toISOString()
  return `${s.slice(0, 10)} ${s.slice(11, 16)} JST`
}

/** Format a Unix timestamp as a JST `datetime-local` value "YYYY-MM-DDTHH:mm". */
export function toJstInputValue(unix: number): string {
  const d = new Date((unix + 9 * 3600) * 1000)
  const p = (n: number) => String(n).padStart(2, '0')
  return `${d.getUTCFullYear()}-${p(d.getUTCMonth() + 1)}-${p(d.getUTCDate())}T${p(d.getUTCHours())}:${p(d.getUTCMinutes())}`
}

/** Parse a JST `datetime-local` value "YYYY-MM-DDTHH:mm" to a Unix timestamp
 * (seconds, UTC). The input is interpreted as JST regardless of browser TZ. */
export function fromJstInputValue(value: string): number {
  const [date, time] = value.split('T')
  const [y, m, d] = date.split('-').map(Number)
  const [hh, mm] = (time ?? '00:00').split(':').map(Number)
  return Math.floor(Date.UTC(y, m - 1, d, hh, mm) / 1000) - 9 * 3600
}

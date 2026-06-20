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

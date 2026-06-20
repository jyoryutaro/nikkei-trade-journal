const API_BASE = 'http://localhost:8080'

export interface Candle {
  contract: string
  timeframe: string
  time: string  // ISO 8601
  open: number
  high: number
  low: number
  close: number
  volume: number
}

export async function fetchMarketData(contract: string, timeframe: string): Promise<Candle[]> {
  const params = new URLSearchParams({ contract, timeframe })
  const res = await fetch(`${API_BASE}/api/market-data?${params}`)
  if (!res.ok) throw new Error(`API error: ${res.status}`)
  const data: Candle[] | null = await res.json()
  return data ?? []
}

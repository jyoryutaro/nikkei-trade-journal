const API_BASE = 'http://localhost:8080'

export type Side = '' | 'long' | 'short'
export type TradeType = '' | 'open' | 'close'

export interface JournalEntryInput {
  contract: string
  time: number // Unix seconds (UTC)
  side: Side
  tradeType: TradeType
  price: number | null
  comment: string
}

export interface JournalEntry {
  id: number
  contract: string
  time: number // Unix seconds (UTC)
  side: Side
  tradeType: TradeType
  price: number | null
  comment: string
  createdAt: number
}

/** Create a journal entry (position record or comment-only note). */
export async function createJournalEntry(input: JournalEntryInput): Promise<void> {
  const res = await fetch(`${API_BASE}/api/journal-entries`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text.trim() || `API error: ${res.status}`)
  }
}

/** Fetch the journal entries recorded for a contract. */
export async function fetchJournalEntries(contract: string): Promise<JournalEntry[]> {
  const params = new URLSearchParams({ contract })
  const res = await fetch(`${API_BASE}/api/journal-entries?${params}`)
  if (!res.ok) throw new Error(`API error: ${res.status}`)
  const data: JournalEntry[] | null = await res.json()
  return data ?? []
}

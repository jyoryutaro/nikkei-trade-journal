import { useCallback, useEffect, useState } from 'react'
import { fetchMarketData, fetchContracts, type Candle } from '../../api/marketData'
import { fetchJournalEntries, type JournalEntry } from '../../api/journal'
import { defaultWindow } from '../../constants/timeframes'
import { CONTRACTS } from '../../constants/contracts'
import type { Option } from '../atoms/Select'
import { AppHeader } from '../organisms/AppHeader'
import { ChartToolbar } from '../organisms/ChartToolbar'
import { CandlestickChart } from '../organisms/CandlestickChart'
import { CommentPanel } from '../organisms/CommentPanel'
import { DashboardTemplate } from '../templates/DashboardTemplate'
import { colors } from '../../theme'

function contractLabel(c: string): string {
  return c.startsWith('^') ? `${c} (指数)` : `${c} (先物)`
}

function EmptyChart() {
  return (
    <div style={{ height: '420px', display: 'flex', alignItems: 'center', justifyContent: 'center', color: colors.textGhost, border: `1px solid ${colors.border}`, borderRadius: '8px' }}>
      <span style={{ fontSize: '0.9rem' }}>
        データがありません — 取引時間外の可能性があります
      </span>
    </div>
  )
}

/** The dashboard page: owns state and data fetching, composes organisms into
 * the dashboard template. */
export function DashboardPage() {
  const [contract, setContract] = useState('^N225')
  const [contracts, setContracts] = useState<Option[]>(CONTRACTS)
  const [timeframe, setTimeframe] = useState('1m')
  const [windowHours, setWindowHours] = useState<number | null>(defaultWindow('1m'))
  const [candles, setCandles] = useState<Candle[]>([])
  const [entries, setEntries] = useState<JournalEntry[]>([])
  const [selectedCandle, setSelectedCandle] = useState<Candle | null>(null)
  const [error, setError] = useState<string | null>(null)

  const handleTimeframeChange = (tf: string) => {
    setTimeframe(tf)
    setWindowHours(defaultWindow(tf))
    setSelectedCandle(null)
  }

  const loadEntries = useCallback(() => {
    fetchJournalEntries(contract)
      .then(setEntries)
      .catch(() => setEntries([]))
  }, [contract])

  useEffect(() => {
    setSelectedCandle(null)
    fetchMarketData(contract, timeframe)
      .then(data => {
        setCandles(data)
        setError(null)
      })
      .catch(() => setError('APIに接続できません'))
  }, [contract, timeframe])

  useEffect(() => {
    loadEntries()
  }, [loadEntries])

  // populate the symbol selector from what's actually stored in the DB
  useEffect(() => {
    fetchContracts()
      .then(list => {
        if (list.length === 0) return
        setContracts(list.map(c => ({ value: c, label: contractLabel(c) })))
        setContract(prev => (list.includes(prev) ? prev : list[0]))
      })
      .catch(() => {})
  }, [])

  return (
    <DashboardTemplate
      header={<AppHeader contract={contract} contracts={contracts} onContractChange={setContract} />}
      error={error ? <p style={{ color: colors.down, fontSize: '0.85rem', marginBottom: '12px' }}>{error}</p> : null}
      toolbar={
        <ChartToolbar
          timeframe={timeframe}
          windowHours={windowHours}
          onTimeframeChange={handleTimeframeChange}
          onWindowChange={setWindowHours}
        />
      }
      chart={
        candles.length === 0 && !error ? (
          <EmptyChart />
        ) : (
          <CandlestickChart candles={candles} entries={entries} timeframe={timeframe} windowHours={windowHours} onSelect={setSelectedCandle} />
        )
      }
      aside={<CommentPanel contract={contract} timeframe={timeframe} selectedCandle={selectedCandle} onSubmitted={loadEntries} />}
    />
  )
}

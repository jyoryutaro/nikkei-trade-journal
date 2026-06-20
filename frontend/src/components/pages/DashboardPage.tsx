import { useEffect, useState } from 'react'
import { fetchMarketData, type Candle } from '../../api/marketData'
import { defaultWindow } from '../../constants/timeframes'
import { AppHeader } from '../organisms/AppHeader'
import { ChartToolbar } from '../organisms/ChartToolbar'
import { CandlestickChart } from '../organisms/CandlestickChart'
import { CommentPanel } from '../organisms/CommentPanel'
import { PriceTable } from '../organisms/PriceTable'
import { DashboardTemplate } from '../templates/DashboardTemplate'
import { colors } from '../../theme'

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
  const [timeframe, setTimeframe] = useState('1m')
  const [windowHours, setWindowHours] = useState<number | null>(defaultWindow('1m'))
  const [candles, setCandles] = useState<Candle[]>([])
  const [selectedCandle, setSelectedCandle] = useState<Candle | null>(null)
  const [error, setError] = useState<string | null>(null)

  const handleTimeframeChange = (tf: string) => {
    setTimeframe(tf)
    setWindowHours(defaultWindow(tf))
    setSelectedCandle(null)
  }

  useEffect(() => {
    setSelectedCandle(null)
    fetchMarketData(contract, timeframe)
      .then(data => {
        setCandles(data)
        setError(null)
      })
      .catch(() => setError('APIに接続できません'))
  }, [contract, timeframe])

  return (
    <DashboardTemplate
      header={<AppHeader contract={contract} onContractChange={setContract} />}
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
          <CandlestickChart candles={candles} windowHours={windowHours} onSelect={setSelectedCandle} />
        )
      }
      aside={<CommentPanel contract={contract} selectedCandle={selectedCandle} />}
      table={<PriceTable candles={candles} />}
    />
  )
}

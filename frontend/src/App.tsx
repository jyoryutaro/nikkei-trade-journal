import { useEffect, useState } from 'react'
import { fetchMarketData, type Candle } from './api/marketData'
import { defaultWindow } from './constants/timeframes'
import { ContractSelector } from './components/ContractSelector'
import { ChartToolbar } from './components/ChartToolbar'
import { CandlestickChart } from './components/CandlestickChart'
import { CommentPanel } from './components/CommentPanel'
import { PriceTable } from './components/PriceTable'

export default function App() {
  const [contract, setContract] = useState('2609')
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
      .then(data => { setCandles(data); setError(null) })
      .catch(() => setError('APIに接続できません'))
  }, [contract, timeframe])

  return (
    <div style={{ background: '#0a1628', minHeight: '100vh', color: '#e2e8f0', fontFamily: 'sans-serif', padding: '20px 24px' }}>

      <div style={{ display: 'flex', alignItems: 'center', gap: '16px', marginBottom: '20px', borderBottom: '1px solid #1e293b', paddingBottom: '16px' }}>
        <h1 style={{ fontSize: '1rem', fontWeight: 600, color: '#94a3b8', margin: 0 }}>
          日経225先物 トレードジャーナル
        </h1>
        <div style={{ width: '1px', height: '16px', background: '#334155' }} />
        <ContractSelector contract={contract} onContractChange={setContract} />
      </div>

      {error && <p style={{ color: '#ef4444', fontSize: '0.85rem', marginBottom: '12px' }}>{error}</p>}

      <ChartToolbar
        timeframe={timeframe}
        windowHours={windowHours}
        onTimeframeChange={handleTimeframeChange}
        onWindowChange={setWindowHours}
      />

      <div style={{ display: 'flex', gap: '16px', alignItems: 'flex-start' }}>
        <div style={{ flex: 1, minWidth: 0 }}>
          {candles.length === 0 && !error ? (
            <div style={{ height: '420px', display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#334155', border: '1px solid #1e293b', borderRadius: '8px' }}>
              <span style={{ fontSize: '0.9rem' }}>データがありません — <code>make seed</code> で投入してください</span>
            </div>
          ) : (
            <CandlestickChart
              candles={candles}
              windowHours={windowHours}
              onSelect={setSelectedCandle}
            />
          )}
        </div>

        <CommentPanel selectedCandle={selectedCandle} />
      </div>

      <PriceTable candles={candles} />
    </div>
  )
}

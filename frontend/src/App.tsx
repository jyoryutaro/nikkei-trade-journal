import { useEffect, useRef, useState } from 'react'
import { createChart, CandlestickSeries } from 'lightweight-charts'

interface Candle {
  contract: string
  time: string
  open: number
  high: number
  low: number
  close: number
  volume: number
}

const API = 'http://localhost:8080'

function App() {
  const chartRef = useRef<HTMLDivElement>(null)
  const [candles, setCandles] = useState<Candle[]>([])
  const [contract, setContract] = useState('2506')
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetch(`${API}/api/market-data?contract=${contract}`)
      .then(r => r.json())
      .then((data: Candle[]) => {
        setCandles(data ?? [])
        setError(null)
      })
      .catch(() => setError('APIに接続できません（バックエンドが起動しているか確認してください）'))
  }, [contract])

  useEffect(() => {
    if (!chartRef.current || candles.length === 0) return

    const chart = createChart(chartRef.current, {
      width: chartRef.current.clientWidth,
      height: 400,
      layout: { background: { color: '#0f172a' }, textColor: '#e2e8f0' },
      grid: { vertLines: { color: '#1e293b' }, horzLines: { color: '#1e293b' } },
      timeScale: { timeVisible: true },
    })

    const series = chart.addSeries(CandlestickSeries, {
      upColor: '#22c55e',
      downColor: '#ef4444',
      borderVisible: false,
      wickUpColor: '#22c55e',
      wickDownColor: '#ef4444',
    })

    series.setData(
      candles.map(c => ({
        time: c.time.slice(0, 10),
        open: c.open,
        high: c.high,
        low: c.low,
        close: c.close,
      }))
    )

    chart.timeScale().fitContent()

    const handleResize = () => chart.applyOptions({ width: chartRef.current!.clientWidth })
    window.addEventListener('resize', handleResize)
    return () => {
      window.removeEventListener('resize', handleResize)
      chart.remove()
    }
  }, [candles])

  return (
    <div style={{ background: '#0f172a', minHeight: '100vh', color: '#e2e8f0', fontFamily: 'sans-serif', padding: '24px' }}>
      <h1 style={{ fontSize: '1.5rem', marginBottom: '16px' }}>日経225先物 トレードジャーナル</h1>

      <div style={{ marginBottom: '16px' }}>
        <label style={{ marginRight: '8px' }}>限月:</label>
        <input
          value={contract}
          onChange={e => setContract(e.target.value)}
          style={{ background: '#1e293b', color: '#e2e8f0', border: '1px solid #334155', borderRadius: '4px', padding: '4px 8px' }}
        />
      </div>

      {error && <p style={{ color: '#ef4444' }}>{error}</p>}

      <div ref={chartRef} style={{ width: '100%' }} />

      {candles.length > 0 && (
        <table style={{ marginTop: '24px', borderCollapse: 'collapse', fontSize: '0.85rem', width: '100%' }}>
          <thead>
            <tr style={{ borderBottom: '1px solid #334155' }}>
              {['日時', '始値', '高値', '安値', '終値', '出来高'].map(h => (
                <th key={h} style={{ padding: '6px 12px', textAlign: 'right' }}>{h}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {candles.slice(-20).reverse().map(c => (
              <tr key={c.time} style={{ borderBottom: '1px solid #1e293b' }}>
                <td style={{ padding: '4px 12px', color: '#94a3b8' }}>{c.time.slice(0, 10)}</td>
                <td style={{ padding: '4px 12px', textAlign: 'right' }}>{c.open.toLocaleString()}</td>
                <td style={{ padding: '4px 12px', textAlign: 'right', color: '#22c55e' }}>{c.high.toLocaleString()}</td>
                <td style={{ padding: '4px 12px', textAlign: 'right', color: '#ef4444' }}>{c.low.toLocaleString()}</td>
                <td style={{ padding: '4px 12px', textAlign: 'right', fontWeight: 'bold' }}>{c.close.toLocaleString()}</td>
                <td style={{ padding: '4px 12px', textAlign: 'right', color: '#94a3b8' }}>{c.volume.toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {candles.length === 0 && !error && (
        <p style={{ color: '#94a3b8' }}>データがありません。`go run ./cmd/seed` でテストデータを投入してください。</p>
      )}
    </div>
  )
}

export default App

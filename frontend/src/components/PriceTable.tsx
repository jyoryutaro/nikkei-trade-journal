import type { Candle } from '../api/marketData'
import { formatTimeJST } from '../api/marketData'

interface Props {
  candles: Candle[]
  maxRows?: number
}

const th: React.CSSProperties = { padding: '6px 12px', textAlign: 'right', color: '#94a3b8', fontWeight: 500, borderBottom: '1px solid #334155' }
const td: React.CSSProperties = { padding: '4px 12px', textAlign: 'right', borderBottom: '1px solid #1e293b' }

export function PriceTable({ candles, maxRows = 20 }: Props) {
  if (candles.length === 0) return null

  const rows = candles.slice(-maxRows).reverse()

  return (
    <div style={{ marginTop: '24px', overflowX: 'auto' }}>
      <table style={{ borderCollapse: 'collapse', fontSize: '0.82rem', width: '100%' }}>
        <thead>
          <tr>
            <th style={{ ...th, textAlign: 'left' }}>日時 (JST)</th>
            <th style={th}>始値</th>
            <th style={th}>高値</th>
            <th style={th}>安値</th>
            <th style={th}>終値</th>
            <th style={th}>出来高</th>
          </tr>
        </thead>
        <tbody>
          {rows.map(c => {
            const up = c.close >= c.open
            return (
              <tr key={c.time}>
                <td style={{ ...td, textAlign: 'left', color: '#94a3b8' }}>{formatTimeJST(c.time)}</td>
                <td style={td}>{c.open.toLocaleString()}</td>
                <td style={{ ...td, color: '#22c55e' }}>{c.high.toLocaleString()}</td>
                <td style={{ ...td, color: '#ef4444' }}>{c.low.toLocaleString()}</td>
                <td style={{ ...td, fontWeight: 'bold', color: up ? '#22c55e' : '#ef4444' }}>{c.close.toLocaleString()}</td>
                <td style={{ ...td, color: '#64748b' }}>{c.volume.toLocaleString()}</td>
              </tr>
            )
          })}
        </tbody>
      </table>
    </div>
  )
}

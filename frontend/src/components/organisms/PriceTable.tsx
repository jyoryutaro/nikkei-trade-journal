import type { CSSProperties } from 'react'
import type { Candle } from '../../api/marketData'
import { formatTimeJST } from '../../api/marketData'
import { colors } from '../../theme'

interface Props {
  candles: Candle[]
  maxRows?: number
}

const th: CSSProperties = { padding: '6px 12px', textAlign: 'right', color: colors.textMuted, fontWeight: 500, borderBottom: `1px solid ${colors.borderStrong}` }
const td: CSSProperties = { padding: '4px 12px', textAlign: 'right', borderBottom: `1px solid ${colors.border}` }

/** Recent candles as a table (most recent first). */
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
                <td style={{ ...td, textAlign: 'left', color: colors.textMuted }}>{formatTimeJST(c.time)}</td>
                <td style={td}>{c.open.toLocaleString()}</td>
                <td style={{ ...td, color: colors.up }}>{c.high.toLocaleString()}</td>
                <td style={{ ...td, color: colors.down }}>{c.low.toLocaleString()}</td>
                <td style={{ ...td, fontWeight: 'bold', color: up ? colors.up : colors.down }}>{c.close.toLocaleString()}</td>
                <td style={{ ...td, color: colors.textFaint }}>{c.volume.toLocaleString()}</td>
              </tr>
            )
          })}
        </tbody>
      </table>
    </div>
  )
}

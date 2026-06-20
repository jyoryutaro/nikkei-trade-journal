import type { Candle } from '../../api/marketData'
import { formatTimeJST } from '../../api/marketData'
import { colors } from '../../theme'

interface Props {
  candle: Candle
}

/** Compact OHLCV readout for a single candle. */
export function OhlcvSummary({ candle }: Props) {
  return (
    <div style={{ background: colors.surface, borderRadius: '6px', padding: '10px 12px', fontSize: '0.82rem', lineHeight: 1.8 }}>
      <p style={{ color: colors.textFaint, marginBottom: '4px' }}>{formatTimeJST(candle.time)}</p>
      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '2px 12px' }}>
        <span style={{ color: colors.textMuted }}>始値</span>
        <span style={{ textAlign: 'right' }}>{candle.open.toLocaleString()}</span>
        <span style={{ color: colors.up }}>高値</span>
        <span style={{ textAlign: 'right', color: colors.up }}>{candle.high.toLocaleString()}</span>
        <span style={{ color: colors.down }}>安値</span>
        <span style={{ textAlign: 'right', color: colors.down }}>{candle.low.toLocaleString()}</span>
        <span style={{ color: colors.textMuted }}>終値</span>
        <span style={{ textAlign: 'right', fontWeight: 600 }}>{candle.close.toLocaleString()}</span>
      </div>
    </div>
  )
}

import type { CSSProperties } from 'react'
import type { Candle } from '../../api/marketData'
import { OhlcvSummary } from '../molecules/OhlcvSummary'
import { TradeEntryForm } from './TradeEntryForm'
import { colors } from '../../theme'

interface Props {
  contract: string
  timeframe: string
  selectedCandle: Candle | null
  onSubmitted?: () => void
}

const panelStyle: CSSProperties = {
  width: '280px',
  flexShrink: 0,
  background: colors.panel,
  border: `1px solid ${colors.border}`,
  borderRadius: '8px',
  padding: '16px',
  display: 'flex',
  flexDirection: 'column',
  gap: '12px',
}

const labelStyle: CSSProperties = {
  fontSize: '0.75rem',
  color: colors.textFaint,
  textTransform: 'uppercase',
  letterSpacing: '0.05em',
}

/** Side panel: shows the selected candle and the entry form (position or
 * comment-only). */
export function CommentPanel({ contract, timeframe, selectedCandle, onSubmitted }: Props) {
  return (
    <div style={panelStyle}>
      <p style={{ ...labelStyle, marginBottom: '4px' }}>記録</p>

      {!selectedCandle ? (
        <div style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: '8px', color: colors.textGhost, textAlign: 'center', padding: '24px 0' }}>
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
            <path d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" strokeLinecap="round" strokeLinejoin="round" />
          </svg>
          <span style={{ fontSize: '0.82rem', lineHeight: 1.5 }}>
            チャートをクリックして<br />時点を選択してください
          </span>
        </div>
      ) : (
        <>
          <OhlcvSummary candle={selectedCandle} />
          {/* key remounts the form (resetting fields/default time) when the
              selected candle or timeframe changes */}
          <TradeEntryForm
            key={`${selectedCandle.time}-${timeframe}`}
            contract={contract}
            candle={selectedCandle}
            timeframe={timeframe}
            onSubmitted={onSubmitted}
          />
        </>
      )}
    </div>
  )
}

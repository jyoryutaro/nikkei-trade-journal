import type { CSSProperties } from 'react'
import type { Candle } from '../../api/marketData'
import { OhlcvSummary } from '../molecules/OhlcvSummary'
import { colors } from '../../theme'

interface Props {
  selectedCandle: Candle | null
}

const panelStyle: CSSProperties = {
  width: '260px',
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

/** Side panel showing the selected candle and a (WIP) comment editor. */
export function CommentPanel({ selectedCandle }: Props) {
  return (
    <div style={panelStyle}>
      <p style={{ ...labelStyle, marginBottom: '4px' }}>コメント</p>

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

          <textarea
            placeholder="この時点にコメントを追加..."
            rows={5}
            style={{
              background: colors.surface,
              color: colors.text,
              border: `1px solid ${colors.borderStrong}`,
              borderRadius: '6px',
              padding: '8px 10px',
              fontSize: '0.82rem',
              resize: 'vertical',
              fontFamily: 'inherit',
              lineHeight: 1.6,
            }}
          />

          <button
            disabled
            title="コメント保存機能は準備中です"
            style={{
              padding: '8px',
              borderRadius: '6px',
              border: 'none',
              background: colors.accentBlueDeep,
              color: '#93c5fd',
              fontSize: '0.82rem',
              cursor: 'not-allowed',
              opacity: 0.5,
            }}
          >
            保存（準備中）
          </button>
        </>
      )}
    </div>
  )
}

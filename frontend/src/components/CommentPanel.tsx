import type { Candle } from '../api/marketData'

interface Props {
  selectedCandle: Candle | null
}

function formatTime(iso: string) {
  const d = new Date(iso)
  const jst = new Date(d.getTime() + 9 * 60 * 60 * 1000)
  return jst.toISOString().slice(0, 16).replace('T', ' ') + ' JST'
}

const panelStyle: React.CSSProperties = {
  width: '260px',
  flexShrink: 0,
  background: '#0f1f35',
  border: '1px solid #1e293b',
  borderRadius: '8px',
  padding: '16px',
  display: 'flex',
  flexDirection: 'column',
  gap: '12px',
}

const labelStyle: React.CSSProperties = {
  fontSize: '0.75rem',
  color: '#64748b',
  textTransform: 'uppercase',
  letterSpacing: '0.05em',
}

export function CommentPanel({ selectedCandle }: Props) {
  return (
    <div style={panelStyle}>
      <p style={{ ...labelStyle, marginBottom: '4px' }}>コメント</p>

      {!selectedCandle ? (
        <div style={{ flex: 1, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', gap: '8px', color: '#334155', textAlign: 'center', padding: '24px 0' }}>
          <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
            <path d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" strokeLinecap="round" strokeLinejoin="round"/>
          </svg>
          <span style={{ fontSize: '0.82rem', lineHeight: 1.5 }}>
            チャートをクリックして<br />時点を選択してください
          </span>
        </div>
      ) : (
        <>
          <div style={{ background: '#1e293b', borderRadius: '6px', padding: '10px 12px', fontSize: '0.82rem', lineHeight: 1.8 }}>
            <p style={{ color: '#64748b', marginBottom: '4px' }}>{formatTime(selectedCandle.time)}</p>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '2px 12px' }}>
              <span style={{ color: '#94a3b8' }}>始値</span>
              <span style={{ textAlign: 'right' }}>{selectedCandle.open.toLocaleString()}</span>
              <span style={{ color: '#22c55e' }}>高値</span>
              <span style={{ textAlign: 'right', color: '#22c55e' }}>{selectedCandle.high.toLocaleString()}</span>
              <span style={{ color: '#ef4444' }}>安値</span>
              <span style={{ textAlign: 'right', color: '#ef4444' }}>{selectedCandle.low.toLocaleString()}</span>
              <span style={{ color: '#94a3b8' }}>終値</span>
              <span style={{ textAlign: 'right', fontWeight: 600 }}>{selectedCandle.close.toLocaleString()}</span>
            </div>
          </div>

          <textarea
            placeholder="この時点にコメントを追加..."
            rows={5}
            style={{
              background: '#1e293b',
              color: '#e2e8f0',
              border: '1px solid #334155',
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
              background: '#1e40af',
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

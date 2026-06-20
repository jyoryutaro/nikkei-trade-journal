import { TIMEFRAMES, WINDOW_OPTIONS } from '../constants/timeframes'

interface Props {
  timeframe: string
  windowHours: number | null
  onTimeframeChange: (v: string) => void
  onWindowChange: (v: number | null) => void
}

const selectStyle: React.CSSProperties = {
  background: '#1e293b',
  color: '#e2e8f0',
  border: '1px solid #334155',
  borderRadius: '4px',
  padding: '4px 10px',
  fontSize: '0.82rem',
  cursor: 'pointer',
}

export function ChartToolbar({ timeframe, windowHours, onTimeframeChange, onWindowChange }: Props) {
  const windows = WINDOW_OPTIONS[timeframe] ?? []

  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px', flexWrap: 'wrap' }}>
      {/* timeframe selector */}
      <div style={{ display: 'flex', gap: '4px' }}>
        {TIMEFRAMES.map(t => {
          const active = t.value === timeframe
          return (
            <button
              key={t.value}
              onClick={() => onTimeframeChange(t.value)}
              style={{
                padding: '4px 10px',
                fontSize: '0.82rem',
                borderRadius: '4px',
                border: 'none',
                cursor: 'pointer',
                background: active ? '#0f766e' : '#1e293b',
                color: active ? '#fff' : '#94a3b8',
                fontWeight: active ? 600 : 400,
              }}
            >
              {t.label}
            </button>
          )
        })}
      </div>

      <div style={{ width: '1px', height: '20px', background: '#334155', margin: '0 4px' }} />

      {/* time window selector (adapts to selected timeframe) */}
      <div style={{ display: 'flex', gap: '4px' }}>
        {windows.map(w => {
          const active = w.hours === windowHours
          return (
            <button
              key={w.label}
              onClick={() => onWindowChange(w.hours)}
              style={{
                padding: '4px 10px',
                fontSize: '0.82rem',
                borderRadius: '4px',
                border: 'none',
                cursor: 'pointer',
                background: active ? '#3b82f6' : '#1e293b',
                color: active ? '#fff' : '#94a3b8',
                fontWeight: active ? 600 : 400,
              }}
            >
              {w.label}
            </button>
          )
        })}
      </div>
    </div>
  )
}

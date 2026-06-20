const WINDOWS: { label: string; hours: number | null }[] = [
  { label: '30分', hours: 0.5 },
  { label: '1時間', hours: 1 },
  { label: '4時間', hours: 4 },
  { label: '全体', hours: null },
]

const TIMEFRAMES = [
  { value: '1m', label: '1分' },
  { value: '5m', label: '5分' },
  { value: '1d', label: '日足' },
]

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
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px', flexWrap: 'wrap' }}>
      <div style={{ display: 'flex', gap: '4px' }}>
        {WINDOWS.map(w => {
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

      <div style={{ width: '1px', height: '20px', background: '#334155', margin: '0 4px' }} />

      <label style={{ color: '#64748b', fontSize: '0.82rem' }}>足種</label>
      <select value={timeframe} onChange={e => onTimeframeChange(e.target.value)} style={selectStyle}>
        {TIMEFRAMES.map(t => (
          <option key={t.value} value={t.value}>{t.label}</option>
        ))}
      </select>
    </div>
  )
}

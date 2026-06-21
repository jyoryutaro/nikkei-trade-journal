import { colors } from '../../theme'

interface Props {
  value: string // "HH:mm"
  onChange: (value: string) => void
}

/** A styled time-of-day input (HH:mm, minute precision). */
export function TimeInput({ value, onChange }: Props) {
  return (
    <input
      type="time"
      value={value}
      onChange={e => onChange(e.target.value)}
      style={{
        background: colors.surface,
        color: colors.text,
        border: `1px solid ${colors.borderStrong}`,
        borderRadius: '6px',
        padding: '8px 10px',
        fontSize: '0.82rem',
        width: '100%',
        boxSizing: 'border-box',
      }}
    />
  )
}

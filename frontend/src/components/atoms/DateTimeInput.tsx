import { colors } from '../../theme'

interface Props {
  value: string // "YYYY-MM-DDTHH:mm"
  onChange: (value: string) => void
}

/** A styled datetime-local input (minute precision). */
export function DateTimeInput({ value, onChange }: Props) {
  return (
    <input
      type="datetime-local"
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

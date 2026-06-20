import { colors } from '../../theme'

interface Props {
  value: string
  onChange: (value: string) => void
  placeholder?: string
  min?: number
  step?: number
}

/** A styled numeric input. Value is kept as a string for controlled editing. */
export function NumberInput({ value, onChange, placeholder, min = 0, step = 1 }: Props) {
  return (
    <input
      type="number"
      inputMode="decimal"
      value={value}
      min={min}
      step={step}
      placeholder={placeholder}
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

import { colors } from '../../theme'

interface Props {
  value: string
  onChange: (value: string) => void
  placeholder?: string
  min?: number
  step?: number
}

/** A styled numeric input that only accepts integers. Value is kept as a string for controlled editing. */
export function NumberInput({ value, onChange, placeholder, min = 0, step = 1 }: Props) {
  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === '.' || e.key === ',' || e.key === 'e' || e.key === 'E') {
      e.preventDefault()
    }
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const raw = e.target.value
    if (raw === '' || raw === '-') {
      onChange(raw)
      return
    }
    if (/^-?\d+$/.test(raw)) onChange(raw)
  }

  return (
    <input
      type="number"
      inputMode="numeric"
      value={value}
      min={min}
      step={step}
      placeholder={placeholder}
      onKeyDown={handleKeyDown}
      onChange={handleChange}
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

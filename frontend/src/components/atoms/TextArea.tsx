import { colors } from '../../theme'

interface Props {
  value: string
  onChange: (value: string) => void
  placeholder?: string
  rows?: number
}

/** A styled multiline text input. */
export function TextArea({ value, onChange, placeholder, rows = 5 }: Props) {
  return (
    <textarea
      value={value}
      placeholder={placeholder}
      rows={rows}
      onChange={e => onChange(e.target.value)}
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
        width: '100%',
        boxSizing: 'border-box',
      }}
    />
  )
}

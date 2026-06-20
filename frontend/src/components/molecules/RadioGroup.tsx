import { colors } from '../../theme'

export interface RadioOption<T> {
  value: T
  label: string
}

interface Props<T> {
  label: string
  name: string
  options: RadioOption<T>[]
  value: T
  onChange: (value: T) => void
}

/** A labelled group of mutually-exclusive radio buttons. */
export function RadioGroup<T extends string>({ label, name, options, value, onChange }: Props<T>) {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
      <span style={{ fontSize: '0.75rem', color: colors.textMuted }}>{label}</span>
      <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap' }}>
        {options.map(opt => (
          <label key={opt.value} style={{ display: 'flex', alignItems: 'center', gap: '4px', fontSize: '0.82rem', cursor: 'pointer' }}>
            <input
              type="radio"
              name={name}
              checked={value === opt.value}
              onChange={() => onChange(opt.value)}
            />
            {opt.label}
          </label>
        ))}
      </div>
    </div>
  )
}

import { ToggleButton } from '../atoms/ToggleButton'

export interface ButtonItem<T> {
  value: T
  label: string
}

interface Props<T> {
  items: ButtonItem<T>[]
  value: T
  onChange: (value: T) => void
  /** Active background colour for the selected button. */
  activeColor?: string
}

/** A horizontal group of mutually-exclusive toggle buttons. */
export function ButtonGroup<T extends string | number | null>({
  items,
  value,
  onChange,
  activeColor,
}: Props<T>) {
  return (
    <div style={{ display: 'flex', gap: '4px' }}>
      {items.map((item, i) => (
        <ToggleButton
          key={`${item.label}-${i}`}
          active={item.value === value}
          activeColor={activeColor}
          onClick={() => onChange(item.value)}
        >
          {item.label}
        </ToggleButton>
      ))}
    </div>
  )
}

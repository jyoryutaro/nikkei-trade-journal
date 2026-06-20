import type { ReactNode } from 'react'
import { colors } from '../../theme'

interface Props {
  active: boolean
  onClick: () => void
  /** Background colour when active (defaults to teal). */
  activeColor?: string
  children: ReactNode
}

/** A single selectable pill button. */
export function ToggleButton({ active, onClick, activeColor = colors.accentTeal, children }: Props) {
  return (
    <button
      onClick={onClick}
      style={{
        padding: '4px 10px',
        fontSize: '0.82rem',
        borderRadius: '4px',
        border: 'none',
        cursor: 'pointer',
        background: active ? activeColor : colors.surface,
        color: active ? '#fff' : colors.textMuted,
        fontWeight: active ? 600 : 400,
      }}
    >
      {children}
    </button>
  )
}

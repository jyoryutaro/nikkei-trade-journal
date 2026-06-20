import type { ReactNode } from 'react'
import { colors } from '../../theme'

interface Props {
  children: ReactNode
  onClick?: () => void
  disabled?: boolean
  type?: 'button' | 'submit'
}

/** Primary action button with a disabled state. */
export function Button({ children, onClick, disabled = false, type = 'button' }: Props) {
  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled}
      style={{
        padding: '8px',
        borderRadius: '6px',
        border: 'none',
        background: disabled ? colors.accentBlueDeep : colors.accentBlue,
        color: disabled ? '#93c5fd' : '#fff',
        fontSize: '0.82rem',
        fontWeight: 600,
        cursor: disabled ? 'not-allowed' : 'pointer',
        opacity: disabled ? 0.5 : 1,
        width: '100%',
      }}
    >
      {children}
    </button>
  )
}

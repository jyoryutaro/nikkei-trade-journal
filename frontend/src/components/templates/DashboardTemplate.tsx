import type { ReactNode } from 'react'
import { colors } from '../../theme'

interface Props {
  header: ReactNode
  toolbar: ReactNode
  chart: ReactNode
  aside: ReactNode
  table: ReactNode
  error?: ReactNode
}

/** Page layout for the dashboard. Pure presentation: it arranges slots and
 * holds no state. */
export function DashboardTemplate({ header, toolbar, chart, aside, table, error }: Props) {
  return (
    <div style={{ background: colors.bg, minHeight: '100vh', color: colors.text, fontFamily: 'sans-serif', padding: '20px 24px' }}>
      {header}
      {error}
      {toolbar}
      <div style={{ display: 'flex', gap: '16px', alignItems: 'flex-start' }}>
        <div style={{ flex: 1, minWidth: 0 }}>{chart}</div>
        {aside}
      </div>
      {table}
    </div>
  )
}

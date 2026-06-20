import { ButtonGroup, type ButtonItem } from '../molecules/ButtonGroup'
import { Divider } from '../atoms/Divider'
import { TIMEFRAMES, WINDOW_OPTIONS } from '../../constants/timeframes'
import { colors } from '../../theme'

interface Props {
  timeframe: string
  windowHours: number | null
  onTimeframeChange: (tf: string) => void
  onWindowChange: (hours: number | null) => void
}

/** Timeframe + time-window selectors for the chart. */
export function ChartToolbar({ timeframe, windowHours, onTimeframeChange, onWindowChange }: Props) {
  const windowItems: ButtonItem<number | null>[] = (WINDOW_OPTIONS[timeframe] ?? []).map(w => ({
    value: w.hours,
    label: w.label,
  }))

  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px', flexWrap: 'wrap' }}>
      <ButtonGroup
        items={TIMEFRAMES}
        value={timeframe}
        onChange={onTimeframeChange}
        activeColor={colors.accentTeal}
      />
      <Divider height={20} />
      <ButtonGroup
        items={windowItems}
        value={windowHours}
        onChange={onWindowChange}
        activeColor={colors.accentBlue}
      />
    </div>
  )
}

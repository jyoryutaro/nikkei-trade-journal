import { colors } from '../../theme'

interface Props {
  /** Height of the vertical rule in pixels. */
  height?: number
}

/** A thin vertical separator. */
export function Divider({ height = 16 }: Props) {
  return <div style={{ width: '1px', height: `${height}px`, background: colors.borderStrong }} />
}

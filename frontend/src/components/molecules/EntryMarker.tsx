import type { JournalEntry } from '../../api/journal'
import { formatTimeJST } from '../../api/marketData'
import { colors } from '../../theme'

interface Props {
  entry: JournalEntry
  /** Pixel coordinates within the chart overlay. */
  x: number
  y: number
  hovered: boolean
  onHover: () => void
  onLeave: () => void
}

const sideLabel = (s: string) => (s === 'long' ? '買い' : s === 'short' ? '売り' : '')
const typeLabel = (t: string) => (t === 'open' ? '新規' : t === 'close' ? '決済' : '')

/** A chart marker placed at an entry's (time, price) coordinate. On hover it
 * shows a speech-bubble tooltip to the right. */
export function EntryMarker({ entry, x, y, hovered, onHover, onLeave }: Props) {
  const isLong = entry.side === 'long'
  const color = isLong ? colors.up : colors.down

  return (
    <div
      style={{ position: 'absolute', left: `${x}px`, top: `${y}px`, transform: 'translate(-50%, -50%)', pointerEvents: 'auto' }}
      onMouseEnter={onHover}
      onMouseLeave={onLeave}
    >
      <div
        style={{
          width: '14px',
          height: '14px',
          borderRadius: '50%',
          background: color,
          border: '2px solid #fff',
          boxShadow: '0 0 0 1px rgba(0,0,0,0.35)',
          cursor: 'pointer',
        }}
      />

      {hovered && (
        <div
          style={{
            position: 'absolute',
            left: '16px',
            top: '50%',
            transform: 'translateY(-50%)',
            background: colors.surface,
            border: `1px solid ${colors.borderStrong}`,
            borderRadius: '8px',
            padding: '8px 10px',
            width: '180px',
            fontSize: '0.72rem',
            color: colors.text,
            lineHeight: 1.6,
            boxShadow: '0 6px 16px rgba(0,0,0,0.45)',
            zIndex: 20,
            pointerEvents: 'none',
          }}
        >
          {/* tail pointing left toward the marker */}
          <div
            style={{
              position: 'absolute',
              left: '-6px',
              top: '50%',
              transform: 'translateY(-50%)',
              width: 0,
              height: 0,
              borderTop: '6px solid transparent',
              borderBottom: '6px solid transparent',
              borderRight: `6px solid ${colors.surface}`,
            }}
          />
          <div style={{ color: colors.textFaint, marginBottom: '2px' }}>{formatTimeJST(entry.time)}</div>
          <div style={{ fontWeight: 600, color }}>
            {sideLabel(entry.side)} ・ {typeLabel(entry.tradeType)}
          </div>
          {entry.price != null && <div>金額: {entry.price.toLocaleString()}</div>}
          {entry.comment && (
            <div style={{ color: colors.textMuted, marginTop: '2px', whiteSpace: 'pre-wrap' }}>{entry.comment}</div>
          )}
        </div>
      )}
    </div>
  )
}

import { useState } from 'react'
import type { Candle } from '../../api/marketData'
import { createJournalEntry, type Side, type TradeType } from '../../api/journal'
import { RadioGroup } from '../molecules/RadioGroup'
import { NumberInput } from '../atoms/NumberInput'
import { TextArea } from '../atoms/TextArea'
import { Button } from '../atoms/Button'
import { colors } from '../../theme'

interface Props {
  contract: string
  candle: Candle
  onSubmitted?: () => void
}

const SIDE_OPTIONS: { value: Side; label: string }[] = [
  { value: '', label: '未選択' },
  { value: 'long', label: '買い' },
  { value: 'short', label: '売り' },
]

const TRADE_TYPE_OPTIONS: { value: TradeType; label: string }[] = [
  { value: 'open', label: '新規' },
  { value: 'close', label: '決済' },
]

/** Form for recording a position (side / trade type / price) or a comment-only
 * note. Position fields appear only once a side is chosen, and the submit
 * button is disabled until the entry is valid. */
export function TradeEntryForm({ contract, candle, onSubmitted }: Props) {
  const [side, setSide] = useState<Side>('')
  const [tradeType, setTradeType] = useState<TradeType>('')
  const [price, setPrice] = useState('')
  const [comment, setComment] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const hasPosition = side !== ''
  const priceNum = Number(price)
  const priceValid = price.trim() !== '' && !Number.isNaN(priceNum) && priceNum > 0

  // Validation: a position needs trade type + price; a comment-only entry needs
  // a non-empty comment.
  const valid = hasPosition ? tradeType !== '' && priceValid : comment.trim() !== ''

  const handleSideChange = (next: Side) => {
    setSide(next)
    if (next === '') {
      setTradeType('')
      setPrice('')
    }
  }

  const reset = () => {
    setSide('')
    setTradeType('')
    setPrice('')
    setComment('')
  }

  const handleSubmit = async () => {
    if (!valid || submitting) return
    setSubmitting(true)
    setError(null)
    try {
      await createJournalEntry({
        contract,
        time: candle.time,
        side,
        tradeType: hasPosition ? tradeType : '',
        price: hasPosition ? priceNum : null,
        comment,
      })
      reset()
      onSubmitted?.()
    } catch (e) {
      setError(e instanceof Error ? e.message : '保存に失敗しました')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
      <RadioGroup label="売買" name="side" options={SIDE_OPTIONS} value={side} onChange={handleSideChange} />

      {hasPosition && (
        <>
          <RadioGroup label="種別" name="tradeType" options={TRADE_TYPE_OPTIONS} value={tradeType} onChange={setTradeType} />
          <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
            <span style={{ fontSize: '0.75rem', color: colors.textMuted }}>金額</span>
            <NumberInput value={price} onChange={setPrice} placeholder="例: 39000" step={5} />
          </div>
        </>
      )}

      <div style={{ display: 'flex', flexDirection: 'column', gap: '6px' }}>
        <span style={{ fontSize: '0.75rem', color: colors.textMuted }}>コメント{hasPosition ? '（任意）' : ''}</span>
        <TextArea value={comment} onChange={setComment} placeholder="この時点にコメントを追加..." rows={4} />
      </div>

      {error && <p style={{ color: colors.down, fontSize: '0.78rem', margin: 0 }}>{error}</p>}

      <Button onClick={handleSubmit} disabled={!valid || submitting} type="button">
        {submitting ? '保存中...' : '保存'}
      </Button>
    </div>
  )
}

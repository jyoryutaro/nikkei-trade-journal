import { Select } from '../atoms/Select'
import { CONTRACTS } from '../../constants/contracts'
import { colors } from '../../theme'

interface Props {
  contract: string
  onContractChange: (contract: string) => void
}

/** Label + select for choosing the futures contract month. */
export function ContractSelector({ contract, onContractChange }: Props) {
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
      <span style={{ fontSize: '0.8rem', color: colors.textMuted }}>限月</span>
      <Select value={contract} onChange={onContractChange} options={CONTRACTS} />
    </div>
  )
}

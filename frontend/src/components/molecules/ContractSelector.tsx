import { Select, type Option } from '../atoms/Select'
import { colors } from '../../theme'

interface Props {
  contract: string
  options: Option[]
  onContractChange: (contract: string) => void
}

/** Label + select for choosing the symbol/contract. Options are provided by the
 * caller (fetched from the backend). */
export function ContractSelector({ contract, options, onContractChange }: Props) {
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
      <span style={{ fontSize: '0.8rem', color: colors.textMuted }}>銘柄</span>
      <Select value={contract} onChange={onContractChange} options={options} />
    </div>
  )
}

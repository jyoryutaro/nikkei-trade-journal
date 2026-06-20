import { Divider } from '../atoms/Divider'
import { ContractSelector } from '../molecules/ContractSelector'
import { colors } from '../../theme'

interface Props {
  contract: string
  onContractChange: (contract: string) => void
}

/** Top bar: app title and contract selector. */
export function AppHeader({ contract, onContractChange }: Props) {
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: '16px', marginBottom: '20px', borderBottom: `1px solid ${colors.border}`, paddingBottom: '16px' }}>
      <h1 style={{ fontSize: '1rem', fontWeight: 600, color: colors.textMuted, margin: 0 }}>
        日経225先物 トレードジャーナル
      </h1>
      <Divider />
      <ContractSelector contract={contract} onContractChange={onContractChange} />
    </div>
  )
}

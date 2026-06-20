const CONTRACTS = [
  { value: '2609', label: 'NKD Sep-2026' },
]

interface Props {
  contract: string
  onContractChange: (v: string) => void
}

export function ContractSelector({ contract, onContractChange }: Props) {
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
      <span style={{ color: '#64748b', fontSize: '0.75rem', textTransform: 'uppercase', letterSpacing: '0.05em' }}>限月</span>
      <select
        value={contract}
        onChange={e => onContractChange(e.target.value)}
        style={{
          background: 'transparent',
          color: '#e2e8f0',
          border: 'none',
          fontSize: '1rem',
          fontWeight: 600,
          cursor: 'pointer',
          padding: '0',
        }}
      >
        {CONTRACTS.map(c => (
          <option key={c.value} value={c.value} style={{ background: '#1e293b' }}>{c.label}</option>
        ))}
      </select>
    </div>
  )
}

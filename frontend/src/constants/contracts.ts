import type { Option } from '../components/atoms/Select'

// Selectable contract months (major-SQ: Mar/Jun/Sep/Dec). Data availability
// depends on what has been seeded.
export const CONTRACTS: Option[] = [
  { value: '2606', label: '2606 (Jun-2026)' },
  { value: '2609', label: '2609 (Sep-2026)' },
  { value: '2612', label: '2612 (Dec-2026)' },
  { value: '2703', label: '2703 (Mar-2027)' },
]

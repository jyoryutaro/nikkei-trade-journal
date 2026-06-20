import { useCallback, useEffect, useRef, useState } from 'react'
import {
  createChart,
  CandlestickSeries,
  type IChartApi,
  type ISeriesApi,
  type UTCTimestamp,
} from 'lightweight-charts'
import type { Candle } from '../../api/marketData'
import type { JournalEntry } from '../../api/journal'
import { EntryMarker } from '../molecules/EntryMarker'

interface Props {
  candles: Candle[]
  entries: JournalEntry[]
  windowHours: number | null // null = fit all content
  onSelect: (candle: Candle | null) => void
}

interface MarkerPos {
  entry: JournalEntry
  x: number
  y: number
}

function applyWindow(chart: IChartApi, candles: Candle[], windowHours: number | null) {
  if (candles.length === 0) return
  if (windowHours === null) {
    chart.timeScale().fitContent()
    return
  }
  const last = candles[candles.length - 1].time as UTCTimestamp
  const from = (last - windowHours * 3600) as UTCTimestamp
  chart.timeScale().setVisibleRange({ from, to: last })
}

export function CandlestickChart({ candles, entries, windowHours, onSelect }: Props) {
  const containerRef = useRef<HTMLDivElement>(null)
  const chartHostRef = useRef<HTMLDivElement>(null)
  const chartRef = useRef<IChartApi | null>(null)
  const seriesRef = useRef<ISeriesApi<'Candlestick'> | null>(null)
  const candlesRef = useRef<Candle[]>(candles)
  const entriesRef = useRef<JournalEntry[]>(entries)
  candlesRef.current = candles
  entriesRef.current = entries

  const [markers, setMarkers] = useState<MarkerPos[]>([])
  const [hoveredId, setHoveredId] = useState<number | null>(null)

  // Project each priced entry onto pixel coordinates at its (time, price).
  const recomputeMarkers = useCallback(() => {
    const chart = chartRef.current
    const series = seriesRef.current
    if (!chart || !series) return
    const next: MarkerPos[] = []
    for (const e of entriesRef.current) {
      if (e.price == null || e.side === '') continue
      const x = chart.timeScale().timeToCoordinate(e.time as UTCTimestamp)
      const y = series.priceToCoordinate(e.price)
      if (x == null || y == null) continue
      next.push({ entry: e, x, y })
    }
    setMarkers(next)
  }, [])

  // create chart once on mount
  useEffect(() => {
    if (!chartHostRef.current) return

    const chart = createChart(chartHostRef.current, {
      width: chartHostRef.current.clientWidth,
      height: 420,
      layout: { background: { color: '#0f172a' }, textColor: '#e2e8f0' },
      grid: { vertLines: { color: '#1e293b' }, horzLines: { color: '#1e293b' } },
      crosshair: { mode: 1 },
      timeScale: { timeVisible: true, secondsVisible: false },
    })

    const series = chart.addSeries(CandlestickSeries, {
      upColor: '#22c55e',
      downColor: '#ef4444',
      borderVisible: false,
      wickUpColor: '#22c55e',
      wickDownColor: '#ef4444',
    })

    chart.subscribeClick(param => {
      if (!param.time) {
        onSelect(null)
        return
      }
      const clicked = param.time as number
      const found = candlesRef.current.find(c => c.time === clicked) ?? null
      onSelect(found)
    })

    // markers must follow pan/zoom and resize
    chart.timeScale().subscribeVisibleTimeRangeChange(recomputeMarkers)

    const handleResize = () => {
      if (chartHostRef.current) chart.applyOptions({ width: chartHostRef.current.clientWidth })
      recomputeMarkers()
    }
    window.addEventListener('resize', handleResize)

    chartRef.current = chart
    seriesRef.current = series

    return () => {
      window.removeEventListener('resize', handleResize)
      chart.timeScale().unsubscribeVisibleTimeRangeChange(recomputeMarkers)
      chart.remove()
      chartRef.current = null
      seriesRef.current = null
    }
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  // update series when candles change
  useEffect(() => {
    if (!seriesRef.current || !chartRef.current) return
    seriesRef.current.setData(
      candles.map(c => ({
        time: c.time as UTCTimestamp,
        open: c.open,
        high: c.high,
        low: c.low,
        close: c.close,
      }))
    )
    applyWindow(chartRef.current, candles, windowHours)
    recomputeMarkers()
  }, [candles]) // eslint-disable-line react-hooks/exhaustive-deps

  // update visible window when windowHours changes independently
  useEffect(() => {
    if (!chartRef.current) return
    applyWindow(chartRef.current, candlesRef.current, windowHours)
    recomputeMarkers()
  }, [windowHours, recomputeMarkers])

  // recompute when entries change
  useEffect(() => {
    recomputeMarkers()
  }, [entries, recomputeMarkers])

  return (
    <div ref={containerRef} style={{ width: '100%', position: 'relative' }}>
      <div ref={chartHostRef} style={{ width: '100%' }} />
      <div style={{ position: 'absolute', inset: 0, pointerEvents: 'none' }}>
        {markers.map(m => (
          <EntryMarker
            key={m.entry.id}
            entry={m.entry}
            x={m.x}
            y={m.y}
            hovered={hoveredId === m.entry.id}
            onHover={() => setHoveredId(m.entry.id)}
            onLeave={() => setHoveredId(null)}
          />
        ))}
      </div>
    </div>
  )
}

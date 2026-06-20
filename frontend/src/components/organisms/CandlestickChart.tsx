import { useCallback, useEffect, useRef, useState } from 'react'
import {
  createChart,
  CandlestickSeries,
  type IChartApi,
  type ISeriesApi,
  type UTCTimestamp,
  type Logical,
} from 'lightweight-charts'
import type { Candle } from '../../api/marketData'
import type { JournalEntry } from '../../api/journal'
import { INTERVAL_SECONDS } from '../../constants/timeframes'
import { EntryMarker } from '../molecules/EntryMarker'

interface Props {
  candles: Candle[]
  entries: JournalEntry[]
  timeframe: string
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

export function CandlestickChart({ candles, entries, timeframe, windowHours, onSelect }: Props) {
  const containerRef = useRef<HTMLDivElement>(null)
  const chartHostRef = useRef<HTMLDivElement>(null)
  const chartRef = useRef<IChartApi | null>(null)
  const seriesRef = useRef<ISeriesApi<'Candlestick'> | null>(null)
  const candlesRef = useRef<Candle[]>(candles)
  const entriesRef = useRef<JournalEntry[]>(entries)
  const intervalRef = useRef<number>(INTERVAL_SECONDS[timeframe] ?? 60)
  const rafRef = useRef<number | null>(null)
  candlesRef.current = candles
  entriesRef.current = entries
  intervalRef.current = INTERVAL_SECONDS[timeframe] ?? 60

  const [markers, setMarkers] = useState<MarkerPos[]>([])
  const [markerSize, setMarkerSize] = useState(10)
  const [hoveredId, setHoveredId] = useState<number | null>(null)

  // Project each priced entry onto pixel coordinates at its (time, price), and
  // size the marker relative to the candle width. Markers outside the data
  // pane (axis gutters / off-screen) are dropped.
  const recomputeMarkers = useCallback(() => {
    const chart = chartRef.current
    const series = seriesRef.current
    if (!chart || !series) return
    const timeScale = chart.timeScale()
    const paneWidth = timeScale.width()
    const paneHeight = (chartHostRef.current?.clientHeight ?? 0) - timeScale.height()

    // candle pitch in px (uniform across the view) → marker a bit smaller
    const c0 = timeScale.logicalToCoordinate(0 as Logical)
    const c1 = timeScale.logicalToCoordinate(1 as Logical)
    if (c0 != null && c1 != null) {
      const spacing = Math.abs(c1 - c0)
      setMarkerSize(Math.max(2, Math.min(10, Math.round(spacing * 0.32))))
    }

    const interval = intervalRef.current
    const next: MarkerPos[] = []
    for (const e of entriesRef.current) {
      if (e.price == null || e.side === '') continue
      // snap the entry time to the candle bucket of the current timeframe so it
      // aligns with a bar on aggregated timeframes (matches the backend buckets)
      const bucketTime = (Math.floor(e.time / interval) * interval) as UTCTimestamp
      const x = timeScale.timeToCoordinate(bucketTime)
      const y = series.priceToCoordinate(e.price)
      if (x == null || y == null) continue
      if (x < 0 || x > paneWidth || y < 0 || y > paneHeight) continue
      next.push({ entry: e, x, y })
    }
    setMarkers(next)
  }, [])

  // Defer recompute to the next frame so the chart's price autoscale and bar
  // spacing have settled before we read coordinates.
  const scheduleRecompute = useCallback(() => {
    if (rafRef.current != null) cancelAnimationFrame(rafRef.current)
    rafRef.current = requestAnimationFrame(() => {
      rafRef.current = null
      recomputeMarkers()
    })
  }, [recomputeMarkers])

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
      // Scroll/zoom are disabled: the range is driven by the toolbar's window
      // buttons. This keeps the HTML entry markers perfectly aligned with the
      // canvas (no per-frame lag during panning).
      handleScroll: false,
      handleScale: false,
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

    chart.timeScale().subscribeVisibleTimeRangeChange(scheduleRecompute)

    const handleResize = () => {
      if (chartHostRef.current) chart.applyOptions({ width: chartHostRef.current.clientWidth })
      scheduleRecompute()
    }
    window.addEventListener('resize', handleResize)

    chartRef.current = chart
    seriesRef.current = series

    return () => {
      window.removeEventListener('resize', handleResize)
      if (rafRef.current != null) cancelAnimationFrame(rafRef.current)
      chart.timeScale().unsubscribeVisibleTimeRangeChange(scheduleRecompute)
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
    scheduleRecompute()
  }, [candles]) // eslint-disable-line react-hooks/exhaustive-deps

  // update visible window when windowHours changes independently
  useEffect(() => {
    if (!chartRef.current) return
    applyWindow(chartRef.current, candlesRef.current, windowHours)
    scheduleRecompute()
  }, [windowHours, scheduleRecompute])

  // recompute when entries change
  useEffect(() => {
    scheduleRecompute()
  }, [entries, scheduleRecompute])

  return (
    <div ref={containerRef} style={{ width: '100%', position: 'relative' }}>
      <div ref={chartHostRef} style={{ width: '100%' }} />
      <div style={{ position: 'absolute', inset: 0, pointerEvents: 'none', zIndex: 3 }}>
        {markers.map(m => (
          <EntryMarker
            key={m.entry.id}
            entry={m.entry}
            x={m.x}
            y={m.y}
            size={markerSize}
            hovered={hoveredId === m.entry.id}
            onHover={() => setHoveredId(m.entry.id)}
            onLeave={() => setHoveredId(null)}
          />
        ))}
      </div>
    </div>
  )
}

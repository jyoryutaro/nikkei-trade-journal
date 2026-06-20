import { useEffect, useRef } from 'react'
import {
  createChart,
  CandlestickSeries,
  type IChartApi,
  type ISeriesApi,
  type UTCTimestamp,
} from 'lightweight-charts'
import type { Candle } from '../api/marketData'

interface Props {
  candles: Candle[]
  windowHours: number | null  // null = fit all content
  onSelect: (candle: Candle | null) => void
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

export function CandlestickChart({ candles, windowHours, onSelect }: Props) {
  const containerRef = useRef<HTMLDivElement>(null)
  const chartRef     = useRef<IChartApi | null>(null)
  const seriesRef    = useRef<ISeriesApi<'Candlestick'> | null>(null)
  const candlesRef   = useRef<Candle[]>(candles)
  candlesRef.current = candles

  // create chart once on mount
  useEffect(() => {
    if (!containerRef.current) return

    const chart = createChart(containerRef.current, {
      width: containerRef.current.clientWidth,
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
      if (!param.time) { onSelect(null); return }
      const clicked = param.time as number
      const found = candlesRef.current.find(c => c.time === clicked) ?? null
      onSelect(found)
    })

    const handleResize = () => {
      if (containerRef.current) chart.applyOptions({ width: containerRef.current.clientWidth })
    }
    window.addEventListener('resize', handleResize)

    chartRef.current  = chart
    seriesRef.current = series

    return () => {
      window.removeEventListener('resize', handleResize)
      chart.remove()
      chartRef.current  = null
      seriesRef.current = null
    }
  }, []) // eslint-disable-line react-hooks/exhaustive-deps

  // update series when candles change
  useEffect(() => {
    if (!seriesRef.current || !chartRef.current) return
    seriesRef.current.setData(
      candles.map(c => ({
        time:  c.time as UTCTimestamp,
        open:  c.open,
        high:  c.high,
        low:   c.low,
        close: c.close,
      }))
    )
    applyWindow(chartRef.current, candles, windowHours)
  }, [candles]) // eslint-disable-line react-hooks/exhaustive-deps

  // update visible window when windowHours changes independently
  useEffect(() => {
    if (!chartRef.current) return
    applyWindow(chartRef.current, candlesRef.current, windowHours)
  }, [windowHours])

  return <div ref={containerRef} style={{ width: '100%' }} />
}

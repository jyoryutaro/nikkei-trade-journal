# UC-01: 市場データ取得・保存

Yahoo Finance API から日経先物（NKD=F）の OHLCV データを取得し、永続化する。

---

## トリガー

フロントエンドからのチャートデータリクエスト時（オンデマンド）

---

## 基本フロー

1. 以下のエンドポイントへリクエストする

   ```
   GET https://query2.finance.yahoo.com/v8/finance/chart/NKD%3DF
       ?interval={interval}
       &range={range}
   ```

   | パラメータ | 値の例 | 説明 |
   |---|---|---|
   | `interval` | `1d`, `1h`, `5m` | 足の種類 |
   | `range` | `1mo`, `6mo`, `1y` | 取得期間 |

2. レスポンスから以下を抽出する

   | JSONパス | 型 | 説明 |
   |---|---|---|
   | `chart.result[0].timestamp[]` | `int64[]` | UNIX 時刻（秒） |
   | `chart.result[0].indicators.quote[0].open[]` | `float64[]` | 始値 |
   | `chart.result[0].indicators.quote[0].high[]` | `float64[]` | 高値 |
   | `chart.result[0].indicators.quote[0].low[]` | `float64[]` | 安値 |
   | `chart.result[0].indicators.quote[0].close[]` | `float64[]` | 終値 |
   | `chart.result[0].indicators.quote[0].volume[]` | `int64[]` | 出来高 |

   全配列は同一長。`null` 要素は欠損値としてスキップする。

3. 抽出した値を `MarketData` エンティティ（symbol, timestamp, open, high, low, close, volume）として構築し、`MarketDataRepository` 経由で保存する。`(symbol, timestamp)` が一致する既存レコードは上書きしない。

4. 保存済みの `MarketData` リストを返す。

---

## 代替フロー

| 条件 | 処理 |
|---|---|
| HTTP 429 | エラーを返す。リトライは呼び出し元に委ねる |
| HTTP 4xx / 5xx | エラーを返す |
| `chart.result` が `null` | 空リストを返す |

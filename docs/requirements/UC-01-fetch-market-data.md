# UC-01: 市場データ取得・保存

Yahoo Finance API から日経先物（NKD=F）の OHLCV データを取得し、永続化する。

---

## トリガー

フロントエンドからのチャートデータリクエスト時（オンデマンド）

---

## 基本フロー

1. 取得ウィンドウ（`period1` / `period2`）を決定する

   - **初回（当該 `symbol` × `interval` の保存済みデータが無い）**: `period1` = その先物の取引開始日。開始日が不明な場合は `0`（提供元が保持する最古から）とする。`period2` = 現在時刻。→ 開始日から全件取得。
   - **2回目以降（既存データあり）**: `period1` = 保存済みの最新 `timestamp` + 1 秒、`period2` = 現在時刻。→ 未取得分のみを追加取得（差分取得）。

2. 以下のエンドポイントへリクエストする

   ```
   GET https://query2.finance.yahoo.com/v8/finance/chart/NKD%3DF
       ?interval={interval}
       &period1={period1}
       &period2={period2}
   ```

   | パラメータ | 値の例 | 説明 |
   |---|---|---|
   | `interval` | `1m`, `5m`, `1h`, `1d` | 足の種類 |
   | `period1` | `0`, `1718841600` | 取得開始（UNIX 秒）。`range` の代わりに使用 |
   | `period2` | `1718841600` | 取得終了（UNIX 秒）。通常は現在時刻 |

3. レスポンスから以下を抽出する

   | JSONパス | 型 | 説明 |
   |---|---|---|
   | `chart.result[0].timestamp[]` | `int64[]` | UNIX 時刻（秒） |
   | `chart.result[0].indicators.quote[0].open[]` | `float64[]` | 始値 |
   | `chart.result[0].indicators.quote[0].high[]` | `float64[]` | 高値 |
   | `chart.result[0].indicators.quote[0].low[]` | `float64[]` | 安値 |
   | `chart.result[0].indicators.quote[0].close[]` | `float64[]` | 終値 |
   | `chart.result[0].indicators.quote[0].volume[]` | `int64[]` | 出来高 |

   全配列は同一長。`null` 要素は欠損値としてスキップする。

4. 抽出した値を `MarketData` エンティティ（symbol, interval, timestamp, open, high, low, close, volume）として構築し、`MarketDataRepository` 経由で保存する。`(symbol, interval, timestamp)` が一致する既存レコードは追加しない（重複排除し、未取得分のみを追加する）。

5. 保存済みの `MarketData` リストを返す。

---

## 代替フロー

| 条件 | 処理 |
|---|---|
| HTTP 429 | エラーを返す。リトライは呼び出し元に委ねる |
| HTTP 4xx / 5xx | エラーを返す |
| `chart.result` が `null` | 空リストを返す |
| 短い足（特に `1m`）で要求期間が提供元の上限を超える | 提供元が返せる範囲のみ取得される（`1m` は直近数日〜約30日、その他の分足・時間足にも期間制限あり）。全期間が必要な場合は `period1`/`period2` を分割してページングする |
| 差分取得で新規データが無い（`period1` 以降に新しい足が無い） | 空リストを返す（既存データは保持） |

# 7. データ要件

主なデータ項目（エンティティ）:

| エンティティ | 主な属性 |
| --- | --- |
| MarketData（相場） | id, 限月(contract), 日時(ts), 始値, 高値, 安値, 終値, 出来高 |
| Trade（トレード） | id, 限月, 方向(long/short), 枚数, IN日時, IN価格, OUT日時, OUT価格, 損益, ステータス |
| Comment（コメント） | id, 対象(trade/時点), 日時, 本文, タグ |
| User（ユーザー）※将来 | id, メールアドレス, 表示名, 作成日時（マルチユーザー化時に追加。各レコードに userId を付与） |

- 対象は日経225先物。相場データはメジャーSQ（3/6/9/12月限）の各限月を可能な限り全て保持する（取得不可なら取得できる分で可。[12-open-issues.md](12-open-issues.md) Q-05）。
- Trade は限月を保持する。数量の単位は「枚（contract）」とする。
- 商品種別（ラージ / ミニ / マイクロ、OSE / CME 等）は限定せず、データ取得先に依存する（[12-open-issues.md](12-open-issues.md) Q-06）。
- データソース: Yahoo Finance chart API（`NKD=F`）から取得する（[要件定義 UC-01](../spec/use-cases/UC-01-fetch-market-data.md)、[12-open-issues.md](12-open-issues.md) Q-01）。

## MySQL テーブル構成

- ストレージは MySQL。当面はローカル単一ユーザーのため userId は持たない（マルチユーザー化時に各テーブルへ userId を追加）。
- 実装済みスキーマ: [`db/init/001_schema.sql`](../../db/init/001_schema.sql)

| テーブル | 状態 | 概要 |
| --- | --- | --- |
| `market_data` | 実装済 | 限月別の OHLCV。`UNIQUE(contract, ts)` で重複排除 |
| `trades` | 未実装（予定） | ポジションの IN / OUT |
| `comments` | 未実装（予定） | トレード / 時点へのコメント |

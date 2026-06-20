# 7. データ要件

主なデータ項目（エンティティ）:

| エンティティ | 主な属性 |
| --- | --- |
| MarketData（相場） | 限月, 足種, 日時(UTC), 始値, 高値, 安値, 終値, 出来高 |
| JournalEntry（記録） | id, 限月, 対象時点(UTC), 売買(long/short/なし), 種別(open/close), 金額, コメント, 登録日時 |
| User（ユーザー）※将来 | id, メールアドレス, 表示名, 作成日時（マルチユーザー化時に追加。各レコードに userId を付与） |

- 対象は日経225先物。相場データはメジャーSQ（3/6/9/12月限）の各限月を可能な限り全て保持する（取得不可なら取得できる分で可。[12-open-issues.md](12-open-issues.md) Q-05）。
- ポジション記録（JournalEntry）は限月・対象時点・売買方向・種別・金額・コメントを保持する。1レコードがポジション記録またはコメントのみを表す（[設計 03-position-recording](../design/03-position-recording.md)）。
- 商品種別（ラージ / ミニ / マイクロ、OSE / CME 等）は限定せず、データ取得先に依存する（[12-open-issues.md](12-open-issues.md) Q-06）。
- データソース: <日経225先物の価格データの取得元。例：API名 / 手入力 / CSV取込。TBD（[12-open-issues.md](12-open-issues.md) Q-01）>

## MarketData テーブル設計（現行実装）

```sql
CREATE TABLE market_data (
  id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  contract   VARCHAR(10)    NOT NULL COMMENT '限月 e.g. 2609',
  timeframe  VARCHAR(5)     NOT NULL DEFAULT '1m' COMMENT 'ベース足種。派生足はオンザフライ集計',
  ts         DATETIME       NOT NULL COMMENT '足の開始日時 (UTC)',
  open       DECIMAL(12,2)  NOT NULL,
  high       DECIMAL(12,2)  NOT NULL,
  low        DECIMAL(12,2)  NOT NULL,
  close      DECIMAL(12,2)  NOT NULL,
  volume     BIGINT UNSIGNED NOT NULL DEFAULT 0,
  created_at DATETIME       NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uq_contract_tf_ts (contract, timeframe, ts)
);
```

### 足種の保存方針

| 足種 | DB 保存 | 取得方法 |
| --- | --- | --- |
| 1分足 (`1m`) | ✅ 保存 | DB 直読み |
| 5分足 (`5m`) | ❌ 保存しない | バックエンドが 1m を集計 |
| 30分足 (`30m`) | ❌ 保存しない | バックエンドが 1m を集計 |
| 1時間足 (`1h`) | ❌ 保存しない | バックエンドが 1m を集計 |
| 日足 (`1d`) | ❌ 保存しない | バックエンドが 1m を集計 |

1m をベースとしてバックエンドで集計することで、DB スキーマをシンプルに保つ。
将来的にデータ量が増えマテリアライズが必要になれば再検討する。

## JournalEntry テーブル設計（現行実装）

```sql
CREATE TABLE journal_entries (
  id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  contract   VARCHAR(10)   NOT NULL COMMENT '限月 e.g. 2609',
  ts         DATETIME      NOT NULL COMMENT '対象時点 (UTC)',
  side       VARCHAR(5)    NULL     COMMENT 'long | short (NULL = コメントのみ)',
  trade_type VARCHAR(5)    NULL     COMMENT 'open(新規) | close(決済)',
  price      DECIMAL(12,2) NULL     COMMENT '約定金額',
  comment    TEXT          NULL,
  created_at DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_contract_ts (contract, ts)
);
```

- `side` が NULL のレコードはコメントのみの記録を表す。検証ルールは[設計 03-position-recording](../design/03-position-recording.md)を参照。
- マルチユーザー化（[FR-00](04-functional/FR-00-auth.md) / Q-03）の際は `user_id` 列を追加してユーザーごとに分離する。

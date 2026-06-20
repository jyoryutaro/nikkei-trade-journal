# アーキテクチャ

実装の構造方針。バックエンドはドメイン駆動設計（DDD）、フロントエンドはアトミックデザインに従う。

## バックエンド（DDD）

`backend/internal/` を4層に分離し、依存方向は**外側 → 内側**（infrastructure / interfaces → application → domain）に統一する。ドメインは他層に依存しない。

```
backend/
├── cmd/
│   ├── server/      # composition root（依存を組み立ててHTTPサーバ起動）
│   └── seed/        # composition root（取得JSON → Import ユースケース）
└── internal/
    ├── domain/marketdata/      # エンティティ/値オブジェクト/ドメインサービス/ポート
    │   ├── candle.go           #   Candle（OHLCV 値オブジェクト）
    │   ├── timeframe.go        #   Timeframe（値オブジェクト）
    │   ├── aggregator.go       #   Aggregate（純粋なドメインサービス）
    │   └── repository.go       #   Repository（永続化ポート）
    ├── application/            # ユースケース（MarketDataService）
    ├── infrastructure/         # アダプタ
    │   ├── db/                 #   MySQL 接続
    │   ├── persistence/mysql/  #   Repository 実装
    │   └── yahoo/              #   Yahoo Finance JSON パーサ
    └── interfaces/http/        # 配信層（ハンドラ/DTO/ルータ/CORS）
```

### 層の責務

| 層 | 責務 | 依存先 |
| --- | --- | --- |
| domain | 業務ルール。永続化・通信を知らない | なし |
| application | ユースケースの調整。ポート経由でドメインを操作 | domain |
| infrastructure | DB・外部APIなどの技術的詳細。ポートを実装 | domain |
| interfaces | HTTP 入出力。DTO ↔ ドメインの変換 | application, domain |

- リポジトリは domain に**インターフェース（ポート）**として定義し、実装は infrastructure に置く（依存性逆転）。
- 例: `var _ marketdata.Repository = (*MarketDataRepository)(nil)` で実装の整合をコンパイル時に保証。

## フロントエンド（Atomic Design）

`frontend/src/components/` を5階層に分離する。下位は表示に専念し、状態とデータ取得は page に集約する。

```
frontend/src/
├── components/
│   ├── atoms/       # 最小部品（ToggleButton / Select / Divider）
│   ├── molecules/   # atoms の組合せ（ButtonGroup / ContractSelector / OhlcvSummary）
│   ├── organisms/   # 意味のあるまとまり（AppHeader / ChartToolbar / CandlestickChart / CommentPanel / PriceTable）
│   ├── templates/   # レイアウトのみ（DashboardTemplate）
│   └── pages/       # 状態・データ取得を集約（DashboardPage）
├── api/             # API クライアント（marketData.ts）
├── constants/       # 定数（timeframes / contracts）
└── theme.ts         # 配色トークン
```

### 方針

| 階層 | 責務 |
| --- | --- |
| atoms | 単一の表示部品。状態を持たない |
| molecules | atoms を組み合わせた小単位 |
| organisms | ドメイン的に意味のあるUIブロック |
| templates | 配置（レイアウト）のみ。状態を持たない |
| pages | 状態・副作用（データ取得）を担い、organisms を template に流し込む |

- 配色などのデザイントークンは `theme.ts` に集約し、atoms から参照する。
- `App.tsx` は `DashboardPage` を描画するだけの薄いシェル。

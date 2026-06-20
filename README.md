# nikkei-trade-journal

日経225先物のチャート上で、ポジションの IN / OUT やコメントを記録・振り返りするためのトレードジャーナルアプリ。

## ドキュメント

[docs/](docs/README.md) を起点に、抽象度ごとに3階層へ分けています。

- [要求定義](docs/requirements/README.md) — 何を実現したいか（機能・セクション単位に分割）
- [要件定義](docs/spec/README.md) — システムが満たすべき振る舞いの仕様（ユースケース）
- [設計書](docs/design/README.md) — どう実現するか（技術選定・アーキテクチャ）

## ステータス

要求定義フェーズ（ドラフト）

---

## 環境構築

### 前提条件

| ツール | バージョン |
|---|---|
| Go | 1.24 以上 |
| Node.js | 20 以上 |
| Docker | 起動済みであること |

### 手順

**1. リポジトリをクローン**

```bash
git clone <repo-url>
cd nikkei-trade-journal
```

**2. MySQL を起動**

```bash
make up
```

Docker で MySQL 8.4 コンテナを起動し、`db/init/001_schema.sql` のスキーマが自動適用されます。

**3. Go バックエンドを起動**（ターミナル 1）

```bash
make server
# → http://localhost:8080
```

**4. フロントエンドを起動**（ターミナル 2）

```bash
make frontend
# → http://localhost:5173
```

**5. テストデータを投入**（初回のみ）

```bash
make seed
```

限月 `2506` のダミー日足データ（約 44 本）が DB に保存されます。

ブラウザで http://localhost:5173 を開くとチャートとテーブルが表示されます。

### 停止

```bash
make down   # MySQL コンテナを停止
# Go / Vite は各ターミナルで Ctrl-C
```

### 環境変数

デフォルト値のまま動作しますが、必要に応じて変更できます。

| 変数 | デフォルト | 用途 |
|---|---|---|
| `DB_DSN` | `app:app@tcp(localhost:3306)/nikkei_trade?parseTime=true&loc=Asia%2FTokyo` | Go からの DB 接続文字列 |
| `ADDR` | `:8080` | Go サーバーのリッスンアドレス |

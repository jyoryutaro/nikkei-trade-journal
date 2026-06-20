# nikkei-trade-journal

日経225先物のチャート上で、ポジションの IN / OUT やコメントを記録・振り返りするためのトレードジャーナルアプリ。

## ドキュメント

[docs/](docs/README.md) を起点に、抽象度ごとに3階層へ分けています。

- [要求定義](docs/requirements/README.md) — 何を実現したいか（機能・セクション単位に分割）
- [要件定義](docs/spec/README.md) — システムが満たすべき振る舞いの仕様（ユースケース）
- [設計書](docs/design/README.md) — どう実現するか（技術選定・アーキテクチャ）

---

## アーキテクチャ概要

```
フロントエンド (React/Vite)
    │
    │ GET /api/market-data?contract=&timeframe=
    ▼
バックエンド (Go)
    ├── GET /api/market-data        ─→ MySQL (ローソク足読み取り・集計)
    └── POST /api/market-data/fetch ─→ Yahoo Finance API (取得・保存)
                 ↑
         InternalOnly + RateLimiter
         (外部からのアクセス不可)
```

### API エンドポイント

| メソッド | パス | 説明 | 保護 |
|---|---|---|---|
| `GET` | `/api/market-data` | DB からローソク足を取得してフロントに返す | なし |
| `POST` | `/api/market-data/fetch` | Yahoo Finance から取得して DB に保存する | `Authorization: Bearer <INTERNAL_SECRET>` + レート制限 (10 req/min) |

#### GET /api/market-data

| クエリパラメータ | 必須 | 説明 |
|---|---|---|
| `contract` | 任意 | DB 上の銘柄コード（例: `^N225`）。省略時は全件 |
| `timeframe` | 任意 | 集計足種 (`1m` / `5m` / `30m` / `1h` / `1d`)。省略時は `1m` |

レスポンス例:
```json
[
  { "contract": "^N225", "timeframe": "5m", "time": 1750000000, "open": 38000, "high": 38100, "low": 37900, "close": 38050, "volume": 0 }
]
```

#### POST /api/market-data/fetch

| クエリパラメータ | 必須 | 説明 |
|---|---|---|
| `symbol` | 必須 | Yahoo Finance ティッカー（例: `^N225`） |

ヘッダー: `Authorization: Bearer <INTERNAL_SECRET>`

レスポンス例:
```json
{ "saved": 1342 }
```

呼び出し例:
```bash
curl -X POST "http://localhost:8080/api/market-data/fetch?symbol=%5EN225" \
  -H "Authorization: Bearer dev-secret"
```

---

## 環境構築

### 前提条件

| ツール | バージョン |
|---|---|
| Go | 1.24 以上 |
| Node.js | **20.19+ または 22.12+**（推奨: 22 LTS。Vite 8 の要件） |
| Docker ランタイム | Docker Desktop または [colima](https://github.com/abiosoft/colima) 等（`docker` / `docker compose` が使えること） |

> nvm を使う場合は `frontend/.nvmrc`（`22`）があるので、`frontend` ディレクトリで `nvm use` すれば対応バージョンに切り替わります。

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
# INTERNAL_SECRET を明示する場合: INTERNAL_SECRET=my-secret make server
```

**4. フロントエンドを起動**（ターミナル 2）

```bash
nvm use            # frontend/.nvmrc により Node 22 へ（nvm 利用時）
make frontend
# → http://localhost:5173
```

**5. データを取得**（初回・更新時）

```bash
curl -X POST "http://localhost:8080/api/market-data/fetch?symbol=%5EN225" \
  -H "Authorization: Bearer dev-secret"
# → {"saved": 1342} のように保存件数が返る
```

ブラウザで http://localhost:5173 を開くとチャートとテーブルが表示されます。

> 指数だけでなく先物も取得できます（`symbol` を変えるだけ）。例: `NKD=F`（CME 日経先物）は `symbol=NKD%3DF`。先物は限月コード（例 `2609`）で保存されます。

### 定期取得（ローカル・launchd で15分毎）

`make server` が動いている前提で、[`scripts/fetch-market-data.sh`](scripts/fetch-market-data.sh) を launchd で定期実行できます（既定で `^N225` と `NKD=F` を取得）。

```bash
# 1) テンプレートをコピーし、ABSOLUTE_REPO_PATH と YOUR_HOME を自分の値に置換
cp scripts/com.nikkei-trade-journal.fetch.plist.example \
   ~/Library/LaunchAgents/com.nikkei-trade-journal.fetch.plist

# 2) 登録（15分毎 + 登録時に1回実行）
launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/com.nikkei-trade-journal.fetch.plist

# 解除
launchctl bootout gui/$(id -u)/com.nikkei-trade-journal.fetch

# ログ
tail -f ~/Library/Logs/nikkei-trade-journal-fetch.log
```

取得対象や接続先は plist の `EnvironmentVariables`（`SYMBOLS` / `BASE_URL` / `INTERNAL_SECRET`）で変更できます。スクリプト単体実行も可: `SYMBOLS="^N225" ./scripts/fetch-market-data.sh`。

> 注意: これは**ローカルのスケジューラ**です。Mac とバックエンドが起動している間だけ動きます。クラウドから定期実行したい場合は、バックエンドを公開デプロイした上で GitHub Actions 等の cron に切り替えてください。

### 停止

```bash
make down   # MySQL コンテナを停止
# Go / Vite は各ターミナルで Ctrl-C
```

### 環境変数

| 変数 | デフォルト | 用途 |
|---|---|---|
| `DB_DSN` | `app:app@tcp(localhost:3306)/nikkei_trade?parseTime=true&loc=Asia%2FTokyo` | Go からの DB 接続文字列 |
| `ADDR` | `:8080` | Go サーバーのリッスンアドレス |
| `INTERNAL_SECRET` | `dev-secret`（`make server` 経由時） | `/api/market-data/fetch` の認証トークン。**本番では必ず変更すること** |

#### フロントエンド（Vite）

| 変数 | デフォルト | 用途 |
|---|---|---|
| `VITE_API_BASE` | `http://localhost:8080` | バックエンド API のベース URL（末尾スラッシュ無し）。デプロイ時は公開 URL を設定。`frontend/.env.example` をコピーして `frontend/.env.local` に記載 |

### make コマンド一覧

| コマンド | 説明 |
|---|---|
| `make up` | MySQL コンテナ起動 |
| `make down` | MySQL コンテナ停止 |
| `make server` | Go バックエンド起動 |
| `make frontend` | フロントエンド開発サーバー起動 |
| `make db-console` | MySQL コンソールに接続 |
| `make test` | バックエンドテスト実行 |

---

## トラブルシューティング

**`INTERNAL_SECRET env var is required` でサーバーが起動しない**

`make server` を使うと `dev-secret` がデフォルトで設定されます。直接 `go run` する場合は環境変数を明示してください。

```bash
INTERNAL_SECRET=dev-secret go run ./cmd/server
```

**`Access denied for user 'app'@'localhost'`**

ホストの 3306 をローカルの MySQL（Homebrew 等）が占有していて、コンテナの MySQL に届いていない可能性があります。

```bash
lsof -nP -iTCP:3306 -sTCP:LISTEN   # 何が 3306 を bind しているか確認
brew services stop mysql           # ローカル MySQL を停止（コンテナを使う場合）
```

**`Cannot find native binding` / `@rolldown/binding-*`（make frontend）**

```bash
cd frontend && rm -rf node_modules package-lock.json && npm install
```

**`TypeError [ERR_INVALID_ARG_VALUE] ... styleText`（make frontend）**

Node が古い（20.19 未満）と Vite 8 / rolldown が起動できません。Node を 20.19+ / 22.12+ に上げてください（`nvm install 22 && nvm use 22`）。

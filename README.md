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
| Node.js | **20.19+ または 22.12+**（推奨: 22 LTS。Vite 8 の要件） |
| Docker ランタイム | Docker Desktop または [colima](https://github.com/abiosoft/colima) 等（`docker` / `docker compose` が使えること） |

> nvm を使う場合は `frontend/.nvmrc`（`22`）があるので、`frontend` ディレクトリで `nvm use` すれば対応バージョンに切り替わります。古い Node（例: 20.12）では Vite 8 が起動しません。

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
nvm use            # frontend/.nvmrc により Node 22 へ（nvm 利用時）
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

### トラブルシューティング

**`Access denied for user 'app'@'localhost'`（make seed / make server）**

ホストの 3306 をローカルの MySQL（Homebrew 等）が占有していて、コンテナの MySQL に届いていない可能性があります。`127.0.0.1:3306` への接続は、より具体的に bind しているローカル MySQL に吸われます。

```bash
lsof -nP -iTCP:3306 -sTCP:LISTEN   # 何が 3306 を bind しているか確認
brew services stop mysql           # ローカル MySQL を停止（コンテナを使う場合）
```

**`Cannot find native binding` / `@rolldown/binding-*`（make frontend）**

npm の optional dependencies のバグで Vite 8 のネイティブバインディングが入らないことがあります。クリーン再インストールで解消します。

```bash
cd frontend && rm -rf node_modules package-lock.json && npm install
```

**`TypeError [ERR_INVALID_ARG_VALUE] ... styleText`（make frontend）**

Node が古い（20.19 未満）と Vite 8 / rolldown が起動できません。Node を 20.19+ / 22.12+ に上げてください（`nvm install 22 && nvm use 22`）。

# 要求定義書

> 本書は機能・セクションごとにファイルを分割しています。**複数人で同時編集してもコンフリクトしにくい**よう、自分の担当ファイルだけを編集してください。
> 未確定の項目は `TBD`（To Be Determined）と記載し、決定後に更新します。

## 文書情報

| 項目 | 内容 |
| --- | --- |
| プロジェクト名 | nikkei-trade-journal |
| 文書バージョン | 0.9 |
| 作成日 | 2026-06-20 |
| 作成者 | 城竜太郎 |
| ステータス | ドラフト |

## 目次

| # | セクション | ファイル |
| --- | --- | --- |
| 1 | 概要 / 背景 | [01-overview.md](01-overview.md) |
| 2 | スコープ | [02-scope.md](02-scope.md) |
| 3 | 利用者 / ペルソナ | [03-personas.md](03-personas.md) |
| 4 | 機能要件（FR） | [04-functional/README.md](04-functional/README.md) |
| 5 | 非機能要件（NFR） | [05-non-functional.md](05-non-functional.md) |
| 6 | 画面 / 画面遷移 | [06-screens.md](06-screens.md) |
| 7 | データ要件 | [07-data-model.md](07-data-model.md) |
| 8 | 外部連携 / 制約 | [08-integrations.md](08-integrations.md) |
| 9 | 前提 / 制約条件 | [09-assumptions.md](09-assumptions.md) |
| 10 | リリース計画 / マイルストーン | [10-roadmap.md](10-roadmap.md) |
| 11 | 受け入れ基準 | [11-acceptance.md](11-acceptance.md) |
| 12 | 未決事項 / 課題 | [12-open-issues.md](12-open-issues.md) |
| - | 改訂履歴 | [CHANGELOG.md](CHANGELOG.md) |

## 編集ルール（コンフリクト回避）

- **1ファイル＝1セクション / 1機能**。原則、自分の担当ファイルのみ編集する。
- 機能を追加するときは `04-functional/FR-XX-<name>.md` を新規作成し、[04-functional/README.md](04-functional/README.md) の一覧に**末尾追記**する（表の途中行を書き換えない）。
- 一覧・目次など共有ファイルへの追記は**末尾に1行追加**を基本とし、既存行の並べ替えを避ける。
- 文書全体に影響する変更をしたら [CHANGELOG.md](CHANGELOG.md) に**末尾追記**する。

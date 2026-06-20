# FR-00 認証

| 項目 | 内容 |
| --- | --- |
| ID | FR-00 |
| 機能名 | 認証 |
| 優先度 | 将来課題（当面スコープ外） |
| 担当 | <担当者> |
| ステータス | 保留 |

## 概要

当面はローカル単一ユーザーでの利用を前提とし、**認証は導入しない**。マルチユーザー化する場合に改めて認証を設計する。

> 旧方針（Firebase Authentication による Google ログイン）は、技術構成を Go + MySQL に変更したため撤回（[12-open-issues.md](../12-open-issues.md) Q-03 / Q-07）。

## 将来導入する場合の方針（参考）

- 認証方式は別途決定する（例: Google OAuth 2.0 を Go で自前実装し、users テーブルを MySQL に保持）。
- 未ログイン時の記録機能の扱い、ログイン / ログアウト導線を設計する。

## 関連

- 非機能: [05-non-functional.md](../05-non-functional.md)（セキュリティ）
- データ: [07-data-model.md](../07-data-model.md)（User: 将来）
- 機能: [FR-07-data-isolation.md](FR-07-data-isolation.md)（データ分離: 将来）
- 課題: [12-open-issues.md](../12-open-issues.md) Q-03

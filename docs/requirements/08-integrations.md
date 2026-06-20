# 8. 外部連携 / 制約

## 外部API・サービス

- **Yahoo Finance chart API（日経225先物の価格取得）** — `NKD=F` の OHLCV を取得（[要件定義 UC-01](../spec/use-cases/UC-01-fetch-market-data.md)）。利用規約・レート制限に注意。intraday（特に 1m）は取得可能期間に上限がある（[12-open-issues.md](12-open-issues.md) Q-01）。

## 技術的制約

- 自前構成の Web アプリ: フロントエンド（React + TypeScript + Vite）＋ バックエンド（Go）＋ DB（MySQL）。技術選定は[設計書](../design/01-tech-stack.md)を参照。
- 当面はローカル環境で稼働（Docker で MySQL を起動）。本番ホスティング先は TBD。

## 法務・コンプライアンス

- 投資助言に当たらないよう「記録ツール」に限定する。
- 認証導入（マルチユーザー化）で個人情報を扱う場合はプライバシーポリシーの整備が必要（[12-open-issues.md](12-open-issues.md) Q-04）。

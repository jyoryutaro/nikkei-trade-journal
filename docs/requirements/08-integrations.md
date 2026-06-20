# 8. 外部連携 / 制約

## 外部API・サービス

- **Firebase（Authentication / Firestore / Hosting）** — 認証・データ保存・ホスティング基盤。Firebase プロジェクトの作成と Google ログインプロバイダの有効化が必要。
- **株価データAPI（日経225先物の価格取得）** — 利用規約・レート制限・費用に注意。先物は限月ごとに銘柄が分かれる点に留意。提供元は TBD（[12-open-issues.md](12-open-issues.md) Q-01）。

## 技術的制約

- Webアプリ（フロントエンド + Firebase）。フロントエンドの使用言語・フレームワークは TBD。
- Firebase 無料枠（Spark）の上限に留意。

## 法務・コンプライアンス

- 投資助言に当たらないよう「記録ツール」に限定する。
- 個人情報（メールアドレス等）を扱うためプライバシーポリシーの整備が必要（[12-open-issues.md](12-open-issues.md) Q-04）。

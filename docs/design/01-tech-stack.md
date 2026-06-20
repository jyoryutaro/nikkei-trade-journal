# 技術選定

> 設計書。要件（[../spec/](../spec/README.md)）を実現するための技術選定。

| 層 | 技術 |
|---|---|
| バックエンド | Go |
| フロントエンド | React + TypeScript + Vite |
| チャート | Lightweight Charts |
| DB | MySQL（ローカル） |

> ⚠️ **要整合**: ここでの選定（Go + MySQL ローカル）は、要求定義で確定済みの **Firebase（Authentication / Firestore / Hosting）** と矛盾する。認証・データ保存・ホスティングの方針を要求定義側と突き合わせて整合させる必要がある（[../requirements/12-open-issues.md](../requirements/12-open-issues.md) Q-07）。

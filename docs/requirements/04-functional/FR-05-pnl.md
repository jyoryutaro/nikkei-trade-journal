# FR-05 損益表示

| 項目 | 内容 |
| --- | --- |
| ID | FR-05 |
| 機能名 | 損益表示 |
| 優先度 | Should |
| 担当 | <担当者> |
| ステータス | ドラフト |

## 概要

IN / OUT から損益（pips / 金額 / %）を自動計算して表示する。

## 詳細仕様 / 振る舞い

- 方向（long / short）・数量・IN価格・OUT価格から損益を計算する。
- 表示単位（pips / 金額 / %）は TBD。
- OUT が未登録（建玉中）の場合は含み損益として扱う（要否は TBD）。

## 受け入れ基準

- [ ] IN / OUT 価格から損益が正しく計算される。
- [ ] long / short で符号が正しい。

## 関連

- データ: [07-data-model.md](../07-data-model.md)（Trade）
- 関連機能: [FR-02-positions.md](FR-02-positions.md)

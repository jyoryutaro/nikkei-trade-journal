# 7. データ要件

主なデータ項目（エンティティ）:

| エンティティ | 主な属性 |
| --- | --- |
| User（ユーザー） | uid(Firebase UID), メールアドレス, 表示名, アバターURL, 作成日時 |
| Trade（トレード） | id, userId, 限月, 方向(long/short), 枚数, IN日時, IN価格, OUT日時, OUT価格, 損益, ステータス |
| Comment（コメント） | id, userId, 対象(trade/時点), 日時, 本文, タグ |
| MarketData（相場） | 限月(または連続足区分), 日時, 始値, 高値, 安値, 終値, 出来高 |

- 対象は日経225先物。先物は限月（SQ）があるため、Trade は限月を保持し、相場データは限月別または連続足として扱う（[12-open-issues.md](12-open-issues.md) Q-05）。
- 数量の単位は「枚（contract）」とする。
- データソース: <日経225先物の価格データの取得元。例：API名 / 手入力 / CSV取込。TBD（[12-open-issues.md](12-open-issues.md) Q-01）>

## Firestore コレクション構成（案）

```
users/{uid}
  trades/{tradeId}
    comments/{commentId}
```

- ユーザーごとのデータはサブコレクションとして保持し、セキュリティルールで `request.auth.uid == uid` のみ許可する。
- 相場データ（MarketData）はユーザー横断の共有データ。保存要否・取得方法は [12-open-issues.md](12-open-issues.md) Q-01 参照。

-- Journal entries: a record attached to a point in time, which is either a
-- position record (side + trade type + price) or a comment-only note.
CREATE TABLE IF NOT EXISTS journal_entries (
  id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  contract    VARCHAR(10)   NOT NULL COMMENT '限月 e.g. 2609',
  ts          DATETIME      NOT NULL COMMENT '対象時点 (UTC)',
  side        VARCHAR(5)    NULL     COMMENT 'long | short (NULL = コメントのみ)',
  trade_type  VARCHAR(5)    NULL     COMMENT 'open(新規) | close(決済)',
  price       DECIMAL(12,2) NULL     COMMENT '約定金額',
  comment     TEXT          NULL,
  created_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_contract_ts (contract, ts)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

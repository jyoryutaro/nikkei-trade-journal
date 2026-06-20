CREATE TABLE IF NOT EXISTS market_data (
  id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  contract    VARCHAR(10)  NOT NULL COMMENT '限月 e.g. 2506',
  ts          DATETIME     NOT NULL COMMENT '足の日時 (JST)',
  open        DECIMAL(10,1) NOT NULL,
  high        DECIMAL(10,1) NOT NULL,
  low         DECIMAL(10,1) NOT NULL,
  close       DECIMAL(10,1) NOT NULL,
  volume      BIGINT UNSIGNED NOT NULL DEFAULT 0,
  created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uq_contract_ts (contract, ts)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

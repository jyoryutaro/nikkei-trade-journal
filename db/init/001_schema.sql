CREATE TABLE IF NOT EXISTS market_data (
  id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  contract    VARCHAR(10)   NOT NULL COMMENT '限月 e.g. 2609',
  timeframe   VARCHAR(5)    NOT NULL DEFAULT '1m' COMMENT '足種 e.g. 1m, 5m, 30m, 1h, 1d',
  ts          DATETIME      NOT NULL COMMENT '足の開始日時 (UTC)',
  open        DECIMAL(12,2) NOT NULL,
  high        DECIMAL(12,2) NOT NULL,
  low         DECIMAL(12,2) NOT NULL,
  close       DECIMAL(12,2) NOT NULL,
  volume      BIGINT UNSIGNED NOT NULL DEFAULT 0,
  created_at  DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY uq_contract_tf_ts (contract, timeframe, ts)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

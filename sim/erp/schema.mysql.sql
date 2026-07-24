-- ============================================================================
-- MindSet fake ERP — schema (MySQL 8 / MariaDB 10 syntax)
-- Mirrors a small but realistic mid-sized manufacturing ERP.
-- Matches the canonical MindSet model (docs/mysql_connector.md §6b).
-- ============================================================================

CREATE TABLE IF NOT EXISTS products (
  product_code   VARCHAR(64) PRIMARY KEY,
  name           VARCHAR(255) NOT NULL,
  target_rate    INT NULL,                          -- units per hour target
  recipe_id      VARCHAR(64) NULL,
  hourly_margin  DECIMAL(10,2) NULL                 -- EUR margin per hour of runtime
);

CREATE TABLE IF NOT EXISTS operators (
  operator_id    VARCHAR(32) PRIMARY KEY,
  name           VARCHAR(255) NOT NULL,
  shift          ENUM('MORNING','AFTERNOON','NIGHT') NOT NULL
);

CREATE TABLE IF NOT EXISTS work_orders (
  of_number      VARCHAR(64) PRIMARY KEY,
  product_code   VARCHAR(64) NOT NULL,
  work_center    VARCHAR(64) NOT NULL,
  planned_qty    INT NOT NULL,
  actual_qty     INT NOT NULL DEFAULT 0,
  status         ENUM('PLANNED','RUNNING','DONE','ABORTED') NOT NULL DEFAULT 'PLANNED',
  started_at     DATETIME NULL,
  finished_at    DATETIME NULL,
  operator_id    VARCHAR(32) NULL,
  due_date       DATE NULL,
  customer_id    VARCHAR(64) NULL,
  INDEX idx_wc_status (work_center, status)
);

CREATE TABLE IF NOT EXISTS batches (
  batch_id       VARCHAR(64) PRIMARY KEY,
  of_number      VARCHAR(64) NOT NULL,
  started_at     DATETIME NOT NULL,
  finished_at    DATETIME NULL,
  quality_status ENUM('PASS','FAIL','REWORK') NULL,
  FOREIGN KEY (of_number) REFERENCES work_orders(of_number),
  INDEX idx_batches_of (of_number)
);

CREATE TABLE IF NOT EXISTS schedules (
  id             INT AUTO_INCREMENT PRIMARY KEY,
  work_center    VARCHAR(64) NOT NULL,
  of_number      VARCHAR(64) NOT NULL,
  planned_start  DATETIME NOT NULL,
  planned_end    DATETIME NOT NULL,
  FOREIGN KEY (of_number) REFERENCES work_orders(of_number),
  INDEX idx_sched_wc_start (work_center, planned_start)
);

CREATE TABLE IF NOT EXISTS quality_results (
  id             INT AUTO_INCREMENT PRIMARY KEY,
  batch_id       VARCHAR(64) NOT NULL,
  measured_at    DATETIME NOT NULL,
  metric         VARCHAR(64) NOT NULL,
  value          DECIMAL(10,4) NOT NULL,
  spec_min       DECIMAL(10,4) NULL,
  spec_max       DECIMAL(10,4) NULL,
  FOREIGN KEY (batch_id) REFERENCES batches(batch_id),
  INDEX idx_quality_batch (batch_id, measured_at)
);

-- ============================================================================
-- MindSet fake ERP — grants
-- Runs LAST (03-*), after schema.mysql.sql and seed.mysql.sql.
--
-- We DO NOT rely on MYSQL_USER/MYSQL_PASSWORD in docker-compose because that
-- flow auto-grants ALL PRIVILEGES on the target database — REVOKE ALL then
-- fails on MySQL 8.4 because of stricter grant matching. Instead we create
-- both users explicitly here with the exact privileges we want.
--
--  * mindset_readonly — used by the MindSet SQL connector. SELECT ONLY.
--  * mindset_writer   — used by cmd/erpsim only, to simulate ERP activity.
-- ============================================================================

CREATE USER IF NOT EXISTS 'mindset_readonly'@'%' IDENTIFIED BY 'readonly_dev';
GRANT  SELECT ON fake_erp.* TO 'mindset_readonly'@'%';

CREATE USER IF NOT EXISTS 'mindset_writer'@'%'   IDENTIFIED BY 'writer_dev';
GRANT  SELECT, INSERT, UPDATE ON fake_erp.* TO 'mindset_writer'@'%';

FLUSH PRIVILEGES;

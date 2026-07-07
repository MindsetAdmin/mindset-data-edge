-- ============================================================================
-- MindSet fake ERP — seed data
-- Relative timestamps (NOW() - INTERVAL ...) so the seed stays realistic
-- no matter when the container is first brought up.
-- Work-centers "machine1/machine2/machine3" match the OPC-UA simulator.
-- ============================================================================

-- ----- Products (mix of agrifood, pharma, cosmetics) --------------------------
INSERT INTO products (product_code, name, target_rate, recipe_id, hourly_margin) VALUES
('PROD-A01', 'Yaourt nature 125g',        1200, 'REC-Y-001',   45.00),
('PROD-A02', 'Yaourt fraise 125g',        1200, 'REC-Y-002',   48.00),
('PROD-A03', 'Yaourt vanille 125g',       1150, 'REC-Y-003',   47.00),
('PROD-A04', 'Fromage blanc 200g',         800, 'REC-F-001',   62.00),
('PROD-A05', 'Creme dessert chocolat',     900, 'REC-C-001',   55.00),
('PROD-A06', 'Compote pomme 100g',        1400, 'REC-COM-01',  32.00),
('PROD-A07', 'Compote fraise 100g',       1350, 'REC-COM-02',  34.00),
('PROD-A08', 'Boisson lactee framboise',  1000, 'REC-BL-01',   51.00),
('PROD-B01', 'Sirop antitussif 100ml',     600, 'REC-PH-001', 180.00),
('PROD-B02', 'Solution buvable 50ml',      800, 'REC-PH-002', 210.00),
('PROD-C01', 'Creme hydratante 50ml',      400, 'REC-CO-001', 125.00),
('PROD-C02', 'Serum jour 30ml',            350, 'REC-CO-002', 175.00);

-- ----- Operators (8 across 3 shifts) -----------------------------------------
INSERT INTO operators (operator_id, name, shift) VALUES
('OP-001', 'Alice Martin',    'MORNING'),
('OP-002', 'Bertrand Leroy',  'MORNING'),
('OP-003', 'Camille Dubois',  'MORNING'),
('OP-004', 'David Lefevre',   'AFTERNOON'),
('OP-005', 'Emma Robert',     'AFTERNOON'),
('OP-006', 'Franck Petit',    'AFTERNOON'),
('OP-007', 'Gabrielle Simon', 'NIGHT'),
('OP-008', 'Hugo Bernard',    'NIGHT');

-- ----- Work Orders — historical (last 30 days, mostly DONE, some ABORTED) ---
INSERT INTO work_orders (of_number, product_code, work_center, planned_qty, actual_qty, status, started_at, finished_at, operator_id) VALUES
('WO-2026-8001', 'PROD-A01', 'machine1', 5000, 4987, 'DONE',    NOW() - INTERVAL 30 DAY, NOW() - INTERVAL 30 DAY + INTERVAL 4 HOUR, 'OP-001'),
('WO-2026-8002', 'PROD-A02', 'machine1', 4800, 4772, 'DONE',    NOW() - INTERVAL 28 DAY, NOW() - INTERVAL 28 DAY + INTERVAL 4 HOUR, 'OP-002'),
('WO-2026-8003', 'PROD-A03', 'machine1', 4600, 4601, 'DONE',    NOW() - INTERVAL 27 DAY, NOW() - INTERVAL 27 DAY + INTERVAL 4 HOUR, 'OP-004'),
('WO-2026-8004', 'PROD-A04', 'machine2', 3200, 3175, 'DONE',    NOW() - INTERVAL 26 DAY, NOW() - INTERVAL 26 DAY + INTERVAL 4 HOUR, 'OP-005'),
('WO-2026-8005', 'PROD-B01', 'machine2', 2400, 1200, 'ABORTED', NOW() - INTERVAL 25 DAY, NOW() - INTERVAL 25 DAY + INTERVAL 2 HOUR, 'OP-005'),
('WO-2026-8006', 'PROD-A05', 'machine3', 3600, 3542, 'DONE',    NOW() - INTERVAL 24 DAY, NOW() - INTERVAL 24 DAY + INTERVAL 4 HOUR, 'OP-006'),
('WO-2026-8007', 'PROD-A06', 'machine3', 5600, 5588, 'DONE',    NOW() - INTERVAL 23 DAY, NOW() - INTERVAL 23 DAY + INTERVAL 4 HOUR, 'OP-006'),
('WO-2026-8008', 'PROD-A01', 'machine1', 5000, 4972, 'DONE',    NOW() - INTERVAL 22 DAY, NOW() - INTERVAL 22 DAY + INTERVAL 4 HOUR, 'OP-001'),
('WO-2026-8009', 'PROD-C01', 'machine2', 1600, 1591, 'DONE',    NOW() - INTERVAL 21 DAY, NOW() - INTERVAL 21 DAY + INTERVAL 4 HOUR, 'OP-004'),
('WO-2026-8010', 'PROD-A07', 'machine3', 5400, 5387, 'DONE',    NOW() - INTERVAL 20 DAY, NOW() - INTERVAL 20 DAY + INTERVAL 4 HOUR, 'OP-006'),
('WO-2026-8011', 'PROD-A02', 'machine1', 4800, 4788, 'DONE',    NOW() - INTERVAL 19 DAY, NOW() - INTERVAL 19 DAY + INTERVAL 4 HOUR, 'OP-002'),
('WO-2026-8012', 'PROD-B02', 'machine2', 3200, 3195, 'DONE',    NOW() - INTERVAL 18 DAY, NOW() - INTERVAL 18 DAY + INTERVAL 4 HOUR, 'OP-004'),
('WO-2026-8013', 'PROD-A08', 'machine3', 4000, 3956, 'DONE',    NOW() - INTERVAL 17 DAY, NOW() - INTERVAL 17 DAY + INTERVAL 4 HOUR, 'OP-005'),
('WO-2026-8014', 'PROD-A03', 'machine1', 4600, 4588, 'DONE',    NOW() - INTERVAL 16 DAY, NOW() - INTERVAL 16 DAY + INTERVAL 4 HOUR, 'OP-002'),
('WO-2026-8015', 'PROD-C02', 'machine2', 1400, 1382, 'DONE',    NOW() - INTERVAL 15 DAY, NOW() - INTERVAL 15 DAY + INTERVAL 4 HOUR, 'OP-005'),
('WO-2026-8016', 'PROD-A04', 'machine2', 3200, 2100, 'ABORTED', NOW() - INTERVAL 14 DAY, NOW() - INTERVAL 14 DAY + INTERVAL 2 HOUR, 'OP-004'),
('WO-2026-8017', 'PROD-A05', 'machine3', 3600, 3577, 'DONE',    NOW() - INTERVAL 13 DAY, NOW() - INTERVAL 13 DAY + INTERVAL 4 HOUR, 'OP-006'),
('WO-2026-8018', 'PROD-A06', 'machine3', 5600, 5596, 'DONE',    NOW() - INTERVAL 12 DAY, NOW() - INTERVAL 12 DAY + INTERVAL 4 HOUR, 'OP-007'),
('WO-2026-8019', 'PROD-A01', 'machine1', 5000, 4977, 'DONE',    NOW() - INTERVAL 11 DAY, NOW() - INTERVAL 11 DAY + INTERVAL 4 HOUR, 'OP-003'),
('WO-2026-8020', 'PROD-B01', 'machine2', 2400, 2388, 'DONE',    NOW() - INTERVAL 10 DAY, NOW() - INTERVAL 10 DAY + INTERVAL 4 HOUR, 'OP-004'),
('WO-2026-8021', 'PROD-A02', 'machine1', 4800, 4782, 'DONE',    NOW() -  INTERVAL 9 DAY, NOW() -  INTERVAL 9 DAY + INTERVAL 4 HOUR, 'OP-001'),
('WO-2026-8022', 'PROD-A07', 'machine3', 5400, 5391, 'DONE',    NOW() -  INTERVAL 8 DAY, NOW() -  INTERVAL 8 DAY + INTERVAL 4 HOUR, 'OP-005'),
('WO-2026-8023', 'PROD-A08', 'machine3', 4000, 3971, 'DONE',    NOW() -  INTERVAL 7 DAY, NOW() -  INTERVAL 7 DAY + INTERVAL 4 HOUR, 'OP-006'),
('WO-2026-8024', 'PROD-C01', 'machine2', 1600, 1594, 'DONE',    NOW() -  INTERVAL 6 DAY, NOW() -  INTERVAL 6 DAY + INTERVAL 4 HOUR, 'OP-004'),
('WO-2026-8025', 'PROD-A03', 'machine1', 4600, 4592, 'DONE',    NOW() -  INTERVAL 5 DAY, NOW() -  INTERVAL 5 DAY + INTERVAL 4 HOUR, 'OP-002'),
('WO-2026-8026', 'PROD-A05', 'machine3', 3600, 3568, 'DONE',    NOW() -  INTERVAL 4 DAY, NOW() -  INTERVAL 4 DAY + INTERVAL 4 HOUR, 'OP-007'),
('WO-2026-8027', 'PROD-B02', 'machine2', 3200, 3191, 'DONE',    NOW() -  INTERVAL 3 DAY, NOW() -  INTERVAL 3 DAY + INTERVAL 4 HOUR, 'OP-005'),
('WO-2026-8028', 'PROD-A04', 'machine2', 3200, 3182, 'DONE',    NOW() -  INTERVAL 2 DAY, NOW() -  INTERVAL 2 DAY + INTERVAL 4 HOUR, 'OP-004'),
('WO-2026-8029', 'PROD-C02', 'machine2', 1400, 1394, 'DONE',    NOW() -  INTERVAL 2 DAY, NOW() -  INTERVAL 2 DAY + INTERVAL 3 HOUR, 'OP-005'),
('WO-2026-8030', 'PROD-A01', 'machine1', 5000, 4991, 'DONE',    NOW() -  INTERVAL 1 DAY, NOW() -  INTERVAL 1 DAY + INTERVAL 4 HOUR, 'OP-001');

-- ----- Work Orders — currently RUNNING (one per work_center) ----------------
INSERT INTO work_orders (of_number, product_code, work_center, planned_qty, actual_qty, status, started_at, finished_at, operator_id) VALUES
('WO-2026-9001', 'PROD-A02', 'machine1', 5000, 2100, 'RUNNING', NOW() - INTERVAL 3 HOUR,     NULL, 'OP-002'),
('WO-2026-9002', 'PROD-B01', 'machine2', 3000,  800, 'RUNNING', NOW() - INTERVAL 1 HOUR,     NULL, 'OP-004'),
('WO-2026-9003', 'PROD-C01', 'machine3', 2000,  450, 'RUNNING', NOW() - INTERVAL 30 MINUTE,  NULL, 'OP-006');

-- ----- Work Orders — PLANNED (near-future) ----------------------------------
INSERT INTO work_orders (of_number, product_code, work_center, planned_qty, actual_qty, status, started_at, finished_at, operator_id) VALUES
('WO-2026-9004', 'PROD-A03', 'machine1', 4600, 0, 'PLANNED', NULL, NULL, NULL),
('WO-2026-9005', 'PROD-B02', 'machine2', 3200, 0, 'PLANNED', NULL, NULL, NULL),
('WO-2026-9006', 'PROD-A05', 'machine3', 3600, 0, 'PLANNED', NULL, NULL, NULL),
('WO-2026-9007', 'PROD-A01', 'machine1', 5000, 0, 'PLANNED', NULL, NULL, NULL),
('WO-2026-9008', 'PROD-C02', 'machine2', 1400, 0, 'PLANNED', NULL, NULL, NULL);

-- ----- Batches — closed (linked to DONE work orders) ------------------------
INSERT INTO batches (batch_id, of_number, started_at, finished_at, quality_status) VALUES
('B-2026-5001', 'WO-2026-8001', NOW() - INTERVAL 30 DAY, NOW() - INTERVAL 30 DAY + INTERVAL 4 HOUR, 'PASS'),
('B-2026-5002', 'WO-2026-8005', NOW() - INTERVAL 25 DAY, NOW() - INTERVAL 25 DAY + INTERVAL 2 HOUR, 'FAIL'),
('B-2026-5003', 'WO-2026-8007', NOW() - INTERVAL 23 DAY, NOW() - INTERVAL 23 DAY + INTERVAL 4 HOUR, 'PASS'),
('B-2026-5004', 'WO-2026-8012', NOW() - INTERVAL 18 DAY, NOW() - INTERVAL 18 DAY + INTERVAL 4 HOUR, 'REWORK'),
('B-2026-5005', 'WO-2026-8016', NOW() - INTERVAL 14 DAY, NOW() - INTERVAL 14 DAY + INTERVAL 2 HOUR, 'FAIL'),
('B-2026-5006', 'WO-2026-8019', NOW() - INTERVAL 11 DAY, NOW() - INTERVAL 11 DAY + INTERVAL 4 HOUR, 'PASS'),
('B-2026-5007', 'WO-2026-8020', NOW() - INTERVAL 10 DAY, NOW() - INTERVAL 10 DAY + INTERVAL 4 HOUR, 'PASS'),
('B-2026-5008', 'WO-2026-8022', NOW() -  INTERVAL 8 DAY, NOW() -  INTERVAL 8 DAY + INTERVAL 4 HOUR, 'PASS'),
('B-2026-5009', 'WO-2026-8025', NOW() -  INTERVAL 5 DAY, NOW() -  INTERVAL 5 DAY + INTERVAL 4 HOUR, 'PASS'),
('B-2026-5010', 'WO-2026-8027', NOW() -  INTERVAL 3 DAY, NOW() -  INTERVAL 3 DAY + INTERVAL 4 HOUR, 'REWORK'),
('B-2026-5011', 'WO-2026-8028', NOW() -  INTERVAL 2 DAY, NOW() -  INTERVAL 2 DAY + INTERVAL 4 HOUR, 'PASS'),
('B-2026-5012', 'WO-2026-8030', NOW() -  INTERVAL 1 DAY, NOW() -  INTERVAL 1 DAY + INTERVAL 4 HOUR, 'PASS');

-- ----- Batches — in-flight (linked to RUNNING work orders) ------------------
INSERT INTO batches (batch_id, of_number, started_at, finished_at, quality_status) VALUES
('B-2026-6001', 'WO-2026-9001', NOW() - INTERVAL 3 HOUR,     NULL, NULL),
('B-2026-6002', 'WO-2026-9002', NOW() - INTERVAL 1 HOUR,     NULL, NULL),
('B-2026-6003', 'WO-2026-9003', NOW() - INTERVAL 30 MINUTE,  NULL, NULL);

-- ----- Quality results ------------------------------------------------------
INSERT INTO quality_results (batch_id, measured_at, metric, value, spec_min, spec_max) VALUES
-- Historical batches
('B-2026-5001', NOW() - INTERVAL 30 DAY + INTERVAL 1 HOUR, 'viscosity',   850.0, 800.0, 900.0),
('B-2026-5001', NOW() - INTERVAL 30 DAY + INTERVAL 2 HOUR, 'temperature',   4.2,   3.5,   4.5),
('B-2026-5001', NOW() - INTERVAL 30 DAY + INTERVAL 3 HOUR, 'ph',            4.4,   4.2,   4.6),
('B-2026-5002', NOW() - INTERVAL 25 DAY + INTERVAL 1 HOUR, 'viscosity',   920.0, 800.0, 900.0),   -- OUT of spec (FAIL)
('B-2026-5003', NOW() - INTERVAL 23 DAY + INTERVAL 2 HOUR, 'moisture',     22.1,  20.0,  25.0),
('B-2026-5004', NOW() - INTERVAL 18 DAY + INTERVAL 2 HOUR, 'ph',            4.7,   4.2,   4.6),   -- OUT of spec (REWORK)
('B-2026-5005', NOW() - INTERVAL 14 DAY + INTERVAL 1 HOUR, 'temperature',   5.1,   3.5,   4.5),   -- OUT of spec (FAIL)
('B-2026-5006', NOW() - INTERVAL 11 DAY + INTERVAL 2 HOUR, 'viscosity',   843.0, 800.0, 900.0),
('B-2026-5007', NOW() - INTERVAL 10 DAY + INTERVAL 2 HOUR, 'moisture',     22.5,  20.0,  25.0),
('B-2026-5008', NOW() -  INTERVAL 8 DAY + INTERVAL 2 HOUR, 'temperature',   4.1,   3.5,   4.5),
('B-2026-5009', NOW() -  INTERVAL 5 DAY + INTERVAL 2 HOUR, 'viscosity',   861.0, 800.0, 900.0),
('B-2026-5010', NOW() -  INTERVAL 3 DAY + INTERVAL 2 HOUR, 'ph',           4.65,   4.2,   4.6),   -- OUT of spec (REWORK)
('B-2026-5011', NOW() -  INTERVAL 2 DAY + INTERVAL 2 HOUR, 'moisture',     21.8,  20.0,  25.0),
('B-2026-5012', NOW() -  INTERVAL 1 DAY + INTERVAL 2 HOUR, 'viscosity',   852.0, 800.0, 900.0),
-- In-flight batches
('B-2026-6001', NOW() - INTERVAL 2 HOUR,    'temperature',  4.3, 3.5, 4.5),
('B-2026-6001', NOW() - INTERVAL 1 HOUR,    'viscosity',  855.0, 800.0, 900.0),
('B-2026-6002', NOW() - INTERVAL 40 MINUTE, 'ph',           4.5, 4.2, 4.6),
('B-2026-6003', NOW() - INTERVAL 20 MINUTE, 'temperature',  4.2, 3.5, 4.5);

-- ----- Schedules — next 24h across the 3 work_centers -----------------------
INSERT INTO schedules (work_center, of_number, planned_start, planned_end) VALUES
('machine1', 'WO-2026-9004', NOW() + INTERVAL 2 HOUR,      NOW() + INTERVAL  6 HOUR),
('machine2', 'WO-2026-9005', NOW() + INTERVAL 1 HOUR,      NOW() + INTERVAL  4 HOUR),
('machine3', 'WO-2026-9006', NOW() + INTERVAL 30 MINUTE,   NOW() + INTERVAL  3 HOUR),
('machine1', 'WO-2026-9007', NOW() + INTERVAL 8 HOUR,      NOW() + INTERVAL 12 HOUR),
('machine2', 'WO-2026-9008', NOW() + INTERVAL 6 HOUR,      NOW() + INTERVAL  9 HOUR);

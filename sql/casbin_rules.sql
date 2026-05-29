USE medical_info;

CREATE TABLE IF NOT EXISTS casbin_rules (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  ptype VARCHAR(16) NOT NULL,
  v0 VARCHAR(191) DEFAULT '',
  v1 VARCHAR(191) DEFAULT '',
  v2 VARCHAR(191) DEFAULT '',
  v3 VARCHAR(191) DEFAULT '',
  v4 VARCHAR(191) DEFAULT '',
  v5 VARCHAR(191) DEFAULT '',
  PRIMARY KEY (id),
  INDEX idx_casbin_rules_ptype (ptype),
  INDEX idx_casbin_rules_v0 (v0),
  INDEX idx_casbin_rules_v1 (v1)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO casbin_rules (ptype, v0, v1, v2) VALUES
  ('g', 'role_888', 'admin', ''),
  ('g', 'role_777', 'specimen_manager', ''),
  ('g', 'role_999', 'viewer', ''),
  ('p', 'admin', '/api/drugs/add', 'POST'),
  ('p', 'admin', '/api/drugs/get', 'GET'),
  ('p', 'admin', '/api/specimens/add', 'POST'),
  ('p', 'admin', '/api/specimens/get', 'GET'),
  ('p', 'admin', '/api/users/add', 'POST'),
  ('p', 'admin', '/api/users/get', 'GET'),
  ('p', 'admin', '/api/users/delete', 'POST'),
  ('p', 'specimen_manager', '/api/drugs/get', 'GET'),
  ('p', 'specimen_manager', '/api/specimens/add', 'POST'),
  ('p', 'specimen_manager', '/api/specimens/get', 'GET'),
  ('p', 'viewer', '/api/drugs/get', 'GET'),
  ('p', 'viewer', '/api/specimens/get', 'GET');

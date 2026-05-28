CREATE DATABASE IF NOT EXISTS medical_info
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE medical_info;
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS drugs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL COMMENT '药品名称',
  manufacturer VARCHAR(150) DEFAULT '' COMMENT '生产厂家',
  approval_number VARCHAR(100) DEFAULT '' COMMENT '批准文号',
  specification VARCHAR(100) DEFAULT '' COMMENT '规格',
  price DECIMAL(10, 2) NOT NULL COMMENT '价格',
  stock INT NOT NULL COMMENT '库存数量',
  created_at DATETIME(3) DEFAULT NULL,
  updated_at DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (id),
  INDEX idx_drugs_name (name),
  CONSTRAINT chk_drugs_price_positive CHECK (price > 0),
  CONSTRAINT chk_drugs_stock_positive CHECK (stock > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS specimen_applications (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(50) NOT NULL COMMENT '姓名',
  gender VARCHAR(10) NOT NULL COMMENT '性别',
  age INT NOT NULL COMMENT '年龄',
  id_number VARCHAR(40) NOT NULL COMMENT 'ID号',
  sample_type VARCHAR(20) NOT NULL COMMENT '送检标本类型',
  pathology_type VARCHAR(50) NOT NULL COMMENT '病理类型',
  pdl1_expression INT NOT NULL DEFAULT 0 COMMENT 'PD-L1表达百分比',
  driver_gene_mutation VARCHAR(255) DEFAULT '' COMMENT '驱动基因突变',
  stage VARCHAR(10) NOT NULL COMMENT '分期',
  last_treatment VARCHAR(255) DEFAULT '' COMMENT '末次治疗',
  follow_up_treatment VARCHAR(255) DEFAULT '' COMMENT '后续治疗方案',
  doctor VARCHAR(50) NOT NULL COMMENT '送检医师',
  inspection_date VARCHAR(10) NOT NULL COMMENT '送检日期',
  created_at DATETIME(3) DEFAULT NULL,
  updated_at DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (id),
  INDEX idx_specimen_name (name),
  INDEX idx_specimen_id_number (id_number),
  CONSTRAINT chk_specimen_age_positive CHECK (age > 0),
  CONSTRAINT chk_specimen_pdl1_range CHECK (pdl1_expression BETWEEN 0 AND 100),
  CONSTRAINT chk_specimen_gender CHECK (gender IN ('男', '女')),
  CONSTRAINT chk_specimen_sample_type CHECK (sample_type IN ('组织', '血浆')),
  CONSTRAINT chk_specimen_stage CHECK (stage IN ('I', 'II', 'III', 'IV'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

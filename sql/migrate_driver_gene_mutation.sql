USE medical_info;

ALTER TABLE specimen_applications
  ADD COLUMN driver_gene_mutation VARCHAR(255) DEFAULT '' COMMENT '驱动基因突变'
  AFTER pdl1_expression;

-- 如果确认旧图片字段不再使用，可以执行：
-- ALTER TABLE specimen_applications DROP COLUMN driver_gene_image;

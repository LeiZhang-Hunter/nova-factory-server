-- 缺陷管理表
DROP TABLE IF EXISTS `sys_defect`;
CREATE TABLE IF NOT EXISTS `sys_defect` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '缺陷ID',
  `defect_code` varchar(50) NOT NULL COMMENT '缺陷编码',
  `defect_name` varchar(100) NOT NULL COMMENT '缺陷名称',
  `index_type` varchar(50) NOT NULL COMMENT '指标类型',
  `defect_level` varchar(20) NOT NULL COMMENT '缺陷等级',
  `dept_id` bigint(20) NOT NULL COMMENT '部门ID',
  `state` char(1) DEFAULT '0' COMMENT '状态（0正常 1停用）',
    `create_by_id` bigint(20) NOT NULL COMMENT '创建者ID',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
    `update_by_id` bigint(20) NOT NULL COMMENT '更新者ID',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `remark` varchar(500) DEFAULT NULL COMMENT '备注',
  `attr_1` varchar(100) DEFAULT NULL COMMENT '属性1',
  `attr_2` varchar(100) DEFAULT NULL COMMENT '属性2',
  `attr_3` varchar(100) DEFAULT NULL COMMENT '属性3',
  `attr_4` varchar(100) DEFAULT NULL COMMENT '属性4',
  `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp DEFAULT  NULL COMMENT '更新时间',
  `deleted_at` timestamp DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_dept_id` (`dept_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='缺陷管理表';

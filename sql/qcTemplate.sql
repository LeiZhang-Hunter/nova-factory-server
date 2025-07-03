-- 检测模板表
CREATE TABLE IF NOT EXISTS `qc_template` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '检测模板ID',
  `template_code` varchar(64) NOT NULL COMMENT '检测模板编号',
  `template_name` varchar(255) NOT NULL COMMENT '检测模板名称',
  `qc_types` varchar(255) NOT NULL COMMENT '检测种类',
  `enable_flag` char(1) DEFAULT 'Y' COMMENT '是否启用',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `attr1` varchar(64) DEFAULT NULL COMMENT '预留字段1',
  `attr2` varchar(255) DEFAULT NULL COMMENT '预留字段2',
  `attr3` int(11) DEFAULT 0 COMMENT '预留字段3',
  `attr4` int(11) DEFAULT 0 COMMENT '预留字段4',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
    `create_by_id` bigint(20) NOT NULL COMMENT '创建者ID',
    `update_by_id` bigint(20) NOT NULL COMMENT '更新者ID',
    `created_at` timestamp DEFAULT NULL COMMENT '创建时间',
    `updated_at` timestamp DEFAULT  NULL COMMENT '更新时间',
    `deleted_at` timestamp DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='检测模板表';

-- 检测项表
DROP TABLE IF EXISTS `qc_index`;
CREATE TABLE IF NOT EXISTS `qc_index` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '检测项ID',
    `index_code` varchar(64) NOT NULL COMMENT '检测项编码',
    `index_name` varchar(255) NOT NULL COMMENT '检测项名称',
    `index_type` varchar(64) NOT NULL COMMENT '检测项类型',
    `qc_tool` varchar(255) DEFAULT NULL COMMENT '检测工具',
    `qc_result_type` varchar(64) NOT NULL COMMENT '质检值类型',
    `qc_result_spc` varchar(255) DEFAULT NULL COMMENT '值属性',
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
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='检测项表';

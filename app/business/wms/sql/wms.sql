CREATE TABLE IF NOT EXISTS `wms_warehouse_area` (
    `id` bigint(20) NOT NULL COMMENT '主键ID',
    `warehouse_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '仓库 ID',
    `area_name` VARCHAR(255) NOT NULL  COMMENT '库区名称',
    `parent_id` bigint(20) NOT NULL DEFAULT 0  COMMENT '父辈仓库 ID',
    `warehouse_area_property` tinyint(1) NOT NULL DEFAULT 0 COMMENT '库区属性',
    `last_update_time` datetime(0) NULL DEFAULT NULL COMMENT '最后更新时间',
    `is_valid` tinyint(1) NOT NULL DEFAULT 0  COMMENT '是否有效，0/1',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`id`),
    KEY `idx_warehouse_id` (`warehouse_id`),
    KEY `idx_area_name` (`area_name`),
    KEY `idx_parent_id` (`parent_id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='仓库分区';


-- ============================================================
-- 表名: wms_warehouse_location
-- 用途: 库位表，维护仓库内具体货位信息。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE IF NOT EXISTS `wms_warehouse_location` (
    `id` bigint(20) NOT NULL COMMENT '主键ID',
    `warehouse_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '仓库 ID',
    `warehouse_name` varchar(255) NOT NULL DEFAULT '' COMMENT '仓库名称',
    `warehouse_area_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '库区 ID',
    `warehouse_area_name` varchar(255) NOT NULL DEFAULT '' COMMENT '库区名称',
    `warehouse_area_property` tinyint(1) NOT NULL DEFAULT 0 COMMENT '库区属性',
    `location_name` varchar(255) NOT NULL DEFAULT '' COMMENT '库位名称',
    `location_length` decimal(10, 2) NULL DEFAULT NULL COMMENT '库位长度',
    `location_width` decimal(10, 2) NULL DEFAULT NULL COMMENT '库位宽度',
    `location_height` decimal(10, 2) NULL DEFAULT NULL COMMENT '库位高度',
    `location_volume` decimal(10, 2) NULL DEFAULT NULL COMMENT '库位容积',
    `location_load` decimal(10, 2) NULL DEFAULT NULL COMMENT '库位承重',
    `roadway_number` varchar(64) NOT NULL DEFAULT '' COMMENT '巷道编号',
    `shelf_number` varchar(64) NOT NULL DEFAULT '' COMMENT '货架编号',
    `layer_number` varchar(64) NOT NULL DEFAULT '' COMMENT '层号',
    `tag_number` varchar(64) NOT NULL DEFAULT '' COMMENT '标签号',
    `last_update_time` datetime(0) NULL DEFAULT NULL COMMENT '最后更新时间',
    `is_valid` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否有效，0/1',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`id`),
    KEY `idx_warehouse_id` (`warehouse_id`),
    KEY `idx_warehouse_area_id` (`warehouse_area_id`),
    KEY `idx_location_name` (`location_name`),
    KEY `idx_tag_number` (`tag_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='仓库库位';

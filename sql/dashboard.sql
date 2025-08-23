
CREATE TABLE IF NOT EXISTS `sys_dashboard` (
    `id` bigint(20) NOT NULL COMMENT '调度id',
    `name` varchar(255) NOT NULL COMMENT '',
    `type` varchar(255) NOT NULL COMMENT '',
    `description` varchar(2048) not null comment '描述',
    `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
    `create_time` datetime DEFAULT NULL COMMENT '创建时间',
    `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
    `update_time` datetime DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `sys_dashboard_data` (
    `id` bigint(20) NOT NULL COMMENT '面板id',
    `datashboard_id` bigint(20) NOT NULL COMMENT '仪表盘数据id',
    `data` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '数据',
    `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
    `create_time` datetime DEFAULT NULL COMMENT '创建时间',
    `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
    `update_time` datetime DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

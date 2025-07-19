
CREATE TABLE `sys_alert` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增标识',
    `gateway_id` bigint(20) NOT NULL COMMENT '网关id',
    `template_id` bigint(20) NOT NULL COMMENT '模板id',
    `name` varchar(255) not null comment '告警策略名称',
    `additions` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '注解',
    `advanced` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '告警规则',
    `ignore` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '忽略规则',
    `matcher` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '匹配规则',
    `description` varchar(1024) not null comment '配置版本',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    `status` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0禁用 1启用）',
    PRIMARY KEY (`id`),
    INDEX `gateway_id`(`gateway_id`) USING BTREE,
    INDEX `template_id`(`template_id`) USING BTREE
) ENGINE=Innodb   COMMENT='告警策略配置';

CREATE TABLE `sys_alert_sink_template`
(
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增标识',
    `name` varchar(255) not null comment '告警模板名称',
    `addr` varchar(512) not null comment '发送alert的http地址，若为空，则不会发送',
    `timeout` int(11) not null comment '发送alert的http timeout',
    `headers` text  null comment '发送alert的http header',
    `method` varchar(16) not null comment '发送alert的http method, 如果不填put(不区分大小写)，都认为是POST',
    `template` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '告警模板',
    `extension` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '扩展信息',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`id`)
)

CREATE TABLE `sys_alert_log` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增标识',
    `gateway_id` bigint(20) NOT NULL COMMENT '网关id',
    `alert_id` bigint(20) NOT NULL COMMENT '告警id',
    `device_id` bigint(20) NOT NULL COMMENT '设备id',
    `device_template_id` bigint(20) NOT NULL COMMENT '设备模板id',
    `device_data_id` bigint(20) NOT NULL COMMENT '设备数据id',
    `context` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '内容',
    `message` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '消息',
    `data` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '消息快照',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`id`),
    INDEX `gateway_id`(`gateway_id`) USING BTREE,
    INDEX `alert_id`(`alert_id`) USING BTREE,
    INDEX `device_id`(`device_id`) USING BTREE,
    INDEX `device_template_id`(`device_template_id`) USING BTREE,
    INDEX `device_data_id`(`device_data_id`) USING BTREE
) ENGINE=Innodb   COMMENT='告警日志';
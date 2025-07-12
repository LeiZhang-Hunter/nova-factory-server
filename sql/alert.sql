
CREATE TABLE `sys_alert` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增标识',
    `name` varchar(255) not null comment '告警策略名称',
    `description` varchar(64) not null comment '配置版本',
    `type` varchar(64) not null comment '配置版本',
    `rule` TEXT not null comment '配置内容',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 `    -1删除）',
    PRIMARY KEY (`id`),
    INDEX `agent_object_id`(`agent_object_id`) USING BTREE
) ENGINE=Innodb   COMMENT='agent配置';

CREATE TABLE `sys_alert_template`
(
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增标识',
    `name` varchar(255) not null comment '告警模板',
    `notify_object` varchar(1024) not null comment '通知对象',
)
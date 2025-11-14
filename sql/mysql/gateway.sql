
/**
    入口管理
 */
CREATE TABLE `sys_gateway_inbound_config`  (
   `gateway_config_id` bigint(20) NOT NULL COMMENT '设备主键',
   `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '设备名称',
   `protocol_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '协议类型',
   `config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL  COMMENT '配置',
   `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
   `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
   `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
   `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
   `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
   `status` int(11) NULL DEFAULT 0 COMMENT '操作状态（0正常 1异常）',
   `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
   PRIMARY KEY (`gateway_config_id`) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '设备协议管理' ROW_FORMAT = Dynamic;
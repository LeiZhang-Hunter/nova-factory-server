SET NAMES utf8mb4;

SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sys_oper_log
-- ----------------------------
CREATE TABLE `sys_device`  (
    `device_id` bigint(20) NOT NULL COMMENT '设备主键',
    `device_group_id` bigint(20) NOT NULL  COMMENT '设备分组id',
    `device_class_id` bigint(20) NOT NULL  COMMENT '设备分类id',
    `device_protocol_id` bigint(20) NOT NULL  COMMENT '协议id',
    `device_building_id` bigint(20) NOT NULL  COMMENT '建筑物id',
    `device_gateway_id` bigint(20) NOT NULL  COMMENT '网关id',
    `slave` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '从设备地址',
    `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '设备名称',
    `number` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '设备编号',
    `communication_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '通信方式',
    `protocol_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '协议类型',
    `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '设备类型',
    `action` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL  COMMENT '行为字段',
    `extension` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL  COMMENT '扩展字段',
    `control_type` int(11) NULL DEFAULT 0 COMMENT '控制模式 0 是手动 1是自动',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `status` int(11) NULL DEFAULT 0 COMMENT '操作状态（0正常 1异常）',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`device_id`) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '设备管理' ROW_FORMAT = Dynamic;


CREATE TABLE `sys_device_group` (
    `group_id` bigint(20) NOT NULL COMMENT '设备主键',
    `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '设备名称',
    `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL  COMMENT '备注',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`group_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='设备组';




CREATE TABLE `sys_material`  (
   `material_id` bigint(20) NOT NULL COMMENT '物料id',
   `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '物料名',
   `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '编码',
   `model` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '型号',
   `unit` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '单位',
   `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '物料类型',
   `factory` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '厂家',
   `address` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '存放地点',
   `price` decimal(10,2) COLLATE utf8mb4_general_ci NULL DEFAULT 0.0 COMMENT '价格',
   `total` bigint(20) NOT NULL COMMENT '总数',
   `outbound` bigint(20) NOT NULL COMMENT '出库',
   `current_total` bigint(20) NOT NULL COMMENT '当前库存',
   `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
   `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
   `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
   `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
   `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
   `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
   PRIMARY KEY (`material_id`) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '设备管理' ROW_FORMAT = Dynamic;



CREATE TABLE `sys_material_inbound`  (
    `inbound_id` bigint(20) NOT NULL COMMENT '入库id',
    `material_id` bigint(20) NOT NULL COMMENT '物料id',
    `number` bigint(20) NOT NULL COMMENT '入库数量',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`inbound_id`) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '物料入库管理' ROW_FORMAT = Dynamic;

CREATE TABLE `sys_material_outbound`  (
    `outbound_id` bigint(20) NOT NULL COMMENT '出库id',
    `material_id` bigint(20) NOT NULL COMMENT '物料id',
    `number` bigint(20) NOT NULL COMMENT '入库数量',
    `reason` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL  COMMENT '出库理由',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`outbound_id`) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '物料出库管理' ROW_FORMAT = Dynamic;

CREATE TABLE `sys_device_template`
(
    `template_id` bigint(20) NOT NULL COMMENT '设备主键',
    `name` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '设备名称',
    `template_type`  tinyint(1) NULL DEFAULT 0 COMMENT '模板类型0是私有1 是共有',
    `vendor` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '供应商',
    `protocol` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '协议类型',
    `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL  COMMENT '备注',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`template_id`) USING BTREE,
    KEY (`protocol`)
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'modbus数据模板' ROW_FORMAT = Dynamic;

CREATE TABLE `sys_modbus_device_config_data`
(
    `device_config_id` bigint(20) NOT NULL COMMENT '文档id',
    `template_id` bigint(20) NOT NULL COMMENT '模板id',
    `name` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '数据名称',
    `type` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '数据类型',
    `data_type` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '设备节点类型，到底是开关还是数值',
    `register` bigint(20) NOT NULL COMMENT '寄存器/偏移量',
    `storage_type` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '存储策略',
    `unit`  varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '单位',
    `precision` bigint(20) NOT NULL COMMENT '数据精度',
    `function_code` tinyint(1) NULL DEFAULT 0 COMMENT '功能码',
    `agg_function` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '聚合函数，用来计算图表',
    `perturbation_variables` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '扰动变量',
    `graph_enable` tinyint(1) NULL DEFAULT 0 COMMENT '是否开启图表',
    `predict_enable` tinyint(1) NULL DEFAULT 0 COMMENT '是否开启趋势预测',
    `mode` tinyint(1) NULL DEFAULT 0 COMMENT '读写权限',
    `data_format` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '数据格式',
    `annotation` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL  COMMENT '注解',
    `sort` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '数据排序',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`device_config_id`) USING BTREE,
    INDEX `template_id`(`template_id`) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'modbus数据配置' ROW_FORMAT = Dynamic;




CREATE TABLE `sys_iot_agent` (
    `object_id` bigint(20) unsigned NOT NULL COMMENT 'agent uuid',
    `name` varchar(64) not null comment 'agent名字',
    `username` varchar(64) not null comment 'username',
    `password` varchar(64) not null comment 'password',
    `operate_state` int(2) NOT NULL DEFAULT 0 COMMENT '操作状态 1-启动中 2-停止中 3-启动失败 4-停止失败',
    `operate_time` timestamp NULL COMMENT '操作时间',
    `version` varchar(64) not null comment 'agent版本',
    `config_id` bigint(20) unsigned NOT NULL COMMENT '配置id',
    `last_config_id` bigint(20) unsigned NOT NULL COMMENT '最后版本的配置id',
    `ipv4`varchar(32) not null default '' comment  'ipv4地址',
    `ipv6`varchar(64) not null default '' comment  'ipv6地址',
    `last_heartbeat_time` timestamp NULL COMMENT '上次心跳时间',
    `update_config_time` timestamp NULL COMMENT '更新配置时间',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`object_id`)
) ENGINE=Innodb   COMMENT='agent信息';

CREATE TABLE `sys_iot_agent_process` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增标识',
    `agent_object_id` bigint(20) unsigned NOT NULL COMMENT 'agent uuid',
    `status` int(2) NOT NULL COMMENT '状态 1-运行 2-停止',
    `name` varchar(64) not null comment '名字',
    `version` varchar(64) not null comment '版本',
    `start_time` timestamp NULL COMMENT '启动时间',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    INDEX (`agent_object_id`, `name`, `version`, `state`),
    UNIQUE INDEX (`agent_object_id`, `name`),
    PRIMARY KEY (`id`)
) ENGINE=Innodb   COMMENT='agent进程信息';

CREATE TABLE `sys_iot_agent_config` (
   `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增标识',
   `agent_object_id` bigint(20) unsigned NOT NULL COMMENT 'agent id',
   `config_version` varchar(64) not null comment '配置版本',
   `content` TEXT not null comment '配置内容',
   `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
   `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
   `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
   `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
   `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
   `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 `    -1删除）',
   PRIMARY KEY (`id`),
   INDEX `agent_object_id`(`agent_object_id`) USING BTREE
) ENGINE=Innodb   COMMENT='agent配置';

CREATE TABLE `sys_iot_db_dev_map`
(
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增标识',
    `device_id` bigint(20) NOT NULL COMMENT '设备id',
    `template_id` bigint(20) NOT NULL COMMENT '模板id',
    `data_id` bigint(20) NOT NULL COMMENT '测点id',
    `device` varchar(128) not null comment '设备名字',
    `data_name` varchar(128) not null comment '数据名字',
    `unit` varchar(128) not null comment '数据单位',
    `dept_id` bigint(20) NULL DEFAULT 0 COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常-1删除）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_dev` (`device`) USING BTREE
) ENGINE=Innodb   COMMENT='数据和iotdb 对照表';



CREATE TABLE `sys_building`  (
    `id` bigint(20) NOT NULL COMMENT 'id',
    `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '建筑物名称',
    `code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '建筑物编号',
    `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '建筑物类型',
    `status` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '建筑物状态',
    `area` bigint(20) NOT NULL COMMENT '建筑面积',
    `floors` bigint(20) NOT NULL COMMENT '楼层数',
    `year` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '建造年份',
    `lifeYears` bigint(20) NOT NULL COMMENT '使用年限',
    `manager` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '负责人w',
    `phone` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '单位',
    `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '单位',
    `remark` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '单位',
     `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '设备管理' ROW_FORMAT = Dynamic;

CREATE TABLE `sys_building`  (
     `id` bigint(20) NOT NULL COMMENT 'id',
     `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '建筑物名称',
     `code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '建筑物编号',
     `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '建筑物类型',
     `status` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '建筑物状态',
     `area` bigint(20) NOT NULL COMMENT '建筑面积',
     `floors` bigint(20) NOT NULL COMMENT '楼层数',
     `year` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '建造年份',
     `lifeYears` bigint(20) NOT NULL COMMENT '使用年限',
     `manager` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '负责人w',
     `phone` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '单位',
     `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '单位',
     `remark` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '单位',
     `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
     `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
     `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
     `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
     `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
     `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
     PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '设备管理' ROW_FORMAT = Dynamic;

CREATE TABLE `sys_device_electric_setting`  (
     `id` bigint(20) NOT NULL COMMENT 'id',
     `device_id` bigint(20) NOT NULL COMMENT 'device_id',
     `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '配置名称',
     `expression` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL  COMMENT '表达式',
     `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
     `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
     `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
     `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
     `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
     `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
     PRIMARY KEY (`id`) USING BTREE,
     UNIQUE KEY `uniq_device_id` (`device_id`) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '设备电流配置，用来计算设备电流是待机 停机 还是 运行' ROW_FORMAT = Dynamic;

-- ----------------------------
-- 3、设备点检保养项目表
-- ----------------------------
drop table if exists sys_device_subject;
create table sys_device_subject (
    id                  bigint(20)      not null auto_increment    comment '项目ID',
    subject_code                varchar(64)     not null                   comment '项目编码',
    subject_name                varchar(255)                               comment '项目名称',
    subject_type                varchar(64)     default 0                  comment '项目类型',
    subject_content             varchar(500)    not null                   comment '项目内容',
    subject_standard            varchar(255)                               comment '标准',
    remark                      varchar(500)    default ''                 comment '备注',
    attr1                       varchar(64)     default null               comment '预留字段1',
    attr2                       varchar(255)    default null               comment '预留字段2',
    attr3                       int(11)         default 0                  comment '预留字段3',
    attr4                       int(11)         default 0                  comment '预留字段4',
    `status`                    tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0禁用 1启用）',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    primary key (id),
    UNIQUE KEY `uniq_subject_code` (`subject_code`) USING BTREE
) engine=innodb auto_increment=200 comment = '设备点检保养项目表';


-- ----------------------------
-- 4、设备点检保养计划头表
-- ----------------------------
drop table if exists sys_device_check_plan;
create table sys_device_check_plan (
   plan_id                     bigint(20)      not null auto_increment    comment '计划ID',
   plan_code                   varchar(64)     not null                   comment '计划编码',
   plan_name                   varchar(255)                               comment '计划名称',
   plan_type                   varchar(64)     not null                   comment '计划类型',
   start_date                  datetime                                   comment '开始日期',
   end_date                    datetime                                   comment '结束日期',
   cycle_type                  varchar(64)                                comment '频率',
   cycle_count                 int(11)                                    comment '次数',
   status                      varchar(64)                                comment '状态',
   remark                      varchar(500)    default ''                 comment '备注',
   attr1                       varchar(64)     default null               comment '预留字段1',
   attr2                       varchar(255)    default null               comment '预留字段2',
   attr3                       int(11)         default 0                  comment '预留字段3',
   attr4                       int(11)         default 0                  comment '预留字段4',
   `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
   `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
   `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
   `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
   `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
   `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
   primary key (plan_id)
) engine=innodb auto_increment=200 comment = '设备点检保养计划头表';




-- ----------------------------
-- 5、点检保养计划设备表
-- ----------------------------
drop table if exists sys_device_check_machinery;
create table sys_device_check_machinery (
    record_id                   bigint(20)      not null auto_increment    comment '流水号',
    plan_id                     bigint(20)      not null                   comment '计划ID',
    machinery_id                bigint(20)      not null                   comment '设备ID',
    machinery_code              varchar(64)     not null                   comment '设备编码',
    machinery_name              varchar(255)    not null                   comment '设备名称',
    machinery_brand             varchar(255)                               comment '品牌',
    machinery_spec              varchar(255)                               comment '规格型号',
    remark                      varchar(500)    default ''                 comment '备注',
    attr1                       varchar(64)     default null               comment '预留字段1',
    attr2                       varchar(255)    default null               comment '预留字段2',
    attr3                       int(11)         default 0                  comment '预留字段3',
    attr4                       int(11)         default 0                  comment '预留字段4',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    primary key (record_id),
    KEY `plan_id_idx` (`plan_id`) USING BTREE
) engine=innodb auto_increment=200 comment = '点检设备表';


-- ----------------------------
-- 6、点检保养计划项目表
-- ----------------------------
drop table if exists sys_device_check_subject;
create table sys_device_check_subject (
    record_id                   bigint(20)      not null auto_increment    comment '流水号',
    plan_id                     bigint(20)      not null                   comment '计划ID',
    subject_id                  bigint(20)      not null                   comment '项目ID',
    subject_code                varchar(64)     not null                   comment '项目编码',
    subject_name                varchar(255)                               comment '项目名称',
    subject_type                varchar(64)                                comment '项目类型',
    subject_content             varchar(500)    not null                   comment '项目内容',
    subject_standard            varchar(255)                               comment '标准',
    remark                      varchar(500)    default ''                 comment '备注',
    attr1                       varchar(64)     default null               comment '预留字段1',
    attr2                       varchar(255)    default null               comment '预留字段2',
    attr3                       int(11)         default 0                  comment '预留字段3',
    attr4                       int(11)         default 0                  comment '预留字段4',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    primary key (record_id)
) engine=innodb auto_increment=200 comment = '点检项目表';



-- ----------------------------
-- 7、设备点检记录表
-- ----------------------------
drop table if exists sys_device_check_record;
create table sys_device_check_record (
    record_id                   bigint(20)      not null auto_increment    comment '计划ID',
    plan_id                     bigint(20)                                 comment '计划ID',
    plan_code                   varchar(64)                                comment '计划编码',
    plan_name                   varchar(255)                               comment '计划名称',
    plan_type                   varchar(64)                                comment '计划类型',
    machinery_id                bigint(20)      not null                   comment '设备ID',
    machinery_code              varchar(64)     not null                   comment '设备编码',
    machinery_name              varchar(255)    not null                   comment '设备名称',
    machinery_brand             varchar(255)                               comment '品牌',
    machinery_spec              varchar(255)                               comment '规格型号',
    check_time                  datetime        not null                   comment '点检时间',
    user_id                     bigint(20)                                 comment '点检人',
    user_name                   varchar(64)                                comment '点检人用户名',
    nick_name                   varchar(255)                               comment '点检人名称',
    status                      varchar(64)     default 'PREPARE'          comment '状态',
    remark                      varchar(500)    default ''                 comment '备注',
    attr1                       varchar(64)     default null               comment '预留字段1',
    attr2                       varchar(255)    default null               comment '预留字段2',
    attr3                       int(11)         default 0                  comment '预留字段3',
    attr4                       int(11)         default 0                  comment '预留字段4',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    primary key (record_id)
) engine=innodb auto_increment=200 comment = '设备点检记录表';



-- ----------------------------
-- 7、设备保养记录表
-- ----------------------------
drop table if exists sys_device_mainten_record;
create table sys_device_mainten_record (
    record_id                   bigint(20)      not null auto_increment    comment '计划ID',
    plan_id                     bigint(20)                                 comment '计划ID',
    plan_code                   varchar(64)                                comment '计划编码',
    plan_name                   varchar(255)                               comment '计划名称',
    plan_type                   varchar(64)                                comment '计划类型',
    machinery_id                bigint(20)      not null                   comment '设备ID',
    machinery_code              varchar(64)     not null                   comment '设备编码',
    machinery_name              varchar(255)    not null                   comment '设备名称',
    machinery_brand             varchar(255)                               comment '品牌',
    machinery_spec              varchar(255)                               comment '规格型号',
    mainten_time                  datetime        not null                 comment '保养时间',
    user_id                     bigint(20)                                 comment '用户ID',
    user_name                   varchar(64)                                comment '用户名',
    nick_name                   varchar(128)                               comment '昵称',
    status                      varchar(64)     default 'PREPARE'          comment '状态',
    remark                      varchar(500)    default ''                 comment '备注',
    attr1                       varchar(64)     default null               comment '预留字段1',
    attr2                       varchar(255)    default null               comment '预留字段2',
    attr3                       int(11)         default 0                  comment '预留字段3',
    attr4                       int(11)         default 0                  comment '预留字段4',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    primary key (record_id)
) engine=innodb auto_increment=200 comment = '设备保养记录表';


-- ----------------------------
-- 8、设备保养记录行表
-- ----------------------------
drop table if exists sys_device_mainten_record_line;
create table sys_device_mainten_record_line (
    line_id                     bigint(20)      not null auto_increment    comment '计划ID',
    record_id                   bigint(20)      not null                   comment '计划ID',
    subject_id                  bigint(20)      not null                   comment '项目ID',
    subject_code                varchar(64)     not null                   comment '项目编码',
    subject_name                varchar(255)                               comment '项目名称',
    subject_type                varchar(64)                                comment '项目类型',
    subject_content             varchar(500)    not null                   comment '项目内容',
    subject_standard            varchar(255)                               comment '标准',
    mainten_status              varchar(64)     not null                   comment '保养结果',
    mainten_result              varchar(500)                               comment '异常描述',
    attr1                       varchar(64)     default null               comment '预留字段1',
    attr2                       varchar(255)    default null               comment '预留字段2',
    attr3                       int(11)         default 0                  comment '预留字段3',
    attr4                       int(11)         default 0                  comment '预留字段4',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    primary key (line_id)
) engine=innodb auto_increment=200 comment = '设备保养记录行表';





-- ----------------------------
-- 7、设备维修单
-- ----------------------------
drop table if exists sys_device_repair;
create table sys_device_repair (
    repair_id                   bigint(20)      not null auto_increment    comment '维修单ID',
    repair_code                 varchar(64)     not null                   comment '维修单编号',
    repair_name                 varchar(255)                               comment '维修单名称',
    machinery_id                bigint(20)      not null                   comment '设备ID',
    machinery_code              varchar(64)     not null                   comment '设备编码',
    machinery_name              varchar(255)    not null                   comment '设备名称',
    machinery_brand             varchar(255)                               comment '品牌',
    machinery_spec              varchar(255)                               comment '规格型号',
    machinery_type_id           bigint(20)                                 comment '设备类型ID',
    require_date                datetime                                   comment '报修日期',
    finish_date                 datetime                                   comment '维修完成日期',
    confirm_date                datetime                                   comment '验收日期',
    repair_result               varchar(64)                                comment '维修结果',
    accepted_id                 bigint(20)                                 comment '维修人员ID',
    accepted_name               varchar(64)                                comment '维修人员用户名',
    accepted_by                 varchar(64)                                comment '维修人员名称',
    confirm_id                  bigint(20)                                 comment '验收人员ID',
    confirm_name                varchar(64)                                comment '验收人员用户名',
    confirm_by                  varchar(64)                                comment '验收人员名称',
    source_doc_type             varchar(64)                                comment '来源单据类型',
    source_doc_id               bigint(20)                                 comment '来源单据ID',
    source_doc_code             varchar(64)                                comment '来源单据编号',
    status                      varchar(64)     default 'PREPARE'          comment '单据状态',
    remark                      varchar(500)    default ''                 comment '备注',
    attr1                       varchar(64)     default null               comment '预留字段1',
    attr2                       varchar(255)    default null               comment '预留字段2',
    attr3                       int(11)         default 0                  comment '预留字段3',
    attr4                       int(11)         default 0                  comment '预留字段4',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    primary key (repair_id)
) engine=innodb auto_increment=200 comment = '设备维修单';



-- ----------------------------
-- 7、设备维修单行
-- ----------------------------
drop table if exists sys_device_repair_line;
create table sys_device_repair_line (
    line_id                     bigint(20)      not null auto_increment    comment '行ID',
    repair_id                   bigint(20)      not null                   comment '维修单ID',
    subject_id                  bigint(20)                                 comment '项目ID',
    subject_code                varchar(64)                                comment '项目编码',
    subject_name                varchar(255)                               comment '项目名称',
    subject_type                varchar(64)                                comment '项目类型',
    subject_content             varchar(500)                               comment '项目内容',
    subject_standard            varchar(255)                               comment '标准',
    malfunction                 varchar(500)    not null                   comment '故障描述',
    malfunction_url             varchar(255)                               comment '故障描述资源',
    repair_des                  varchar(500)                               comment '维修情况',
    remark                      varchar(500)    default ''                 comment '备注',
    attr1                       varchar(64)     default null               comment '预留字段1',
    attr2                       varchar(255)    default null               comment '预留字段2',
    attr3                       int(11)         default 0                  comment '预留字段3',
    attr4                       int(11)         default 0                  comment '预留字段4',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    primary key (line_id)
) engine=innodb auto_increment=200 comment = '设备维修单行';
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
    `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '设备名称',
    `number` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '设备编号',
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

Create Table: CREATE TABLE `sys_dataset` (
     `dataset_id` bigint NOT NULL COMMENT '出库id',
     `dataset_avatar` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT 'base64 编码的头像。',
    `dataset_chunk_method` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '要创建的数据集的分块方法',
    `dataset_create_date` datetime DEFAULT NULL COMMENT '创建时间',
    `dataset_create_time` bigint NOT NULL COMMENT '创建时间',
    `dataset_created_by` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '创建人',
    `dataset_description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '描述',
    `dataset_document_count` bigint NOT NULL DEFAULT '0' COMMENT '文档数',
    `dataset_embedding_model` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '嵌入模型',
    `dataset_uuid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '知识库id',
    `dataset_language` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '知识库语言',
    `dataset_name` varchar(136) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '要创建的数据集的唯一名称',
    `dataset_pagerank` int DEFAULT '0' COMMENT '页面排名',
    `dataset_parser_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '数据集解析器的配置设置。',
    `dataset_permission` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '指定谁可以访问要创建的数据集。',
    `dataset_similarity_threshold` decimal(10,2) DEFAULT '0.00' COMMENT '最小相似度分数',
    `dataset_token_num` bigint NOT NULL DEFAULT '0' COMMENT 'token num',
    `dataset_chunk_count` bigint NOT NULL DEFAULT '0' COMMENT 'chunk num',
    `dataset_update_date` datetime DEFAULT NULL COMMENT '更新时间',
    `dataset_update_time` bigint NOT NULL COMMENT '更新时间',
    `dataset_vector_similarity_weight` decimal(10,2) DEFAULT '0.00' COMMENT '矢量余弦相似性的权重。',
    `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint DEFAULT NULL COMMENT '创建者',
    `create_time` datetime DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint DEFAULT NULL COMMENT '更新者',
    `update_time` datetime DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`dataset_id`) USING BTREE,
    UNIQUE KEY `dataset_name` (`dataset_name`,`state`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='物料出库管理';

CREATE TABLE `sys_dataset_document`  (
    `document_id` bigint(20) NOT NULL COMMENT '文档id',
    `dataset_id` bigint(20) NOT NULL COMMENT '数据集id',
    `dataset_chunk_method` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '要创建的数据集的分块方法',
    `dataset_created_by` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '创建人',
    `dataset_document_uuid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '文档id',
    `dataset_dataset_uuid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '知识库id',
    `dataset_language` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '知识库语言',
    `dataset_location` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '知识库定位',
    `dataset_name` varchar(136) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '要创建的数据集的唯一名称',
    `dataset_parser_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL  COMMENT '数据集解析器的配置设置。',
    `dataset_run` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '知识库运行状态',
    `dataset_size` bigint(20) NOT NULL DEFAULT 0 COMMENT '知识库尺寸大小',
    `dataset_thumbnail` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '知识库锁略图',
    `dataset_type` varchar(136) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '知识库类型',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`document_id`) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '文档管理' ROW_FORMAT = Dynamic;

alter table add column `dataset_chunk_count` bigint NOT NULL DEFAULT '0' COMMENT 'chunk num';

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
    PRIMARY KEY (`template_id`) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'modbus数据模板' ROW_FORMAT = Dynamic;

CREATE TABLE `sys_modbus_device_config_data`
(
    `device_config_id` bigint(20) NOT NULL COMMENT '文档id',
    `template_id` bigint(20) NOT NULL COMMENT '模板id',
    `name` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '数据名称',
    `type` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '数据类型',
    `slave` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '从设备地址',
    `register` bigint(20) NOT NULL COMMENT '寄存器/偏移量',
    `storage_type` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '存储策略',
    `unit`  varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '单位',
    `precision` bigint(20) NOT NULL COMMENT '数据精度',
    `function_code` tinyint(1) NULL DEFAULT 0 COMMENT '功能码',
    `mode` tinyint(1) NULL DEFAULT 0 COMMENT '读写权限',
    `data_format` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '数据格式',
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
    `config_uuid` varchar(64) not null default '' comment '配置版本',
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
   `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增标识',
   `uuid` VARCHAR(32) NOT NULL COMMENT '配置 uuid',
   `company_uuid` VARCHAR(32) NOT NULL COMMENT '公司 uuid',
   `config_version` varchar(64) not null comment '配置版本',
   `content` TEXT not null comment '配置内容',
   `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
   `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
   `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
   `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
   `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
   `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
   PRIMARY KEY (`id`),
   UNIQUE INDEX (`uuid`),
   UNIQUE INDEX (`company_uuid`, `config_version`)
) ENGINE=Innodb   COMMENT='agent配置';
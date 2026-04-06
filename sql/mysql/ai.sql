
CREATE TABLE IF NOT EXISTS `sys_dataset` (
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

CREATE TABLE IF NOT EXISTS `sys_dataset_document`  (
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
     `dataset_chunk_count` bigint NOT NULL DEFAULT '0' COMMENT 'chunk num',
     `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
     `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
     `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
     `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
     `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
     `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
     PRIMARY KEY (`document_id`) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '文档管理' ROW_FORMAT = Dynamic;

CREATE TABLE IF NOT EXISTS `sys_dataset_role_permission` (
    `id` bigint(20) NOT NULL COMMENT '主键ID',
    `role_id` bigint(20) NOT NULL COMMENT '角色ID(sys_role.role_id)',
    `dataset_id` bigint(20) NOT NULL COMMENT '知识库ID(sys_dataset.dataset_id)',
    `document_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '文档ID(sys_dataset_document.document_id)，为0表示对整个知识库授权',
    `permission` varchar(32) NOT NULL DEFAULT 'read' COMMENT '权限类型：read/write/admin',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_role_dataset_document` (`role_id`, `dataset_id`, `document_id`, `state`) USING BTREE,
    KEY `idx_role_id` (`role_id`) USING BTREE,
    KEY `idx_dataset_id` (`dataset_id`) USING BTREE,
    KEY `idx_document_id` (`document_id`) USING BTREE,
    KEY `idx_dept_id` (`dept_id`) USING BTREE,
    KEY `idx_create_time` (`create_time`) USING BTREE,
    KEY `idx_update_time` (`update_time`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci
  COMMENT = '知识库/文档-角色权限表';

CREATE TABLE IF NOT EXISTS `sys_ai_prediction_list`  (
     `id` bigint(20) NOT NULL COMMENT 'id',
     `reason_id` bigint(20) NOT NULL COMMENT '模型推理id',
     `action_id` bigint(20) NOT NULL COMMENT '处理通知id',
     `threshold` bigint(5) NOT NULL COMMENT 'threshold',
     `name` varchar(255) not null comment '智能预警名称',
     `dev` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '测点名称',
     `advanced` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '告警规则',
     `model` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '预测模型',
     `field` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '预测字段',
     `interval` bigint(20) NOT NULL COMMENT '预测时间段',
     `predict_length` bigint(20) unsigned NULL DEFAULT 0 COMMENT '预测长度',
     `agg_function` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '聚合函数，用来计算图表',
     `perturbation_variables` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '关联变量',
     `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
     `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
     `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
     `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
     `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
     `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
     PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '预测列表' ROW_FORMAT = Dynamic;

CREATE TABLE IF NOT EXISTS `sys_ai_prediction_exception`  (
    `id` bigint(20) NOT NULL COMMENT 'id',
    `reason_id` bigint(20) NOT NULL COMMENT '模型推理id',
    `action_id` bigint(20) NOT NULL COMMENT '处理通知id',
    `threshold` bigint(5) NOT NULL COMMENT 'threshold',
    `name` varchar(255) not null comment '智能预警名称',
    `dev` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '测点名称',
    `model` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '预测模型',
    `interval` bigint(20) NOT NULL COMMENT '预测时间段',
    `predict_length` bigint(20) unsigned NULL DEFAULT 0 COMMENT '预测长度',
    `agg_function` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '聚合函数，用来计算图表',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '异常预警' ROW_FORMAT = Dynamic;

CREATE TABLE IF NOT EXISTS `sys_ai_prediction_control`  (
    `id` bigint(20) NOT NULL COMMENT 'id',
    `device_gateway_id` bigint(20) NOT NULL  COMMENT '网关id',
    `name` varchar(255) not null comment '智能预警名称',
    `parallelism` bigint(5) NOT NULL COMMENT '并发',
    `threshold` bigint(5) NOT NULL COMMENT 'threshold',
    `model` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '预测模型',
    `interval` bigint(20) NOT NULL COMMENT '预测时间段',
    `predict_length` bigint(20) unsigned NULL DEFAULT 0 COMMENT '预测长度',
    `agg_function` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '聚合函数，用来计算图表',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '趋势控制' ROW_FORMAT = Dynamic;


CREATE TABLE IF NOT EXISTS `ai_model_provider` (
    `id` bigint(20) NOT NULL COMMENT '主键ID',
    `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '供应商名称',
    `logo` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '图标地址',
    `tags` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '能力标签，逗号分隔',
    `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1启用，0禁用',
    `rank` int(11) NOT NULL DEFAULT 0 COMMENT '排序值，越大越靠前',
    `dept_id` bigint(20) DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) DEFAULT NULL COMMENT '创建时间(系统)',
    `update_by` bigint(20) DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) DEFAULT NULL COMMENT '更新时间(系统)',
    `state` tinyint(1) DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE KEY `uk_provider_name_state` (`name`,`state`) USING BTREE,
    KEY `idx_status` (`status`) USING BTREE,
    KEY `idx_rank` (`rank`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'AI模型供应商配置' ROW_FORMAT = Dynamic;

CREATE TABLE IF NOT EXISTS `ai_llm` (
    `llm_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '模型名称',
    `model_type` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '模型类型',
    `fid` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '模型唯一标识',
    `max_tokens` int(11) NOT NULL COMMENT '最大Token',
    `tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '模型标签',
    `is_tools` tinyint(1) NOT NULL COMMENT '是否支持工具调用：0否1是',
    `status` varchar(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '状态',
    `dept_id` bigint(20) DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) DEFAULT NULL COMMENT '创建时间(系统)',
    `update_by` bigint(20) DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) DEFAULT NULL COMMENT '更新时间(系统)',
    `state` tinyint(1) DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`fid`, `llm_name`) USING BTREE,
    KEY `idx_llm_name` (`llm_name`) USING BTREE,
    KEY `idx_llm_model_type` (`model_type`) USING BTREE,
    KEY `idx_llm_fid` (`fid`) USING BTREE,
    KEY `idx_llm_tags` (`tags`) USING BTREE,
    KEY `idx_llm_status` (`status`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'LLM模型' ROW_FORMAT = Dynamic;


 CREATE TABLE IF NOT EXISTS `ai_llm_setting` (
     `id` bigint(20) NOT NULL COMMENT '主键ID',
    `name` varchar(100) DEFAULT NULL,
    `public_key` varchar(255) DEFAULT NULL,
    `llm_id` varchar(128) NOT NULL,
    `embd_id` varchar(128) NOT NULL,
    `asr_id` varchar(128) NOT NULL,
    `img2txt_id` varchar(128) NOT NULL,
    `rerank_id` varchar(128) NOT NULL,
    `tts_id` varchar(256) DEFAULT NULL,
    `parser_ids` varchar(256) NOT NULL,
    `credit` int NOT NULL,
    `status` varchar(1) DEFAULT NULL,
    `dept_id` bigint(20) DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) DEFAULT NULL COMMENT '创建时间(系统)',
    `update_by` bigint(20) DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) DEFAULT NULL COMMENT '更新时间(系统)',
    `state` tinyint(1) DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`id`),
    KEY `tenant_name` (`name`),
    KEY `tenant_public_key` (`public_key`),
    KEY `tenant_llm_id` (`llm_id`),
    KEY `tenant_embd_id` (`embd_id`),
    KEY `tenant_asr_id` (`asr_id`),
    KEY `tenant_img2txt_id` (`img2txt_id`),
    KEY `tenant_rerank_id` (`rerank_id`),
    KEY `tenant_tts_id` (`tts_id`),
    KEY `tenant_parser_ids` (`parser_ids`),
    KEY `tenant_credit` (`credit`),
    KEY `tenant_status` (`status`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = 'LLM模型配置' ROW_FORMAT = Dynamic;;

CREATE TABLE IF NOT EXISTS `ai_user_llm` (
    `user_id` bigint(20) NOT NULL COMMENT '用户id',
    `llm_factory` varchar(128) NOT NULL,
    `model_type` varchar(128) DEFAULT NULL,
    `llm_name` varchar(128) NOT NULL,
    `api_key` text,
    `api_base` varchar(255) DEFAULT NULL,
    `max_tokens` int NOT NULL,
    `used_tokens` int NOT NULL,
    `status` varchar(1) NOT NULL,
    PRIMARY KEY (`user_id`,`llm_factory`,`llm_name`),
    KEY `tenantllm_user_id` (`user_id`),
    KEY `tenantllm_llm_factory` (`llm_factory`),
    KEY `tenantllm_model_type` (`model_type`),
    KEY `tenantllm_llm_name` (`llm_name`),
    KEY `tenantllm_max_tokens` (`max_tokens`),
    KEY `tenantllm_used_tokens` (`used_tokens`),
    KEY `tenant_llm_status` (`status`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT = '用户选择的模型厂商' ROW_FORMAT = Dynamic;

CREATE TABLE IF NOT EXISTS ai_gateway (
    id               bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    name             VARCHAR(128)    NOT NULL COMMENT '网关名称',
    base_url         VARCHAR(512)    NOT NULL COMMENT 'API服务器地址',
    api_key          VARCHAR(512)    NOT NULL COMMENT 'API Key',
    enabled   TINYINT(1)      NOT NULL DEFAULT 1 COMMENT '启用状态:1启用,0停用',
    active    TINYINT(1)      NOT NULL DEFAULT 0 COMMENT '在线状态:1在线,0离线',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (id)
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_general_ci
    COMMENT = '智能体网关配置表';

CREATE TABLE IF NOT EXISTS ai_conversations (
    id bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    name varchar(255) NOT NULL DEFAULT '' COMMENT '会话名称',
    message text NULL COMMENT '最后一条消息',
    llm_provider_id varchar(128) NOT NULL DEFAULT '' COMMENT '模型供应商ID',
    llm_model_id varchar(128) NOT NULL DEFAULT '' COMMENT '模型ID',
    enable_thinking tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否开启思考模式',
    chat_mode varchar(32) NOT NULL DEFAULT '' COMMENT '聊天模式',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (id),
    KEY idx_create_time (create_time),
    KEY idx_update_time (update_time)
) ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '会话表';

CREATE TABLE IF NOT EXISTS mcp_servers (
    id bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'MCP服务ID',
    `name` varchar(255) NOT NULL COMMENT 'MCP服务名称',
    `description` text NULL COMMENT 'MCP服务描述',
    `transport` varchar(32) NOT NULL DEFAULT 'stdio' COMMENT '传输方式: stdio | streamableHttp',
    `command` text NULL COMMENT 'stdio 模式启动命令',
    `args` text NULL COMMENT 'stdio 模式参数(JSON数组字符串)',
    `env` text NULL COMMENT 'stdio 模式环境变量(JSON对象字符串)',
    `url` varchar(1024) DEFAULT '' COMMENT 'streamableHttp 模式URL',
    `headers` text NULL COMMENT 'streamableHttp 请求头(JSON对象字符串)',
    `timeout` int(11) NOT NULL DEFAULT 30 COMMENT '超时时间(秒)',
    `is_common` tinyint(1) DEFAULT 0 COMMENT '是否是共用的',
    `enabled` tinyint(1) DEFAULT 1 COMMENT '是否启用',
    `dept_id` bigint(20) DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) DEFAULT NULL COMMENT '创建时间(系统)',
    `update_by` bigint(20) DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) DEFAULT NULL COMMENT '更新时间(系统)',
    `state` tinyint(1) DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`id`),
    KEY `idx_mcp_servers_name` (`name`),
    KEY `idx_mcp_servers_transport` (`transport`),
    KEY `idx_mcp_servers_enabled` (`enabled`)
) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'mcp服务器配置';

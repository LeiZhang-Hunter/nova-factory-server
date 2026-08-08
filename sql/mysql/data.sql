-- Data Pipeline configuration module.

DROP TABLE IF EXISTS `data_dispatch_logs`;
DROP TABLE IF EXISTS `data_extraction_rules`;

CREATE TABLE IF NOT EXISTS `data_service_connections` (
    `id`                    VARCHAR(36)  NOT NULL COMMENT '连接档案ID',
    `name`                  VARCHAR(255) NOT NULL COMMENT '连接名称',
    `description`           TEXT         DEFAULT NULL COMMENT '连接描述',
    `source_type`           VARCHAR(20)  NOT NULL COMMENT '连接类型：mysql|api',
    `status`                VARCHAR(20)  NOT NULL DEFAULT 'enabled' COMMENT '状态：enabled|disabled',
    `config`                JSON         NOT NULL COMMENT '连接参数JSON（不含敏感凭据）',
    `encrypted_credentials` TEXT         DEFAULT NULL COMMENT 'AES-GCM加密凭据',
    `created_at`            DATETIME     DEFAULT NULL COMMENT '创建时间',
    `updated_at`            DATETIME     DEFAULT NULL COMMENT '更新时间',
    `deleted_at`            DATETIME     DEFAULT NULL COMMENT '软删除时间',
    `created_by`            VARCHAR(100) DEFAULT NULL COMMENT '创建人',
    `updated_by`            VARCHAR(100) DEFAULT NULL COMMENT '更新人',
    PRIMARY KEY (`id`),
    INDEX `idx_data_connection_type` (`source_type`),
    INDEX `idx_data_connection_status` (`status`),
    INDEX `idx_data_connection_deleted` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据服务连接档案';

CREATE TABLE IF NOT EXISTS `data_pipeline_rules` (
    `id`          VARCHAR(36)  NOT NULL COMMENT 'Pipeline规则ID',
    `name`        VARCHAR(255) NOT NULL COMMENT '规则名称',
    `description` TEXT         DEFAULT NULL COMMENT '规则描述',
    `source_type` VARCHAR(20)  NOT NULL COMMENT '采集源类型：file|mysql|api',
    `status`      VARCHAR(20)  NOT NULL DEFAULT 'enabled' COMMENT '状态：enabled|disabled',
    `version`     INT          NOT NULL DEFAULT 1 COMMENT '规则版本号',
    `config`      JSON         NOT NULL COMMENT '规范化PipelineBundle配置JSON',
    `created_at`  DATETIME     DEFAULT NULL COMMENT '创建时间',
    `updated_at`  DATETIME     DEFAULT NULL COMMENT '更新时间',
    `deleted_at`  DATETIME     DEFAULT NULL COMMENT '软删除时间',
    `created_by`  VARCHAR(100) DEFAULT NULL COMMENT '创建人',
    `updated_by`  VARCHAR(100) DEFAULT NULL COMMENT '更新人',
    PRIMARY KEY (`id`),
    INDEX `idx_data_pipeline_type` (`source_type`),
    INDEX `idx_data_pipeline_status` (`status`),
    INDEX `idx_data_pipeline_name` (`name`),
    INDEX `idx_data_pipeline_deleted` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据Pipeline规则包';

CREATE TABLE IF NOT EXISTS `data_model_config` (
    `id`                       VARCHAR(36)  NOT NULL COMMENT '配置ID（固定 default）',
    `provider`                 VARCHAR(100) DEFAULT NULL COMMENT '模型供应商',
    `model`                    VARCHAR(200) DEFAULT NULL COMMENT '默认模型',
    `temperature`              DECIMAL(4,2) NOT NULL DEFAULT 0.70 COMMENT 'Temperature',
    `enable_temperature`       TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '是否启用 Temperature',
    `top_p`                    DECIMAL(4,2) NOT NULL DEFAULT 0.90 COMMENT 'Top P',
    `enable_top_p`             TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '是否启用 Top P',
    `max_tokens`               INT          NOT NULL DEFAULT 4096 COMMENT 'Max Tokens',
    `enable_max_tokens`        TINYINT(1)   NOT NULL DEFAULT 0 COMMENT '是否启用 Max Tokens',
    `max_context_count`        INT          NOT NULL DEFAULT 10 COMMENT '上下文轮数',
    `retrieval_top_k`          INT          NOT NULL DEFAULT 5 COMMENT '检索 TopK',
    `retrieval_match_threshold` DECIMAL(4,2) NOT NULL DEFAULT 0.50 COMMENT '检索匹配阈值',
    `created_at`               DATETIME     DEFAULT NULL COMMENT '创建时间',
    `updated_at`               DATETIME     DEFAULT NULL COMMENT '更新时间',
    `created_by`               VARCHAR(100) DEFAULT NULL COMMENT '创建人',
    `updated_by`               VARCHAR(100) DEFAULT NULL COMMENT '更新人',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据平台模型配置';

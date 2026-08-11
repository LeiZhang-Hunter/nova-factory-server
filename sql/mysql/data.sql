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

CREATE TABLE IF NOT EXISTS `data_collectors` (
    `id`                 VARCHAR(36)  NOT NULL COMMENT '采集器ID',
    `name`               VARCHAR(255) NOT NULL COMMENT '设备名称',
    `device_id`          VARCHAR(100) NOT NULL DEFAULT '' COMMENT '采集器上报的设备ID',
    `token`              VARCHAR(255) NOT NULL DEFAULT '' COMMENT '采集器认证Token',
    `status`             VARCHAR(20)  NOT NULL DEFAULT 'offline' COMMENT '状态：online|offline',
    `last_heartbeat`     DATETIME     DEFAULT NULL COMMENT '最后心跳时间',
    `last_connected_at`  DATETIME     DEFAULT NULL COMMENT '最后连接时间',
    `last_checksum`      VARCHAR(64)  DEFAULT NULL COMMENT '最后下发的配置MD5',
    `last_dispatched_at` DATETIME     DEFAULT NULL COMMENT '最后下发时间',
    `version`            VARCHAR(50)  DEFAULT NULL COMMENT '采集器版本',
    `tags`               JSON         DEFAULT NULL COMMENT '标签JSON',
    `created_at`         DATETIME     DEFAULT NULL COMMENT '创建时间',
    `updated_at`         DATETIME     DEFAULT NULL COMMENT '更新时间',
    `deleted_at`         DATETIME     DEFAULT NULL COMMENT '软删除时间',
    `created_by`         VARCHAR(100) DEFAULT NULL COMMENT '创建人',
    `updated_by`         VARCHAR(100) DEFAULT NULL COMMENT '更新人',
    PRIMARY KEY (`id`),
    INDEX `idx_data_collector_device` (`device_id`),
    INDEX `idx_data_collector_status` (`status`),
    INDEX `idx_data_collector_deleted` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据采集器设备';

CREATE TABLE IF NOT EXISTS `data_dispatch_logs` (
    `id`             VARCHAR(36)  NOT NULL COMMENT '记录ID',
    `collector_id`   VARCHAR(36)  NOT NULL COMMENT '采集器ID',
    `collector_name` VARCHAR(255) DEFAULT NULL COMMENT '采集器名称',
    `rule_ids`       JSON         DEFAULT NULL COMMENT '本次下发的规则ID列表',
    `rule_names`     JSON         DEFAULT NULL COMMENT '本次下发的规则名称列表',
    `dispatch_no`    INT          NOT NULL DEFAULT 1 COMMENT '该采集器的下发序号',
    `config`         TEXT         NOT NULL COMMENT '合并后的YAML配置快照',
    `config_md5`     CHAR(32)     NOT NULL COMMENT '合并配置MD5',
    `created_at`     DATETIME     DEFAULT NULL COMMENT '创建时间',
    `created_by`     VARCHAR(100) DEFAULT NULL COMMENT '操作人',
    PRIMARY KEY (`id`),
    INDEX `idx_dl_collector` (`collector_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='采集器规则下发版本记录';



CREATE TABLE IF NOT EXISTS nova_metrics_device ON CLUSTER `default`
(
    `device_id` UInt64 COMMENT '设备id',
    `template_id` UInt64 COMMENT '设备模板id',
    `data_id` bigint(20) NOT NULL COMMENT '数据id',
    series_id UInt64 COMMENT '序列id',
    attributes            Map(String, String)  COMMENT '属性' CODEC (ZSTD(1)),
    start_time_unix         DateTime64(9) COMMENT '开始时间' CODEC (Delta, ZSTD(1))  ,
    time_unix              DateTime64(9) COMMENT '当前时间' CODEC (Delta, ZSTD(1)) ,
    value                 Float64 COMMENT '统计值' CODEC (ZSTD(1)) ,
    INDEX idx_attr_key mapKeys(attributes) TYPE bloom_filter(0.01) GRANULARITY 1
    ) ENGINE = MergeTree
    PARTITION BY toDate(time_unix)
    ORDER BY (device_id, template_id, data_id, series_id)
    TTL toDateTime(time_unix) + toIntervalDay(180)
    SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1;

CREATE TABLE IF NOT EXISTS nova_control_log
(
    `device_id` UInt64 COMMENT '设备id',
    `data_id` bigint(20) NOT NULL COMMENT '数据id',
    `message` Nullable(String),
    `type` Nullable(String),
    series_id UInt64 COMMENT '序列id',
    attributes            Map(String, String)  COMMENT '属性' CODEC (ZSTD(1)),
    start_time_unix         DateTime64(9) COMMENT '开始时间' CODEC (Delta, ZSTD(1))  ,
    time_unix              DateTime64(9) COMMENT '当前时间' CODEC (Delta, ZSTD(1)) ,
    INDEX idx_attr_key mapKeys(attributes) TYPE bloom_filter(0.01) GRANULARITY 1
    ) ENGINE = MergeTree
    PARTITION BY toDate(time_unix)
    ORDER BY (device_id, data_id, series_id)
    TTL toDateTime(time_unix) + toIntervalDay(30)
    SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1;

CREATE TABLE IF NOT EXISTS nova_alert_log
(
    `object_id` UInt64 COMMENT '告警对象id',
    `gateway_id` UInt64 COMMENT '设备id',
    `device_id` UInt64 COMMENT '设备id',
    `device_template_id` UInt64 COMMENT '设备id',
    `device_data_id` UInt64 COMMENT '设备id',
    `alert_id` UInt64 COMMENT '设备id',
    `series_id` UInt64 COMMENT '序列id',
    `context` Nullable(String),
    `reason` Nullable(String),
    `message` Nullable(String),
    `data` Nullable(String),
    start_time_unix         DateTime64(9) COMMENT '开始时间' CODEC (Delta, ZSTD(1))  ,
    time_unix              DateTime64(9) COMMENT '当前时间' CODEC (Delta, ZSTD(1)) ,
    INDEX idx_device_id device_id TYPE minmax GRANULARITY 5,
    ) ENGINE = MergeTree
    PARTITION BY toDate(time_unix)
    ORDER BY (object_id, start_time_unix)
    TTL toDateTime(time_unix) + toIntervalDay(30)
    SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1;


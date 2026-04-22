CREATE TABLE IF NOT EXISTS `erp_sales_order_agent_config` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name` VARCHAR(128) NOT NULL COMMENT '配置名称',
    `agent_id` VARCHAR(128) NOT NULL COMMENT '智能体AgentID',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用，0禁用',
    `remark` VARCHAR(500) DEFAULT NULL COMMENT '备注',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='销售订单智能体配置表';

CREATE TABLE IF NOT EXISTS `erp_integration_config` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `type` VARCHAR(128) NOT NULL COMMENT '接入类型',
    `data` TEXT  COMMENT '配置快照',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用，0禁用',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` bigint(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` bigint(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `type`(`type`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='erp集成系统';

CREATE TABLE IF NOT EXISTS `erp_order_audit` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `tid` VARCHAR(50) NOT NULL COMMENT '网店订单编号',
    `weight` DECIMAL(18,4) DEFAULT NULL COMMENT '订单重量',
    `size` DECIMAL(18,4) DEFAULT NULL COMMENT '订单尺寸/体积',

    `buyer_nick` VARCHAR(50) DEFAULT NULL COMMENT '买家账号',
    `buyer_message` VARCHAR(50) DEFAULT NULL COMMENT '买家留言',
    `seller_memo` VARCHAR(50) DEFAULT NULL COMMENT '卖家备注',

    `total` DECIMAL(20,8) DEFAULT NULL COMMENT '订单总金额（含运费）',
    `privilege` DECIMAL(20,8) NOT NULL DEFAULT 0.00000000 COMMENT '订单优惠金额',
    `post_fee` DECIMAL(20,8) NOT NULL DEFAULT 0.00000000 COMMENT '运费',

    `receiver_name` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '收货人名称',
    `receiver_province` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '收货省编码',
    `receiver_province_name` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '收货省名称',
    `receiver_city` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '收货市编码',
    `receiver_city_name` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '收货市名称',
    `receiver_district` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '收货区编码',
    `receiver_district_name` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '收货区名称',
    `receiver_street` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '收货街道编码',
    `receiver_street_name` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '收货街道名称',
    `receiver_address` VARCHAR(255) NOT NULL COMMENT '收货详细地址',
    `receiver_phone` VARCHAR(50) DEFAULT NULL COMMENT '收货电话',
    `receiver_mobile` VARCHAR(50) NOT NULL COMMENT '收货人手机号',
    `receiver_zip` VARCHAR(20) DEFAULT NULL COMMENT '邮编',

    `status` VARCHAR(200) NOT NULL COMMENT '来源订单状态',
    `order_type` VARCHAR(50) NOT NULL COMMENT '订单类型（Cod/NoCod）',
    `invoice_name` VARCHAR(50) DEFAULT NULL COMMENT '发票抬头',
    `seller_flag` VARCHAR(50) DEFAULT NULL COMMENT '卖家旗帜',
    `pay_time` DATETIME(0) DEFAULT NULL COMMENT '付款时间',

    `logist_b_type_code` VARCHAR(100) DEFAULT NULL COMMENT '物流公司编码',
    `logist_bill_code` VARCHAR(100) DEFAULT NULL COMMENT '物流单号',
    `b_type_code` VARCHAR(50) DEFAULT NULL COMMENT '往来单位编码',

    `details_json` LONGTEXT COMMENT '订单商品明细JSON',
    `accounts_json` LONGTEXT COMMENT '订单账户信息JSON',
    `source_json` LONGTEXT COMMENT '原始订单报文JSON',

    `audit_status` TINYINT NOT NULL DEFAULT 0 COMMENT '审核状态：0待审核 1审核通过 2审核驳回',
    `audit_remark` VARCHAR(500) DEFAULT NULL COMMENT '审核意见',
    `audit_by` BIGINT(20) NULL DEFAULT NULL COMMENT '审核人',
    `audit_time` DATETIME(0) DEFAULT NULL COMMENT '审核时间',

    `transfer_status` TINYINT NOT NULL DEFAULT 0 COMMENT '入正式订单状态：0未入库 1已入库 2入库失败',
    `transfer_message` VARCHAR(500) DEFAULT NULL COMMENT '入正式订单结果描述',
    `transfer_time` DATETIME(0) DEFAULT NULL COMMENT '入正式订单时间',
    `erp_order_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '审核通过后生成的正式订单ID',

    `dept_id` BIGINT(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` BIGINT(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` DATETIME(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` BIGINT(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` DATETIME(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` TINYINT(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',

    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_erp_order_audit_tid` (`tid`) USING BTREE,
    INDEX `idx_erp_order_audit_status` (`audit_status`) USING BTREE,
    INDEX `idx_erp_order_audit_transfer_status` (`transfer_status`) USING BTREE,
    INDEX `idx_erp_order_audit_erp_order_id` (`erp_order_id`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ERP订单审核表';

CREATE TABLE IF NOT EXISTS `erp_order` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `tid` VARCHAR(50) NOT NULL COMMENT '网店订单编号',
    `weight` DECIMAL(18,4) DEFAULT NULL COMMENT '订单重量',
    `size` DECIMAL(18,4) DEFAULT NULL COMMENT '订单尺寸/体积',

    `buyer_nick` VARCHAR(50) DEFAULT NULL COMMENT '买家账号',
    `buyer_message` VARCHAR(50) DEFAULT NULL COMMENT '买家留言',
    `seller_memo` VARCHAR(50) DEFAULT NULL COMMENT '卖家备注',

    `total` DECIMAL(20,8) DEFAULT NULL COMMENT '订单总金额（含运费）',
    `privilege` DECIMAL(20,8) NOT NULL DEFAULT 0.00000000 COMMENT '订单优惠金额',
    `post_fee` DECIMAL(20,8) NOT NULL DEFAULT 0.00000000 COMMENT '运费',

    `receiver_name` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '收货人名称',

    `receiver_province` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '收货省编码',
    `receiver_province_name` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '收货省名称',

    `receiver_city` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '收货市编码',
    `receiver_city_name` VARCHAR(100) NOT NULL DEFAULT '' COMMENT '收货市名称',

    `receiver_district` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '收货区编码',
    `receiver_district_name` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '收货区名称',

    `receiver_street` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '收货街道编码',
    `receiver_street_name` VARCHAR(50) NOT NULL DEFAULT '' COMMENT '收货街道名称',

    `receiver_address` VARCHAR(255) NOT NULL COMMENT '收货详细地址',
    `receiver_phone` VARCHAR(50) DEFAULT NULL COMMENT '收货电话',
    `receiver_mobile` VARCHAR(50) NOT NULL COMMENT '收货人手机号',
    `receiver_zip` VARCHAR(20) DEFAULT NULL COMMENT '邮编',

    `status` VARCHAR(200) NOT NULL COMMENT '订单状态',
    `order_type` VARCHAR(50) NOT NULL COMMENT '订单类型（Cod/NoCod）',
    `invoice_name` VARCHAR(50) DEFAULT NULL COMMENT '发票抬头',
    `seller_flag` VARCHAR(50) DEFAULT NULL COMMENT '卖家旗帜',
    `pay_time` DATETIME(0) DEFAULT NULL COMMENT '付款时间',

    `logist_b_type_code` VARCHAR(100) DEFAULT NULL COMMENT '物流公司编码',
    `logist_bill_code` VARCHAR(100) DEFAULT NULL COMMENT '物流单号',
    `b_type_code` VARCHAR(50) DEFAULT NULL COMMENT '往来单位编码',

    `details_json` LONGTEXT COMMENT '订单商品明细JSON',
    `accounts_json` LONGTEXT COMMENT '订单账户信息JSON',

    `bill_code` VARCHAR(100) DEFAULT NULL COMMENT '管家婆返回订单编号',
    `sync_message` VARCHAR(500) DEFAULT NULL COMMENT '同步结果描述',
    `sync_status` TINYINT NOT NULL DEFAULT 0 COMMENT '同步状态：0待同步 1成功 2失败',
    `sync_time` DATETIME(0) DEFAULT NULL COMMENT '同步时间',

    `dept_id` BIGINT(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` BIGINT(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` DATETIME(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` BIGINT(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` DATETIME(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` TINYINT(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',

    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_erp_order_tid` (`tid`) USING BTREE,
    INDEX `idx_erp_order_status` (`status`) USING BTREE,
    INDEX `idx_erp_order_sync_status` (`sync_status`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ERP订单主表';


CREATE TABLE IF NOT EXISTS `erp_order_detail` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `order_id` BIGINT UNSIGNED NOT NULL COMMENT '订单主表ID',
    `tid` VARCHAR(50) NOT NULL COMMENT '网店订单编号',
    `oid` VARCHAR(50) NOT NULL COMMENT '网店订单明细编号',

    `barcode` VARCHAR(50) DEFAULT NULL COMMENT '商品条码',
    `eshop_goods_id` VARCHAR(100) DEFAULT NULL COMMENT '网店商品ID',
    `outer_iid` VARCHAR(100) DEFAULT NULL COMMENT '网店商家编码',
    `eshop_goods_name` VARCHAR(100) NOT NULL COMMENT '网店商品名称',
    `eshop_sku_id` VARCHAR(100) DEFAULT NULL COMMENT '网店商品SKUID',
    `eshop_sku_name` VARCHAR(100) DEFAULT NULL COMMENT '网店商品SKU名称',

    `num_iid` BIGINT(20) DEFAULT NULL COMMENT '商品ID（方案B）',
    `sku_id` BIGINT(20) DEFAULT NULL COMMENT '规格ID（方案B）',

    `num` DECIMAL(18,4) NOT NULL COMMENT '基本单位数量',
    `payment` DECIMAL(20,8) NOT NULL COMMENT '商品总额',
    `pic_path` VARCHAR(255) DEFAULT NULL COMMENT '商品图片路径',
    `weight` DECIMAL(18,4) DEFAULT NULL COMMENT '重量',
    `size` DECIMAL(18,4) DEFAULT NULL COMMENT '尺寸/体积',

    `unit_id` BIGINT(20) DEFAULT NULL COMMENT '销售单位ID',
    `unit_qty` DECIMAL(18,4) DEFAULT NULL COMMENT '销售单位数量',

    `dept_id` BIGINT(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` BIGINT(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` DATETIME(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` BIGINT(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` DATETIME(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` TINYINT(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',

    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_erp_order_detail_oid` (`oid`) USING BTREE,
    INDEX `idx_erp_order_detail_order_id` (`order_id`) USING BTREE,
    INDEX `idx_erp_order_detail_tid` (`tid`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ERP订单明细表';


CREATE TABLE IF NOT EXISTS `erp_order_account` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `order_id` BIGINT UNSIGNED NOT NULL COMMENT '订单主表ID',
    `tid` VARCHAR(50) NOT NULL COMMENT '网店订单编号',
    `finance_code` VARCHAR(50) NOT NULL COMMENT '账户编码',
    `total` DECIMAL(20,8) NOT NULL COMMENT '收款金额',

    `dept_id` BIGINT(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` BIGINT(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` DATETIME(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` BIGINT(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` DATETIME(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` TINYINT(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',

    PRIMARY KEY (`id`),
    INDEX `idx_erp_order_account_order_id` (`order_id`) USING BTREE,
    INDEX `idx_erp_order_account_tid` (`tid`) USING BTREE,
    INDEX `idx_erp_order_account_finance_code` (`finance_code`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ERP订单账户表';


CREATE TABLE IF NOT EXISTS `erp_logistics_company` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `code` VARCHAR(64) NOT NULL COMMENT '物流公司编码',
    `company` VARCHAR(128) NOT NULL COMMENT '关联公司',
    `name` VARCHAR(128) NOT NULL COMMENT '物流公司名称',
    `short_name` VARCHAR(64) DEFAULT NULL COMMENT '物流公司简称',
    `contact_name` VARCHAR(64) DEFAULT NULL COMMENT '联系人',
    `contact_phone` VARCHAR(32) DEFAULT NULL COMMENT '联系电话',
    `address` VARCHAR(255) DEFAULT NULL COMMENT '联系地址',
    
    `province_code` VARCHAR(32) DEFAULT NULL COMMENT '省编码',
    `province_name` VARCHAR(64) NOT NULL COMMENT '省名称',
    `city_code` VARCHAR(32) DEFAULT NULL COMMENT '市编码',
    `city_name` VARCHAR(64) NOT NULL COMMENT '市名称',
    `district_code` VARCHAR(32) DEFAULT NULL COMMENT '区编码',
    `district_name` VARCHAR(64) NOT NULL COMMENT '区名称',
    `street_code` VARCHAR(32) DEFAULT NULL COMMENT '街道编码',
    `street_name` VARCHAR(64) DEFAULT NULL COMMENT '街道名称',

    `remark` VARCHAR(500) DEFAULT NULL COMMENT '备注',
    `sort` INT NOT NULL DEFAULT 0 COMMENT '排序值',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1启用，0禁用',
    `dept_id` BIGINT(20) NULL DEFAULT NULL COMMENT '部门ID',
    `create_by` BIGINT(20) NULL DEFAULT NULL COMMENT '创建者',
    `create_time` DATETIME(0) NULL DEFAULT NULL COMMENT '创建时间',
    `update_by` BIGINT(20) NULL DEFAULT NULL COMMENT '更新者',
    `update_time` DATETIME(0) NULL DEFAULT NULL COMMENT '更新时间',
    `state` TINYINT(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_erp_logistics_company_code` (`code`) USING BTREE,
    UNIQUE INDEX `uk_erp_logistics_company_name` (`name`) USING BTREE,
    INDEX `idx_erp_logistics_company_status` (`status`) USING BTREE,
    INDEX `idx_erp_logistics_company_sort` (`sort`) USING BTREE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ERP物流公司管理表';

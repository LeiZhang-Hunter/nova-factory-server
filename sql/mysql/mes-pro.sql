
-- ----------------------------
-- 1、工艺路线表
-- ----------------------------
create table sys_craft_route (
   route_id                       bigint(20)      not null auto_increment    comment '工艺路线ID',
   route_code                     varchar(64)     not null                   comment '工艺路线编号',
   loop_time                      int(11)         default 30                  comment '循环执行时间',
   route_name                     varchar(255)    not null                   comment '工艺路线名称',
   route_desc                     varchar(500)                               comment '工艺路线说明',
   remark                         varchar(500)    default ''                 comment '备注',
   attr1                          varchar(64)     default null               comment '预留字段1',
   attr2                          varchar(255)    default null               comment '预留字段2',
   attr3                          int(11)         default 0                  comment '预留字段3',
   attr4                          int(11)         default 0                  comment '预留字段4',
   `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
   create_by                      varchar(64)     default ''                 comment '创建者',
   create_time 	                 datetime                                   comment '创建时间',
   update_by                      varchar(64)     default ''                 comment '更新者',
   `status` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 1异常）',
   update_time                    datetime                                   comment '更新时间',
   `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
   primary key (route_id)
) engine=innodb auto_increment=0 comment = '工艺路线表';

-- ----------------------------
-- 2、生产工序表
-- ----------------------------
create table sys_pro_process (
     process_id                     bigint(20)      not null auto_increment    comment '工序ID',
     process_code                   varchar(64)     not null                   comment '工序编码',
     process_name                   varchar(255)    not null                   comment '工序名称',
     attention                      varchar(1000)                              comment '工艺要求',
     remark                         varchar(500)    default ''                 comment '备注',
     attr1                          varchar(64)     default null               comment '预留字段1',
     attr2                          varchar(255)    default null               comment '预留字段2',
     attr3                          int(11)         default 0                  comment '预留字段3',
     attr4                          int(11)         default 0                  comment '预留字段4',
     `status`                       tinyint(1) NULL DEFAULT 0 COMMENT '是否启用（0禁用 1启用）',
     create_by                      varchar(64)     default ''                 comment '创建者',
     create_time 	                 datetime                                   comment '创建时间',
     update_by                      varchar(64)     default ''                 comment '更新者',
     update_time                    datetime                                   comment '更新时间',
     `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
     `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
     primary key (process_id)
) engine=innodb auto_increment=0 comment = '生产工序表';


-- ----------------------------
-- 3、生产工序内容表
-- ----------------------------
create table sys_pro_process_content (
    content_id                     bigint(20)      not null auto_increment    comment '内容ID',
    process_id                     bigint(20)      not null                   comment '工序ID',
    control_name                   varchar(255)                               comment '控制名称',
    order_num                      int(4)          default 0                  comment '顺序编号',
    content_text                   varchar(500)                               comment '内容说明',
    device                         varchar(255)                               comment '辅助设备',
    material                       varchar(255)                               comment '辅助材料',
    control_type                   varchar(255)                               comment '控制方式',
    doc_url                        varchar(255)                               comment '材料URL',
    remark                         varchar(500)    default ''                 comment '备注',
    attr1                          varchar(64)     default null               comment '预留字段1',
    attr2                          varchar(255)    default null               comment '预留字段2',
    attr3                          int(11)         default 0                  comment '预留字段3',
    attr4                          int(11)         default 0                  comment '预留字段4',
    `extension` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL  COMMENT '扩展字段',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    create_by                      varchar(64)     default ''                 comment '创建者',
    create_time 	                 datetime                                   comment '创建时间',
    update_by                      varchar(64)     default ''                 comment '更新者',
    update_time                    datetime                                   comment '更新时间',
    primary key (content_id)
) engine=innodb auto_increment=0 comment = '生产工序内容表';
alter table sys_pro_process_content add column control_type                   varchar(255)                               comment '控制方式';


-- ----------------------------
-- 4、工艺组成表
-- ----------------------------
create table sys_pro_route_process (
   record_id                      bigint(20)      not null auto_increment    comment '记录ID',
   route_id                       bigint(20)      not null                   comment '工艺路线ID',
   process_id                     bigint(20)      not null                   comment '工序ID',
   process_code                   varchar(64)                                comment '工序编码',
   process_name                   varchar(255)                               comment '工序名称',
   order_num                      int(4)          default 1                  comment '序号',
   next_process_id                bigint(20)      not null                   comment '工序ID',
   next_process_code              varchar(64)                                comment '工序编码',
   next_process_name              varchar(255)                               comment '工序名称',
   link_type                      varchar(64)     default 'SS'               comment '与下一道工序关系',
   default_pre_time               int(11)         default 0                  comment '准备时间',
   default_suf_time               int(11)         default 0                  comment '等待时间',
   color_code                     char(7)         default '#00AEF3'          comment '甘特图显示颜色',
   key_flag                       varchar(64)     default 'N'                comment '关键工序',
   is_check                       char(1)         default 'N'                comment '是否检验',
   remark                         varchar(500)    default ''                 comment '备注',
   attr1                          varchar(64)     default null               comment '预留字段1',
   attr2                          varchar(255)    default null               comment '预留字段2',
   attr3                          int(11)         default 0                  comment '预留字段3',
   attr4                          int(11)         default 0                  comment '预留字段4',
   `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
   `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
   create_by                      varchar(64)     default ''                 comment '创建者',
   create_time 	                 datetime                                   comment '创建时间',
   update_by                      varchar(64)     default ''                 comment '更新者',
   update_time                    datetime                                   comment '更新时间',
   primary key (record_id)
) engine=innodb auto_increment=0 comment = '工艺组成表';

-- ----------------------------
-- 5、产品制程
-- ----------------------------
create table sys_pro_route_product (
    record_id                      bigint(20)      not null auto_increment    comment '记录ID',
    route_id                       bigint(20)      not null                   comment '工艺路线ID',
    item_id                        bigint(20)      not null                   comment '产品物料ID',
    item_code                      varchar(64)     not null                   comment '产品物料编码',
    item_name                      varchar(255)    not null                   comment '产品物料名称',
    specification                  varchar(500)    default null               comment '规格型号',
    unit_of_measure                varchar(64)     not null                   comment '单位',
    unit_name                      varchar(64)                                comment '单位名称',
    quantity                       int(11)         default 1                  comment '生产数量',
    production_time                double(12,2)    default 1                  comment '生产用时',
    time_unit_type                 varchar(64)     default 'MINUTE'           comment '时间单位',
    remark                         varchar(500)    default ''                 comment '备注',
    attr1                          varchar(64)     default null               comment '预留字段1',
    attr2                          varchar(255)    default null               comment '预留字段2',
    attr3                          int(11)         default 0                  comment '预留字段3',
    attr4                          int(11)         default 0                  comment '预留字段4',
    `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
    `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
    create_by                      varchar(64)     default ''                 comment '创建者',
    create_time 	                 datetime                                   comment '创建时间',
    update_by                      varchar(64)     default ''                 comment '更新者',
    update_time                    datetime                                   comment '更新时间',
    primary key (record_id)
) engine=innodb auto_increment=0 comment = '产品制程';

-- ----------------------------
-- 6、产品制程物料BOM表
-- ----------------------------
create table pro_route_product_bom (
   record_id                      bigint(20)      not null auto_increment    comment '记录ID',
   route_id                       bigint(20)      not null                   comment '工艺路线ID',
   process_id                     bigint(20)      not null                   comment '工序ID',
   product_id                     bigint(20)      not null                   comment '产品BOM中的唯一ID',
   item_id                        bigint(20)      not null                   comment '产品物料ID',
   item_code                      varchar(64)     not null                   comment '产品物料编码',
   item_name                      varchar(255)    not null                   comment '产品物料名称',
   specification                  varchar(500)    default null               comment '规格型号',
   unit_of_measure                varchar(64)     not null                   comment '单位',
   unit_name                      varchar(64)                                comment '单位名称',
   quantity                       double(12,2)         default 1             comment '用料比例',
   remark                         varchar(500)    default ''                 comment '备注',
   attr1                          varchar(64)     default null               comment '预留字段1',
   attr2                          varchar(255)    default null               comment '预留字段2',
   attr3                          int(11)         default 0                  comment '预留字段3',
   attr4                          int(11)         default 0                  comment '预留字段4',
   `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
   `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
   create_by                      varchar(64)     default ''                 comment '创建者',
   create_time 	                 datetime                                   comment '创建时间',
   update_by                      varchar(64)     default ''                 comment '更新者',
   update_time                    datetime                                   comment '更新时间',
   primary key (record_id)
) engine=innodb auto_increment=0 comment = '产品制程物料BOM表';


 create table sys_pro_workorder (
   workorder_id                bigint(20)      not null auto_increment    comment '工单ID',
   workorder_code              varchar(64)     not null                   comment '工单编码',
   workorder_name              varchar(255)    not null                   comment '工单名称',
   workorder_type              varchar(64)     default 'SELF'             comment '工单类型',
   order_source                varchar(64)     not null                   comment '来源类型',
   source_code                 varchar(64)                                comment '来源单据',
   product_id                  bigint(20)      not null                   comment '产品ID',
   product_code                varchar(64)     not null                   comment '产品编号',
   product_name                varchar(255)    not null                   comment '产品名称',
   product_spc                 varchar(255)                               comment '规格型号',
   unit_of_measure             varchar(64)     not null                   comment '单位',
   quantity                    double(14,2)    default 0 not null         comment '生产数量',
   quantity_produced           double(14,2)    default 0                  comment '已生产数量',
   quantity_changed            double(14,2)    default 0                  comment '调整数量',
   quantity_scheduled          double(14,2)    default 0                  comment '已排产数量',
   client_id                   bigint(20)                                 comment '客户ID',
   client_code                 varchar(64)                                comment '客户编码',
   client_name                 varchar(255)                               comment '客户名称',
   vendor_id                   bigint(20)                                 comment '供应商ID',
   vendor_code                 varchar(64)                                comment '供应商编号',
   vendor_name                 varchar(255)                               comment '供应商名称',
   batch_code                  varchar(64)                                comment '批次号',
   request_date                datetime        not null                   comment '需求日期',
   parent_id                   bigint(20)      default 0 not null         comment '父工单',
   ancestors                   varchar(500)    not null                   comment '所有父节点ID',
   finish_date                 datetime                                   comment '完成时间',
   status                      varchar(64)     default 'PREPARE'          comment '单据状态',
   remark                      varchar(500)    default ''                 comment '备注',
   attr1                       varchar(64)     default null               comment '预留字段1',
   attr2                       varchar(255)    default null               comment '预留字段2',
   attr3                       int(11)         default 0                  comment '预留字段3',
   attr4                       int(11)         default 0                  comment '预留字段4',
   `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
   `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
   create_by                   varchar(64)     default ''                 comment '创建者',
   create_time 	              datetime                                   comment '创建时间',
   update_by                   varchar(64)     default ''                 comment '更新者',
   update_time                 datetime                                   comment '更新时间',
   primary key (workorder_id)
) engine=innodb auto_increment=200 comment = '生产工单表';

create table sys_pro_task (
  task_id                        bigint(20)      not null auto_increment    comment '任务ID',
  task_code                      varchar(64)     not null                   comment '任务编号',
  task_name                      varchar(255)    not null                   comment '任务名称',
  workorder_id                   bigint(20)      not null                   comment '生产工单ID',
  workorder_code                 varchar(64)     not null                   comment '生产工单编号',
  workorder_name                 varchar(255)    not null                   comment '工单名称',
  workstation_id                 bigint(20)      not null                   comment '工作站ID',
  gateway_id                   bigint(20)      not null                   comment '网关id',
  workstation_code               varchar(64)     not null                   comment '工作站编号',
  workstation_name               varchar(255)    not null                   comment '工作站名称',
  route_id                       bigint(20)      not null                   comment '工艺ID',
  route_code                     varchar(64)                                comment '工艺编号',
  process_id                     bigint(20)      not null                   comment '工序ID',
  process_code                   varchar(64)                                comment '工序编码',
  process_name                   varchar(255)                               comment '工序名称',
  item_id                        bigint(20)      not null                   comment '产品物料ID',
  item_code                      varchar(64)     not null                   comment '产品物料编码',
  item_name                      varchar(255)    not null                   comment '产品物料名称',
  specification                  varchar(500)                               comment '规格型号',
  unit_of_measure                varchar(64)     not null                   comment '单位',
  unit_name                      varchar(64)                                comment '单位名称',
  quantity                       double(14,2)    default 1 not null         comment '排产数量',
  quantity_produced              double(14,2)    default 0                  comment '已生产数量',
  quantity_quanlify              double(14,2)    default 0                  comment '合格品数量',
  quantity_unquanlify            double(14,2)    default 0                  comment '不良品数量',
  quantity_changed               double(14,2)    default 0                  comment '调整数量',
  client_id                      bigint(20)                                 comment '客户ID',
  client_code                    varchar(64)                                comment '客户编码',
  client_name                    varchar(255)                               comment '客户名称',
  client_nick                    varchar(255)                               comment '客户简称',
  start_time                     datetime        default CURRENT_TIMESTAMP  comment '开始生产时间',
  duration                       int(11)         default 1                  comment '生产时长',
  end_time                       datetime                                   comment '完成生产时间',
  color_code                     char(7)         default '#00AEF3'          comment '甘特图显示颜色',
  request_date                   datetime                                   comment '需求日期',
  status                         int     default 0           comment '生产状态 0 是暂停 1 是启动',
  remark                         varchar(500)    default ''                 comment '备注',
  attr1                          varchar(64)     default null               comment '预留字段1',
  attr2                          varchar(255)    default null               comment '预留字段2',
  attr3                          int(11)         default 0                  comment '预留字段3',
  attr4                          int(11)         default 0                  comment '预留字段4',
  `dept_id` bigint(20) NULL DEFAULT NULL COMMENT '部门ID',
  `state` tinyint(1) NULL DEFAULT 0 COMMENT '操作状态（0正常 -1删除）',
  create_by                      varchar(64)     default ''                 comment '创建者',
  create_time 	                 datetime                                   comment '创建时间',
  update_by                      varchar(64)     default ''                 comment '更新者',
  update_time                    datetime                                   comment '更新时间',
  primary key (task_id)
) engine=innodb auto_increment=200 comment = '生产任务表';


CREATE TABLE `sys_craft_route_config` (
  `route_config_id` bigint NOT NULL COMMENT '工艺路线配置id',
  `route_id` bigint NOT NULL COMMENT '工艺路线ID',
  `context` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '原始配置',
  `config` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '下发配置',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`route_config_id`),
  UNIQUE KEY `route_id` (`route_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='工艺路线表';

CREATE TABLE `sys_product_schedule` (
   `id` bigint(20) NOT NULL COMMENT '调度id',
   `gateway_id` bigint(20) NOT NULL  COMMENT '网关id',
   `schedule_name` varchar(255) NOT NULL COMMENT '计划名称',
   `time` varchar(64) NOT NULL COMMENT '时间序列化格式,普通日程,1,2,3,4,5;特殊日程:2025-04-04 ~ 2025-04-04',
   `time_manager` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '时间任务列表',
   `schedule_type` tinyint(1) NOT NULL COMMENT '0为普通日程 1为特殊日程',
   `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '操作状态（0正常 1启动）',
   `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
   `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
   `create_time` datetime DEFAULT NULL COMMENT '创建时间',
   `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
   `update_time` datetime DEFAULT NULL COMMENT '更新时间',
   `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
   PRIMARY KEY (`id`),
   KEY `schedule_type` (`schedule_type`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8;


CREATE TABLE `sys_product_schedule_map` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `schedule_id` bigint(20) NOT NULL COMMENT '计划标识',
    `gateway_id` bigint(20) NOT NULL  COMMENT '网关id',
    `begin_time` int(11) NOT NULL DEFAULT '0' COMMENT '开始时间',
    `end_time` int(11) NOT NULL DEFAULT '0' COMMENT '结束时间',
    `date` tinyint(1) NOT NULL DEFAULT '0' COMMENT '日期 1 2 3 4 5 6 7分别代表周几',
    `craft_route_id` bigint(20) NOT NULL  COMMENT '工艺路线id',
    `schedule_type` tinyint(1) NOT NULL COMMENT '0为 循环日程 1为特殊日程',
    `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '操作状态（0正常 1启动）',
    `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
    `create_time` datetime DEFAULT NULL COMMENT '创建时间',
    `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
    `update_time` datetime DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NOT NULL DEFAULT '1',
    PRIMARY KEY (`id`),
    KEY `schedule_id` (`schedule_id`) USING BTREE,
    KEY `begin_time` (`begin_time`)
) ENGINE=InnoDB AUTO_INCREMENT=39797 DEFAULT CHARSET=utf8;
-- alter table sys_product_schedule_map add column   `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '操作状态（0正常 1启动）';
-- ----------------------------
-- 班次设置
-- ----------------------------

CREATE TABLE `sys_work_shift_setting` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `name`  varchar(255)    not null   comment '班次名称',
    `begin_time` int(11) NOT NULL DEFAULT '0' COMMENT '开始时间',
    `begin_time_str`  varchar(255)    not null   comment '开始时间字符串',
    `end_time` int(11) NOT NULL DEFAULT '0' COMMENT '结束时间',
    `end_time_str`  varchar(255)    not null   comment '结束时间字符串',
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用班次设置',
    `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
    `create_time` datetime DEFAULT NULL COMMENT '创建时间',
    `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
    `update_time` datetime DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;


-- ----------------------------
-- 化验单管理
-- ----------------------------
CREATE TABLE `sys_product_laboratory` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `type` tinyint(1) NOT NULL COMMENT '0 为物料化验 1为 配料化验',
    `material`  varchar(255)    not null   comment '物料',
    `contact`  varchar(255)    not null   comment '采样人',
    `date`  varchar(255)    not null   comment '采样时间',
    `address`  varchar(255)    not null   comment '采样地点',
    `dry_high` decimal(10, 2) NOT NULL DEFAULT '0' COMMENT '干燥基高位',
    `received_low` decimal(10, 2) NOT NULL DEFAULT '0' COMMENT '收到基低位',
    `heat` decimal(10, 2) NOT NULL DEFAULT '0' COMMENT '含热量',
    `sulphur`  decimal(10, 2) NOT NULL DEFAULT '0'   comment '含硫量',
    `volatility`  decimal(10, 2) NOT NULL DEFAULT '0'   comment '挥发性',
    `water`  decimal(10, 2) NOT NULL DEFAULT '0'   comment '含水量',
    `weight` decimal(10, 2) NOT NULL DEFAULT '0'   comment '样重',
    `number` varchar(255)    not null   comment '样号',
    `img`  varchar(512)    not null   comment '图片地址',
    `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
    `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
    `create_time` datetime DEFAULT NULL COMMENT '创建时间',
    `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
    `update_time` datetime DEFAULT NULL COMMENT '更新时间',
    `state` tinyint(1) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY `idx_type` (`type`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;





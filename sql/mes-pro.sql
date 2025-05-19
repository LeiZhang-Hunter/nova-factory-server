
-- ----------------------------
-- 1、工艺路线表
-- ----------------------------
create table sys_craft_route (
   route_id                       bigint(20)      not null auto_increment    comment '工艺路线ID',
   route_code                     varchar(64)     not null                   comment '工艺路线编号',
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
    order_num                      int(4)          default 0                  comment '顺序编号',
    content_text                   varchar(500)                               comment '内容说明',
    device                         varchar(255)                               comment '辅助设备',
    material                       varchar(255)                               comment '辅助材料',
    doc_url                        varchar(255)                               comment '材料URL',
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
    primary key (content_id)
) engine=innodb auto_increment=0 comment = '生产工序内容表';

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
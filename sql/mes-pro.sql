
-- ----------------------------
-- 3、工艺路线表
-- ----------------------------
drop table if exists sys_craft_route;
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
) engine=innodb auto_increment=200 comment = '工艺路线表';


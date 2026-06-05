-- ModernWMS schema generated from current EF entity classes.
-- SQL dialect: SQLite-compatible DDL.
-- 说明：本文件补充了表用途说明、字段注释和建议索引。
PRAGMA foreign_keys = ON;

-- ============================================================
-- 表名: action_log
-- 用途: 操作日志表，记录用户在前端页面发起的业务操作。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "action_log" (
    -- 主键 ID，自增。
                              "id" INTEGER NOT NULL CONSTRAINT "PK_action_log" PRIMARY KEY AUTOINCREMENT
    -- 前端页面路径。
    ,"vue_path" TEXT NOT NULL
    -- 用户名。
    ,"user_name" TEXT NOT NULL
    -- 操作内容。
    ,"action_content" TEXT NOT NULL
    -- 操作时间。
    ,"action_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: action_log
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_action_log_user_name" ON "action_log" ("user_name");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_action_log_action_time" ON "action_log" ("action_time");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_action_log_tenant_id" ON "action_log" ("tenant_id");

-- ============================================================
-- 表名: asnmaster
-- 用途: 入库单主表，保存一张 ASN 单据的主信息。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "asnmaster" (
    -- 主键 ID，自增。
                             "id" INTEGER NOT NULL CONSTRAINT "PK_asnmaster" PRIMARY KEY AUTOINCREMENT
    -- ASN 单号。
    ,"asn_no" TEXT NOT NULL
    -- ASN 批次号。
    ,"asn_batch" TEXT NOT NULL
    -- 预计到货时间。
    ,"estimated_arrival_time" TEXT NOT NULL
    -- ASN 状态。
    ,"asn_status" INTEGER NOT NULL
    -- 重量。
    ,"weight" TEXT NOT NULL
    -- 体积。
    ,"volume" TEXT NOT NULL
    -- 货主 ID。
    ,"goods_owner_id" INTEGER NOT NULL
    -- 货主名称。
    ,"goods_owner_name" TEXT NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: asnmaster
-- 业务单号/编码查询索引。
CREATE INDEX IF NOT EXISTS "IX_asnmaster_asn_no" ON "asnmaster" ("asn_no");
-- 状态过滤索引。
CREATE INDEX IF NOT EXISTS "IX_asnmaster_asn_status" ON "asnmaster" ("asn_status");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_asnmaster_goods_owner_id" ON "asnmaster" ("goods_owner_id");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_asnmaster_tenant_id" ON "asnmaster" ("tenant_id");

-- ============================================================
-- 表名: asnsort
-- 用途: 入库分拣记录表，记录 ASN 明细分拣和上架过程。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "asnsort" (
    -- 主键 ID，自增。
                           "id" INTEGER NOT NULL CONSTRAINT "PK_asnsort" PRIMARY KEY AUTOINCREMENT
    -- ASN 明细 ID。
    ,"asn_id" INTEGER NOT NULL
    -- 已分拣数量。
    ,"sorted_qty" INTEGER NOT NULL
    -- 序列号。
    ,"series_number" TEXT NOT NULL
    -- 已上架数量。
    ,"putaway_qty" INTEGER NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 是否有效，0/1。
    ,"is_valid" INTEGER NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: asnsort
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_asnsort_asn_id" ON "asnsort" ("asn_id");
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_asnsort_series_number" ON "asnsort" ("series_number");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_asnsort_tenant_id" ON "asnsort" ("tenant_id");

-- ============================================================
-- 表名: category
-- 用途: 商品分类表，用于维护 SPU 分类树。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "category" (
    -- 主键 ID，自增。
                            "id" INTEGER NOT NULL CONSTRAINT "PK_category" PRIMARY KEY AUTOINCREMENT
    -- 分类名称。
    ,"category_name" TEXT NOT NULL
    -- 父级 ID。
    ,"parent_id" INTEGER NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 是否有效，0/1。
    ,"is_valid" INTEGER NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: category
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_category_category_name" ON "category" ("category_name");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_category_parent_id" ON "category" ("parent_id");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_category_tenant_id" ON "category" ("tenant_id");

-- ============================================================
-- 表名: company
-- 用途: 公司信息表，保存当前租户或主体公司的基础资料。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "company" (
    -- 主键 ID，自增。
                           "id" INTEGER NOT NULL CONSTRAINT "PK_company" PRIMARY KEY AUTOINCREMENT
    -- 公司名称。
    ,"company_name" TEXT NOT NULL
    -- 城市。
    ,"city" TEXT NOT NULL
    -- 地址。
    ,"address" TEXT NOT NULL
    -- 负责人。
    ,"manager" TEXT NOT NULL
    -- 联系电话。
    ,"contact_tel" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: company
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_company_company_name" ON "company" ("company_name");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_company_tenant_id" ON "company" ("tenant_id");

-- ============================================================
-- 表名: customer
-- 用途: 客户表，保存发货客户资料。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "customer" (
    -- 主键 ID，自增。
                            "id" INTEGER NOT NULL CONSTRAINT "PK_customer" PRIMARY KEY AUTOINCREMENT
    -- 客户名称。
    ,"customer_name" TEXT NOT NULL
    -- 城市。
    ,"city" TEXT NOT NULL
    -- 地址。
    ,"address" TEXT NOT NULL
    -- 邮箱。
    ,"email" TEXT NOT NULL
    -- 负责人。
    ,"manager" TEXT NOT NULL
    -- 联系电话。
    ,"contact_tel" TEXT NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 是否有效，0/1。
    ,"is_valid" INTEGER NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: customer
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_customer_customer_name" ON "customer" ("customer_name");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_customer_tenant_id" ON "customer" ("tenant_id");

-- ============================================================
-- 表名: dispatchlist
-- 用途: 出库单主表，保存发货单头信息和执行状态。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "dispatchlist" (
    -- 主键 ID，自增。
                                "id" INTEGER NOT NULL CONSTRAINT "PK_dispatchlist" PRIMARY KEY AUTOINCREMENT
    -- 发货单号。
    ,"dispatch_no" TEXT NOT NULL
    -- 发货单状态。
    ,"dispatch_status" INTEGER NOT NULL
    -- 客户 ID。
    ,"customer_id" INTEGER NOT NULL
    -- 客户名称。
    ,"customer_name" TEXT NOT NULL
    -- SKU ID。
    ,"sku_id" INTEGER NOT NULL
    -- 数量。
    ,"qty" INTEGER NOT NULL
    -- 重量。
    ,"weight" TEXT NOT NULL
    -- 体积。
    ,"volume" TEXT NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 破损数量。
    ,"damage_qty" INTEGER NOT NULL
    -- 锁定数量。
    ,"lock_qty" INTEGER NOT NULL
    -- 已拣数量。
    ,"picked_qty" INTEGER NOT NULL
    -- 在途数量。
    ,"intrasit_qty" INTEGER NOT NULL
    -- 已打包数量。
    ,"package_qty" INTEGER NOT NULL
    -- 已称重数量。
    ,"weighing_qty" INTEGER NOT NULL
    -- 实际数量。
    ,"actual_qty" INTEGER NOT NULL
    -- 签收数量。
    ,"sign_qty" INTEGER NOT NULL
    -- 包裹号。
    ,"package_no" TEXT NOT NULL
    -- 打包人。
    ,"package_person" TEXT NOT NULL
    -- 打包时间。
    ,"package_time" TEXT NOT NULL
    -- 称重单号。
    ,"weighing_no" TEXT NOT NULL
    -- 称重人。
    ,"weighing_person" TEXT NOT NULL
    -- 称重重量。
    ,"weighing_weight" TEXT NOT NULL
    -- 运单号。
    ,"waybill_no" TEXT NOT NULL
    -- 承运商。
    ,"carrier" TEXT NOT NULL
    -- 运费。
    ,"freightfee" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
    -- 复核人 ID。
    ,"pick_checker_id" INTEGER NOT NULL
    -- 复核人。
    ,"pick_checker" TEXT NOT NULL
);

-- 索引: dispatchlist
-- 业务单号/编码查询索引。
CREATE INDEX IF NOT EXISTS "IX_dispatchlist_dispatch_no" ON "dispatchlist" ("dispatch_no");
-- 状态过滤索引。
CREATE INDEX IF NOT EXISTS "IX_dispatchlist_dispatch_status" ON "dispatchlist" ("dispatch_status");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_dispatchlist_customer_id" ON "dispatchlist" ("customer_id");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_dispatchlist_tenant_id" ON "dispatchlist" ("tenant_id");

-- ============================================================
-- 表名: flowsetmain
-- 用途: 审批流主表，定义某个菜单或业务单据对应的审批流。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "flowsetmain" (
    -- 主键 ID，自增。
                               "id" INTEGER NOT NULL CONSTRAINT "PK_flowsetmain" PRIMARY KEY AUTOINCREMENT
    -- 菜单标识。
    ,"menu" TEXT NOT NULL
    -- 审批流名称。
    ,"flow_name" TEXT NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: flowsetmain
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_flowsetmain_menu" ON "flowsetmain" ("menu");
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_flowsetmain_flow_name" ON "flowsetmain" ("flow_name");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_flowsetmain_tenant_id" ON "flowsetmain" ("tenant_id");

-- ============================================================
-- 表名: freightfee
-- 用途: 运费模板表，维护承运商线路和计费规则。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "freightfee" (
    -- 主键 ID，自增。
                              "id" INTEGER NOT NULL CONSTRAINT "PK_freightfee" PRIMARY KEY AUTOINCREMENT
    -- 承运商。
    ,"carrier" TEXT NOT NULL
    -- 始发城市。
    ,"departure_city" TEXT NOT NULL
    -- 到达城市。
    ,"arrival_city" TEXT NOT NULL
    -- 单位重量运费。
    ,"price_per_weight" TEXT NOT NULL
    -- 单位体积运费。
    ,"price_per_volume" TEXT NOT NULL
    -- 最低收费。
    ,"min_payment" TEXT NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 是否有效，0/1。
    ,"is_valid" INTEGER NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: freightfee
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_freightfee_carrier" ON "freightfee" ("carrier");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_freightfee_departure_city" ON "freightfee" ("departure_city");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_freightfee_arrival_city" ON "freightfee" ("arrival_city");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_freightfee_tenant_id" ON "freightfee" ("tenant_id");

-- ============================================================
-- 表名: global_unique_serial
-- 用途: 全局流水号表，维护业务单号生成规则和当前序号。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "global_unique_serial" (
    -- 主键 ID，自增。
                                        "id" INTEGER NOT NULL CONSTRAINT "PK_global_unique_serial" PRIMARY KEY AUTOINCREMENT
    -- 表名。
    ,"table_name" TEXT NOT NULL
    -- 流水号前缀。
    ,"prefix_char" TEXT NOT NULL
    -- 流水号重置规则。
    ,"reset_rule" TEXT NOT NULL
    -- 当前流水号。
    ,"current_no" INTEGER NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: global_unique_serial
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_global_unique_serial_table_name" ON "global_unique_serial" ("table_name");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_global_unique_serial_tenant_id" ON "global_unique_serial" ("tenant_id");

-- ============================================================
-- 表名: goodslocation
-- 用途: 库位表，维护仓库内具体货位信息。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "goodslocation" (
    -- 主键 ID，自增。
                                 "id" INTEGER NOT NULL CONSTRAINT "PK_goodslocation" PRIMARY KEY AUTOINCREMENT
    -- 仓库 ID。
    ,"warehouse_id" INTEGER NOT NULL
    -- 仓库名称。
    ,"warehouse_name" TEXT NOT NULL
    -- 库区名称。
    ,"warehouse_area_name" TEXT NOT NULL
    -- 库区属性。
    ,"warehouse_area_property" INTEGER NOT NULL
    -- 库位名称。
    ,"location_name" TEXT NOT NULL
    -- 库位长度。
    ,"location_length" TEXT NOT NULL
    -- 库位宽度。
    ,"location_width" TEXT NOT NULL
    -- 库位高度。
    ,"location_heigth" TEXT NOT NULL
    -- 库位容积。
    ,"location_volume" TEXT NOT NULL
    -- 库位承重。
    ,"location_load" TEXT NOT NULL
    -- 巷道编号。
    ,"roadway_number" TEXT NOT NULL
    -- 货架编号。
    ,"shelf_number" TEXT NOT NULL
    -- 层号。
    ,"layer_number" TEXT NOT NULL
    -- 标签号。
    ,"tag_number" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 是否有效，0/1。
    ,"is_valid" INTEGER NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
    -- 库区 ID。
    ,"warehouse_area_id" INTEGER NOT NULL
);

-- 索引: goodslocation
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_goodslocation_warehouse_id" ON "goodslocation" ("warehouse_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_goodslocation_warehouse_area_id" ON "goodslocation" ("warehouse_area_id");
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_goodslocation_location_name" ON "goodslocation" ("location_name");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_goodslocation_tenant_id" ON "goodslocation" ("tenant_id");

-- ============================================================
-- 表名: goodsowner
-- 用途: 货主表，维护库存所属货主资料。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "goodsowner" (
    -- 主键 ID，自增。
                              "id" INTEGER NOT NULL CONSTRAINT "PK_goodsowner" PRIMARY KEY AUTOINCREMENT
    -- 货主名称。
    ,"goods_owner_name" TEXT NOT NULL
    -- 城市。
    ,"city" TEXT NOT NULL
    -- 地址。
    ,"address" TEXT NOT NULL
    -- 负责人。
    ,"manager" TEXT NOT NULL
    -- 联系电话。
    ,"contact_tel" TEXT NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 是否有效，0/1。
    ,"is_valid" INTEGER NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: goodsowner
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_goodsowner_goods_owner_name" ON "goodsowner" ("goods_owner_name");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_goodsowner_tenant_id" ON "goodsowner" ("tenant_id");

-- ============================================================
-- 表名: menu
-- 用途: 菜单表，定义前端功能菜单和路由信息。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "menu" (
    -- 主键 ID，自增。
                        "id" INTEGER NOT NULL CONSTRAINT "PK_menu" PRIMARY KEY AUTOINCREMENT
    -- 菜单名称。
    ,"menu_name" TEXT NOT NULL
    -- 所属模块。
    ,"module" TEXT NOT NULL
    -- 前端页面路径。
    ,"vue_path" TEXT NOT NULL
    -- 前端详情页路径。
    ,"vue_path_detail" TEXT NOT NULL
    -- 前端目录。
    ,"vue_directory" TEXT NOT NULL
    -- 排序值。
    ,"sort" INTEGER NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
    -- 菜单可用动作集合。
    ,"menu_actions" TEXT NOT NULL
);

-- 索引: menu
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_menu_module" ON "menu" ("module");
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_menu_vue_path" ON "menu" ("vue_path");
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_menu_menu_name" ON "menu" ("menu_name");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_menu_tenant_id" ON "menu" ("tenant_id");

-- ============================================================
-- 表名: user_defined_print_solution
-- 用途: 自定义打印方案表，保存用户打印模板和布局配置。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "user_defined_print_solution" (
    -- 主键 ID，自增。
                                               "id" INTEGER NOT NULL CONSTRAINT "PK_user_defined_print_solution" PRIMARY KEY AUTOINCREMENT
    -- 前端页面路径。
    ,"vue_path" TEXT NOT NULL
    -- 标签页标识。
    ,"tab_page" TEXT NOT NULL
    -- 打印方案名称。
    ,"solution_name" TEXT NOT NULL
    -- 打印模板配置 JSON。
    ,"config_json" TEXT NOT NULL
    -- 打印纸张长度。
    ,"report_length" TEXT NOT NULL
    -- 打印纸张宽度。
    ,"report_width" TEXT NOT NULL
    -- 打印方向。
    ,"report_direction" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: user_defined_print_solution
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_user_defined_print_solution_vue_path" ON "user_defined_print_solution" ("vue_path");
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_user_defined_print_solution_tab_page" ON "user_defined_print_solution" ("tab_page");
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_user_defined_print_solution_solution_name" ON "user_defined_print_solution" ("solution_name");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_user_defined_print_solution_tenant_id" ON "user_defined_print_solution" ("tenant_id");

-- ============================================================
-- 表名: rolemenu
-- 用途: 角色菜单权限表，维护角色与菜单的授权关系。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "rolemenu" (
    -- 主键 ID，自增。
                            "id" INTEGER NOT NULL CONSTRAINT "PK_rolemenu" PRIMARY KEY AUTOINCREMENT
    -- 角色 ID。
    ,"userrole_id" INTEGER NOT NULL
    -- 菜单 ID。
    ,"menu_id" INTEGER NOT NULL
    -- 权限值。
    ,"authority" INTEGER NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
    -- 菜单动作权限集合。
    ,"menu_actions_authority" TEXT NOT NULL
);

-- 索引: rolemenu
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_rolemenu_userrole_id" ON "rolemenu" ("userrole_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_rolemenu_menu_id" ON "rolemenu" ("menu_id");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_rolemenu_tenant_id" ON "rolemenu" ("tenant_id");

-- ============================================================
-- 表名: spu
-- 用途: SPU 表，维护商品主数据。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "spu" (
    -- 主键 ID，自增。
                       "id" INTEGER NOT NULL CONSTRAINT "PK_spu" PRIMARY KEY AUTOINCREMENT
    -- SPU 编码。
    ,"spu_code" TEXT NOT NULL
    -- SPU 名称。
    ,"spu_name" TEXT NOT NULL
    -- 分类 ID。
    ,"category_id" INTEGER NOT NULL
    -- SPU 描述。
    ,"spu_description" TEXT NOT NULL
    -- 供应商 ID。
    ,"supplier_id" INTEGER NOT NULL
    -- 供应商名称。
    ,"supplier_name" TEXT NOT NULL
    -- 品牌。
    ,"brand" TEXT NOT NULL
    -- 产地。
    ,"origin" TEXT NOT NULL
    -- 长度单位。
    ,"length_unit" INTEGER NOT NULL
    -- 体积单位。
    ,"volume_unit" INTEGER NOT NULL
    -- 重量单位。
    ,"weight_unit" INTEGER NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 是否有效，0/1。
    ,"is_valid" INTEGER NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: spu
-- 业务单号/编码查询索引。
CREATE INDEX IF NOT EXISTS "IX_spu_spu_code" ON "spu" ("spu_code");
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_spu_spu_name" ON "spu" ("spu_name");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_spu_category_id" ON "spu" ("category_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_spu_supplier_id" ON "spu" ("supplier_id");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_spu_tenant_id" ON "spu" ("tenant_id");

-- ============================================================
-- 表名: stock
-- 用途: 即时库存表，记录货主、SKU、库位维度的当前库存。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "stock" (
    -- 主键 ID，自增。
                         "id" INTEGER NOT NULL CONSTRAINT "PK_stock" PRIMARY KEY AUTOINCREMENT
    -- SKU ID。
    ,"sku_id" INTEGER NOT NULL
    -- 库位 ID。
    ,"goods_location_id" INTEGER NOT NULL
    -- 数量。
    ,"qty" INTEGER NOT NULL
    -- 货主 ID。
    ,"goods_owner_id" INTEGER NOT NULL
    -- 是否冻结，0/1。
    ,"is_freeze" INTEGER NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
    -- 序列号。
    ,"series_number" TEXT NOT NULL
    -- 有效期或失效日期。
    ,"expiry_date" TEXT NOT NULL
    -- 单价。
    ,"price" TEXT NOT NULL
    -- 上架日期。
    ,"putaway_date" TEXT NOT NULL
);

-- 索引: stock
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stock_sku_id" ON "stock" ("sku_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stock_goods_location_id" ON "stock" ("goods_location_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stock_goods_owner_id" ON "stock" ("goods_owner_id");
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_stock_series_number" ON "stock" ("series_number");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_stock_tenant_id" ON "stock" ("tenant_id");

-- ============================================================
-- 表名: stockadjust
-- 用途: 库存调整表，记录盘盈盘亏等库存调整明细。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "stockadjust" (
    -- 主键 ID，自增。
                               "id" INTEGER NOT NULL CONSTRAINT "PK_stockadjust" PRIMARY KEY AUTOINCREMENT
    -- 业务作业单号。
    ,"job_code" TEXT NOT NULL
    -- SKU ID。
    ,"sku_id" INTEGER NOT NULL
    -- 货主 ID。
    ,"goods_owner_id" INTEGER NOT NULL
    -- 库位 ID。
    ,"goods_location_id" INTEGER NOT NULL
    -- 数量。
    ,"qty" INTEGER NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
    -- 是否已回写库存，0/1。
    ,"is_update_stock" INTEGER NOT NULL
    -- 作业类型。
    ,"job_type" INTEGER NOT NULL
    -- 来源业务表主键。
    ,"source_table_id" INTEGER NOT NULL
    -- 序列号。
    ,"series_number" TEXT NOT NULL
    -- 有效期或失效日期。
    ,"expiry_date" TEXT NOT NULL
    -- 单价。
    ,"price" TEXT NOT NULL
    -- 上架日期。
    ,"putaway_date" TEXT NOT NULL
);

-- 索引: stockadjust
-- 业务单号/编码查询索引。
CREATE INDEX IF NOT EXISTS "IX_stockadjust_job_code" ON "stockadjust" ("job_code");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stockadjust_sku_id" ON "stockadjust" ("sku_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stockadjust_goods_owner_id" ON "stockadjust" ("goods_owner_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stockadjust_goods_location_id" ON "stockadjust" ("goods_location_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stockadjust_source_table_id" ON "stockadjust" ("source_table_id");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_stockadjust_tenant_id" ON "stockadjust" ("tenant_id");

-- ============================================================
-- 表名: stockfreeze
-- 用途: 库存冻结解冻表，记录库存冻结业务。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "stockfreeze" (
    -- 主键 ID，自增。
                               "id" INTEGER NOT NULL CONSTRAINT "PK_stockfreeze" PRIMARY KEY AUTOINCREMENT
    -- 业务作业单号。
    ,"job_code" TEXT NOT NULL
    -- 作业类型。
    ,"job_type" INTEGER NOT NULL
    -- SKU ID。
    ,"sku_id" INTEGER NOT NULL
    -- 货主 ID。
    ,"goods_owner_id" INTEGER NOT NULL
    -- 库位 ID。
    ,"goods_location_id" INTEGER NOT NULL
    -- 处理人。
    ,"handler" TEXT NOT NULL
    -- 处理时间。
    ,"handle_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
    -- 序列号。
    ,"series_number" TEXT NOT NULL
);

-- 索引: stockfreeze
-- 业务单号/编码查询索引。
CREATE INDEX IF NOT EXISTS "IX_stockfreeze_job_code" ON "stockfreeze" ("job_code");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stockfreeze_sku_id" ON "stockfreeze" ("sku_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stockfreeze_goods_owner_id" ON "stockfreeze" ("goods_owner_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stockfreeze_goods_location_id" ON "stockfreeze" ("goods_location_id");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_stockfreeze_tenant_id" ON "stockfreeze" ("tenant_id");

-- ============================================================
-- 表名: stockmove
-- 用途: 库内移位表，记录库存从原库位移动到目标库位的过程。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "stockmove" (
    -- 主键 ID，自增。
                             "id" INTEGER NOT NULL CONSTRAINT "PK_stockmove" PRIMARY KEY AUTOINCREMENT
    -- 业务作业单号。
    ,"job_code" TEXT NOT NULL
    -- 移位状态。
    ,"move_status" INTEGER NOT NULL
    -- SKU ID。
    ,"sku_id" INTEGER NOT NULL
    -- 原库位 ID。
    ,"orig_goods_location_id" INTEGER NOT NULL
    -- 目标库位 ID。
    ,"dest_googs_location_id" INTEGER NOT NULL
    -- 数量。
    ,"qty" INTEGER NOT NULL
    -- 货主 ID。
    ,"goods_owner_id" INTEGER NOT NULL
    -- 处理人。
    ,"handler" TEXT NOT NULL
    -- 处理时间。
    ,"handle_time" TEXT NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
    -- 序列号。
    ,"series_number" TEXT NOT NULL
    -- 有效期或失效日期。
    ,"expiry_date" TEXT NOT NULL
    -- 单价。
    ,"price" TEXT NOT NULL
    -- 上架日期。
    ,"putaway_date" TEXT NOT NULL
);

-- 索引: stockmove
-- 业务单号/编码查询索引。
CREATE INDEX IF NOT EXISTS "IX_stockmove_job_code" ON "stockmove" ("job_code");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stockmove_sku_id" ON "stockmove" ("sku_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stockmove_orig_goods_location_id" ON "stockmove" ("orig_goods_location_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stockmove_dest_googs_location_id" ON "stockmove" ("dest_googs_location_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stockmove_goods_owner_id" ON "stockmove" ("goods_owner_id");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_stockmove_tenant_id" ON "stockmove" ("tenant_id");

-- ============================================================
-- 表名: stockprocess
-- 用途: 库存加工主表，记录加工、拆装等库存处理任务。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "stockprocess" (
    -- 主键 ID，自增。
                                "id" INTEGER NOT NULL CONSTRAINT "PK_stockprocess" PRIMARY KEY AUTOINCREMENT
    -- 业务作业单号。
    ,"job_code" TEXT NOT NULL
    -- 作业类型。
    ,"job_type" INTEGER NOT NULL
    -- 处理状态。
    ,"process_status" INTEGER NOT NULL
    -- 处理人。
    ,"processor" TEXT NOT NULL
    -- 处理时间。
    ,"process_time" TEXT NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: stockprocess
-- 业务单号/编码查询索引。
CREATE INDEX IF NOT EXISTS "IX_stockprocess_job_code" ON "stockprocess" ("job_code");
-- 状态过滤索引。
CREATE INDEX IF NOT EXISTS "IX_stockprocess_process_status" ON "stockprocess" ("process_status");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_stockprocess_tenant_id" ON "stockprocess" ("tenant_id");

-- ============================================================
-- 表名: stocktaking
-- 用途: 库存盘点表，记录盘点任务及账实差异。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "stocktaking" (
    -- 主键 ID，自增。
                               "id" INTEGER NOT NULL CONSTRAINT "PK_stocktaking" PRIMARY KEY AUTOINCREMENT
    -- 业务作业单号。
    ,"job_code" TEXT NOT NULL
    -- 作业状态。
    ,"job_status" INTEGER NOT NULL
    -- SKU ID。
    ,"sku_id" INTEGER NOT NULL
    -- 货主 ID。
    ,"goods_owner_id" INTEGER NOT NULL
    -- 库位 ID。
    ,"goods_location_id" INTEGER NOT NULL
    -- 序列号。
    ,"series_number" TEXT NOT NULL
    -- 有效期或失效日期。
    ,"expiry_date" TEXT NOT NULL
    -- 单价。
    ,"price" TEXT NOT NULL
    -- 上架日期。
    ,"putaway_date" TEXT NOT NULL
    -- 账面数量。
    ,"book_qty" INTEGER NOT NULL
    -- 实盘数量。
    ,"counted_qty" INTEGER NOT NULL
    -- 差异数量。
    ,"difference_qty" INTEGER NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
    -- 处理人。
    ,"handler" TEXT NOT NULL
    -- 处理时间。
    ,"handle_time" TEXT NOT NULL
);

-- 索引: stocktaking
-- 业务单号/编码查询索引。
CREATE INDEX IF NOT EXISTS "IX_stocktaking_job_code" ON "stocktaking" ("job_code");
-- 状态过滤索引。
CREATE INDEX IF NOT EXISTS "IX_stocktaking_job_status" ON "stocktaking" ("job_status");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stocktaking_sku_id" ON "stocktaking" ("sku_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stocktaking_goods_owner_id" ON "stocktaking" ("goods_owner_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stocktaking_goods_location_id" ON "stocktaking" ("goods_location_id");
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_stocktaking_series_number" ON "stocktaking" ("series_number");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_stocktaking_tenant_id" ON "stocktaking" ("tenant_id");

-- ============================================================
-- 表名: supplier
-- 用途: 供应商表，维护采购或入库供应商资料。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "supplier" (
    -- 主键 ID，自增。
                            "id" INTEGER NOT NULL CONSTRAINT "PK_supplier" PRIMARY KEY AUTOINCREMENT
    -- 供应商名称。
    ,"supplier_name" TEXT NOT NULL
    -- 城市。
    ,"city" TEXT NOT NULL
    -- 地址。
    ,"address" TEXT NOT NULL
    -- 邮箱。
    ,"email" TEXT NOT NULL
    -- 负责人。
    ,"manager" TEXT NOT NULL
    -- 联系电话。
    ,"contact_tel" TEXT NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 是否有效，0/1。
    ,"is_valid" INTEGER NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: supplier
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_supplier_supplier_name" ON "supplier" ("supplier_name");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_supplier_tenant_id" ON "supplier" ("tenant_id");

-- ============================================================
-- 表名: userrole
-- 用途: 角色表，定义系统角色。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "userrole" (
    -- 主键 ID，自增。
                            "id" INTEGER NOT NULL CONSTRAINT "PK_userrole" PRIMARY KEY AUTOINCREMENT
    -- 角色名称。
    ,"role_name" TEXT NOT NULL
    -- 是否有效，0/1。
    ,"is_valid" INTEGER NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: userrole
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_userrole_role_name" ON "userrole" ("role_name");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_userrole_tenant_id" ON "userrole" ("tenant_id");

-- ============================================================
-- 表名: warehouse
-- 用途: 仓库表，维护仓库基础资料。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "warehouse" (
    -- 主键 ID，自增。
                             "id" INTEGER NOT NULL CONSTRAINT "PK_warehouse" PRIMARY KEY AUTOINCREMENT
    -- 仓库名称。
    ,"warehouse_name" TEXT NOT NULL
    -- 城市。
    ,"city" TEXT NOT NULL
    -- 地址。
    ,"address" TEXT NOT NULL
    -- 邮箱。
    ,"email" TEXT NOT NULL
    -- 负责人。
    ,"manager" TEXT NOT NULL
    -- 联系电话。
    ,"contact_tel" TEXT NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 是否有效，0/1。
    ,"is_valid" INTEGER NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: warehouse
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_warehouse_warehouse_name" ON "warehouse" ("warehouse_name");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_warehouse_tenant_id" ON "warehouse" ("tenant_id");

-- ============================================================
-- 表名: warehousearea
-- 用途: 库区表，维护仓库下的库区结构。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "warehousearea" (
    -- 主键 ID，自增。
     "id" INTEGER NOT NULL CONSTRAINT "PK_warehousearea" PRIMARY KEY AUTOINCREMENT
    -- 仓库 ID。
    ,"warehouse_id" INTEGER NOT NULL
    -- 库区名称。
    ,"area_name" TEXT NOT NULL
    -- 父级 ID。
    ,"parent_id" INTEGER NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 是否有效，0/1。
    ,"is_valid" INTEGER NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
    -- 库区属性。
    ,"area_property" INTEGER NOT NULL
);

-- 索引: warehousearea
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_warehousearea_warehouse_id" ON "warehousearea" ("warehouse_id");
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_warehousearea_area_name" ON "warehousearea" ("area_name");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_warehousearea_parent_id" ON "warehousearea" ("parent_id");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_warehousearea_tenant_id" ON "warehousearea" ("tenant_id");

-- ============================================================
-- 表名: user
-- 用途: 用户表，保存系统账号信息。
-- 关联: 无显式外键。
-- ============================================================
CREATE TABLE "user" (
    -- 主键 ID，自增。
                        "id" INTEGER NOT NULL CONSTRAINT "PK_user" PRIMARY KEY AUTOINCREMENT
    -- 用户编号。
    ,"user_num" TEXT NOT NULL
    -- 用户名。
    ,"user_name" TEXT NOT NULL
    -- 联系电话。
    ,"contact_tel" TEXT NOT NULL
    -- 角色名称或角色标识。
    ,"user_role" TEXT NOT NULL
    -- 性别。
    ,"sex" TEXT NOT NULL
    -- 是否有效，0/1。
    ,"is_valid" INTEGER NOT NULL
    -- 登录密码或认证串。
    ,"auth_string" TEXT NOT NULL
    -- 邮箱。
    ,"email" TEXT NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
);

-- 索引: user
-- 业务单号/编码查询索引。
CREATE INDEX IF NOT EXISTS "IX_user_user_num" ON "user" ("user_num");
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_user_user_name" ON "user" ("user_name");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_user_tenant_id" ON "user" ("tenant_id");

-- ============================================================
-- 表名: asn
-- 用途: asn entity
-- 关联: asnmaster
-- ============================================================
CREATE TABLE "asn" (
    -- 主键 ID，自增。
                       "id" INTEGER NOT NULL CONSTRAINT "PK_asn" PRIMARY KEY AUTOINCREMENT
    -- ASN 主表 ID。
    ,"asnmaster_id" INTEGER NOT NULL
    -- ASN 单号。
    ,"asn_no" TEXT NOT NULL
    -- ASN 状态。
    ,"asn_status" INTEGER NOT NULL
    -- SPU ID。
    ,"spu_id" INTEGER NOT NULL
    -- SKU ID。
    ,"sku_id" INTEGER NOT NULL
    -- ASN 数量。
    ,"asn_qty" INTEGER NOT NULL
    -- 实际数量。
    ,"actual_qty" INTEGER NOT NULL
    -- 实际到货时间。
    ,"arrival_time" TEXT NOT NULL
    -- 卸货时间。
    ,"unload_time" TEXT NOT NULL
    -- 卸货人 ID。
    ,"unload_person_id" INTEGER NOT NULL
    -- 卸货人。
    ,"unload_person" TEXT NOT NULL
    -- 已分拣数量。
    ,"sorted_qty" INTEGER NOT NULL
    -- 短少数量。
    ,"shortage_qty" INTEGER NOT NULL
    -- 超收数量。
    ,"more_qty" INTEGER NOT NULL
    -- 破损数量。
    ,"damage_qty" INTEGER NOT NULL
    -- 重量。
    ,"weight" TEXT NOT NULL
    -- 体积。
    ,"volume" TEXT NOT NULL
    -- 供应商 ID。
    ,"supplier_id" INTEGER NOT NULL
    -- 供应商名称。
    ,"supplier_name" TEXT NOT NULL
    -- 货主 ID。
    ,"goods_owner_id" INTEGER NOT NULL
    -- 货主名称。
    ,"goods_owner_name" TEXT NOT NULL
    -- 创建人。
    ,"creator" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 是否有效，0/1。
    ,"is_valid" INTEGER NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
    -- 有效期或失效日期。
    ,"expiry_date" TEXT NOT NULL
    -- 单价。
    ,"price" TEXT NOT NULL
    -- 外键: asnmaster_id -> asnmaster.id
    ,CONSTRAINT "FK_asn_asnmaster_asnmaster_id" FOREIGN KEY ("asnmaster_id") REFERENCES "asnmaster" ("id") ON DELETE CASCADE
);

-- 索引: asn
-- 外键关联查询索引。
CREATE INDEX IF NOT EXISTS "IX_asn_asnmaster_id" ON "asn" ("asnmaster_id");

-- ============================================================
-- 表名: dispatchpicklist
-- 用途: 拣货明细表，记录发货单对应的拣货明细。
-- 关联: dispatchlist
-- ============================================================
CREATE TABLE "dispatchpicklist" (
    -- 主键 ID，自增。
                                    "id" INTEGER NOT NULL CONSTRAINT "PK_dispatchpicklist" PRIMARY KEY AUTOINCREMENT
    -- 发货单 ID。
    ,"dispatchlist_id" INTEGER NOT NULL
    -- 货主 ID。
    ,"goods_owner_id" INTEGER NOT NULL
    -- 库位 ID。
    ,"goods_location_id" INTEGER NOT NULL
    -- SKU ID。
    ,"sku_id" INTEGER NOT NULL
    -- 待拣数量。
    ,"pick_qty" INTEGER NOT NULL
    -- 已拣数量。
    ,"picked_qty" INTEGER NOT NULL
    -- 是否已回写库存，0/1。
    ,"is_update_stock" INTEGER NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 序列号。
    ,"series_number" TEXT NOT NULL
    -- 拣货人 ID。
    ,"picker_id" INTEGER NOT NULL
    -- 拣货人。
    ,"picker" TEXT NOT NULL
    -- 有效期或失效日期。
    ,"expiry_date" TEXT NOT NULL
    -- 单价。
    ,"price" TEXT NOT NULL
    -- 上架日期。
    ,"putaway_date" TEXT NOT NULL
    -- 外键: dispatchlist_id -> dispatchlist.id
    ,CONSTRAINT "FK_dispatchpicklist_dispatchlist_dispatchlist_id" FOREIGN KEY ("dispatchlist_id") REFERENCES "dispatchlist" ("id") ON DELETE CASCADE
);

-- 索引: dispatchpicklist
-- 外键关联查询索引。
CREATE INDEX IF NOT EXISTS "IX_dispatchpicklist_dispatchlist_id" ON "dispatchpicklist" ("dispatchlist_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_dispatchpicklist_goods_owner_id" ON "dispatchpicklist" ("goods_owner_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_dispatchpicklist_goods_location_id" ON "dispatchpicklist" ("goods_location_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_dispatchpicklist_sku_id" ON "dispatchpicklist" ("sku_id");
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_dispatchpicklist_series_number" ON "dispatchpicklist" ("series_number");

-- ============================================================
-- 表名: flowset
-- 用途: 审批流节点表，定义审批流中的节点。
-- 关联: flowsetmain
-- ============================================================
CREATE TABLE "flowset" (
    -- 主键 ID，自增。
                           "id" INTEGER NOT NULL CONSTRAINT "PK_flowset" PRIMARY KEY AUTOINCREMENT
    -- 审批流主表 ID。
    ,"flowsetmain_id" INTEGER NOT NULL
    -- 是否起始节点，0/1。
    ,"is_origin" INTEGER NOT NULL
    -- 是否结束节点，0/1。
    ,"is_end" INTEGER NOT NULL
    -- 节点 GUID。
    ,"node_guid" TEXT NOT NULL
    -- 节点名称。
    ,"node_name" TEXT NOT NULL
    -- 前一节点 GUID。
    ,"prev_node_guid" TEXT NOT NULL
    -- 外键: flowsetmain_id -> flowsetmain.id
    ,CONSTRAINT "FK_flowset_flowsetmain_flowsetmain_id" FOREIGN KEY ("flowsetmain_id") REFERENCES "flowsetmain" ("id") ON DELETE CASCADE
);

-- 索引: flowset
-- 外键关联查询索引。
CREATE INDEX IF NOT EXISTS "IX_flowset_flowsetmain_id" ON "flowset" ("flowsetmain_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_flowset_node_guid" ON "flowset" ("node_guid");

-- ============================================================
-- 表名: sku
-- 用途: SKU 表，维护具体库存单位信息。
-- 关联: spu
-- ============================================================
CREATE TABLE "sku" (
    -- 主键 ID，自增。
                       "id" INTEGER NOT NULL CONSTRAINT "PK_sku" PRIMARY KEY AUTOINCREMENT
    -- SPU ID。
    ,"spu_id" INTEGER NOT NULL
    -- SKU 编码。
    ,"sku_code" TEXT NOT NULL
    -- SKU 名称。
    ,"sku_name" TEXT NOT NULL
    -- 条码。
    ,"bar_code" TEXT NOT NULL
    -- 重量。
    ,"weight" TEXT NOT NULL
    -- 长度。
    ,"lenght" TEXT NOT NULL
    -- 宽度。
    ,"width" TEXT NOT NULL
    -- 高度。
    ,"height" TEXT NOT NULL
    -- 体积。
    ,"volume" TEXT NOT NULL
    -- 单位。
    ,"unit" TEXT NOT NULL
    -- 成本价。
    ,"cost" TEXT NOT NULL
    -- 单价。
    ,"price" TEXT NOT NULL
    -- 创建时间。
    ,"create_time" TEXT NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 图片地址。
    ,"image_url" TEXT NOT NULL
    -- 外键: spu_id -> spu.id
    ,CONSTRAINT "FK_sku_spu_spu_id" FOREIGN KEY ("spu_id") REFERENCES "spu" ("id") ON DELETE CASCADE
);

-- 索引: sku
-- 外键关联查询索引。
CREATE INDEX IF NOT EXISTS "IX_sku_spu_id" ON "sku" ("spu_id");
-- 业务单号/编码查询索引。
CREATE INDEX IF NOT EXISTS "IX_sku_sku_code" ON "sku" ("sku_code");
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_sku_sku_name" ON "sku" ("sku_name");
-- 业务单号/编码查询索引。
CREATE INDEX IF NOT EXISTS "IX_sku_bar_code" ON "sku" ("bar_code");

-- ============================================================
-- 表名: stockprocessdetail
-- 用途: 库存加工明细表，记录加工任务涉及的明细库存。
-- 关联: stockprocess
-- ============================================================
CREATE TABLE "stockprocessdetail" (
    -- 主键 ID，自增。
                                      "id" INTEGER NOT NULL CONSTRAINT "PK_stockprocessdetail" PRIMARY KEY AUTOINCREMENT
    -- 库存加工主表 ID。
    ,"stock_process_id" INTEGER NOT NULL
    -- SKU ID。
    ,"sku_id" INTEGER NOT NULL
    -- 货主 ID。
    ,"goods_owner_id" INTEGER NOT NULL
    -- 库位 ID。
    ,"goods_location_id" INTEGER NOT NULL
    -- 数量。
    ,"qty" INTEGER NOT NULL
    -- 最后更新时间。
    ,"last_update_time" TEXT NOT NULL
    -- 租户 ID。
    ,"tenant_id" INTEGER NOT NULL
    -- 是否源数据，0/1。
    ,"is_source" INTEGER NOT NULL
    -- 是否已回写库存，0/1。
    ,"is_update_stock" INTEGER NOT NULL
    -- 序列号。
    ,"series_number" TEXT NOT NULL
    -- 有效期或失效日期。
    ,"expiry_date" TEXT NOT NULL
    -- 单价。
    ,"price" TEXT NOT NULL
    -- 上架日期。
    ,"putaway_date" TEXT NOT NULL
    -- 外键: stock_process_id -> stockprocess.id
    ,CONSTRAINT "FK_stockprocessdetail_stockprocess_stock_process_id" FOREIGN KEY ("stock_process_id") REFERENCES "stockprocess" ("id") ON DELETE CASCADE
);

-- 索引: stockprocessdetail
-- 外键关联查询索引。
CREATE INDEX IF NOT EXISTS "IX_stockprocessdetail_stock_process_id" ON "stockprocessdetail" ("stock_process_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stockprocessdetail_sku_id" ON "stockprocessdetail" ("sku_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stockprocessdetail_goods_owner_id" ON "stockprocessdetail" ("goods_owner_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_stockprocessdetail_goods_location_id" ON "stockprocessdetail" ("goods_location_id");
-- 名称或关键检索字段索引。
CREATE INDEX IF NOT EXISTS "IX_stockprocessdetail_series_number" ON "stockprocessdetail" ("series_number");
-- 租户隔离查询索引。
CREATE INDEX IF NOT EXISTS "IX_stockprocessdetail_tenant_id" ON "stockprocessdetail" ("tenant_id");

-- ============================================================
-- 表名: flowsetfilter
-- 用途: 审批流条件表，定义节点流转条件。
-- 关联: flowset
-- ============================================================
CREATE TABLE "flowsetfilter" (
    -- 主键 ID，自增。
                                 "id" INTEGER NOT NULL CONSTRAINT "PK_flowsetfilter" PRIMARY KEY AUTOINCREMENT
    -- 审批流节点 ID。
    ,"flowset_id" INTEGER NOT NULL
    -- 审批流主表 ID。
    ,"flowsetmain_id" INTEGER NOT NULL
    -- 节点 GUID。
    ,"node_guid" TEXT NOT NULL
    -- 逻辑运算符。
    ,"logic" TEXT NOT NULL
    -- 左括号。
    ,"c1" TEXT NOT NULL
    -- 字段显示名。
    ,"col_label" TEXT NOT NULL
    -- 字段名。
    ,"col_name" TEXT NOT NULL
    -- 比较符。
    ,"compare" TEXT NOT NULL
    -- 条件内容。
    ,"content" TEXT NOT NULL
    -- 右括号。
    ,"c2" TEXT NOT NULL
    -- 排序值。
    ,"sort" INTEGER NOT NULL
    -- 条件分组。
    ,"condition_group" TEXT NOT NULL
    -- 公式表达式。
    ,"formulas" TEXT NOT NULL
    -- 断言模式。
    ,"assert_mode" TEXT NOT NULL
    -- 表名。
    ,"table_name" TEXT NOT NULL
    -- 方案名称。
    ,"scheme_name" TEXT NOT NULL
    -- 外键: flowset_id -> flowset.id
    ,CONSTRAINT "FK_flowsetfilter_flowset_flowset_id" FOREIGN KEY ("flowset_id") REFERENCES "flowset" ("id") ON DELETE CASCADE
);

-- 索引: flowsetfilter
-- 外键关联查询索引。
CREATE INDEX IF NOT EXISTS "IX_flowsetfilter_flowset_id" ON "flowsetfilter" ("flowset_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_flowsetfilter_flowsetmain_id" ON "flowsetfilter" ("flowsetmain_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_flowsetfilter_node_guid" ON "flowsetfilter" ("node_guid");

-- ============================================================
-- 表名: flowsetusers
-- 用途: 审批流节点人员表，维护节点审批人。
-- 关联: flowset
-- ============================================================
CREATE TABLE "flowsetusers" (
    -- 主键 ID，自增。
                                "id" INTEGER NOT NULL CONSTRAINT "PK_flowsetusers" PRIMARY KEY AUTOINCREMENT
    -- 审批流节点 ID。
    ,"flowset_id" INTEGER NOT NULL
    -- 审批流主表 ID。
    ,"flowsetmain_id" INTEGER NOT NULL
    -- 节点 GUID。
    ,"node_guid" TEXT NOT NULL
    -- 用户 ID。
    ,"user_id" INTEGER NOT NULL
    -- 外键: flowset_id -> flowset.id
    ,CONSTRAINT "FK_flowsetusers_flowset_flowset_id" FOREIGN KEY ("flowset_id") REFERENCES "flowset" ("id") ON DELETE CASCADE
);

-- 索引: flowsetusers
-- 外键关联查询索引。
CREATE INDEX IF NOT EXISTS "IX_flowsetusers_flowset_id" ON "flowsetusers" ("flowset_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_flowsetusers_flowsetmain_id" ON "flowsetusers" ("flowsetmain_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_flowsetusers_user_id" ON "flowsetusers" ("user_id");

-- ============================================================
-- 表名: sku_safety_stock
-- 用途: SKU 安全库存表，按仓库维护安全库存阈值。
-- 关联: sku
-- ============================================================
CREATE TABLE "sku_safety_stock" (
    -- 主键 ID，自增。
                                    "id" INTEGER NOT NULL CONSTRAINT "PK_sku_safety_stock" PRIMARY KEY AUTOINCREMENT
    -- SKU ID。
    ,"sku_id" INTEGER NOT NULL
    -- 仓库 ID。
    ,"warehouse_id" INTEGER NOT NULL
    -- 安全库存数量。
    ,"safety_stock_qty" INTEGER NOT NULL
    -- 外键: sku_id -> sku.id
    ,CONSTRAINT "FK_sku_safety_stock_sku_sku_id" FOREIGN KEY ("sku_id") REFERENCES "sku" ("id") ON DELETE CASCADE
);

-- 索引: sku_safety_stock
-- 外键关联查询索引。
CREATE INDEX IF NOT EXISTS "IX_sku_safety_stock_sku_id" ON "sku_safety_stock" ("sku_id");
-- 常用筛选与检索索引。
CREATE INDEX IF NOT EXISTS "IX_sku_safety_stock_warehouse_id" ON "sku_safety_stock" ("warehouse_id");

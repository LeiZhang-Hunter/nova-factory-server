-- MySQL dump 10.13  Distrib 8.0.45, for Linux (x86_64)
--
-- Host: localhost    Database: nova
-- ------------------------------------------------------
-- Server version	8.0.45-0ubuntu0.24.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `gen_table`
--

DROP TABLE IF EXISTS `gen_table`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gen_table` (
  `table_id` bigint NOT NULL AUTO_INCREMENT COMMENT '编号',
  `table_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '表名称',
  `parent_menu_id` bigint DEFAULT NULL COMMENT '父菜单ID',
  `table_comment` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '表描述',
  `sub_table_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '关联子表的表名',
  `sub_table_fk_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '子表关联的外键名',
  `struct_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '结构体名称',
  `tpl_category` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'crud' COMMENT '使用的模板（crud单表操作 tree树表操作 sub主子表操作）',
  `package_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '生成包路径',
  `module_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '生成模块名',
  `business_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '生成业务名',
  `function_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '生成功能名',
  `function_author` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '生成功能作者',
  `options` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '其它生成选项',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`table_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=221635920305065985 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='代码生成业务表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gen_table`
--

LOCK TABLES `gen_table` WRITE;
/*!40000 ALTER TABLE `gen_table` DISABLE KEYS */;
/*!40000 ALTER TABLE `gen_table` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gen_table_column`
--

DROP TABLE IF EXISTS `gen_table_column`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gen_table_column` (
  `column_id` bigint NOT NULL AUTO_INCREMENT COMMENT '编号',
  `table_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '归属表编号',
  `column_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '列名称',
  `column_comment` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '列描述',
  `column_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '列类型',
  `go_type` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'go类型',
  `go_field` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'go字段名',
  `is_pk` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '是否主键（1是）',
  `is_required` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '是否必填（1是）',
  `is_insert` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '是否为插入字段（1是）',
  `is_edit` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '是否编辑字段（1是）',
  `is_list` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '是否列表字段（1是）',
  `is_query` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '是否查询字段（1是）',
  `query_type` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'EQ' COMMENT '查询方式（等于、不等于、大于、小于、范围）',
  `html_type` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '显示类型（文本框、文本域、下拉框、复选框、单选框、日期控件）',
  `dict_type` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '字典类型',
  `sort` int DEFAULT NULL COMMENT '排序',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `html_field` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  PRIMARY KEY (`column_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=221635920317648906 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='代码生成业务表字段';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gen_table_column`
--

LOCK TABLES `gen_table_column` WRITE;
/*!40000 ALTER TABLE `gen_table_column` DISABLE KEYS */;
/*!40000 ALTER TABLE `gen_table_column` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `pro_route_product_bom`
--

DROP TABLE IF EXISTS `pro_route_product_bom`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `pro_route_product_bom` (
  `record_id` bigint NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `route_id` bigint NOT NULL COMMENT '工艺路线ID',
  `process_id` bigint NOT NULL COMMENT '工序ID',
  `product_id` bigint NOT NULL COMMENT '产品BOM中的唯一ID',
  `item_id` bigint NOT NULL COMMENT '产品物料ID',
  `item_code` varchar(64) NOT NULL COMMENT '产品物料编码',
  `item_name` varchar(255) NOT NULL COMMENT '产品物料名称',
  `specification` varchar(500) DEFAULT NULL COMMENT '规格型号',
  `unit_of_measure` varchar(64) NOT NULL COMMENT '单位',
  `unit_name` varchar(64) DEFAULT NULL COMMENT '单位名称',
  `quantity` double(12,2) DEFAULT '1.00' COMMENT '用料比例',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `attr1` varchar(64) DEFAULT NULL COMMENT '预留字段1',
  `attr2` varchar(255) DEFAULT NULL COMMENT '预留字段2',
  `attr3` int DEFAULT '0' COMMENT '预留字段3',
  `attr4` int DEFAULT '0' COMMENT '预留字段4',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`record_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='产品制程物料BOM表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pro_route_product_bom`
--

LOCK TABLES `pro_route_product_bom` WRITE;
/*!40000 ALTER TABLE `pro_route_product_bom` DISABLE KEYS */;
/*!40000 ALTER TABLE `pro_route_product_bom` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_ai_prediction_control`
--

DROP TABLE IF EXISTS `sys_ai_prediction_control`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_ai_prediction_control` (
  `id` bigint NOT NULL COMMENT 'id',
  `device_gateway_id` bigint NOT NULL COMMENT '网关id',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '智能预警名称',
  `parallelism` bigint NOT NULL COMMENT '并发',
  `threshold` bigint NOT NULL COMMENT 'threshold',
  `model` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '预测模型',
  `interval` bigint NOT NULL COMMENT '预测时间段',
  `predict_length` bigint unsigned DEFAULT '0' COMMENT '预测长度',
  `agg_function` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '聚合函数，用来计算图表',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='趋势控制';
/*!40101 SET character_set_client = @saved_cs_client */;

-- Table structure for table `sys_ai_prediction_exception`
--

DROP TABLE IF EXISTS `sys_ai_prediction_exception`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_ai_prediction_exception` (
  `id` bigint NOT NULL COMMENT 'id',
  `reason_id` bigint NOT NULL COMMENT '模型推理id',
  `action_id` bigint NOT NULL COMMENT '处理通知id',
  `threshold` bigint NOT NULL COMMENT 'threshold',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '智能预警名称',
  `dev` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '测点名称',
  `model` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '预测模型',
  `interval` bigint NOT NULL COMMENT '预测时间段',
  `predict_length` bigint unsigned DEFAULT '0' COMMENT '预测长度',
  `agg_function` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '聚合函数，用来计算图表',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='预测列表';
/*!40101 SET character_set_client = @saved_cs_client */;



--
-- Table structure for table `sys_ai_prediction_list`
--

DROP TABLE IF EXISTS `sys_ai_prediction_list`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_ai_prediction_list` (
  `id` bigint NOT NULL COMMENT 'id',
  `reason_id` bigint NOT NULL COMMENT '模型推理id',
  `action_id` bigint NOT NULL COMMENT '处理通知id',
  `threshold` bigint NOT NULL COMMENT 'threshold',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '智能预警名称',
  `dev` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '测点名称',
  `advanced` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '告警规则',
  `model` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '预测模型',
  `field` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '预测字段',
  `interval` bigint NOT NULL COMMENT '预测时间段',
  `predict_length` bigint unsigned DEFAULT '0' COMMENT '预测长度',
  `agg_function` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '聚合函数，用来计算图表',
  `perturbation_variables` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '关联变量',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='预测列表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_ai_prediction_list`
--



--
-- Table structure for table `sys_alert`
--

DROP TABLE IF EXISTS `sys_alert`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_alert` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增标识',
  `gateway_id` bigint NOT NULL COMMENT '网关id',
  `template_id` bigint NOT NULL COMMENT '模板id',
  `name` varchar(255) NOT NULL COMMENT '告警策略名称',
  `additions` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '注解',
  `advanced` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '告警规则',
  `ignore` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '忽略规则',
  `matcher` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '匹配规则',
  `description` varchar(64) NOT NULL COMMENT '配置版本',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  `status` tinyint(1) DEFAULT '0' COMMENT '操作状态（0禁用 1启用）',
  `reason_id` bigint NOT NULL COMMENT '模型推理id',
  `action_id` bigint NOT NULL COMMENT '模型推理id',
  PRIMARY KEY (`id`),
  KEY `gateway_id` (`gateway_id`) USING BTREE,
  KEY `template_id` (`template_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=397461697243123713 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='告警策略配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_alert`
--


--
-- Table structure for table `sys_alert_action`
--

DROP TABLE IF EXISTS `sys_alert_action`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_alert_action` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(255) NOT NULL COMMENT '告警策略名称',
  `period` varchar(255) DEFAULT NULL COMMENT '通知周期，逗号分隔',
  `user_notify` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '通知用户',
  `api_notify` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '回调通知',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='告警处理动作';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_alert_ai_reason`
--

DROP TABLE IF EXISTS `sys_alert_ai_reason`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_alert_ai_reason` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(255) NOT NULL COMMENT '告警策略名称',
  `prompt` text COMMENT '提示词设置',
  `message` text COMMENT '提问消息模板',
  `agent_id` varchar(64) DEFAULT '' COMMENT '大模型agentid',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=398504110250266625 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='告警ai 推理配置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_alert_ai_reason`
--


--
-- Table structure for table `sys_alert_log`
--

DROP TABLE IF EXISTS `sys_alert_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_alert_log` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增标识',
  `gateway_id` bigint NOT NULL COMMENT '网关id',
  `alert_id` bigint NOT NULL COMMENT '告警id',
  `device_id` bigint NOT NULL COMMENT '设备id',
  `device_template_id` bigint NOT NULL COMMENT '设备模板id',
  `device_data_id` bigint NOT NULL COMMENT '设备数据id',
  `context` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '内容',
  `message` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '消息',
  `data` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '消息快照',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  `reason` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '推理数据',
  PRIMARY KEY (`id`),
  KEY `gateway_id` (`gateway_id`) USING BTREE,
  KEY `alert_id` (`alert_id`) USING BTREE,
  KEY `device_id` (`device_id`) USING BTREE,
  KEY `device_template_id` (`device_template_id`) USING BTREE,
  KEY `device_data_id` (`device_data_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=434023369021591553 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='告警日志';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_alert_sink_template`
--

DROP TABLE IF EXISTS `sys_alert_sink_template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_alert_sink_template` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增标识',
  `name` varchar(255) NOT NULL COMMENT '告警模板名称',
  `addr` varchar(512) NOT NULL COMMENT '发送alert的http地址，若为空，则不会发送',
  `timeout` int NOT NULL COMMENT '发送alert的http timeout',
  `headers` text COMMENT '发送alert的http header',
  `method` varchar(16) NOT NULL COMMENT '发送alert的http method, 如果不填put(不区分大小写)，都认为是POST',
  `template` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '告警模板',
  `extension` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '扩展信息',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=393670255097942017 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_alert_sink_template`
--


--
-- Table structure for table `sys_building`
--

DROP TABLE IF EXISTS `sys_building`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_building` (
  `id` bigint NOT NULL COMMENT 'id',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '建筑物名称',
  `code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '建筑物编号',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '建筑物类型',
  `status` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '建筑物状态',
  `area` bigint NOT NULL COMMENT '建筑面积',
  `floors` bigint NOT NULL COMMENT '楼层数',
  `year` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '建造年份',
  `lifeYears` bigint NOT NULL COMMENT '使用年限',
  `manager` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '负责人w',
  `phone` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '单位',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '单位',
  `remark` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '单位',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='设备管理';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sys_config`
--

DROP TABLE IF EXISTS `sys_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_config` (
  `config_id` bigint NOT NULL AUTO_INCREMENT COMMENT '参数主键',
  `config_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '参数名称',
  `config_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '参数键名',
  `config_value` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '参数键值',
  `config_type` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'N' COMMENT '系统内置（Y是 N否）',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`config_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='参数配置表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_config`
--

LOCK TABLES `sys_config` WRITE;
/*!40000 ALTER TABLE `sys_config` DISABLE KEYS */;
INSERT INTO `sys_config` VALUES (1,'主框架页-默认皮肤样式名称','sys.index.skinName','skin-blue','Y',1,'2024-02-08 04:10:56',1,NULL,'蓝色 skin-blue、绿色 skin-green、紫色 skin-purple、红色 skin-red、黄色 skin-yellow'),(2,'用户管理-账号初始密码','sys.user.initPassword','123456','Y',1,'2024-02-08 04:10:56',1,NULL,'初始化密码 123456'),(3,'主框架页-侧边栏主题','sys.index.sideTheme','theme-dark','Y',1,'2024-02-08 04:10:56',1,NULL,'深色主题theme-dark，浅色主题theme-light'),(4,'账号自助-验证码开关','sys.account.captchaEnabled','false','Y',1,'2024-02-08 04:10:56',1,'2024-03-09 13:17:45','是否开启验证码功能（true开启，false关闭）'),(5,'账号自助-是否开启用户注册功能','sys.account.registerUser','true','Y',1,'2024-02-08 04:10:56',1,'2024-03-09 13:24:50','是否开启注册用户功能（true开启，false关闭）');
/*!40000 ALTER TABLE `sys_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_configuration`
--

DROP TABLE IF EXISTS `sys_configuration`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_configuration` (
  `id` bigint NOT NULL COMMENT '组态id',
  `name` varchar(255) NOT NULL COMMENT '组态名称',
  `tag` varchar(255) NOT NULL COMMENT '标签',
  `annotation` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '注解',
  `description` varchar(2048) NOT NULL COMMENT '描述',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_craft_route`
--

DROP TABLE IF EXISTS `sys_craft_route`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_craft_route` (
  `route_id` bigint NOT NULL AUTO_INCREMENT COMMENT '工艺路线ID',
  `route_code` varchar(64) NOT NULL COMMENT '工艺路线编号',
  `route_name` varchar(255) NOT NULL COMMENT '工艺路线名称',
  `route_desc` varchar(500) DEFAULT NULL COMMENT '工艺路线说明',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `attr1` varchar(64) DEFAULT NULL COMMENT '预留字段1',
  `attr2` varchar(255) DEFAULT NULL COMMENT '预留字段2',
  `attr3` int DEFAULT '0' COMMENT '预留字段3',
  `attr4` int DEFAULT '0' COMMENT '预留字段4',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `status` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 1异常）',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`route_id`)
) ENGINE=InnoDB AUTO_INCREMENT=430357899168976897 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='工艺路线表';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_craft_route_config`
--

DROP TABLE IF EXISTS `sys_craft_route_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
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
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_dashboard`
--

DROP TABLE IF EXISTS `sys_dashboard`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dashboard` (
  `id` bigint NOT NULL COMMENT '调度id',
  `name` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  `description` varchar(2048) NOT NULL COMMENT '描述',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_dashboard_data`
--

DROP TABLE IF EXISTS `sys_dashboard_data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dashboard_data` (
  `id` bigint NOT NULL COMMENT '面板id',
  `datashboard_id` bigint NOT NULL COMMENT '仪表盘数据id',
  `data` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '数据',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `datashboard_id` (`datashboard_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_dataset`
--

DROP TABLE IF EXISTS `sys_dataset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dataset` (
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
  `dataset_update_date` datetime DEFAULT NULL COMMENT '更新时间',
  `dataset_update_time` bigint NOT NULL COMMENT '更新时间',
  `dataset_vector_similarity_weight` decimal(10,2) DEFAULT '0.00' COMMENT '矢量余弦相似性的权重。',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  `dataset_chunk_count` bigint NOT NULL DEFAULT '0' COMMENT 'chunk num',
  PRIMARY KEY (`dataset_id`) USING BTREE,
  UNIQUE KEY `dataset_name` (`dataset_name`,`state`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='物料出库管理';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_dataset_document`
--

DROP TABLE IF EXISTS `sys_dataset_document`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dataset_document` (
  `document_id` bigint NOT NULL COMMENT '文档id',
  `dataset_id` bigint NOT NULL COMMENT '数据集id',
  `dataset_chunk_method` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '要创建的数据集的分块方法',
  `dataset_created_by` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '创建人',
  `dataset_document_uuid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '文档id',
  `dataset_dataset_uuid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '知识库id',
  `dataset_language` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '知识库语言',
  `dataset_location` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '知识库定位',
  `dataset_name` varchar(136) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '要创建的数据集的唯一名称',
  `dataset_parser_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '数据集解析器的配置设置。',
  `dataset_run` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '知识库运行状态',
  `dataset_size` bigint NOT NULL DEFAULT '0' COMMENT '知识库尺寸大小',
  `dataset_thumbnail` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '知识库锁略图',
  `dataset_type` varchar(136) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '知识库类型',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`document_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='文档管理';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_dept`
--

DROP TABLE IF EXISTS `sys_dept`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dept` (
  `dept_id` bigint NOT NULL AUTO_INCREMENT COMMENT '部门id',
  `parent_id` bigint DEFAULT '0' COMMENT '父部门id',
  `ancestors` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '祖级列表',
  `dept_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '部门名称',
  `order_num` int DEFAULT '0' COMMENT '显示顺序',
  `leader` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '负责人',
  `phone` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '联系电话',
  `email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '邮箱',
  `status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '0' COMMENT '部门状态（0正常 1停用）',
  `del_flag` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '0' COMMENT '删除标志（0代表存在 2代表删除）',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`dept_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=110 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='部门表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dept`
--

LOCK TABLES `sys_dept` WRITE;
/*!40000 ALTER TABLE `sys_dept` DISABLE KEYS */;
INSERT INTO `sys_dept` VALUES (100,0,'0','白泽科技',0,'白泽','15888888888','bz@qq.com','0','0',1,'2024-02-08 04:10:55',1,'2025-10-04 20:14:07'),(101,100,'0,100','深圳总公司',1,'白泽','15888888888','bz@qq.com','0','0',1,'2024-02-08 04:10:55',1,NULL),(102,100,'0,100','长沙分公司',2,'白泽','15888888888','bz@qq.com','0','0',1,'2024-02-08 04:10:55',1,NULL),(103,101,'0,100,101','研发部门',1,'白泽','15888888888','bz@qq.com','0','0',1,'2024-02-08 04:10:55',1,NULL),(104,101,'0,100,101','市场部门',2,'白泽','15888888888','bz@qq.com','0','0',1,'2024-02-08 04:10:55',1,NULL),(105,101,'0,100,101','测试部门',3,'白泽','15888888888','bz@qq.com','0','0',1,'2024-02-08 04:10:55',1,NULL),(106,101,'0,100,101','财务部门',4,'白泽','15888888888','bz@qq.com','0','0',1,'2024-02-08 04:10:55',1,NULL),(107,101,'0,100,101','运维部门',5,'白泽','15888888888','bz@qq.com','0','0',1,'2024-02-08 04:10:55',1,NULL),(108,102,'0,100,102','市场部门',1,'白泽','15888888888','bz@qq.com','0','0',1,'2024-02-08 04:10:55',1,NULL),(109,102,'0,100,102','财务部门',2,'白泽','15888888888','bz@qq.com','0','0',1,'2024-02-08 04:10:55',1,NULL);
/*!40000 ALTER TABLE `sys_dept` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_device`
--

DROP TABLE IF EXISTS `sys_device`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_device` (
  `device_id` bigint NOT NULL COMMENT '设备主键',
  `device_group_id` bigint NOT NULL COMMENT '设备分组id',
  `device_class_id` bigint NOT NULL COMMENT '设备分类id',
  `device_protocol_id` bigint NOT NULL COMMENT '协议id',
  `device_building_id` bigint NOT NULL COMMENT '建筑物id',
  `device_gateway_id` bigint NOT NULL COMMENT '网关id',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '设备名称',
  `number` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '设备编号',
  `communication_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '通信方式',
  `protocol_type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '协议类型',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '设备类型',
  `action` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '行为字段',
  `extension` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '扩展字段',
  `control_type` int DEFAULT '0' COMMENT '控制模式 0 是手动 1是自动',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `status` int DEFAULT '0' COMMENT '操作状态（0正常 1异常）',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  `enable` tinyint(1) DEFAULT '0' COMMENT '启用状态（0正常 1异常）',
  `model_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '模型图片',
  PRIMARY KEY (`device_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='设备管理';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_device_electric_setting`
--

DROP TABLE IF EXISTS `sys_device_electric_setting`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_device_electric_setting` (
  `id` bigint NOT NULL COMMENT 'id',
  `device_id` bigint NOT NULL COMMENT 'device_id',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '配置名称',
  `expression` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '表达式',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uniq_device_id` (`device_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='设备电流配置，用来计算设备电流是待机 停机 还是 运行';
/*!40101 SET character_set_client = @saved_cs_client */;



--
-- Table structure for table `sys_device_group`
--

DROP TABLE IF EXISTS `sys_device_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_device_group` (
  `group_id` bigint NOT NULL COMMENT '设备主键',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '设备名称',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '备注',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`group_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='设备组';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_device_group`
--


--
-- Table structure for table `sys_device_template`
--

DROP TABLE IF EXISTS `sys_device_template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_device_template` (
  `template_id` bigint NOT NULL COMMENT '设备主键',
  `name` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '设备名称',
  `template_type` tinyint(1) DEFAULT '0' COMMENT '模板类型0是私有1 是共有',
  `vendor` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '供应商',
  `protocol` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '协议类型',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '备注',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`template_id`) USING BTREE,
  KEY `idx_protocol` (`protocol`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='modbus数据模板';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_dict_data`
--

DROP TABLE IF EXISTS `sys_dict_data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dict_data` (
  `dict_code` bigint NOT NULL AUTO_INCREMENT COMMENT '字典编码',
  `dict_sort` int DEFAULT '0' COMMENT '字典排序',
  `dict_label` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '字典标签',
  `dict_value` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '字典键值',
  `dict_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '字典类型',
  `css_class` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '样式属性（其他样式扩展）',
  `list_class` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '表格回显样式',
  `is_default` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'N' COMMENT '是否默认（Y是 N否）',
  `status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '0' COMMENT '状态（0正常 1停用）',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '备注',
  PRIMARY KEY (`dict_code`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=467291272223133697 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='字典数据表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dict_data`
--

LOCK TABLES `sys_dict_data` WRITE;
/*!40000 ALTER TABLE `sys_dict_data` DISABLE KEYS */;
INSERT INTO `sys_dict_data` VALUES (1,1,'男','0','sys_user_sex','','','Y','0',1,'2024-02-08 04:10:56',1,NULL,'性别男'),(2,2,'女','1','sys_user_sex','','','N','0',1,'2024-02-08 04:10:56',1,NULL,'性别女'),(3,3,'未知','2','sys_user_sex','','','N','0',1,'2024-02-08 04:10:56',1,NULL,'性别未知'),(6,1,'正常','0','sys_normal_disable','','primary','Y','0',1,'2024-02-08 04:10:56',1,NULL,'正常状态'),(7,2,'停用','1','sys_normal_disable','','danger','N','0',1,'2024-02-08 04:10:56',1,NULL,'停用状态'),(8,1,'正常','0','sys_job_status','','primary','Y','0',1,'2024-02-08 04:10:56',1,NULL,'正常状态'),(9,2,'暂停','1','sys_job_status','','danger','N','0',1,'2024-02-08 04:10:56',1,NULL,'停用状态'),(10,1,'默认','DEFAULT','sys_job_group','','','Y','0',1,'2024-02-08 04:10:56',1,NULL,'默认分组'),(11,2,'系统','SYSTEM','sys_job_group','','','N','0',1,'2024-02-08 04:10:56',1,NULL,'系统分组'),(12,1,'是','Y','sys_yes_no','','primary','Y','0',1,'2024-02-08 04:10:56',1,NULL,'系统默认是'),(13,2,'否','N','sys_yes_no','','danger','N','0',1,'2024-02-08 04:10:56',1,NULL,'系统默认否'),(14,1,'通知','1','sys_notice_type','','warning','Y','0',1,'2024-02-08 04:10:56',1,NULL,'通知'),(15,2,'公告','2','sys_notice_type','','success','N','0',1,'2024-02-08 04:10:56',1,NULL,'公告'),(18,0,'其他','0','sys_oper_type','','info','N','0',1,'2024-02-08 04:10:56',1,NULL,'其他操作'),(19,1,'新增','1','sys_oper_type','','info','N','0',1,'2024-02-08 04:10:56',1,NULL,'新增操作'),(20,2,'修改','2','sys_oper_type','','info','N','0',1,'2024-02-08 04:10:56',1,NULL,'修改操作'),(21,3,'删除','3','sys_oper_type','','danger','N','0',1,'2024-02-08 04:10:56',1,NULL,'删除操作'),(22,4,'强退','4','sys_oper_type','','danger','N','0',1,'2024-02-08 04:10:56',1,'2024-04-12 14:33:31','强退操作'),(23,5,'清空数据','5','sys_oper_type','','danger','N','0',1,'2024-02-08 04:10:56',1,NULL,'清空操作'),(28,1,'成功','0','sys_common_status','','primary','N','0',1,'2024-02-08 04:10:56',1,NULL,'正常状态'),(29,2,'失败','1','sys_common_status','','danger','N','0',1,'2024-02-08 04:10:56',1,NULL,'停用状态'),(374007595364519936,0,'数据属性','data','data_type','','default','','0',1,'2025-06-05 09:32:39',1,'2025-06-05 09:32:39',''),(374007672149643264,1,'配置属性','config','data_type','','default','','0',1,'2025-06-05 09:32:58',1,'2025-06-05 09:32:58',''),(374008691638145024,0,'数据属性','data','数据类型','','default','','0',1,'2025-06-05 09:37:01',1,'2025-06-05 09:37:01',''),(374008751314702336,0,'配置属性','config','数据类型','','default','','0',1,'2025-06-05 09:37:15',1,'2025-06-05 09:37:15',''),(374012607507468288,0,' ModBus TCP','modbus-tcp','protocol','','default','','0',1,'2025-06-05 09:52:34',1,'2025-06-10 11:03:33',''),(374016522445656064,0,'ABCD','ABCD','device_template_data_sort','','default','','0',1,'2025-06-05 10:08:08',1,'2025-06-05 10:08:08',''),(374016715719184384,0,'CDAB','CDAB','device_template_data_sort','','default','','0',1,'2025-06-05 10:08:54',1,'2025-06-05 10:08:54',''),(374049816121970688,0,'变化储存','change_storage','device_template_data_storage_type','','default','','0',1,'2025-06-05 12:20:26',1,'2025-06-05 12:20:26',''),(374050461042348032,0,'实时存储','peer_storage','device_template_data_storage_type','','default','','0',1,'2025-06-05 12:22:59',1,'2025-06-05 12:22:59',''),(374050558908043264,0,'30s','peer_30s_storage','device_template_data_storage_type','','default','','0',1,'2025-06-05 12:23:23',1,'2025-06-05 12:23:23',''),(374050635844161536,0,'1min','peer_1min_storage','device_template_data_storage_type','','default','','0',1,'2025-06-05 12:23:41',1,'2025-06-05 12:23:41',''),(374050696531546112,0,'5min','peer_5min_storage','device_template_data_storage_type','','default','','0',1,'2025-06-05 12:23:56',1,'2025-06-05 12:23:56',''),(374050777607442432,0,'10min','peer_10min_storage','device_template_data_storage_type','','default','','0',1,'2025-06-05 12:24:15',1,'2025-06-05 12:24:15',''),(374051245578522624,0,'30min','peer_30min_storage','device_template_data_storage_type','','default','','0',1,'2025-06-05 12:26:06',1,'2025-06-05 12:26:06',''),(374051339556098048,0,'1h','peer_1h_storage','device_template_data_storage_type','','default','','0',1,'2025-06-05 12:26:29',1,'2025-06-05 12:26:29',''),(374051411689738240,0,'1day','peer_1day_storage','device_template_data_storage_type','','default','','0',1,'2025-06-05 12:26:46',1,'2025-06-05 12:26:46',''),(374053154020397056,0,'只读','0','device_template_data_mode','','default','','0',1,'2025-06-05 12:33:41',1,'2025-06-09 20:47:08',''),(374053235754799104,0,'读写','1','device_template_data_mode','','default','','0',1,'2025-06-05 12:34:01',1,'2025-06-09 20:47:11',''),(374053539539849216,0,'0','0','device_template_data_precision','','default','','0',1,'2025-06-05 12:35:13',1,'2025-06-05 12:35:13',''),(374053568153391104,0,'1','1','device_template_data_precision','','default','','0',1,'2025-06-05 12:35:20',1,'2025-06-05 12:35:20',''),(374053581700993024,0,'2','2','device_template_data_precision','','default','','0',1,'2025-06-05 12:35:23',1,'2025-06-05 12:35:23',''),(374053592912367616,0,'3','3','device_template_data_precision','','default','','0',1,'2025-06-05 12:35:26',1,'2025-06-05 12:35:26',''),(374053605205872640,0,'4','4','device_template_data_precision','','default','','0',1,'2025-06-05 12:35:29',1,'2025-06-05 12:35:29',''),(374053618648616960,0,'5','5','device_template_data_precision','','default','','0',1,'2025-06-05 12:35:32',1,'2025-06-05 12:35:32',''),(374053641729871872,0,'6','6','device_template_data_precision','','default','','0',1,'2025-06-05 12:35:38',1,'2025-06-05 12:35:38',''),(374067336803520512,0,'读线圈状态','1','device_template_data_function_code','','default','','0',1,'2025-06-05 13:30:03',1,'2025-06-05 13:30:03',''),(374067438557335552,0,'读离散输入状态','2','device_template_data_function_code','','default','','0',1,'2025-06-05 13:30:27',1,'2025-06-05 13:30:27',''),(374067477384007680,0,'读保持寄存器','3','device_template_data_function_code','','default','','0',1,'2025-06-05 13:30:36',1,'2025-06-05 13:30:36',''),(374067509378158592,0,'读输入寄存器','4','device_template_data_function_code','','default','','0',1,'2025-06-05 13:30:44',1,'2025-06-05 13:30:44',''),(374067536741797888,0,'写单个线圈','5','device_template_data_function_code','','default','','0',1,'2025-06-05 13:30:51',1,'2025-06-05 13:30:51',''),(374067645424603136,0,'写单个保持寄存器','6','device_template_data_function_code','','default','','0',1,'2025-06-05 13:31:17',1,'2025-06-05 13:31:17',''),(374068012463951872,0,'写多个线圈','15','device_template_data_function_code','','default','','0',1,'2025-06-05 13:32:44',1,'2025-06-05 13:32:44',''),(374068446033350656,0,'写多个保持寄存器','16','device_template_data_function_code','','default','','0',1,'2025-06-05 13:34:27',1,'2025-06-05 13:34:27',''),(375625676335616000,0,'最小的数据单位（0或1）','bit','device_template_data_format','','default','','0',1,'2025-06-09 20:42:20',1,'2025-07-01 15:36:32',''),(375625704278069248,0,'uint8','uint8','device_template_data_format','','default','','0',1,'2025-06-09 20:42:27',1,'2025-06-09 20:42:27',''),(375625726545629184,0,'uint16','uint16','device_template_data_format','','default','','0',1,'2025-06-09 20:42:32',1,'2025-06-09 20:42:32',''),(375625748553142272,0,'int16','int16','device_template_data_format','','default','','0',1,'2025-06-09 20:42:37',1,'2025-06-09 20:42:37',''),(375625769906343936,0,'uint32','uint32','device_template_data_format','','default','','0',1,'2025-06-09 20:42:42',1,'2025-06-09 20:42:42',''),(375625789929951232,0,'int32','int32','device_template_data_format','','default','','0',1,'2025-06-09 20:42:47',1,'2025-06-09 20:42:47',''),(375625813279641600,0,'float32','float32','device_template_data_format','','default','','0',1,'2025-06-09 20:42:53',1,'2025-06-09 20:42:53',''),(375625837870845952,0,'float64','float64','device_template_data_format','','default','','0',1,'2025-06-09 20:42:59',1,'2025-06-09 20:42:59',''),(375827192547905536,0,'局域网','local','device_communication_type','','default','','0',1,'2025-06-10 10:03:05',1,'2025-06-10 10:03:05',''),(375827224370089984,0,'4G','4G','device_communication_type','','default','','0',1,'2025-06-10 10:03:13',1,'2025-06-10 10:03:13',''),(376200464347172864,0,'在线','1','agent_active','','default','','0',1,'2025-06-11 10:46:20',1,'2025-06-11 10:46:38',''),(376200567371862016,0,'离线','0','agent_active','','default','','0',1,'2025-06-11 10:46:45',1,'2025-06-11 10:46:45',''),(376200596056707072,0,'故障','2','agent_active','','default','','0',1,'2025-06-11 10:46:52',1,'2025-06-11 10:46:52',''),(376373915414433792,0,'初始化中','0','agent_process_status','','default','','0',1,'2025-06-11 22:15:34',1,'2025-06-11 22:17:07',''),(376373959488180224,0,'启动中','1','agent_process_status','','default','','0',1,'2025-06-11 22:15:45',1,'2025-06-11 22:17:15',''),(376374005935902720,0,'运行中','2','agent_process_status','','default','','0',1,'2025-06-11 22:15:56',1,'2025-06-11 22:17:25',''),(376374038882160640,0,'停止中','3','agent_process_status','','default','','0',1,'2025-06-11 22:16:04',1,'2025-06-11 22:17:37',''),(376374112265703424,0,'中断','4','agent_process_status','','default','','0',1,'2025-06-11 22:16:21',1,'2025-06-11 22:17:49',''),(376374151176261632,0,'完成','5','agent_process_status','','default','','0',1,'2025-06-11 22:16:30',1,'2025-06-11 22:17:58',''),(376374203953188864,0,'致命错误','6','agent_process_status','','default','','0',1,'2025-06-11 22:16:43',1,'2025-06-11 22:18:13',''),(376374258172956672,0,'未知','7','agent_process_status','','default','','0',1,'2025-06-11 22:16:56',1,'2025-06-11 22:18:21',''),(382065020508311552,0,'当前工序开始生产，下一道工序才开始生产','s-to-s','craft_process_relation','','default','','0',1,'2025-06-27 15:09:59',1,'2025-06-27 15:09:59',''),(382065236678545408,0,'当前工序结束生产，下一道工序才结束生产','f-to-f','craft_process_relation','','default','','0',1,'2025-06-27 15:10:51',1,'2025-06-27 15:11:01',''),(382065413204217856,0,'当前工序开始生产，下一道工序才结束生产','s-to-f','craft_process_relation','','default','','0',1,'2025-06-27 15:11:33',1,'2025-06-27 15:11:33',''),(382065538358054912,0,'当前工序结束生产，下一道工序才开始生产','f-to-s','craft_process_relation','','default','','0',1,'2025-06-27 15:12:03',1,'2025-06-27 15:12:03',''),(383253308694859776,0,'数值类型','value','node_data_type','','default','','0',1,'2025-06-30 21:51:49',1,'2025-06-30 22:14:05',''),(383259277462081536,0,'状态类型','status','node_data_type','','default','','0',1,'2025-06-30 22:15:32',1,'2025-06-30 22:15:32',''),(383259349276954624,0,'开关类型','switch','node_data_type','','default','','0',1,'2025-06-30 22:15:49',1,'2025-06-30 22:15:49',''),(383259486212591616,0,'GPS类型','gps','node_data_type','','default','','0',1,'2025-06-30 22:16:22',1,'2025-06-30 22:16:22',''),(383259572770443264,0,'显示型','display','node_data_type','','default','','0',1,'2025-06-30 22:16:43',1,'2025-06-30 22:16:43',''),(383259703053914112,0,'动态型','dynamic','node_data_type','','default','','0',1,'2025-06-30 22:17:14',1,'2025-06-30 22:17:14',''),(386790751595401216,0,'SUM','SUM','agg_function','','default','','0',1,'2025-07-10 16:08:21',1,'2025-07-10 16:08:49',''),(386790915278114816,0,'COUNT','COUNT','agg_function','','default','','0',1,'2025-07-10 16:09:00',1,'2025-07-10 16:09:00',''),(386790941073084416,0,'AVG','AVG','agg_function','','default','','0',1,'2025-07-10 16:09:07',1,'2025-07-10 16:09:07',''),(386791000200187904,0,'MAX_VALUE','MAX_VALUE','agg_function','','default','','0',1,'2025-07-10 16:09:21',1,'2025-07-10 16:09:21',''),(386791026263592960,0,'MIN_VALUE','MIN_VALUE','agg_function','','default','','0',1,'2025-07-10 16:09:27',1,'2025-07-10 16:09:27',''),(393668319107878912,0,'http地址','http://81.71.98.26:8080','http','','default','','0',1,'2025-07-29 15:37:21',1,'2025-07-29 15:37:21',''),(393737757744173056,0,'匹配告警','regexp','alarm_mode','','default','','0',1,'2025-07-29 20:13:17',1,'2025-08-09 03:01:17',''),(393737758549479424,0,'无数据告警','noData','alarm_mode','','default','','0',1,'2025-07-29 20:13:17',1,'2025-08-09 03:00:59',''),(393738501008396288,0,'全部成立','all','alarm_match_type','','default','','0',1,'2025-07-29 20:16:14',1,'2025-08-09 03:02:08',''),(393738550815756288,0,'成立任意一个','any','alarm_match_type','','default','','0',1,'2025-07-29 20:16:26',1,'2025-08-09 03:02:23',''),(394759485693890560,0,'邮件','email','alert_action_channel','','default','','0',1,'2025-08-01 15:53:16',1,'2025-08-01 15:53:16',''),(394759606292713472,1,'短信','sms','alert_action_channel','','default','','0',1,'2025-08-01 15:53:44',1,'2025-08-01 16:16:37',''),(394759765021954048,2,'微信','wechat','alert_action_channel','','default','','0',1,'2025-08-01 15:54:22',1,'2025-08-01 16:16:40',''),(394759840594923520,3,'企业微信','work-wechat','alert_action_channel','','default','','0',1,'2025-08-01 15:54:40',1,'2025-08-01 16:16:44',''),(394765148495024128,4,'电话','phone','alert_action_channel','','default','','0',1,'2025-08-01 16:15:46',1,'2025-08-01 16:16:48',''),(396884699391201280,0,'告警信息','{{.Message}}','ai_reason_var','','default','','0',1,'2025-08-07 12:38:06',1,'2025-08-07 12:38:06',''),(396884769142476800,0,'设备名称','{{.Device}}','ai_reason_var','','default','','0',1,'2025-08-07 12:38:23',1,'2025-08-07 12:38:23',''),(397464996579119104,0,'>','gt','alert_operation','','default','','0',1,'2025-08-09 03:04:00',1,'2025-08-09 03:05:31','超过'),(397465163294314496,0,'=','eq','alert_operation','','default','','0',1,'2025-08-09 03:04:39',1,'2025-08-09 03:05:39','等于'),(397465226804465664,0,'<','lt','alert_operation','','default','','0',1,'2025-08-09 03:04:55',1,'2025-08-09 03:05:48','小于'),(397997102694666240,0,'url','http://localhost:8080','alarm_http_server','','default','','0',1,'2025-08-10 14:18:24',1,'2025-08-10 14:18:24',''),(399636645231464448,0,'生产车间','workshop','building_type','','default','','0',1,'2025-08-15 02:53:21',1,'2025-08-15 02:53:21',''),(399636689963716608,0,'仓库','warehouse','building_type','','default','','0',1,'2025-08-15 02:53:32',1,'2025-08-15 02:53:32',''),(399636794582241280,0,'办公楼','office','building_type','','default','','0',1,'2025-08-15 02:53:57',1,'2025-08-15 02:53:57',''),(399636837481582592,0,'宿舍楼','dormitory','building_type','','default','','0',1,'2025-08-15 02:54:07',1,'2025-08-15 02:54:07',''),(399636872965394432,0,'食堂','canteen','building_type','','default','','0',1,'2025-08-15 02:54:15',1,'2025-08-15 02:54:15',''),(399636910554746880,0,'其他','other','building_type','','default','','0',1,'2025-08-15 02:54:24',1,'2025-08-15 02:54:24',''),(399637065563639808,0,'正常使用','normal','building_status','','default','','0',1,'2025-08-15 02:55:01',1,'2025-08-15 02:55:01',''),(399637105661186048,0,'维修中','maintenance','building_status','','default','','0',1,'2025-08-15 02:55:11',1,'2025-08-15 02:55:11',''),(399637140310331392,0,'停用','disabled','building_status','','default','','0',1,'2025-08-15 02:55:19',1,'2025-08-15 02:55:19',''),(399637175399878656,0,'建设中','constructing','building_status','','default','','0',1,'2025-08-15 02:55:28',1,'2025-08-15 02:55:28',''),(400559835352928256,0,'业务监控','business','dashboard_type','','default','','0',1,'2025-08-17 16:01:47',1,'2025-08-17 16:01:47',''),(400559885915262976,0,'系统监控','system','dashboard_type','','default','','0',1,'2025-08-17 16:01:59',1,'2025-08-17 16:01:59',''),(400559927287877632,0,'设备监控','device','dashboard_type','','default','','0',1,'2025-08-17 16:02:09',1,'2025-08-17 16:02:09',''),(400559963954483200,0,'数据分析','analysis','dashboard_type','','default','','0',1,'2025-08-17 16:02:17',1,'2025-08-17 16:02:17',''),(400560001376063488,0,'报表统计','report','dashboard_type','','default','','0',1,'2025-08-17 16:02:26',1,'2025-08-17 16:02:26',''),(400560041364557824,0,'其他','other','dashboard_type','','default','','0',1,'2025-08-17 16:02:36',1,'2025-08-17 16:02:36',''),(401924240942567424,0,'预测-Arima','_Arima','predict_model','','default','','0',1,'2025-08-21 10:23:26',1,'2025-08-21 12:21:33','预测'),(401924277823082496,0,'预测-STLForecaster','_STLForecaster','predict_model','','default','','0',1,'2025-08-21 10:23:35',1,'2025-08-21 12:21:39','预测'),(401924343736569856,0,'预测-NaiveForecaster','_NaiveForecaster','predict_model','','default','','0',1,'2025-08-21 10:23:51',1,'2025-08-21 12:21:47','预测'),(401924370072604672,0,'预测--ExponentialSmoothing','_ExponentialSmoothing','predict_model','','default','','0',1,'2025-08-21 10:23:57',1,'2025-08-21 12:21:55','预测'),(401924394135326720,0,'标注-GaussianHMM','_GaussianHMM','predict_model','','default','','0',1,'2025-08-21 10:24:03',1,'2025-08-21 12:24:03','标注'),(401924423491260416,0,'标注-GMMHMM','_GMMHMM','predict_model','','default','','0',1,'2025-08-21 10:24:10',1,'2025-08-21 12:24:19','标注'),(401924459193176064,0,'异常检测-Stray','_Stray','predict_model','','default','','0',1,'2025-08-21 10:24:18',1,'2025-08-21 12:24:27','异常检测'),(401925937395929088,0,'1天前','86400','dashboard_interval','','default','','0',1,'2025-08-21 10:30:11',1,'2025-09-19 19:50:12',''),(401926002835460096,0,'1周前','604800','dashboard_interval','','default','','0',1,'2025-08-21 10:30:27',1,'2025-09-19 19:50:25',''),(401926067499044864,0,'1个月前','2592000','dashboard_interval','','default','','0',1,'2025-08-21 10:30:42',1,'2025-09-19 19:50:40',''),(401926523537330176,0,'3个月前','7776000','dashboard_interval','','default','','0',1,'2025-08-21 10:32:31',1,'2025-09-19 19:50:51',''),(401926586196037632,0,'半年前','15552000','dashboard_interval','','default','','0',1,'2025-08-21 10:32:46',1,'2025-09-19 19:51:04',''),(401926651308412928,0,'1年前','31104000','dashboard_interval','','default','','0',1,'2025-08-21 10:33:01',1,'2025-09-19 19:51:18',''),(412431620511895552,0,'预测-Arima','arima','prediction_model','','default','','0',1,'2025-09-19 10:16:01',1,'2025-11-03 19:59:28',''),(412431680876318720,0,'预测-STLForecaster','stl_forecaster','prediction_model','','default','','0',1,'2025-09-19 10:16:15',1,'2025-11-03 19:59:16',''),(412431735377104896,0,'预测-NaiveForecaster','naive_forecaster','prediction_model','','default','','0',1,'2025-09-19 10:16:28',1,'2025-11-03 19:59:42',''),(412431796022546432,0,'预测--ExponentialSmoothing','exponential_smoothing','prediction_model','','default','','0',1,'2025-09-19 10:16:43',1,'2025-11-03 19:59:53',''),(412969127079055360,0,'Stray-异常检测','_Stray','prediction_exception_model','','default','','0',1,'2025-09-20 21:51:52',1,'2025-09-20 21:51:52',''),(421706109854683136,0,'Mqtt','mqtt','protocol','','default','','0',1,'2025-10-15 00:29:32',1,'2025-10-15 00:29:50',''),(423126593976668160,0,'电流电压采集器','electric_dev','device_type','','default','','0',1,'2025-10-18 22:34:01',1,'2025-10-18 22:34:01',''),(423126767566327808,0,'数控设备','cnc_dev','device_type','','default','','0',1,'2025-10-18 22:34:43',1,'2025-10-18 22:34:43',''),(423185153506938880,0,'设备离线','0','device_run_status','','default','','0',1,'2025-10-19 02:26:43',1,'2025-10-19 02:27:51',''),(423185255109758976,0,'设备待机','1','device_run_status','','default','','0',1,'2025-10-19 02:27:07',1,'2025-10-19 02:27:07',''),(423185302765441024,0,'设备运行','2','device_run_status','','default','','0',1,'2025-10-19 02:27:19',1,'2025-10-19 02:27:43',''),(424103059698749440,0,'address','localhost:6000','grpc_sink_address','','default','','0',1,'2025-10-21 15:14:09',1,'2025-10-21 15:14:09',''),(433558951187976192,0,'服务器主机','http://localhost:8080','server_host','','default','','0',1,'2025-11-16 17:28:29',1,'2026-01-23 17:01:23',''),(441645267779850240,0,'煤矸石','gangue','laboratory_material','','default','','0',1,'2025-12-09 01:00:37',1,'2025-12-12 17:08:52',''),(441645313241911296,0,'煤灰','coom','laboratory_material','','default','','0',1,'2025-12-09 01:00:48',1,'2025-12-09 01:00:48',''),(441645354505474048,0,'炉渣','slag','laboratory_material','','default','','0',1,'2025-12-09 01:00:58',1,'2026-01-22 10:46:15',''),(448414513125920768,0,'粉混料','mixed_material','assaying_material','','default','','0',1,'2025-12-27 17:19:11',1,'2025-12-27 17:19:11',''),(451717319635243008,0,'环保资料','environmental','resource_file_category','','default','','0',1,'2026-01-05 20:03:22',1,'2026-01-05 20:03:22',''),(451717367135735808,0,'销售资料','sale','resource_file_category','','default','','0',1,'2026-01-05 20:03:33',1,'2026-01-05 20:03:33',''),(457737783843229696,0,'污泥','sludge','laboratory_material','','default','','0',1,'2026-01-22 10:46:32',1,'2026-01-22 10:46:32',''),(461488519865438208,0,'address','localhost:10050','grpc_control_server_address','','default','','0',1,'2026-02-01 19:10:37',1,'2026-02-01 19:10:37',''),(467288849161129984,0,'阈值算法','threshold','control_type','','default','','0',1,'2026-02-17 19:19:03',1,'2026-02-17 19:19:03',''),(467288882132553728,0,'PID算法','pid','control_type','','default','','0',1,'2026-02-17 19:19:11',1,'2026-02-17 19:19:11',''),(467288936381681664,0,'模型预测算法','mpc','control_type','','default','','0',1,'2026-02-17 19:19:24',1,'2026-02-17 19:19:24',''),(467291229168603136,0,'即时控制','concurrent_control','control_mode','','default','','0',1,'2026-02-17 19:28:31',1,'2026-02-17 19:28:31',''),(467291272223133696,0,'延时控制','delay_control','control_mode','','default','','0',1,'2026-02-17 19:28:41',1,'2026-02-17 19:28:41','');
/*!40000 ALTER TABLE `sys_dict_data` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_dict_type`
--

DROP TABLE IF EXISTS `sys_dict_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dict_type` (
  `dict_id` bigint NOT NULL AUTO_INCREMENT COMMENT '字典主键',
  `dict_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '字典名称',
  `dict_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '字典类型',
  `status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '0' COMMENT '状态（0正常 1停用）',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`dict_id`) USING BTREE,
  UNIQUE KEY `dict_type` (`dict_type`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=467290806022049793 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='字典类型表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dict_type`
--

LOCK TABLES `sys_dict_type` WRITE;
/*!40000 ALTER TABLE `sys_dict_type` DISABLE KEYS */;
INSERT INTO `sys_dict_type` VALUES (1,'用户性别','sys_user_sex','0',1,'2024-02-08 04:10:55',1,NULL,'用户性别列表'),(3,'系统开关','sys_normal_disable','0',1,'2024-02-08 04:10:55',1,NULL,'系统开关列表'),(4,'任务状态','sys_job_status','0',1,'2024-02-08 04:10:55',1,NULL,'任务状态列表'),(5,'任务分组','sys_job_group','0',1,'2024-02-08 04:10:55',1,NULL,'任务分组列表'),(6,'系统是否','sys_yes_no','0',1,'2024-02-08 04:10:55',1,NULL,'系统是否列表'),(7,'通知类型','sys_notice_type','0',1,'2024-02-08 04:10:55',1,NULL,'通知类型列表'),(9,'操作类型','sys_oper_type','0',1,'2024-02-08 04:10:55',1,NULL,'操作类型列表'),(10,'系统状态','sys_common_status','0',1,'2024-02-08 04:10:55',1,NULL,'登录状态列表'),(374008168541327360,'数据类型','data_type','0',1,'2025-06-05 09:34:56',1,'2025-06-05 09:36:15',''),(374012104740442112,'协议类型','protocol','0',1,'2025-06-05 09:50:35',1,'2025-06-05 09:50:35','设备模版'),(374016186398019584,'数据顺序','device_template_data_sort','0',1,'2025-06-05 10:06:48',1,'2025-06-05 10:06:48','设备模板数据顺序'),(374047324155940864,'存储策略','device_template_data_storage_type','0',1,'2025-06-05 12:10:32',1,'2025-06-05 12:10:32','存储策略'),(374047797793525760,'数据精度','device_template_data_precision','0',1,'2025-06-05 12:12:24',1,'2025-06-05 12:12:24','数据精度'),(374048117458210816,'读写方式','device_template_data_mode','0',1,'2025-06-05 12:13:41',1,'2025-06-05 12:13:41','读写方式'),(374048618883059712,'数据格式','device_template_data_format','0',1,'2025-06-05 12:15:40',1,'2025-06-05 12:15:40','数据格式'),(374066372566585344,'功能码','device_template_data_function_code','0',1,'2025-06-05 13:26:13',1,'2025-06-05 13:26:20','功能码'),(375827071483514880,'通信方式','device_communication_type','0',1,'2025-06-10 10:02:36',1,'2025-06-10 10:02:46','通信方式'),(376200332595695616,'agent状态','agent_active','0',1,'2025-06-11 10:45:49',1,'2025-06-11 10:45:49',''),(376373575113773056,'agent进程状态','agent_process_status','0',1,'2025-06-11 22:14:13',1,'2025-06-11 22:14:13',''),(382064790723366912,'工艺工序关系','craft_process_relation','0',1,'2025-06-27 15:09:05',1,'2025-06-27 15:09:05',''),(383252987469893632,'设备节点类型','node_data_type','0',1,'2025-06-30 21:50:33',1,'2025-06-30 21:50:33','设备模版 -详情 -设备节点类型'),(386790261985906688,'图表聚合函数','agg_function','0',1,'2025-07-10 16:06:25',1,'2025-07-10 16:06:28','图表聚合函数'),(393667656319766528,'http地址','alarm_http_server','0',1,'2025-07-29 15:34:43',1,'2025-08-10 14:16:38',''),(393737592379543552,'模式','alarm_mode','0',1,'2025-07-29 20:12:37',1,'2025-07-29 20:14:56','告警模式'),(393738419290771456,'匹配类型','alarm_match_type','0',1,'2025-07-29 20:15:54',1,'2025-07-29 20:15:54','告警匹配类型'),(394759394698465280,'告警通知渠道','alert_action_channel','0',1,'2025-08-01 15:52:54',1,'2025-08-01 15:52:54',''),(396883892335808512,'告警AI推理变量','ai_reason_var','0',1,'2025-08-07 12:34:54',1,'2025-08-07 12:34:54',''),(397464900613443584,'告警操作符','alert_operation','0',1,'2025-08-09 03:03:37',1,'2025-08-09 03:03:49','告警操作符'),(399636315072630784,'建筑物类型','building_type','0',1,'2025-08-15 02:52:02',1,'2025-08-15 02:52:02',''),(399636458484273152,'建筑物状态','building_status','0',1,'2025-08-15 02:52:37',1,'2025-08-15 02:52:37',''),(400559327108141056,'仪表盘类型','dashboard_type','0',1,'2025-08-17 15:59:46',1,'2025-08-17 15:59:46','仪表盘类型'),(401923807444471808,'预测模型','predict_model','0',1,'2025-08-21 10:21:43',1,'2025-08-21 10:21:43','预测模型'),(401925303229747200,'仪表盘时间间隔','dashboard_interval','0',1,'2025-08-21 10:27:40',1,'2025-08-21 10:27:52','仪表盘时间间隔'),(412430838978842624,'预警模型','prediction_model','0',1,'2025-09-19 10:12:55',1,'2025-09-19 10:12:55','预警模型'),(412968251866222592,'异常预警模型','prediction_exception_model','0',1,'2025-09-20 21:48:24',1,'2025-09-20 21:48:24',''),(423126262676983808,'设备类型','device_type','0',1,'2025-10-18 22:32:42',1,'2025-10-19 02:24:35','设备类型，电流电压采集器标注后会生成机器 稼动率采集配置'),(423185002306473984,'设备运行状态','device_run_status','0',1,'2025-10-19 02:26:07',1,'2025-10-19 02:26:07','设备运行状态'),(424102839531343872,'网关数据写入地址','grpc_sink_address','0',1,'2025-10-21 15:13:16',1,'2025-10-21 15:13:22','网关数据写入地址'),(433558833080569856,'服务器地址','server_host','0',1,'2025-11-16 17:28:01',1,'2025-11-16 17:28:01',''),(441645194777989120,'物料化验种类','laboratory_material','0',1,'2025-12-09 01:00:20',1,'2025-12-09 12:23:27',''),(448414420125618176,'配料化验种类','assaying_material','0',1,'2025-12-27 17:18:49',1,'2025-12-27 17:18:49','配料化验种类'),(451717254015356928,'资源业务场景','resource_file_category','0',1,'2026-01-05 20:03:06',1,'2026-01-05 20:03:06',''),(461487484686045184,'控制服务器','grpc_control_server_address','0',1,'2026-02-01 19:06:30',1,'2026-02-01 19:06:30',''),(467288730919505920,'控制算法类型','control_type','0',1,'2026-02-17 19:18:35',1,'2026-02-17 19:18:35',''),(467290806022049792,'控制方式','control_mode','0',1,'2026-02-17 19:26:50',1,'2026-02-17 19:26:50','');
/*!40000 ALTER TABLE `sys_dict_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_gateway_inbound_config`
--

DROP TABLE IF EXISTS `sys_gateway_inbound_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_gateway_inbound_config` (
  `gateway_config_id` bigint NOT NULL COMMENT '设备主键',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '设备名称',
  `protocol_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '协议类型',
  `config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '配置',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `status` int DEFAULT '0' COMMENT '操作状态（0正常 1异常）',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`gateway_config_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='设备协议管理';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_gateway_inbound_config`
--

LOCK TABLES `sys_gateway_inbound_config` WRITE;
/*!40000 ALTER TABLE `sys_gateway_inbound_config` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_gateway_inbound_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_iot_agent`
--

DROP TABLE IF EXISTS `sys_iot_agent`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_iot_agent` (
  `object_id` bigint unsigned NOT NULL COMMENT 'agent uuid',
  `name` varchar(64) NOT NULL COMMENT 'agent名字',
  `username` varchar(64) NOT NULL COMMENT 'username',
  `password` varchar(64) NOT NULL COMMENT 'password',
  `operate_state` int NOT NULL DEFAULT '0' COMMENT '操作状态 1-启动中 2-停止中 3-启动失败 4-停止失败',
  `operate_time` timestamp NULL DEFAULT NULL COMMENT '操作时间',
  `version` varchar(64) NOT NULL COMMENT 'agent版本',
  `config_uuid` varchar(64) NOT NULL DEFAULT '' COMMENT '配置版本',
  `ipv4` varchar(32) NOT NULL DEFAULT '' COMMENT 'ipv4地址',
  `ipv6` varchar(64) NOT NULL DEFAULT '' COMMENT 'ipv6地址',
  `last_heartbeat_time` timestamp NULL DEFAULT NULL COMMENT '上次心跳时间',
  `update_config_time` timestamp NULL DEFAULT NULL COMMENT '更新配置时间',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  `last_config_id` bigint unsigned NOT NULL COMMENT '最后版本的配置id',
  `config_id` bigint unsigned NOT NULL COMMENT '版本的配置id',
  PRIMARY KEY (`object_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='agent信息';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sys_iot_agent_config`
--

DROP TABLE IF EXISTS `sys_iot_agent_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_iot_agent_config` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增标识',
  `agent_object_id` bigint unsigned NOT NULL COMMENT 'agent id',
  `config_version` varchar(64) NOT NULL COMMENT '配置版本',
  `content` text NOT NULL COMMENT '配置内容',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`id`),
  KEY `agent_object_id` (`agent_object_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=462280058967429121 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='agent配置';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_iot_agent_process`
--

DROP TABLE IF EXISTS `sys_iot_agent_process`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_iot_agent_process` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增标识',
  `agent_object_id` bigint unsigned NOT NULL COMMENT 'agent uuid',
  `status` int NOT NULL COMMENT '状态 1-运行 2-停止',
  `name` varchar(64) NOT NULL COMMENT '名字',
  `version` varchar(64) NOT NULL COMMENT '版本',
  `start_time` timestamp NULL DEFAULT NULL COMMENT '启动时间',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `agent_object_id_2` (`agent_object_id`,`name`),
  KEY `agent_object_id` (`agent_object_id`,`name`,`version`,`state`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='agent进程信息';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_iot_agent_process`
--

LOCK TABLES `sys_iot_agent_process` WRITE;
/*!40000 ALTER TABLE `sys_iot_agent_process` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_iot_agent_process` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_iot_db_dev_map`
--

DROP TABLE IF EXISTS `sys_iot_db_dev_map`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_iot_db_dev_map` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增标识',
  `device_id` bigint NOT NULL COMMENT '设备id',
  `template_id` bigint NOT NULL COMMENT '模板id',
  `data_id` bigint NOT NULL COMMENT '测点id',
  `device` varchar(128) NOT NULL COMMENT '设备名字',
  `data_name` varchar(128) NOT NULL COMMENT '数据名字',
  `unit` varchar(128) NOT NULL COMMENT '数据单位',
  `dept_id` bigint DEFAULT '0' COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常-1删除）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_dev` (`device`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=459628233344684033 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='数据和iotdb 对照表';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_job`
--

DROP TABLE IF EXISTS `sys_job`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_job` (
  `job_id` bigint NOT NULL AUTO_INCREMENT COMMENT '任务ID',
  `job_name` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_unicode_ci NOT NULL DEFAULT '' COMMENT '任务名称',
  `job_params` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_unicode_ci DEFAULT NULL COMMENT '参数',
  `invoke_target` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_unicode_ci NOT NULL COMMENT '调用目标字符串',
  `cron_expression` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_unicode_ci DEFAULT '' COMMENT 'cron执行表达式',
  `status` char(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_unicode_ci DEFAULT '0' COMMENT '状态（0正常 1暂停）',
  `create_by` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_unicode_ci DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_unicode_ci DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`job_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=417654295689695233 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='定时任务调度表';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_job_log`
--

DROP TABLE IF EXISTS `sys_job_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_job_log` (
  `job_log_id` bigint NOT NULL AUTO_INCREMENT COMMENT '任务日志ID',
  `job_id` bigint NOT NULL COMMENT '任务ID',
  `job_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务名称',
  `job_group` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '任务组名',
  `invoke_target` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '调用目标字符串',
  `status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '0' COMMENT '执行状态（0正常 1失败）',
  `exception_info` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '异常信息',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `cost_time` bigint DEFAULT '0' COMMENT '耗时（毫秒）',
  PRIMARY KEY (`job_log_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='定时任务调度日志表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_job_log`
--

LOCK TABLES `sys_job_log` WRITE;
/*!40000 ALTER TABLE `sys_job_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_job_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_logininfor`
--

DROP TABLE IF EXISTS `sys_logininfor`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_logininfor` (
  `info_id` bigint NOT NULL AUTO_INCREMENT COMMENT '访问ID',
  `user_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '用户账号',
  `ipaddr` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '登录IP地址',
  `login_location` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '登录地点',
  `browser` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '浏览器类型',
  `os` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '操作系统',
  `status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '0' COMMENT '登录状态（0成功 1失败）',
  `msg` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '提示消息',
  `login_time` datetime DEFAULT NULL COMMENT '访问时间',
  PRIMARY KEY (`info_id`) USING BTREE,
  KEY `idx_sys_logininfor_s` (`status`) USING BTREE,
  KEY `idx_sys_logininfor_lt` (`login_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=469696468198887425 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='系统访问记录';
/*!40101 SET character_set_client = @saved_cs_client */;


DROP TABLE IF EXISTS `sys_material`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_material` (
  `material_id` bigint NOT NULL COMMENT '设备主键',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '物料名',
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '编码',
  `model` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '型号',
  `unit` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '单位',
  `factory` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '厂家',
  `address` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '存放地点',
  `price` decimal(10,2) DEFAULT '0.00' COMMENT '价格',
  `total` bigint NOT NULL COMMENT '总数',
  `outbound` bigint NOT NULL COMMENT '出库',
  `current_total` bigint NOT NULL COMMENT '当前库存',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '物料类型',
  PRIMARY KEY (`material_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='设备管理';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_material`
--

--
-- Table structure for table `sys_material_inbound`
--

DROP TABLE IF EXISTS `sys_material_inbound`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_material_inbound` (
  `inbound_id` bigint NOT NULL COMMENT '入库id',
  `material_id` bigint NOT NULL COMMENT '物料id',
  `number` bigint NOT NULL COMMENT '入库数量',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`inbound_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='物料入库管理';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_material_outbound`
--

DROP TABLE IF EXISTS `sys_material_outbound`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_material_outbound` (
  `outbound_id` bigint NOT NULL COMMENT '出库id',
  `material_id` bigint NOT NULL COMMENT '物料id',
  `number` bigint NOT NULL COMMENT '入库数量',
  `reason` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '出库理由',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`outbound_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='物料出库管理';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_material_outbound`
--

LOCK TABLES `sys_material_outbound` WRITE;
/*!40000 ALTER TABLE `sys_material_outbound` DISABLE KEYS */;
INSERT INTO `sys_material_outbound` VALUES (361044707347795968,361009835388440576,1,NULL,105,1,'2025-04-30 15:02:46',1,'2025-04-30 15:02:46',0),(361045342587719680,361009835388440576,1,'无出库原因',105,1,'2025-04-30 15:05:18',1,'2025-04-30 15:05:18',0),(361492515972452352,360819632518467584,1,'卖掉连了',105,1,'2025-05-01 20:42:12',1,'2025-05-01 20:42:12',0),(361498070875115520,361497126322049024,1,'坏了',105,1,'2025-05-01 21:04:16',1,'2025-05-01 21:04:16',0),(370923976907558912,361497126322049024,10,'测试出库',105,1,'2025-05-27 21:19:28',1,'2025-05-27 21:19:28',0);
/*!40000 ALTER TABLE `sys_material_outbound` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_modbus_device_config_data`
--

DROP TABLE IF EXISTS `sys_modbus_device_config_data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_modbus_device_config_data` (
  `device_config_id` bigint NOT NULL COMMENT '文档id',
  `template_id` bigint NOT NULL COMMENT '模板id',
  `name` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '数据名称',
  `type` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '数据类型',
  `slave` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '从设备地址',
  `register` bigint NOT NULL COMMENT '寄存器/偏移量',
  `storage_type` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '存储策略',
  `unit` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '单位',
  `precision` bigint NOT NULL COMMENT '数据精度',
  `function_code` tinyint(1) DEFAULT '0' COMMENT '功能码',
  `mode` tinyint(1) DEFAULT '0' COMMENT '功能码',
  `data_format` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '读写方式',
  `sort` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '数据排序',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  `data_type` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '设备节点类型，到底是开关还是数值',
  `agg_function` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '聚合函数，用来计算图表',
  `perturbation_variables` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '扰动变量',
  `graph_enable` tinyint(1) DEFAULT '0' COMMENT '是否开启图表',
  `predict_enable` tinyint(1) DEFAULT '0' COMMENT '是否开启趋势预测',
  `annotation` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '注解',
  `expression` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '转换表达式',
  `configuration_enable` tinyint(1) DEFAULT '0' COMMENT '组态属性',
  PRIMARY KEY (`device_config_id`) USING BTREE,
  KEY `template_id` (`template_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='modbus数据配置';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_notice`
--

DROP TABLE IF EXISTS `sys_notice`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_notice` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '公告ID',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '公告标题',
  `type` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '公告类型（1通知 2公告）',
  `txt` longblob COMMENT '公告内容',
  `dept_id` bigint DEFAULT NULL COMMENT '发件人所在部门',
  `dept_ids` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '发送部门',
  `create_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '创建者名称',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=418350850675576833 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='通知公告表';
/*!40101 SET character_set_client = @saved_cs_client */;



--
-- Table structure for table `sys_notice_user`
--

DROP TABLE IF EXISTS `sys_notice_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_notice_user` (
  `user_id` bigint NOT NULL,
  `notice_id` bigint NOT NULL,
  `status` char(1) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  PRIMARY KEY (`user_id`,`notice_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='消息通知关联用户表';
/*!40101 SET character_set_client = @saved_cs_client */;



--
-- Table structure for table `sys_oper_log`
--

DROP TABLE IF EXISTS `sys_oper_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_oper_log` (
  `oper_id` bigint NOT NULL AUTO_INCREMENT COMMENT '日志主键',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '模块标题',
  `business_type` int DEFAULT '0' COMMENT '业务类型（0其它 1新增 2修改 3删除）',
  `method` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '方法名称',
  `request_method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '请求方式',
  `oper_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '操作人员',
  `oper_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '请求URL',
  `oper_ip` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '主机地址',
  `oper_param` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '请求参数',
  `json_result` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '返回参数',
  `status` int DEFAULT '0' COMMENT '操作状态（0正常 1异常）',
  `user_id` bigint DEFAULT NULL COMMENT '用户ID',
  `oper_time` datetime DEFAULT NULL COMMENT '操作时间',
  `cost_time` bigint DEFAULT '0' COMMENT '消耗时间',
  PRIMARY KEY (`oper_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='操作日志记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_oper_log`
--

LOCK TABLES `sys_oper_log` WRITE;
/*!40000 ALTER TABLE `sys_oper_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_oper_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_permission`
--

DROP TABLE IF EXISTS `sys_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_permission` (
  `permission_id` bigint NOT NULL COMMENT '主键',
  `permission_name` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '权限名称',
  `parent_id` bigint DEFAULT NULL COMMENT '父ID',
  `permission` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '权限标识符',
  `sort` int DEFAULT NULL COMMENT '排序',
  `status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '0' COMMENT '状态',
  `create_by` bigint DEFAULT NULL COMMENT '创建人',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新人',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`permission_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='权限表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_permission`
--

LOCK TABLES `sys_permission` WRITE;
/*!40000 ALTER TABLE `sys_permission` DISABLE KEYS */;
INSERT INTO `sys_permission` VALUES (1,'系统管理',0,'system',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(2,'系统监控',0,'monitor',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(3,'系统工具',0,'tool',3,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(100,'用户管理',1,'system:user',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(101,'角色管理',1,'system:role',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(102,'权限管理',1,'system:permission',3,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(103,'部门管理',1,'system:dept',4,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(104,'岗位管理',1,'system:post',5,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(105,'字典管理',1,'system:dict',6,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(106,'参数设置',1,'system:config',7,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(107,'通知公告',1,'system:notice',8,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(108,'日志管理',1,'system:monitor',9,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(109,'在线用户',2,'monitor:online',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(110,'定时任务',2,'monitor:job',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(112,'服务监控',2,'monitor:server',3,'0',1,'2025-02-28 13:38:05',1,'2025-03-26 22:37:27'),(115,'表单构建',3,'tool:build',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(116,'代码生成',3,'tool:gen',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(117,'系统接口',3,'tool:swagger',3,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(500,'操作日志',108,'system:monitor:operlog',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(501,'登录日志',108,'system:monitor:logininfor',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1000,'用户查询',100,'system:user:query',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1001,'用户新增',100,'system:user:add',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1002,'用户修改',100,'system:user:edit',3,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1003,'用户删除',100,'system:user:remove',4,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1004,'用户导出',100,'system:user:export',5,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1005,'用户导入',100,'system:user:import',6,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1006,'重置密码',100,'system:user:resetPwd',7,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1007,'角色查询',101,'system:role:query',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1008,'角色新增',101,'system:role:add',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1009,'角色修改',101,'system:role:edit',3,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1010,'角色删除',101,'system:role:remove',4,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1011,'角色导出',101,'system:role:export',5,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1012,'菜单查询',102,'system:permission:query',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1013,'菜单新增',102,'system:permission:add',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1014,'菜单修改',102,'system:permission:edit',3,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1015,'菜单删除',102,'system:permission:remove',4,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1016,'部门查询',103,'system:dept:query',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1017,'部门新增',103,'system:dept:add',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1018,'部门修改',103,'system:dept:edit',3,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1019,'部门删除',103,'system:dept:remove',4,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1020,'岗位查询',104,'system:post:query',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1021,'岗位新增',104,'system:post:add',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1022,'岗位修改',104,'system:post:edit',3,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1023,'岗位删除',104,'system:post:remove',4,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1024,'岗位导出',104,'system:post:export',5,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1025,'字典查询',105,'system:dict:query',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1026,'字典新增',105,'system:dict:add',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1027,'字典修改',105,'system:dict:edit',3,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1028,'字典删除',105,'system:dict:remove',4,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1029,'字典导出',105,'system:dict:export',5,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1030,'参数查询',106,'system:config:query',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1031,'参数新增',106,'system:config:add',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1032,'参数修改',106,'system:config:edit',3,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1033,'参数删除',106,'system:config:remove',4,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1034,'参数导出',106,'system:config:export',5,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1035,'公告查询',107,'system:notice:query',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1036,'公告新增',107,'system:notice:add',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1037,'公告修改',107,'system:notice:edit',3,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1038,'公告删除',107,'system:notice:remove',4,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1039,'操作查询',500,'monitor:operlog:query',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1040,'操作删除',500,'monitor:operlog:remove',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1041,'日志导出',500,'monitor:operlog:export',3,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1042,'登录查询',501,'monitor:logininfor:query',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1043,'登录删除',501,'monitor:logininfor:remove',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1044,'日志导出',501,'monitor:logininfor:export',3,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1046,'在线查询',109,'monitor:online:query',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1047,'批量强退',109,'monitor:online:batchLogout',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1048,'单条强退',109,'monitor:online:forceLogout',3,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1049,'任务查询',110,'monitor:job:query',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1050,'任务新增',110,'monitor:job:add',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1051,'任务修改',110,'monitor:job:edit',3,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1052,'任务删除',110,'monitor:job:remove',4,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1053,'状态修改',110,'monitor:job:changeStatus',5,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1054,'任务导出',110,'monitor:job:export',6,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1055,'生成查询',116,'tool:gen:query',1,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1056,'生成修改',116,'tool:gen:edit',2,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1057,'生成删除',116,'tool:gen:remove',3,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1058,'导入代码',116,'tool:gen:import',4,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1059,'预览代码',116,'tool:gen:preview',5,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(1060,'生成代码',116,'tool:gen:code',6,'0',1,'2025-02-28 13:38:05',1,'2025-02-28 13:38:05'),(359684443914375168,'资产管理',0,'asset',4,'0',1,'2025-04-26 20:57:34',1,'2025-04-26 21:06:04'),(359686535840272384,'设备列表',359684443914375168,'asset:device',0,'0',1,'2025-04-26 21:05:53',1,'2025-04-26 21:05:53'),(359686725640916992,'添加设备',359686535840272384,'asset:device:set',0,'0',1,'2025-04-26 21:06:38',1,'2025-04-26 21:06:38'),(359900397353046016,'设备分组',359684443914375168,'asset:deviceGroup',0,'0',1,'2025-04-27 11:15:41',1,'2025-04-27 12:56:17'),(359926159607074816,'设置设备信息',359900397353046016,'asset:deviceGroup:set',0,'0',1,'2025-04-27 12:58:04',1,'2025-04-27 12:58:04'),(360335414705262592,'删除设备分组',359900397353046016,'asset:deviceGroup:remove',0,'0',1,'2025-04-28 16:04:18',1,'2025-04-28 16:04:18'),(360342197171654656,'删除设备',359686535840272384,'asset:device:remove',0,'0',1,'2025-04-28 16:31:15',1,'2025-04-28 16:31:15'),(360816738536787968,'物料管理',359684443914375168,'asset:material',0,'0',1,'2025-04-29 23:56:54',1,'2025-04-29 23:56:54'),(360816857738907648,'登记物料',360816738536787968,'asset:material:set',0,'0',1,'2025-04-29 23:57:23',1,'2025-04-29 23:57:23'),(360816923304267776,'删除物料',360816738536787968,'asset:material:remove',0,'0',1,'2025-04-29 23:57:38',1,'2025-04-29 23:57:38'),(361010500642803712,'物料入库列表',360816738536787968,'asset:material:inbound',0,'0',1,'2025-04-30 12:46:51',1,'2025-05-01 19:14:47'),(361010584994451456,'添加入库单',360816738536787968,'asset:material:inbound:add',0,'0',1,'2025-04-30 12:47:11',1,'2025-05-01 19:15:17'),(361044431924629504,'物料出库列表',360816738536787968,'asset:material:outbound',0,'0',1,'2025-04-30 15:01:40',1,'2025-05-01 19:15:06'),(361044530234920960,'添加出库单',360816738536787968,'asset:material:outbound:add',0,'0',1,'2025-04-30 15:02:04',1,'2025-05-01 19:14:54'),(361474631036571648,'工业智能体',0,'ai',5,'0',1,'2025-05-01 19:31:08',1,'2025-05-01 19:39:47'),(361475137612025856,'知识库',361474631036571648,'ai:dataset',0,'0',1,'2025-05-01 19:33:09',1,'2025-05-01 19:33:09'),(361475259188121600,'添加知识库',361475137612025856,'ai:dataset:create',0,'0',1,'2025-05-01 19:33:38',1,'2025-05-01 19:33:38'),(361787734412496896,'移除知识库',361475137612025856,'ai:dataset:remove',0,'0',1,'2025-05-02 16:15:18',1,'2025-05-02 16:15:18'),(361787808001560576,'更新知识库',361475137612025856,'ai:dataset:update',0,'0',1,'2025-05-02 16:15:35',1,'2025-05-02 16:15:35'),(362089395790024704,'上传文档',361475137612025856,'ai:dataset:upload:document',0,'0',1,'2025-05-03 12:13:59',1,'2025-05-03 12:13:59'),(362247033651728384,'更新知识库文档',361475137612025856,'ai:dataset:update:document',0,'0',1,'2025-05-03 22:40:23',1,'2025-05-03 22:40:23'),(362521217439109120,'删除知识库文档',361475137612025856,'ai:dataset:remove:document',0,'0',1,'2025-05-04 16:49:54',1,'2025-05-04 16:49:54'),(362576536529801216,'停止解析文档',361475137612025856,'ai:dataset:stop:document',0,'0',1,'2025-05-04 20:29:43',1,'2025-05-04 20:29:43'),(362576598731329536,'解析文档',361475137612025856,'ai:dataset:start:document',0,'0',1,'2025-05-04 20:29:57',1,'2025-05-04 20:29:57'),(362651361076056064,'chunk列表',361475137612025856,'ai:dataset:chunk',0,'0',1,'2025-05-05 01:27:02',1,'2025-05-05 01:27:02'),(362651427786461184,'添加chunk',361475137612025856,'ai:dataset:chunk:add',0,'0',1,'2025-05-05 01:27:18',1,'2025-05-05 01:27:18'),(362651499408396288,'移除chunk',361475137612025856,'ai:dataset:chunk:remove',0,'0',1,'2025-05-05 01:27:35',1,'2025-05-05 01:27:35'),(362651569637822464,'更新chunk',361475137612025856,'ai:dataset:chunk:update',0,'0',1,'2025-05-05 01:27:52',1,'2025-05-05 01:27:52'),(362651630966935552,'检索chunk',361475137612025856,'ai:dataset:chunk:retrieval',0,'0',1,'2025-05-05 01:28:07',1,'2025-05-05 01:28:07'),(363936074398961664,'生产工艺管理',0,'craft',0,'0',1,'2025-05-08 14:32:02',1,'2025-05-19 15:49:09'),(363936143873413120,'工艺管理',363936074398961664,'craft:route',0,'0',1,'2025-05-08 14:32:18',1,'2025-05-19 15:48:36'),(363936228376055808,'设置工艺',363936143873413120,'craft:route:set',0,'0',1,'2025-05-08 14:32:38',1,'2025-05-19 15:48:47'),(363936291261255680,'删除工艺',363936143873413120,'craft:route:remove',0,'0',1,'2025-05-08 14:32:53',1,'2025-05-19 15:48:59'),(364059932867170304,'AI聊天助理',361474631036571648,'ai:dataset:assistant',0,'0',1,'2025-05-08 22:44:12',1,'2025-05-09 22:09:50'),(364059987762221056,'创建聊天助理',364059932867170304,'ai:dataset:chats:assistant:create',0,'0',1,'2025-05-08 22:44:25',1,'2025-05-08 22:49:08'),(364413536425742336,'会话管理',361474631036571648,'ai:dataset:session',0,'0',1,'2025-05-09 22:09:18',1,'2025-05-09 22:09:18'),(364413837832622080,'使用助理创建会话',364413536425742336,'ai:dataset:session:create',0,'0',1,'2025-05-09 22:10:29',1,'2025-05-09 22:10:29'),(364421764240904192,'使用助理更新会话',364413536425742336,'ai:dataset:session:update',0,'0',1,'2025-05-09 22:41:59',1,'2025-05-09 22:41:59'),(364433815591981056,'会话列表',364413536425742336,'ai:dataset:session:list',0,'0',1,'2025-05-09 23:29:52',1,'2025-05-09 23:29:52'),(364442854312906752,'删除助理会话',364413536425742336,'ai:dataset:session:remove',0,'0',1,'2025-05-10 00:05:47',1,'2025-05-10 00:05:47'),(364652283557842944,'与助手聊天',364413536425742336,'ai:dataset:session:charts:completions',0,'0',1,'2025-05-10 13:57:59',1,'2025-05-10 13:57:59'),(365685718661468160,'工序管理',363936074398961664,'craft:process',0,'0',1,'2025-05-13 10:24:29',1,'2025-05-19 15:46:44'),(365685913927290880,'设置工序',365685718661468160,'craft:process:set',0,'0',1,'2025-05-13 10:25:16',1,'2025-05-19 15:46:55'),(365685977705877504,'移除工序',365685718661468160,'craft:process:remove',0,'0',1,'2025-05-13 10:25:31',1,'2025-05-19 15:47:14'),(366056868528787456,'工序内容管理',363936074398961664,'craft:process:context',0,'0',1,'2025-05-14 10:59:18',1,'2025-05-19 15:46:05'),(366056957984903168,'设置工序内容',366056868528787456,'craft:process:context:set',0,'0',1,'2025-05-14 10:59:40',1,'2025-05-19 15:45:57'),(366057043401904128,'移除工序内容',366056868528787456,'craft:process:context:remove',0,'0',1,'2025-05-14 11:00:00',1,'2025-05-19 15:46:28'),(366230539859922944,'移除助理',364059932867170304,'ai:dataset:chats:assistant:remove',0,'0',1,'2025-05-14 22:29:25',1,'2025-05-14 22:29:25'),(366568187384303616,'更新聊天助理',364059932867170304,'ai:dataset:assistant:update',0,'0',1,'2025-05-15 20:51:06',1,'2025-05-15 20:51:06'),(366777330686758912,'读取知识库',361475137612025856,'ai:dataset:info',0,'0',1,'2025-05-16 10:42:10',1,'2025-05-16 10:42:10'),(366863885531090944,'组成工序管理',363936074398961664,'craft:route:process',0,'0',1,'2025-05-16 16:26:06',1,'2025-05-19 15:47:27'),(366863942808506368,'设置组成工序内容',366863885531090944,'craft:route:process:set',0,'0',1,'2025-05-16 16:26:20',1,'2025-05-19 15:48:08'),(366863990531297280,'移除组成工序内容',366863885531090944,'craft:route:process:remove',0,'0',1,'2025-05-16 16:26:31',1,'2025-05-19 15:48:18'),(367663286088372224,'Agent管理',361474631036571648,'ai:dataset:agents:list',0,'0',1,'2025-05-18 21:22:38',1,'2025-05-18 21:22:38'),(367664829529329664,'创建agent会话',367663286088372224,'ai:dataset:session:agents:create',0,'0',1,'2025-05-18 21:28:46',1,'2025-05-18 21:28:46'),(367664941072650240,'和Agent聊天',367663286088372224,'ai:dataset:session:agents:completions',0,'0',1,'2025-05-18 21:29:13',1,'2025-05-18 21:29:13'),(367665204437192704,'Agent会话列表',367663286088372224,'ai:dataset:session:agents:list',0,'0',1,'2025-05-18 21:30:16',1,'2025-05-18 21:30:16'),(367665320741048320,'删除agent会话',367663286088372224,'ai:dataset:session:agents:delete',0,'0',1,'2025-05-18 21:30:43',1,'2025-05-18 21:30:43'),(367665410381713408,'相关提问',364413536425742336,'ai:dataset:session:conversation:related_questions',0,'0',1,'2025-05-18 21:31:05',1,'2025-05-18 21:31:05'),(367940361919664128,'产品制程管理',363936074398961664,'craft:route:product',0,'0',1,'2025-05-19 15:43:38',1,'2025-05-19 15:44:00'),(367940533215039488,'设置产品制程',367940361919664128,'craft:route:product:set',0,'0',1,'2025-05-19 15:44:19',1,'2025-05-19 15:44:19'),(367940616467779584,'删除产品制程',367940361919664128,'craft:route:product:remove',0,'0',1,'2025-05-19 15:44:39',1,'2025-05-19 15:44:39'),(370897941847609344,'智能回答',364413536425742336,'ai:dataset:session:ask',0,'0',1,'2025-05-27 19:36:00',1,'2025-05-27 19:36:00'),(373343224594436096,'设备模板管理',359684443914375168,'asset:device:template',0,'0',1,'2025-06-03 13:32:41',1,'2025-06-03 13:32:41'),(373343302742708224,'获取模板列表',373343224594436096,'asset:device:template:list',0,'0',1,'2025-06-03 13:33:00',1,'2025-06-03 13:33:00'),(373343417628889088,'设置模板',373343224594436096,'asset:device:template:set',0,'0',1,'2025-06-03 13:33:27',1,'2025-06-03 13:33:27'),(373343563192209408,'删除设备模板',373343224594436096,'asset:device:template:remove',0,'0',1,'2025-06-03 13:34:02',1,'2025-06-03 13:34:02'),(373374786467794944,'模板协议列表',373343224594436096,'asset:device:template:protocols',0,'0',1,'2025-06-03 15:38:06',1,'2025-06-03 15:39:36'),(373706277987028992,'模板数据列表',359686535840272384,'asset:device:template:data:list',0,'0',1,'2025-06-04 13:35:20',1,'2025-06-04 13:35:20'),(373706373285810176,'设置模板数据',359686535840272384,'asset:device:template:data:set',0,'0',1,'2025-06-04 13:35:43',1,'2025-06-04 13:35:43'),(373706447705346048,'移除模板数据',359686535840272384,'asset:device:template:data:remove',0,'0',1,'2025-06-04 13:36:00',1,'2025-06-04 13:36:00'),(374788609506545664,'网关管理',0,'daemonize',0,'0',1,'2025-06-07 13:16:08',1,'2025-06-07 13:16:08'),(374788724715687936,'agent管理',374788609506545664,'daemonize:agent',0,'0',1,'2025-06-07 13:16:35',1,'2025-06-07 13:16:49'),(374788976810135552,'agent列表',374788724715687936,'daemonize:agent:list',0,'0',1,'2025-06-07 13:17:35',1,'2025-06-07 13:17:35'),(374789035895296000,'设置agent',374788724715687936,'daemonize:agent:set',0,'0',1,'2025-06-07 13:17:49',1,'2025-06-07 13:17:49'),(374789081369939968,'移除agent',374788724715687936,'daemonize:agent:remove',0,'0',1,'2025-06-07 13:18:00',1,'2025-06-07 13:18:00'),(375148581360766976,'启动进程',374788724715687936,'daemonize:agent:process:start',0,'0',1,'2025-06-08 13:06:32',1,'2025-06-08 13:06:32'),(375148631210070016,'停止进程',374788724715687936,'daemonize:agent:process:stop',0,'0',1,'2025-06-08 13:06:44',1,'2025-06-08 13:06:44'),(375973793090244608,'生成agent配置',374788724715687936,'gateway:agent:config:generate',0,'0',1,'2025-06-10 19:45:38',1,'2025-06-10 19:45:38'),(376624128414715904,'agent详情',374788724715687936,'daemonize:agent:info',0,'0',1,'2025-06-12 14:49:50',1,'2025-06-12 14:49:50'),(376627010648150016,'agent配置列表',374788724715687936,'gateway:agent:config:list',0,'0',1,'2025-06-12 15:01:17',1,'2025-06-12 15:01:17'),(376627074644840448,'保存agent配置',374788724715687936,'gateway:agent:config:set',0,'0',1,'2025-06-12 15:01:32',1,'2025-06-12 15:01:32'),(377388132074524672,'设备监控',0,'device',0,'0',1,'2025-06-14 17:25:42',1,'2025-06-14 17:26:09'),(377388341726810112,'设备监控',377388132074524672,'device:monitor',0,'0',1,'2025-06-14 17:26:32',1,'2025-06-14 17:26:32'),(377388458961801216,'设备监控',377388341726810112,'device:monitor:list',0,'0',1,'2025-06-14 17:27:00',1,'2025-06-14 17:27:00'),(377678742174044160,'设备指标',377388341726810112,'device:monitor:metric',0,'0',1,'2025-06-15 12:40:29',1,'2025-06-15 12:40:29'),(379927109822320640,'工艺路线详情',363936143873413120,'craft:route:detail',0,'0',1,'2025-06-21 17:34:42',1,'2025-06-21 17:34:42'),(382423962505711616,'生产工单管理',363936074398961664,'craft:route:work_order',0,'0',1,'2025-06-28 14:56:18',1,'2025-06-28 14:56:18'),(382424273555296256,'生产工单列表',382423962505711616,'craft:route:work_order:list',0,'0',1,'2025-06-28 14:57:32',1,'2025-06-28 14:57:32'),(382424371806867456,'保存生产工单',382423962505711616,'craft:route:work_order:set',0,'0',1,'2025-06-28 14:57:55',1,'2025-06-28 14:57:55'),(382424450659782656,'移除生产工单',382423962505711616,'craft:route:work_order:remove',0,'0',1,'2025-06-28 14:58:14',1,'2025-06-28 14:58:14'),(383853950169780224,'保存工艺制图',363936143873413120,'craft:route:config:save',0,'0',1,'2025-07-02 13:38:33',1,'2025-07-02 13:38:33'),(387615595530555392,'设备数据列表',388601769464172544,'device:monitor:data:dev:list',0,'0',1,'2025-07-12 22:46:00',1,'2025-07-15 16:08:17'),(387615666456236032,'实时数据',388601769464172544,'device:monitor:data:list',0,'0',1,'2025-07-12 22:46:16',1,'2025-07-15 16:05:43'),(387813840298971136,'实时数据导出',388601769464172544,'device:monitor:data:export',0,'0',1,'2025-07-13 11:53:45',1,'2025-07-15 16:05:24'),(388601769464172544,'设备数据',377388132074524672,'device:monitor:data',0,'0',1,'2025-07-15 16:04:42',1,'2025-07-15 16:04:42'),(389399690874982400,'告警管理',0,'alert',0,'0',1,'2025-07-17 20:55:21',1,'2025-07-17 20:55:55'),(389399774467461120,'告警模板',389399690874982400,'alert:template',0,'0',1,'2025-07-17 20:55:41',1,'2025-07-17 20:56:11'),(389399983243137024,'告警模板列表',389399774467461120,'alert:template:list',0,'0',1,'2025-07-17 20:56:31',1,'2025-07-17 20:56:31'),(389400076717395968,'设置告警模板',389399774467461120,'alert:template:set',0,'0',1,'2025-07-17 20:56:53',1,'2025-07-17 20:56:53'),(389400378266882048,'告警规则',389399690874982400,'alert:rule',0,'0',1,'2025-07-17 20:58:05',1,'2025-07-17 20:58:05'),(389400707612020736,'告警规则列表',389400378266882048,'alert:rule:list',0,'0',1,'2025-07-17 20:59:23',1,'2025-07-17 20:59:23'),(389400784246149120,'设置告警规则',389400378266882048,'alert:rule:set',0,'0',1,'2025-07-17 20:59:42',1,'2025-07-17 20:59:42'),(389400891922321408,'删除告警规则',389400378266882048,'alert:rule:remove',0,'0',1,'2025-07-17 21:00:07',1,'2025-07-17 21:00:07'),(390837150549020672,'调度管理',363936074398961664,'craft:route:schedule',0,'0',1,'2025-07-21 20:07:18',1,'2025-07-21 20:07:18'),(390837560819060736,'当月调度列表',390837150549020672,'craft:route:product:month:list',0,'0',1,'2025-07-21 20:08:56',1,'2025-07-21 20:08:56'),(391066329458675712,'告警数据管理',389399690874982400,'alert:log',0,'0',1,'2025-07-22 11:17:59',1,'2025-07-22 11:18:05'),(391066480818524160,'告警列表',391066329458675712,'alert:log:list',0,'0',1,'2025-07-22 11:18:35',1,'2025-07-22 11:18:35'),(392296367583662080,'保存调度任务',390837150549020672,'craft:route:schedule:set',0,'0',1,'2025-07-25 20:45:43',1,'2025-07-25 20:45:43'),(393428243878776832,'调度列表',390837150549020672,'craft:route:schedule:list',0,'0',1,'2025-07-28 23:43:23',1,'2025-07-28 23:43:23'),(393433880197074944,'移除告警模板',389399774467461120,'alert:template:remove',0,'0',1,'2025-07-29 00:05:47',1,'2025-07-29 00:05:47'),(393737216251138048,'调度详情',390837150549020672,'craft:route:schedule:detail',0,'0',1,'2025-07-29 20:11:08',1,'2025-07-29 20:11:08'),(394512351845421056,'告警处理',389399690874982400,'alert:action',0,'0',1,'2025-07-31 23:31:14',1,'2025-07-31 23:31:14'),(394512431788855296,'告警处理列表',394512351845421056,'alert:action:list',0,'0',1,'2025-07-31 23:31:33',1,'2025-07-31 23:31:33'),(394512539913818112,'保存处理规则',394512351845421056,'alert:action:set',0,'0',1,'2025-07-31 23:31:59',1,'2025-07-31 23:31:59'),(394512644633006080,'移除告警处理规则',394512351845421056,'alert:action:remove',0,'0',1,'2025-07-31 23:32:24',1,'2025-07-31 23:32:24'),(395211005040267264,'告警推理',389399690874982400,'alert:reason',0,'0',1,'2025-08-02 21:47:26',1,'2025-08-02 21:47:26'),(395444252689043456,'告警AI推理列表',395211005040267264,'alert:reason:list',0,'0',1,'2025-08-03 13:14:17',1,'2025-08-03 13:14:17'),(395444328442368000,'设置告警AI推理发送配置',395211005040267264,'alert:reason:set',0,'0',1,'2025-08-03 13:14:35',1,'2025-08-03 13:14:35'),(395444655237369856,'删除告警AI推理发送配置',395211005040267264,'alert:reason:remove',0,'0',1,'2025-08-03 13:15:53',1,'2025-08-03 13:15:53'),(399261236966985728,'建筑管理',359684443914375168,'asset:building',0,'0',1,'2025-08-14 02:01:37',1,'2025-08-14 02:01:37'),(399634680992763904,'建筑物列表',399261236966985728,'asset:building:list',0,'0',1,'2025-08-15 02:45:33',1,'2025-08-15 02:45:33'),(399634796147380224,'保存建筑物',399261236966985728,'asset:building:set',0,'0',1,'2025-08-15 02:46:00',1,'2025-08-15 02:46:00'),(399634921259274240,'移除建筑物',399261236966985728,'asset:building:remove',0,'0',1,'2025-08-15 02:46:30',1,'2025-08-15 02:46:30'),(400238451561074688,'仪表盘',0,'dashboard',8,'0',1,'2025-08-16 18:44:43',1,'2025-08-16 20:41:45'),(400255117632212992,'仪表盘管理',400238451561074688,'dashboard:manager',0,'0',1,'2025-08-16 19:50:56',1,'2025-08-16 21:37:06'),(400266689276547072,'仪表盘列表',400255117632212992,'dashboard:manager:list',0,'0',1,'2025-08-16 20:36:55',1,'2025-08-16 21:37:19'),(400299197661712384,'移除调度任务',390837150549020672,'craft:route:schedule:remove',0,'0',1,'2025-08-16 22:46:06',1,'2025-08-16 22:53:33'),(400555685768597504,'保存仪表盘',400255117632212992,'dashboard:manager:set',0,'0',1,'2025-08-17 15:45:17',1,'2025-08-17 15:45:17'),(400555947891625984,'移除仪表盘',400255117632212992,'dashboard:manager:remove',0,'0',1,'2025-08-17 15:46:20',1,'2025-08-17 15:46:20'),(400646833996566528,'仪表盘数据管理',400238451561074688,'dashboard:data',0,'0',1,'2025-08-17 21:47:29',1,'2025-08-17 21:47:29'),(400646920785104896,'仪表盘数据详情',400646833996566528,'dashboard:data:detail',0,'0',1,'2025-08-17 21:47:50',1,'2025-08-17 21:47:50'),(402100740207677440,'时序指标列表',377388341726810112,'metric:time:seq:list',0,'0',1,'2025-08-21 22:04:47',1,'2025-08-21 22:12:04'),(402182963023843328,'查询面板',400646833996566528,'dashboard:manager:metric:query',0,'0',1,'2025-08-22 03:31:31',1,'2025-08-22 03:31:31'),(402903819106652160,'保存面板',400646833996566528,'dashboard:data:set',0,'0',1,'2025-08-24 03:15:56',1,'2025-08-24 03:15:56'),(402903898223808512,'移除面板',400646833996566528,'dashboard:data:remove',0,'0',1,'2025-08-24 03:16:15',1,'2025-08-24 03:16:15'),(402921721436311552,'读取面板',400646833996566528,'dashboard:data:info',0,'0',1,'2025-08-24 04:27:04',1,'2025-08-24 04:27:04'),(412105593227055104,'智能预警',361474631036571648,'ai:prediction',0,'0',1,'2025-09-18 12:40:30',1,'2025-09-18 12:40:30'),(412106176428249088,'异常检测',361474631036571648,'ai:exception',0,'0',1,'2025-09-18 12:42:49',1,'2025-09-18 12:42:49'),(412563092380061696,'智能预警列表',412105593227055104,'ai:prediction:list',0,'0',1,'2025-09-19 18:58:26',1,'2025-09-19 18:58:26'),(412563179353149440,'设置智能预警',412105593227055104,'ai:prediction:set',0,'0',1,'2025-09-19 18:58:47',1,'2025-09-19 18:58:47'),(412563252694749184,'删除智能预警',412105593227055104,'ai:prediction:remove',0,'0',1,'2025-09-19 18:59:04',1,'2025-09-19 22:34:47'),(412960189709291520,'异常预警列表',412106176428249088,'ai:exception:list',0,'0',1,'2025-09-20 21:16:22',1,'2025-09-20 21:16:22'),(412960265802354688,'设置异常预警',412106176428249088,'ai:exception:set',0,'0',1,'2025-09-20 21:16:40',1,'2025-09-20 21:16:40'),(412960335926923264,'删除异常预警',412106176428249088,'ai:exception:remove',0,'0',1,'2025-09-20 21:16:57',1,'2025-09-20 21:16:57'),(413276830359883776,'预测控制',361474631036571648,'ai:control',0,'0',1,'2025-09-21 18:14:35',1,'2025-09-21 18:14:35'),(413276929043468288,'预测控制列表',413276830359883776,'ai:control:list',0,'0',1,'2025-09-21 18:14:58',1,'2025-09-21 18:14:58'),(413276999113510912,'设置智能控制',413276830359883776,'ai:control:set',0,'0',1,'2025-09-21 18:15:15',1,'2025-09-21 18:15:15'),(413277070659948544,'删除预测控制',413276830359883776,'ai:control:remove',0,'0',1,'2025-09-21 18:15:32',1,'2025-09-21 18:15:32'),(423051602132209664,'电流参数配置',1,'system:electric',0,'0',1,'2025-10-18 17:36:02',1,'2025-10-18 17:36:02'),(423052338089955328,'设备电流配置列表',423051602132209664,'system:electric:set',0,'0',1,'2025-10-18 17:38:57',1,'2025-10-18 17:40:35'),(423052435095818240,'删除设备电流配置',423051602132209664,'system:electric:remove',0,'0',1,'2025-10-18 17:39:20',1,'2025-10-18 17:39:20'),(423052699966115840,'设备电流配置列表',423051602132209664,'system:electric:list',0,'0',1,'2025-10-18 17:40:24',1,'2025-10-18 17:40:49'),(424116834506117120,'通过设备查询下面的模板数据',359686535840272384,'asset:device:template:data:device',0,'0',1,'2025-10-21 16:08:53',1,'2025-10-21 16:08:53'),(424253010877616128,'设备测点',377388132074524672,'metric:time:seq',0,'0',1,'2025-10-22 01:10:00',1,'2025-10-22 01:10:00'),(424253067232284672,'设置测点',424253010877616128,'metric:time:seq:set',0,'0',1,'2025-10-22 01:10:13',1,'2025-10-22 01:10:13'),(424253139395284992,'移除测点',424253010877616128,'metric:time:seq:remove',0,'0',1,'2025-10-22 01:10:31',1,'2025-10-22 01:10:31'),(424549398710587392,'稼动率统计',377388341726810112,'device:monitor:utilization:stat',0,'0',1,'2025-10-22 20:47:44',1,'2025-10-22 20:47:44'),(424602568920928256,'班次管理',1,'system:shift',0,'0',1,'2025-10-23 00:19:01',1,'2025-10-23 00:19:01'),(424602649103437824,'班次管理列表',424602568920928256,'system:shift:list',0,'0',1,'2025-10-23 00:19:20',1,'2025-10-23 00:19:20'),(424602778455773184,'设置班次配置',424602568920928256,'system:shift:set',0,'0',1,'2025-10-23 00:19:51',1,'2025-10-23 00:19:51'),(424602843849166848,'移除班次配置',424602568920928256,'system:shift:remove',0,'0',1,'2025-10-23 00:20:07',1,'2025-10-23 00:20:07'),(441451927406907392,'生产管理',0,'product',0,'0',1,'2025-12-08 12:12:21',1,'2025-12-08 12:12:21'),(441452001889357824,'物料化验单',441451927406907392,'product:laboratory',0,'0',1,'2025-12-08 12:12:39',1,'2025-12-10 16:06:20'),(441452113885663232,'化验单列表',441452001889357824,'product:laboratory:list',0,'0',1,'2025-12-08 12:13:06',1,'2025-12-08 12:13:06'),(441452189274083328,'化验单信息',441452001889357824,'product:laboratory:info',0,'0',1,'2025-12-08 12:13:24',1,'2025-12-08 12:13:24'),(441452285436891136,'设置化验单',441452001889357824,'product:laboratory:set',0,'0',1,'2025-12-08 12:13:47',1,'2025-12-08 12:13:47'),(441452343255371776,'删除化验单',441452001889357824,'product:laboratory:remove',0,'0',1,'2025-12-08 12:14:00',1,'2025-12-08 12:14:00'),(441452404488015872,'用户化验单列表',441452001889357824,'product:laboratory:our:list',0,'0',1,'2025-12-08 12:14:15',1,'2025-12-08 12:14:15'),(451695491684503552,'资料管理',359684443914375168,'asset:resource',0,'0',1,'2026-01-05 18:36:37',1,'2026-01-05 18:36:37'),(451695616918032384,'资料列表',451695491684503552,'asset:resource:list',0,'0',1,'2026-01-05 18:37:07',1,'2026-01-05 18:37:07'),(451695685792698368,'设置资料',451695491684503552,'asset:resource:set',0,'0',1,'2026-01-05 18:37:24',1,'2026-01-05 18:37:24'),(451695797944193024,'移除资料',451695491684503552,'asset:resource:remove',0,'0',1,'2026-01-05 18:37:50',1,'2026-01-05 18:37:50'),(455744359673892864,'首页',0,'home',0,'0',1,'2026-01-16 22:45:23',1,'2026-01-16 22:45:23'),(455744483879817216,'首页统计',455744359673892864,'home:stats',0,'0',1,'2026-01-16 22:45:52',1,'2026-01-16 22:45:52'),(459178189533483008,'移除agent配置',374788724715687936,'gateway:agent:config:remove',0,'0',1,'2026-01-26 10:10:12',1,'2026-01-26 10:10:12'),(463342988819435520,'设备数据预测',377388341726810112,'device:monitor:predict',0,'0',1,'2026-02-06 21:59:37',1,'2026-02-06 21:59:37'),(467225738282536960,'组态管理',0,'configuration',0,'0',1,'2026-02-17 15:08:17',1,'2026-02-17 15:08:17'),(467226020349480960,'组态管理',467225738282536960,'configuration:manager',0,'0',1,'2026-02-17 15:09:24',1,'2026-02-17 15:09:24'),(467226264915152896,'组态列表',467226020349480960,'configuration:manager:list',0,'0',1,'2026-02-17 15:10:22',1,'2026-02-17 15:10:22'),(467226346293039104,'设置组态',467226020349480960,'configuration:manager:set',0,'0',1,'2026-02-17 15:10:42',1,'2026-02-17 15:10:42'),(467226738544349184,'删除组态',467226020349480960,'configuration:manager:remove',0,'0',1,'2026-02-17 15:12:15',1,'2026-02-17 15:12:15'),(469432251474513920,'设备控制',377388341726810112,'device:monitor:control',0,'0',1,'2026-02-23 17:16:10',1,'2026-02-23 18:31:31'),(469432334664339456,'设备控制日志',377388341726810112,'device:monitor:control:log:list',0,'0',1,'2026-02-23 17:16:30',1,'2026-02-23 18:31:38');
/*!40000 ALTER TABLE `sys_permission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_post`
--

DROP TABLE IF EXISTS `sys_post`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_post` (
  `post_id` bigint NOT NULL AUTO_INCREMENT COMMENT '岗位ID',
  `post_code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '岗位编码',
  `post_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '岗位名称',
  `post_sort` int NOT NULL COMMENT '显示顺序',
  `status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '状态（0正常 1停用）',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`post_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='岗位信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_post`
--

LOCK TABLES `sys_post` WRITE;
/*!40000 ALTER TABLE `sys_post` DISABLE KEYS */;
INSERT INTO `sys_post` VALUES (1,'ceo','董事长',1,'0',1,'2024-02-08 04:10:55',1,NULL,''),(2,'se','项目经理',2,'0',1,'2024-02-08 04:10:55',1,NULL,''),(3,'hr','人力资源',3,'0',1,'2024-02-08 04:10:55',1,NULL,''),(4,'user','普通员工',4,'0',1,'2024-02-08 04:10:55',1,'2024-02-25 10:40:41','');
/*!40000 ALTER TABLE `sys_post` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_pro_process`
--

DROP TABLE IF EXISTS `sys_pro_process`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_pro_process` (
  `process_id` bigint NOT NULL AUTO_INCREMENT COMMENT '工序ID',
  `process_code` varchar(64) NOT NULL COMMENT '工序编码',
  `process_name` varchar(255) NOT NULL COMMENT '工序名称',
  `attention` varchar(1000) DEFAULT NULL COMMENT '工艺要求',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `attr1` varchar(64) DEFAULT NULL COMMENT '预留字段1',
  `attr2` varchar(255) DEFAULT NULL COMMENT '预留字段2',
  `attr3` int DEFAULT '0' COMMENT '预留字段3',
  `attr4` int DEFAULT '0' COMMENT '预留字段4',
  `status` tinyint(1) DEFAULT '0' COMMENT '是否启用（0禁用 1启用）',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`process_id`)
) ENGINE=InnoDB AUTO_INCREMENT=467287660776394753 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='生产工序表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sys_pro_process_content`
--

DROP TABLE IF EXISTS `sys_pro_process_content`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_pro_process_content` (
  `content_id` bigint NOT NULL AUTO_INCREMENT COMMENT '内容ID',
  `process_id` bigint NOT NULL COMMENT '工序ID',
  `order_num` int DEFAULT '0' COMMENT '顺序编号',
  `content_text` varchar(500) DEFAULT NULL COMMENT '内容说明',
  `device` varchar(255) DEFAULT NULL COMMENT '辅助设备',
  `material` varchar(255) DEFAULT NULL COMMENT '辅助材料',
  `doc_url` varchar(255) DEFAULT NULL COMMENT '材料URL',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `attr1` varchar(64) DEFAULT NULL COMMENT '预留字段1',
  `attr2` varchar(255) DEFAULT NULL COMMENT '预留字段2',
  `attr3` int DEFAULT '0' COMMENT '预留字段3',
  `attr4` int DEFAULT '0' COMMENT '预留字段4',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `extension` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '扩展字段',
  `control_type` varchar(255) DEFAULT NULL COMMENT '控制方式',
  PRIMARY KEY (`content_id`)
) ENGINE=InnoDB AUTO_INCREMENT=430358804064899073 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='生产工序内容表';
/*!40101 SET character_set_client = @saved_cs_client */;


DROP TABLE IF EXISTS `sys_pro_route_process`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_pro_route_process` (
  `record_id` bigint NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `route_id` bigint NOT NULL COMMENT '工艺路线ID',
  `process_id` bigint NOT NULL COMMENT '工序ID',
  `process_code` varchar(64) DEFAULT NULL COMMENT '工序编码',
  `process_name` varchar(255) DEFAULT NULL COMMENT '工序名称',
  `order_num` int DEFAULT '1' COMMENT '序号',
  `next_process_id` bigint NOT NULL COMMENT '工序ID',
  `next_process_code` varchar(64) DEFAULT NULL COMMENT '工序编码',
  `next_process_name` varchar(255) DEFAULT NULL COMMENT '工序名称',
  `link_type` varchar(64) DEFAULT 'SS' COMMENT '与下一道工序关系',
  `default_pre_time` int DEFAULT '0' COMMENT '准备时间',
  `default_suf_time` int DEFAULT '0' COMMENT '等待时间',
  `color_code` char(7) DEFAULT '#00AEF3' COMMENT '甘特图显示颜色',
  `key_flag` varchar(64) DEFAULT 'N' COMMENT '关键工序',
  `is_check` char(1) DEFAULT 'N' COMMENT '是否检验',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `attr1` varchar(64) DEFAULT NULL COMMENT '预留字段1',
  `attr2` varchar(255) DEFAULT NULL COMMENT '预留字段2',
  `attr3` int DEFAULT '0' COMMENT '预留字段3',
  `attr4` int DEFAULT '0' COMMENT '预留字段4',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`record_id`)
) ENGINE=InnoDB AUTO_INCREMENT=380554321055453185 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='工艺组成表';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_pro_route_product`
--

DROP TABLE IF EXISTS `sys_pro_route_product`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_pro_route_product` (
  `record_id` bigint NOT NULL AUTO_INCREMENT COMMENT '记录ID',
  `route_id` bigint NOT NULL COMMENT '工艺路线ID',
  `item_id` bigint NOT NULL COMMENT '产品物料ID',
  `item_code` varchar(64) NOT NULL COMMENT '产品物料编码',
  `item_name` varchar(255) NOT NULL COMMENT '产品物料名称',
  `specification` varchar(500) DEFAULT NULL COMMENT '规格型号',
  `unit_of_measure` varchar(64) NOT NULL COMMENT '单位',
  `unit_name` varchar(64) DEFAULT NULL COMMENT '单位名称',
  `quantity` int DEFAULT '1' COMMENT '生产数量',
  `production_time` double(12,2) DEFAULT '1.00' COMMENT '生产用时',
  `time_unit_type` varchar(64) DEFAULT 'MINUTE' COMMENT '时间单位',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `attr1` varchar(64) DEFAULT NULL COMMENT '预留字段1',
  `attr2` varchar(255) DEFAULT NULL COMMENT '预留字段2',
  `attr3` int DEFAULT '0' COMMENT '预留字段3',
  `attr4` int DEFAULT '0' COMMENT '预留字段4',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`record_id`)
) ENGINE=InnoDB AUTO_INCREMENT=367942272861343745 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='产品制程';
/*!40101 SET character_set_client = @saved_cs_client */;



--
-- Table structure for table `sys_pro_task`
--

DROP TABLE IF EXISTS `sys_pro_task`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_pro_task` (
  `task_id` bigint NOT NULL AUTO_INCREMENT COMMENT '任务ID',
  `task_code` varchar(64) NOT NULL COMMENT '任务编号',
  `task_name` varchar(255) NOT NULL COMMENT '任务名称',
  `workorder_id` bigint NOT NULL COMMENT '生产工单ID',
  `workorder_code` varchar(64) NOT NULL COMMENT '生产工单编号',
  `workorder_name` varchar(255) NOT NULL COMMENT '工单名称',
  `workstation_id` bigint NOT NULL COMMENT '工作站ID',
  `gateway_id` bigint NOT NULL COMMENT '网关id',
  `workstation_code` varchar(64) NOT NULL COMMENT '工作站编号',
  `workstation_name` varchar(255) NOT NULL COMMENT '工作站名称',
  `route_id` bigint NOT NULL COMMENT '工艺ID',
  `route_code` varchar(64) DEFAULT NULL COMMENT '工艺编号',
  `process_id` bigint NOT NULL COMMENT '工序ID',
  `process_code` varchar(64) DEFAULT NULL COMMENT '工序编码',
  `process_name` varchar(255) DEFAULT NULL COMMENT '工序名称',
  `item_id` bigint NOT NULL COMMENT '产品物料ID',
  `item_code` varchar(64) NOT NULL COMMENT '产品物料编码',
  `item_name` varchar(255) NOT NULL COMMENT '产品物料名称',
  `specification` varchar(500) DEFAULT NULL COMMENT '规格型号',
  `unit_of_measure` varchar(64) NOT NULL COMMENT '单位',
  `unit_name` varchar(64) DEFAULT NULL COMMENT '单位名称',
  `quantity` double(14,2) NOT NULL DEFAULT '1.00' COMMENT '排产数量',
  `quantity_produced` double(14,2) DEFAULT '0.00' COMMENT '已生产数量',
  `quantity_quanlify` double(14,2) DEFAULT '0.00' COMMENT '合格品数量',
  `quantity_unquanlify` double(14,2) DEFAULT '0.00' COMMENT '不良品数量',
  `quantity_changed` double(14,2) DEFAULT '0.00' COMMENT '调整数量',
  `client_id` bigint DEFAULT NULL COMMENT '客户ID',
  `client_code` varchar(64) DEFAULT NULL COMMENT '客户编码',
  `client_name` varchar(255) DEFAULT NULL COMMENT '客户名称',
  `client_nick` varchar(255) DEFAULT NULL COMMENT '客户简称',
  `start_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '开始生产时间',
  `duration` int DEFAULT '1' COMMENT '生产时长',
  `end_time` datetime DEFAULT NULL COMMENT '完成生产时间',
  `color_code` char(7) DEFAULT '#00AEF3' COMMENT '甘特图显示颜色',
  `request_date` datetime DEFAULT NULL COMMENT '需求日期',
  `status` int DEFAULT '0' COMMENT '生产状态 0 是暂停 1 是启动',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `attr1` varchar(64) DEFAULT NULL COMMENT '预留字段1',
  `attr2` varchar(255) DEFAULT NULL COMMENT '预留字段2',
  `attr3` int DEFAULT '0' COMMENT '预留字段3',
  `attr4` int DEFAULT '0' COMMENT '预留字段4',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`task_id`)
) ENGINE=InnoDB AUTO_INCREMENT=200 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='生产任务表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_pro_task`
--

LOCK TABLES `sys_pro_task` WRITE;
/*!40000 ALTER TABLE `sys_pro_task` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_pro_task` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_pro_workorder`
--

DROP TABLE IF EXISTS `sys_pro_workorder`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_pro_workorder` (
  `workorder_id` bigint NOT NULL AUTO_INCREMENT COMMENT '工单ID',
  `workorder_code` varchar(64) NOT NULL COMMENT '工单编码',
  `workorder_name` varchar(255) NOT NULL COMMENT '工单名称',
  `workorder_type` varchar(64) DEFAULT 'SELF' COMMENT '工单类型',
  `order_source` varchar(64) NOT NULL COMMENT '来源类型',
  `source_code` varchar(64) DEFAULT NULL COMMENT '来源单据',
  `router_id` bigint NOT NULL COMMENT '工艺路线id',
  `product_id` bigint NOT NULL COMMENT '产品ID',
  `product_code` varchar(64) NOT NULL COMMENT '产品编号',
  `product_name` varchar(255) NOT NULL COMMENT '产品名称',
  `product_spc` varchar(255) DEFAULT NULL COMMENT '规格型号',
  `unit_of_measure` varchar(64) NOT NULL COMMENT '单位',
  `quantity` double(14,2) NOT NULL DEFAULT '0.00' COMMENT '生产数量',
  `quantity_produced` double(14,2) DEFAULT '0.00' COMMENT '已生产数量',
  `quantity_changed` double(14,2) DEFAULT '0.00' COMMENT '调整数量',
  `quantity_scheduled` double(14,2) DEFAULT '0.00' COMMENT '已排产数量',
  `client_id` bigint DEFAULT NULL COMMENT '客户ID',
  `client_code` varchar(64) DEFAULT NULL COMMENT '客户编码',
  `client_name` varchar(255) DEFAULT NULL COMMENT '客户名称',
  `vendor_id` bigint DEFAULT NULL COMMENT '供应商ID',
  `vendor_code` varchar(64) DEFAULT NULL COMMENT '供应商编号',
  `vendor_name` varchar(255) DEFAULT NULL COMMENT '供应商名称',
  `batch_code` varchar(64) DEFAULT NULL COMMENT '批次号',
  `request_date` datetime NOT NULL COMMENT '需求日期',
  `parent_id` bigint NOT NULL DEFAULT '0' COMMENT '父工单',
  `ancestors` varchar(500) NOT NULL COMMENT '所有父节点ID',
  `finish_date` datetime DEFAULT NULL COMMENT '完成时间',
  `status` varchar(64) DEFAULT 'PREPARE' COMMENT '单据状态',
  `remark` varchar(500) DEFAULT '' COMMENT '备注',
  `attr1` varchar(64) DEFAULT NULL COMMENT '预留字段1',
  `attr2` varchar(255) DEFAULT NULL COMMENT '预留字段2',
  `attr3` int DEFAULT '0' COMMENT '预留字段3',
  `attr4` int DEFAULT '0' COMMENT '预留字段4',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`workorder_id`)
) ENGINE=InnoDB AUTO_INCREMENT=200 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='生产工单表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_pro_workorder`
--

LOCK TABLES `sys_pro_workorder` WRITE;
/*!40000 ALTER TABLE `sys_pro_workorder` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_pro_workorder` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_product_laboratory`
--

DROP TABLE IF EXISTS `sys_product_laboratory`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_product_laboratory` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `material` varchar(255) NOT NULL COMMENT '物料',
  `contact` varchar(255) NOT NULL COMMENT '采样人',
  `date` varchar(255) NOT NULL COMMENT '采样时间',
  `address` varchar(255) NOT NULL COMMENT '采样地点',
  `dry_high` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '干燥基高位',
  `received_low` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '收到基低位',
  `heat` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '含热量',
  `sulphur` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '含硫量',
  `volatility` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '挥发性',
  `water` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '含水量',
  `weight` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '样重',
  `number` varchar(255) NOT NULL COMMENT '样号',
  `img` varchar(512) NOT NULL COMMENT '图片地址',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) NOT NULL DEFAULT '0',
  `type` tinyint(1) NOT NULL COMMENT '0 为物料化验 1为 配料化验',
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`)
) ENGINE=InnoDB AUTO_INCREMENT=464713292607131649 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_product_schedule`
--

DROP TABLE IF EXISTS `sys_product_schedule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_product_schedule` (
  `id` bigint NOT NULL COMMENT '调度id',
  `gateway_id` bigint NOT NULL COMMENT '网关id',
  `schedule_name` varchar(255) NOT NULL COMMENT '计划名称',
  `time` varchar(64) NOT NULL COMMENT '时间序列化格式,普通日程,1,2,3,4,5;特殊日程:2025-04-04 ~ 2025-04-04',
  `time_manager` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '时间任务列表',
  `schedule_type` tinyint(1) NOT NULL COMMENT '0为普通日程 1为特殊日程',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '操作状态（0正常 1启动）',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`id`),
  KEY `schedule_type` (`schedule_type`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;



--
-- Table structure for table `sys_product_schedule_map`
--

DROP TABLE IF EXISTS `sys_product_schedule_map`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_product_schedule_map` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `schedule_id` bigint NOT NULL COMMENT '计划标识',
  `gateway_id` bigint NOT NULL COMMENT '网关id',
  `begin_time` int NOT NULL DEFAULT '0' COMMENT '开始时间',
  `end_time` int NOT NULL DEFAULT '0' COMMENT '结束时间',
  `date` tinyint(1) NOT NULL DEFAULT '0' COMMENT '日期 1 2 3 4 5 6 7分别代表周几',
  `craft_route_id` bigint NOT NULL COMMENT '工艺路线id',
  `schedule_type` tinyint(1) NOT NULL COMMENT '0为 循环日程 1为特殊日程',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) NOT NULL DEFAULT '1',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '操作状态（0正常 1启动）',
  PRIMARY KEY (`id`),
  KEY `schedule_id` (`schedule_id`) USING BTREE,
  KEY `begin_time` (`begin_time`)
) ENGINE=InnoDB AUTO_INCREMENT=467319757331238919 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_resource_file`
--

DROP TABLE IF EXISTS `sys_resource_file`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_resource_file` (
  `resource_id` bigint NOT NULL COMMENT '主键ID',
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '资源名称',
  `original_name` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '原始文件名（仅文件有）',
  `extension` varchar(50) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '文件扩展名（仅文件有）',
  `mime_type` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'MIME类型（仅文件有）',
  `type` enum('FILE','FOLDER') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'FILE' COMMENT '类型：FILE-文件，FOLDER-文件夹',
  `path` varchar(1000) COLLATE utf8mb4_general_ci NOT NULL COMMENT '完整路径（不含域名）',
  `file_size` bigint unsigned DEFAULT '0' COMMENT '文件大小（字节）',
  `md5_hash` varchar(32) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '文件MD5哈希值',
  `storage_key` varchar(500) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '存储服务中的唯一标识',
  `storage_type` enum('LOCAL','OSS','S3','COS') COLLATE utf8mb4_general_ci DEFAULT 'LOCAL' COMMENT '存储类型',
  `parent_id` bigint unsigned DEFAULT '0' COMMENT '父级ID（0表示根目录）',
  `level` tinyint unsigned DEFAULT '0' COMMENT '层级深度',
  `lineage` varchar(1000) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '层级路径，格式：/父ID/祖父ID/...',
  `permission_mask` int unsigned DEFAULT '0' COMMENT '权限掩码',
  `category` varchar(50) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '分类：IMAGE, DOCUMENT, VIDEO, AUDIO等',
  `status` tinyint(1) DEFAULT '1' COMMENT '状态：0-删除，1-正常，2-隐藏',
  `description` text COLLATE utf8mb4_general_ci COMMENT '描述',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) DEFAULT '0' COMMENT '操作状态（0正常 -1删除）',
  PRIMARY KEY (`resource_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_type_category` (`type`,`category`),
  KEY `idx_md5_hash` (`md5_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='系统资源文件表';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `sys_role`
--

DROP TABLE IF EXISTS `sys_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_role` (
  `role_id` bigint NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `role_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '角色名称',
  `role_sort` int NOT NULL COMMENT '显示顺序',
  `status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '角色状态（0正常 1停用）',
  `del_flag` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '0' COMMENT '删除标志（0代表存在 2代表删除）',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`role_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=441489218926022657 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='角色信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_role`
--

LOCK TABLES `sys_role` WRITE;
/*!40000 ALTER TABLE `sys_role` DISABLE KEYS */;
INSERT INTO `sys_role` VALUES (1,'admin',1,'0','0',1,'2024-02-08 04:10:55',1,NULL,'超级管理员'),(2,'common',2,'0','0',1,'2024-02-08 04:10:55',359920431269941248,'2025-06-11 21:52:36','普通角色'),(441489218926022656,'化验员',0,'0','0',1,'2025-12-08 14:40:32',1,'2025-12-08 14:40:32','');
/*!40000 ALTER TABLE `sys_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_role_permission`
--

DROP TABLE IF EXISTS `sys_role_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_role_permission` (
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `permission_id` bigint NOT NULL COMMENT '权限ID',
  PRIMARY KEY (`role_id`,`permission_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='角色和权限关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_role_permission`
--

LOCK TABLES `sys_role_permission` WRITE;
/*!40000 ALTER TABLE `sys_role_permission` DISABLE KEYS */;
INSERT INTO `sys_role_permission` VALUES (2,1),(2,2),(2,3),(2,100),(2,101),(2,102),(2,103),(2,104),(2,105),(2,106),(2,107),(2,108),(2,109),(2,110),(2,112),(2,115),(2,116),(2,117),(2,500),(2,501),(2,1000),(2,1001),(2,1002),(2,1003),(2,1004),(2,1005),(2,1006),(2,1007),(2,1008),(2,1009),(2,1010),(2,1011),(2,1012),(2,1013),(2,1014),(2,1015),(2,1016),(2,1017),(2,1018),(2,1019),(2,1020),(2,1021),(2,1022),(2,1023),(2,1024),(2,1025),(2,1026),(2,1027),(2,1028),(2,1029),(2,1030),(2,1031),(2,1032),(2,1033),(2,1034),(2,1035),(2,1036),(2,1037),(2,1038),(2,1039),(2,1040),(2,1041),(2,1042),(2,1043),(2,1044),(2,1046),(2,1047),(2,1048),(2,1049),(2,1050),(2,1051),(2,1052),(2,1053),(2,1054),(2,1055),(2,1056),(2,1057),(2,1058),(2,1059),(2,1060),(2,359684443914375168),(2,359686535840272384),(2,359686725640916992),(2,359900397353046016),(2,359926159607074816),(2,360335414705262592),(2,360342197171654656),(2,360816738536787968),(2,360816857738907648),(2,360816923304267776),(2,361010500642803712),(2,361010584994451456),(2,361044431924629504),(2,361044530234920960),(2,361474631036571648),(2,361475137612025856),(2,361475259188121600),(2,361787734412496896),(2,361787808001560576),(2,362089395790024704),(2,362247033651728384),(2,362521217439109120),(2,362576536529801216),(2,362576598731329536),(2,362651361076056064),(2,362651427786461184),(2,362651499408396288),(2,362651569637822464),(2,362651630966935552),(2,363936074398961664),(2,363936143873413120),(2,363936228376055808),(2,363936291261255680),(2,364059932867170304),(2,364059987762221056),(2,364413536425742336),(2,364413837832622080),(2,364421764240904192),(2,364433815591981056),(2,364442854312906752),(2,364652283557842944),(2,365685718661468160),(2,365685913927290880),(2,365685977705877504),(2,366056868528787456),(2,366056957984903168),(2,366057043401904128),(2,366230539859922944),(2,366568187384303616),(2,366777330686758912),(2,366863885531090944),(2,366863942808506368),(2,366863990531297280),(2,367663286088372224),(2,367664829529329664),(2,367664941072650240),(2,367665204437192704),(2,367665320741048320),(2,367665410381713408),(2,367940361919664128),(2,367940533215039488),(2,367940616467779584),(2,370897941847609344),(2,373343224594436096),(2,373343302742708224),(2,373343417628889088),(2,373343563192209408),(2,373374786467794944),(2,373706277987028992),(2,373706373285810176),(2,373706447705346048),(2,374788609506545664),(2,374788724715687936),(2,374788976810135552),(2,374789035895296000),(2,374789081369939968),(2,375148581360766976),(2,375148631210070016),(2,375973793090244608),(441489218926022656,441451927406907392),(441489218926022656,441452001889357824),(441489218926022656,441452113885663232),(441489218926022656,441452189274083328),(441489218926022656,441452285436891136),(441489218926022656,441452343255371776),(441489218926022656,441452404488015872);
/*!40000 ALTER TABLE `sys_role_permission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user`
--

DROP TABLE IF EXISTS `sys_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user` (
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `user_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户账号',
  `nick_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户昵称',
  `email` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '用户邮箱',
  `phonenumber` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '手机号码',
  `sex` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '0' COMMENT '用户性别（0男 1女 2未知）',
  `avatar` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '头像地址',
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '密码',
  `data_scope` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '5' COMMENT '数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限,无任何）权限',
  `status` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '0' COMMENT '帐号状态（0正常 1停用）',
  `del_flag` char(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '0' COMMENT '删除标志（0代表存在 2代表删除）',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '备注',
  PRIMARY KEY (`user_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user`
--

LOCK TABLES `sys_user` WRITE;
/*!40000 ALTER TABLE `sys_user` DISABLE KEYS */;
INSERT INTO `sys_user` VALUES (1,105,'admin','白泽','bz@163.com','15888888887','0','http://81.71.98.26:11801/profile/1/2025-07-11/0197f797-f316-7c68-9bad-0a836d08165a.jpg','$2a$14$LX4INaYU1kCGISEJyST9Qe3JR/r1hddldg40TcSuvWSH8AIWdnTiS','1','0','0',1,'2024-02-08 04:10:55',1,'2026-02-21 18:11:35','管理员');
/*!40000 ALTER TABLE `sys_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user_dept_scope`
--

DROP TABLE IF EXISTS `sys_user_dept_scope`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user_dept_scope` (
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `dept_id` bigint NOT NULL COMMENT '部门ID',
  PRIMARY KEY (`user_id`,`dept_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户和部门关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user_dept_scope`
--

LOCK TABLES `sys_user_dept_scope` WRITE;
/*!40000 ALTER TABLE `sys_user_dept_scope` DISABLE KEYS */;
INSERT INTO `sys_user_dept_scope` VALUES (2,101),(2,102),(2,103),(2,104),(2,108),(2,109);
/*!40000 ALTER TABLE `sys_user_dept_scope` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user_post`
--

DROP TABLE IF EXISTS `sys_user_post`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user_post` (
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `post_id` bigint NOT NULL COMMENT '岗位ID',
  PRIMARY KEY (`user_id`,`post_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户与岗位关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user_post`
--

LOCK TABLES `sys_user_post` WRITE;
/*!40000 ALTER TABLE `sys_user_post` DISABLE KEYS */;
INSERT INTO `sys_user_post` VALUES (1,1),(2,2),(359920431269941248,2),(417514047974412288,4);
/*!40000 ALTER TABLE `sys_user_post` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user_role`
--

DROP TABLE IF EXISTS `sys_user_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user_role` (
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`user_id`,`role_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户和角色关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user_role`
--

LOCK TABLES `sys_user_role` WRITE;
/*!40000 ALTER TABLE `sys_user_role` DISABLE KEYS */;
INSERT INTO `sys_user_role` VALUES (1,1),(2,2),(359920431269941248,1),(359920431269941248,2),(417514047974412288,2),(441488943934869504,441489218926022656);
/*!40000 ALTER TABLE `sys_user_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_work_shift_setting`
--

DROP TABLE IF EXISTS `sys_work_shift_setting`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_work_shift_setting` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL COMMENT '班次名称',
  `begin_time` int NOT NULL DEFAULT '0' COMMENT '开始时间',
  `begin_time_str` varchar(255) NOT NULL COMMENT '开始时间字符串',
  `end_time` int NOT NULL DEFAULT '0' COMMENT '结束时间',
  `end_time_str` varchar(255) NOT NULL COMMENT '结束时间字符串',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用班次设置',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `state` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=424626908886470657 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;



/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-02-24 15:19:44

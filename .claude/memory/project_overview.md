---
name: nova-factory-server 项目概览
description: 白泽工厂管理系统 - 基于 Go + Gin 的企业级快速开发平台
type: project
---

# 项目概述
**项目名称**: nova-factory-server（白泽工厂服务器）
**技术栈**: Go 1.26.1 + Gin + Zap + GORM + Wire

## 架构模式
- 清晰分层架构: Controller → Service → DAO
- 依赖注入: Google Wire
- API文档: Swagger

## 核心业务模块

### admin 后台管理
- **system**: 用户、角色、菜单、部门、岗位、字典、参数、通知
- **monitor**: 定时任务、操作日志、登录日志、在线用户、服务监控
- **tool**: 代码生成器（一键生成CRUD）
- **product**: 产品管理

### ai 智能体模块
- **agent/dataset**: AI数据集管理（集成RAGFlow）
- **gateway**: LLM网关路由
- **core**: LLM客户端抽象

### iot 物联网模块
- **asset**: 设备资产、摄像头、建筑、物料
- **metric**: 设备指标采集（IoTDB/ClickHouse）
- **alert**: 告警管理
- **dashboard**: 监控仪表盘
- **craft**: 工艺路线
- **devicemonitor**: 设备实时监控
- **daemonize**: 设备守护进程通信
- **home**: 首页概览

### erp 企业资源计划
- **order**: 订单管理
- **setting**: 系统设置

### shop 商城模块
- **product**: 商品管理
- **order**: 订单管理
- **user**: 用户管理
- **activity**: 活动管理
- **home**: 商城首页

## 数据存储
- **MySQL**: 主数据库（业务数据）
- **Redis**: 缓存 + Session共享
- **IoTDB/ClickHouse**: 时序指标数据
- **S3/Local**: 文件存储

## 关键特性
- RBAC 细粒度权限控制
- 动态权限菜单
- 多终端认证（Session共享）
- 代码生成器（Go/SQL）
- MCP（Model Context Protocol）支持
- 定时任务调度

## 构建命令
```bash
wire gen -tags="ai iot"    # 生成依赖注入
go build -tags="ai iot"    # 编译指定模块
swag init                  # 生成API文档
```

**默认账号**: admin/admin123
**API文档**: http://localhost:8080/swagger/doc.json

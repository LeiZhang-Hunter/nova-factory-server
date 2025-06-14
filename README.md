
## 平台简介

白泽是一套全部开源的快速开发平台，毫无保留给个人及企业免费使用。

## 前端代码地址
[baize-react](https://gitee.com/baizeplus/baize-react) https://gitee.com/baizeplus/baize-react  , https://github.com/baizeplus/baize-react
<br>
[baize-vue](https://gitee.com/baizeplus/baize-vue) https://gitee.com/baizeplus/baize-vue  , https://github.com/baizeplus/baize-vue
<br>

# Baize: A Comprehensive Web Application Framework with Role-Based Access Control

Baize is a modern Go web application framework that provides a robust foundation for building secure, scalable, and maintainable web applications with comprehensive role-based access control (RBAC) and monitoring capabilities.

The framework implements a clean architecture pattern with clear separation of concerns between controllers, services, and data access layers. It provides built-in support for user management, authentication, authorization, configuration management, monitoring, and code generation.

## Repository Structure
```
.
├── app/                            # Main application code
│   ├── nova-factory-server/                     # Core framework utilities and base types
│   ├── business/                  # Business logic modules
│   │   ├── monitor/              # Monitoring functionality (jobs, logs, etc.)
│   │   ├── system/               # Core system functionality (users, roles, etc.)
│   │   └── tool/                 # Development tools and generators
│   ├── constant/                 # Application constants
│   ├── datasource/              # Data source implementations (MySQL, Redis, S3)
│   ├── docs/                    # API documentation (Swagger)
│   ├── middlewares/             # HTTP middleware components
│   ├── routes/                  # Route definitions
│   ├── setting/                 # Application configuration
│   └── utils/                   # Utility functions
├── config/                      # Configuration files
├── template/                    # Code generation templates
└── sql/                        # Database schema definitions
```


## Data Flow
The application follows a clean architecture pattern with clear data flow:

1. HTTP Request → Middleware (Auth/Logging) → Controller
2. Controller → Service Layer (Business Logic)
3. Service Layer → Data Access Layer (DAO)
4. DAO → Database/Cache/Storage

```ascii
Request → [Middleware] → [Controller] → [Service] → [DAO] → [Database]
                                                         ↓
                                                    [Cache/Redis]
                                                         ↓
                                                  [Storage (S3/Local)]
```

Key Component Interactions:
- Controllers handle HTTP requests and response formatting
- Services implement business logic and transaction management
- DAOs handle data persistence and retrieval
- Middleware provides cross-cutting concerns (auth, logging)
- Cache layer improves performance for frequently accessed data
- Storage service handles file uploads and downloads





- 后端采用Gin、Zap、sqly(sqlx升级)。
- 权限认证使用共享Session(单机支持本地缓存,集群需要redis)，支持多终端认证系统。
- 支持加载动态权限菜单，多方式轻松权限控制。
- 高效率开发，使用代码生成器可以一键生成后端代码。(前端正在发开)
- 特别鸣谢：[ruoyi-vue](https://gitee.com/y_project/RuoYi-Vue?_from=gitee_search )，

## <p>随手 star ⭐是一种美德。 你们的star就是我的动力</p>

## 也期待小伙伴一起加入白泽   

## 内置功能

1. 用户管理：用户是系统操作者，该功能主要完成系统用户配置。
2. 部门管理：配置系统组织机构（公司、部门、小组），树结构展现支持数据权限。
3. 岗位管理：配置系统用户所属担任职务。
4. 菜单管理：配置系统菜单，操作权限，按钮权限标识等。
5. 角色管理：角色菜单权限分配、设置角色按机构进行数据范围权限划分。
6. 字典管理：对系统中经常使用的一些较为固定的数据进行维护。
7. 参数管理：对系统动态配置常用参数。
8. 通知公告：系统通知公告信息发布维护。
9. 登录日志：系统登录日志记录查询包含登录异常。
10. 在线用户：当前系统中活跃用户状态监控。
11. 服务监控：监视当前系统CPU、内存、磁盘等相关信息。
12. 操作日志：系统正常操作日志记录和查询；系统异常信息日志记录和查询。
13. 系统接口：根据业务代码注释自动生成相关的api接口文档。
14. 代码生成：前后端代码的生成（Go、sql）支持CRUD下载 。
15. 定时任务：在线（添加、修改、删除)任务调度。
## 版本规则
v5.6.7<br>
1位为主版本号（5）：当功能模块有较大的变动，比如增加多个模块或者整体架构发生变化,数据库结构发生变化。此版本号由项目决定是否修改。
<br>
2为次版本号（6）：当功能有一定的增加或变化，比如修改了API接口。此版本号由项目决定是否修改。
<br>
3为阶段版本号(7)：一般是 Bug 修复或是一些小的变动，要经常发布修订版，时间间隔不限。
<br>
主版本号升级请参考更新说明更新修改或添加相应的数据表。
次版本号升级请参考更新说明查看API接口修改情况。
阶段版本号不会影响数据库与api接口，除修复重大bug不更新说明文档



## 在线体验

- admin/admin123

react演示地址：http://react.ibaize.vip
<br>
vue3演示地址：http://vue.ibaize.vip
<br>
文档地址：https://doc.ibaize.vip
<br>

## 演示图

<table>
    <tr>
        <td><img src="https://gitee.com/smell2/nova-factory-server/raw/imgs/202110241805797.jpg"/></td>
        <td><img src="https://gitee.com/smell2/nova-factory-server/raw/imgs/202110241806256.jpg"/></td>
    </tr>
    <tr>
        <td><img src="https://gitee.com/smell2/nova-factory-server/raw/imgs/202110242322137.png"/></td>
        <td><img src="https://gitee.com/smell2/nova-factory-server/raw/imgs/202110242323820.png"/></td>
    </tr>  
    <tr>
        <td><img src="https://gitee.com/smell2/nova-factory-server/raw/imgs/202112082243214.png"/></td>
        <td><img src="https://gitee.com/smell2/nova-factory-server/raw/imgs/202112082242154.png"/></td>
    </tr>

</table>

## 重新生成wire

```
wire gen
```

## 重新生成swag

```
swag init
```

swag 地址:

```
http://localhost:8080/swagger/doc.json
```


## 白泽管理系统交流群


QQ群： [![加入QQ群](https://img.shields.io/badge/83064682-blue.svg)](https://qm.qq.com/cgi-bin/qm/qr?k=rAIw_VQ_blbSQu0J6fApnm5RbAc2CHbp&jump_from=webapi) 点击按钮入群。
##欢迎加入

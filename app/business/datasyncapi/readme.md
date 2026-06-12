# datasyncapi

管家婆开放平台 API 对接模块，提供 OAuth 授权、商品同步、库存更新等服务。

## 目录结构

```
datasyncapi/
├── controller/qqd/       # 路由控制器
│   ├── authorize.go       # OAuth 授权回调处理
│   ├── token.go           # Token 签发与刷新
│   ├── api.go             # 通用 API 网关（method 分发）
│   └── helpers.go         # 请求参数解析、错误响应封装
├── dao/                   # 数据访问层
│   ├── i_qqd_goods_dao.go       # 商品 DAO 接口
│   ├── i_qqd_goods_sku_dao.go   # 商品 SKU DAO 接口
│   └── impl/                    # DAO 实现
├── models/                # 数据模型
│   └── qqd.go             # QQDConfig / 商品表 / SKU 表结构
├── service/               # 业务逻辑层
│   ├── qqd/i_qqd_service.go          # Service 接口定义
│   └── impl/qqd/                     # Service 实现（含签名算法）
├── description.go         # 构建标签声明
├── register.go            # Wire 依赖注入注册
└── routes.go              # 路由注册（/api/v1/erp-api/qqd）
```

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/erp-api/qqd/oauth/authorize` | OAuth 授权，生成回调地址 |
| POST | `/api/v1/erp-api/qqd/oauth/token` | 换取 / 刷新 access_token |
| POST | `/api/v1/erp-api/qqd/api` | 通用 API 入口，通过 `method` 参数分发 |

### method 支持值

- `selfmall.product.list.get` — 商品列表查询
- `selfmall.product.add` — 批量添加商品
- `selfmall.product.stock.update` — 商品库存更新

## 签名算法

参数按 key 升序排列，拼接规则：

```
sign = MD5(appSecret + key1value1 + key2value2 + ... + body + appSecret)
```

详见 `service/impl/qqd/sign.go`。

## 构建标签

本模块通过 `erp` 构建标签控制条件编译，需使用 `-tags erp` 编译。

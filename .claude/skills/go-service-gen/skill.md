---
name: "go-service-gen"
description: "根据 Service 接口定义自动生成对应的 ServiceImpl、DAO 接口及 DAOImpl 基础代码。当用户定义好一个新的 Service 接口或要求生成服务实现层时调用。"
---

# Go Service Generator

该 Skill 用于自动化生成符合本项目四层架构（Interface -> ServiceImpl -> DAO Interface -> DAOImpl）的 Go 代码。

## 使用场景
- 当用户在 `aidatasetservice` 或类似的包中定义了新的 `interface` 后。
- 当用户要求“实现该服务”或“生成配套的 DAO 代码”时。

## 重要规则

### 1. 每个文件只能有一个结构体
- 一个 `.go` 文件只包含一个结构体的定义、构造函数和方法实现
- 多个相关的 DAO/Service 应分别放在独立文件中
- 文件名应与结构体名对应，如 `i_shop_order_dao_impl.go` 对应 `IShopOrderDaoImpl`

### 2. DAO 结构体必须包含 db 和 tableName 字段
```go
type IShopOrderDaoImpl struct {
    db        *gorm.DB
    tableName string
}
```

### 3. 所有数据库操作必须使用 .Table(d.tableName)
- **错误**: `d.db.WithContext(c).Where(...)` 或 `d.db.WithContext(c).Model(&models.Order{}).Where(...)`
- **正确**: `d.db.WithContext(c).Table(d.tableName).Where(...)`

原因：GORM 默认使用结构体名的复数形式作为表名（如 `Order` → `orders`），必须显式指定实际表名。

### 4. DAO 方法中使用 d.db 而非 c.Value("db")
- **正确**: `d.db.WithContext(c).Table(d.tableName).Create(order)`
- **错误**: `c.Value("db").(*gorm.DB).WithContext(c).Create(order)`

### 5. 构造函数必须注入 db 并设置 tableName
```go
func NewIShopOrderDaoImpl(db *gorm.DB) dao.IShopOrderDao {
    return &IShopOrderDaoImpl{
        db:        db,
        tableName: "shop_order",  // 必须指定实际表名
    }
}
```

## 注释要求
- 生成的代码必须带上必要注释，不能只给裸代码。
- 导出类型、构造函数、导出方法都要补充 Go 风格注释，便于后续维护。
- 注释内容要说明职责、关键依赖和方法用途，避免无意义重复描述。
- 如果是 DAO 实现，优先说明表用途、查询条件和软删除/部门隔离等行为。

## 代码生成规则

### 1. Service 实现 (ServiceImpl)

#### Admin 层（shop/xxx/service/）
- **接口路径**: `.../service/i_shop_<name>_service.go`
- **实现路径**: `.../service/impl/i_shop_<name>_service_impl.go`
- **文件名**: 一个文件一个结构体，如 `i_shop_order_service_impl.go`
- **结构**:
    - 包含 `type I<Name>ServiceImpl struct { dao dao.I<Name>Dao }`
    - 包含构造函数 `func NewI<Name>ServiceImpl(dao dao.I<Name>Dao) service.I<Name>Service`
    - 实现接口中的所有方法，默认逻辑为直接调用 `dao` 的对应方法。
    - `struct`、构造函数、导出方法都要带注释。

#### API 层（shop/api/service/）
- **接口路径**: `.../service/i_api_shop_<name>_service.go`
- **实现路径**: `.../service/impl/i_api_shop_<name>_service_impl.go`
- **文件名**: 例如 `i_api_shop_cart_service.go`（接口）, `i_api_shop_cart_service_impl.go`（实现）
- **结构**:
    - 包含 `type IApiShop<Name>ServiceImpl struct { dao dao.IApiShop<Name>Dao }`
    - 包含构造函数 `func NewIApiShop<Name>ServiceImpl(dao dao.IApiShop<Name>Dao) service.IApiShop<Name>Service`
    - 实现接口中的所有方法，默认逻辑为直接调用 `dao` 的对应方法。
    - `struct`、构造函数、导出方法都要带注释。

### 2. DAO 接口 (DAO Interface)

#### Admin 层（shop/xxx/dao/）
- **路径**: `.../dao/i_shop_<name>_dao.go`
- **文件名**: 一个文件一个接口，如 `i_shop_order_dao.go`

#### API 层（shop/api/dao/）
- **路径**: `.../dao/i_api_shop_<name>_dao.go`
- **文件名**: 一个文件一个接口，如 `i_api_shop_cart_dao.go`

### 3. DAO 实现 (DAOImpl)

#### Admin 层（shop/xxx/dao/impl/）
- **路径**: `.../dao/impl/i_shop_<name>_dao_impl.go`
- **文件名**: 一个文件一个结构体，如 `i_shop_order_dao_impl.go`
- **结构**:
    - 包含 `type I<Name>DaoImpl struct { db *gorm.DB, tableName string }`
    - 包含构造函数 `func NewI<Name>DaoImpl(db *gorm.DB) dao.I<Name>Dao`
    - **必须指定 tableName**
    - 实现基础的 CRUD 逻辑
    - **使用 d.db.WithContext(c).Table(d.tableName) 进行所有数据库操作**
    - DAO 实现类型、构造函数、CRUD 方法都要带注释。

#### API 层（shop/api/dao/impl/）
- **路径**: `.../dao/impl/i_api_shop_<name>_dao_impl.go`
- **文件名**: 一个文件一个结构体，如 `i_api_shop_cart_dao_impl.go`
- **结构**:
    - 包含 `type IApiShop<Name>DaoImpl struct { db *gorm.DB, tableName string }`
    - 包含构造函数 `func NewIApiShop<Name>DaoImpl(db *gorm.DB) dao.IApiShop<Name>Dao`
    - **必须指定 tableName**
    - 实现基础的 CRUD 逻辑
    - **使用 d.db.WithContext(c).Table(d.tableName) 进行所有数据库操作**

### 4. Provider 集中注册
- **路径**: `.../dao/impl/provider.go` 或 `.../service/impl/provider.go`
- 在 Provider 中集中注册所有同包下的构造函数

## 示例模板

### Admin 层 ServiceImpl 模板
```go
package impl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
)

// IShopOrderServiceImpl 提供订单相关的业务实现。
type IShopOrderServiceImpl struct {
	orderDao dao.IShopOrderDao
}

// NewIShopOrderServiceImpl 创建订单服务实现。
func NewIShopOrderServiceImpl(orderDao dao.IShopOrderDao) service.IShopOrderService {
	return &IShopOrderServiceImpl{
		orderDao: orderDao,
	}
}

// Create 创建订单。
func (s *IShopOrderServiceImpl) Create(c *gin.Context, username string, req *models.OrderSetReq) (*models.Order, error) {
	// 实现逻辑...
}
```

### API 层 ServiceImpl 模板
```go
package impl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
)

// IApiShopCartServiceImpl 提供购物车相关的业务实现。
type IApiShopCartServiceImpl struct {
	cartDao dao.IApiShopCartDao
}

// NewIApiShopCartServiceImpl 创建购物车服务实现。
func NewIApiShopCartServiceImpl(cartDao dao.IApiShopCartDao) service.IApiShopCartService {
	return &IApiShopCartServiceImpl{
		cartDao: cartDao,
	}
}

// Add 添加商品到购物车。
func (s *IApiShopCartServiceImpl) Add(c *gin.Context, userID int64, req *models.CartAddReq) error {
	// 实现逻辑...
}
```

**API 层命名要点**：
- 接口文件名：`i_api_shop_<name>_service.go`（如 `i_api_shop_cart_service.go`）
- 实现文件名：`i_api_shop_<name>_service_impl.go`（如 `i_api_shop_cart_service_impl.go`）
- 结构体名：`IApiShop<Name>ServiceImpl`（如 `IApiShopCartServiceImpl`）
- 构造函数：`NewIApiShop<Name>ServiceImpl`
- 接口返回类型：`service.IApiShop<Name>Service`

### DAOImpl 模板（API 层示例）
```go
package impl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
)

// IApiShopGoodsDaoImpl 提供商品的数据库访问能力。
type IApiShopGoodsDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewIApiShopGoodsDaoImpl 创建商品 DAO 实现。
func NewIApiShopGoodsDaoImpl(db *gorm.DB) dao.IApiShopGoodsDao {
	return &IApiShopGoodsDaoImpl{
		db:        db,
		tableName: "shop_goods",
	}
}

// GetByID 根据ID获取商品。
func (d *IApiShopGoodsDaoImpl) GetByID(c *gin.Context, id int64) (*models.Goods, error) {
	var item models.Goods
	err := d.db.WithContext(c).Table(d.tableName).Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// List 查询商品列表，支持分页和条件筛选。
func (d *IApiShopGoodsDaoImpl) List(c *gin.Context, query *models.GoodsQuery) (*models.GoodsListData, error) {
	var total int64
	var items []*models.Goods

	// 构建查询条件
	q := d.db.WithContext(c).Table(d.tableName)

	if query.GoodsName != "" {
		q = q.Where("goods_name LIKE ?", "%"+query.GoodsName+"%")
	}
	// ... 其他条件

	// 统计总数
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (query.Page - 1) * query.Size
	if err := q.Order("id DESC").Offset(offset).Limit(query.Size).Find(&items).Error; err != nil {
		return nil, err
	}

	return &models.GoodsListData{Rows: items, Total: total}, nil
}
```

### Provider 模板

#### Admin 层 Provider
```go
package impl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewIShopOrderDaoImpl, NewIShopOrderItemDaoImpl, NewIShopOrderServiceImpl)
```

#### API 层 Provider
```go
package impl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewIApiShopGoodsDaoImpl,
	NewIApiShopGoodsServiceImpl,
	NewIApiShopCartDaoImpl,
	NewIApiShopCartServiceImpl,
)
```

## 常见错误

| 错误写法 | 正确写法 | 说明 |
|---------|---------|------|
| `d.db.WithContext(c).Where(...)` | `d.db.WithContext(c).Table(d.tableName).Where(...)` | 必须指定表名 |
| `d.db.WithContext(c).Model(&models.Order{}).Where(...)` | `d.db.WithContext(c).Table(d.tableName).Where(...)` | 必须指定表名 |
| `c.Value("db").(*gorm.DB).WithContext(c)` | `d.db.WithContext(c)` | db 已注入到结构体 |
| 结构体只有 `db` 字段 | `struct { db *gorm.DB; tableName string }` | 必须包含 tableName |
| 一个文件多个结构体 | 一个文件一个结构体 | 便于维护 |

## Shop API 模块命名规范（仅适用于 `shop/api/`）

**适用范围**：小程序/前端接口 `shop/api/`

### 命名对照表

| 层 | 文件格式 | 接口/结构体格式 | 示例 |
|----|----------|-----------------|------|
| Service 接口 | `i_api_shop_{entity}_service.go` | `IApiShop{Entity}Service` | `IApiShopGoodsService` |
| Service Impl | `i_api_shop_{entity}_service_impl.go` | `IApiShop{Entity}ServiceImpl` | `IApiShopGoodsServiceImpl` |
| DAO 接口 | `i_api_shop_{entity}_dao.go` | `IApiShop{Entity}Dao` | `IApiShopGoodsDao` |
| DAO Impl | `i_api_shop_{entity}_dao_impl.go` | `IApiShop{Entity}DaoImpl` | `IApiShopGoodsDaoImpl` |

### 构造函数对照

| 结构体 | 构造函数 |
|--------|----------|
| Controller | `New{Entity}` |
| Service Impl | `NewIApiShop{Entity}ServiceImpl(dao IApiShop{Entity}Dao)` |
| DAO Impl | `NewIApiShop{Entity}DaoImpl(db *gorm.DB)` |

### 代码生成示例

**生成 goods 模块**：
```
输入: entity=goods, scope=api, layer=service
输出:
  - 文件: i_api_shop_goods_service.go
  - 接口: IApiShopGoodsService
  - Impl文件: i_api_shop_goods_service_impl.go
  - Impl结构: IApiShopGoodsServiceImpl
  - 构造函数: NewIApiShopGoodsServiceImpl

输入: entity=goods, scope=api, layer=dao
输出:
  - 文件: i_api_shop_goods_dao.go
  - 接口: IApiShopGoodsDao
  - Impl文件: i_api_shop_goods_dao_impl.go
  - Impl结构: IApiShopGoodsDaoImpl
  - 构造函数: NewIApiShopGoodsDaoImpl
  - tableName: "shop_goods"
```

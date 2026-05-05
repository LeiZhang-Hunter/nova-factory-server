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
- **路径**: `.../service/impl/i_<name>_service_impl.go`
- **文件名**: 一个文件一个结构体，如 `i_shop_order_service_impl.go`
- **结构**:
    - 包含 `type I<Name>ServiceImpl struct { dao dao.I<Name>Dao }`
    - 包含构造函数 `func NewI<Name>ServiceImpl(dao dao.I<Name>Dao) service.I<Name>Service`
    - 实现接口中的所有方法，默认逻辑为直接调用 `dao` 的对应方法。
    - `struct`、构造函数、导出方法都要带注释。

### 2. DAO 接口 (DAO Interface)
- **路径**: `.../dao/i_<name>_dao.go`
- **文件名**: 一个文件一个接口，如 `i_shop_order_dao.go`
- **结构**:
    - 接口方法通常与 Service 接口保持一致，或者根据底层存储需求调整。
    - 接口与方法要补充职责说明。
    - 如果有多个相关接口（如订单 + 订单明细），应放在不同文件中

### 3. DAO 实现 (DAOImpl)
- **路径**: `.../dao/impl/i_<name>_dao_impl.go`
- **文件名**: 一个文件一个结构体，如 `i_shop_order_dao_impl.go`
- **结构**:
    - 包含 `type I<Name>DaoImpl struct { db *gorm.DB, tableName string }`
    - 包含构造函数 `func NewI<Name>DaoImpl(db *gorm.DB) dao.I<Name>Dao`
    - **必须指定 tableName**
    - 实现基础的 CRUD 逻辑
    - **使用 d.db.WithContext(c).Table(d.tableName) 进行所有数据库操作**
    - DAO 实现类型、构造函数、CRUD 方法都要带注释。

### 4. Provider 集中注册
- **路径**: `.../dao/impl/provider.go` 或 `.../service/impl/provider.go`
- 在 Provider 中集中注册所有同包下的构造函数

## 示例模板

### ServiceImpl 模板
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

### DAOImpl 模板（重要）
```go
package impl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
)

// IShopOrderDaoImpl 提供订单的数据库访问能力。
type IShopOrderDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewIShopOrderDaoImpl 创建订单 DAO 实现。
func NewIShopOrderDaoImpl(db *gorm.DB) dao.IShopOrderDao {
	return &IShopOrderDaoImpl{
		db:        db,
		tableName: "shop_order",
	}
}

// Create 新增订单记录。
func (d *IShopOrderDaoImpl) Create(c *gin.Context, order *models.Order) (*models.Order, error) {
	if err := d.db.WithContext(c).Table(d.tableName).Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

// GetByID 根据ID获取订单。
func (d *IShopOrderDaoImpl) GetByID(c *gin.Context, id int64) (*models.Order, error) {
	var order models.Order
	err := d.db.WithContext(c).Table(d.tableName).Where("id = ? AND state = 0", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// List 查询订单列表，支持分页和条件筛选。
func (d *IShopOrderDaoImpl) List(c *gin.Context, query *models.OrderQuery) (*models.OrderListData, error) {
	var total int64
	var orders []*models.Order

	// 构建查询条件：仅查询未删除记录
	q := d.db.WithContext(c).Table(d.tableName).Where("state = 0")

	if query.UserID > 0 {
		q = q.Where("user_id = ?", query.UserID)
	}
	// ... 其他条件

	// 统计总数
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	if err := q.Order("create_time DESC").Offset(offset).Limit(size).Find(&orders).Error; err != nil {
		return nil, err
	}

	return &models.OrderListData{Rows: orders, Total: total}, nil
}

// UpdateStatus 更新订单状态。
func (d *IShopOrderDaoImpl) UpdateStatus(c *gin.Context, id int64, status int32, version int32) (int64, error) {
	result := d.db.WithContext(c).Table(d.tableName).
		Where("id = ? AND version = ?", id, version).
		Updates(map[string]interface{}{"status": status})
	return result.RowsAffected, result.Error
}
```

### Provider 模板
```go
package impl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewIShopOrderDaoImpl, NewIShopOrderItemDaoImpl)
```

## 常见错误

| 错误写法 | 正确写法 | 说明 |
|---------|---------|------|
| `d.db.WithContext(c).Where(...)` | `d.db.WithContext(c).Table(d.tableName).Where(...)` | 必须指定表名 |
| `d.db.WithContext(c).Model(&models.Order{}).Where(...)` | `d.db.WithContext(c).Table(d.tableName).Where(...)` | 必须指定表名 |
| `c.Value("db").(*gorm.DB).WithContext(c)` | `d.db.WithContext(c)` | db 已注入到结构体 |
| 结构体只有 `db` 字段 | `struct { db *gorm.DB; tableName string }` | 必须包含 tableName |
| 一个文件多个结构体 | 一个文件一个结构体 | 便于维护 |

---
name: "go-service-gen"
description: "根据 Service 接口定义自动生成对应的 ServiceImpl、DAO 接口及 DAOImpl 基础代码。当用户定义好一个新的 Service 接口或要求生成服务实现层时调用。"
---

# Go Service Generator

该 Skill 用于自动化生成符合本项目四层架构（Interface -> ServiceImpl -> DAO Interface -> DAOImpl）的 Go 代码。

## 使用场景
- 当用户在 `aidatasetservice` 或类似的包中定义了新的 `interface` 后。
- 当用户要求“实现该服务”或“生成配套的 DAO 代码”时。

## 注释要求
- 生成的代码必须带上必要注释，不能只给裸代码。
- 导出类型、构造函数、导出方法都要补充 Go 风格注释，便于后续维护。
- 注释内容要说明职责、关键依赖和方法用途，避免无意义重复描述。
- 如果是 DAO 实现，优先说明表用途、查询条件和软删除/部门隔离等行为。

## 代码生成规则

### 1. Service 实现 (ServiceImpl)
- **路径**: `.../aidatasetservice/aidatasetserviceimpl/i_<name>_service_impl.go`
- **结构**: 
    - 包含 `type <Name>ServiceImpl struct { dao aidatasetdao.I<Name>Dao }`
    - 包含构造函数 `func NewI<Name>ServiceImpl(dao aidatasetdao.I<Name>Dao) aidatasetservice.I<Name>Service`
    - 实现接口中的所有方法，默认逻辑为直接调用 `dao` 的对应方法。
    - `struct`、构造函数、导出方法都要带注释。

### 2. DAO 接口 (DAO Interface)
- **路径**: `.../aidatasetdao/i_<name>_dao.go`
- **结构**: 
    - 接口方法通常与 Service 接口保持一致，或者根据底层存储需求调整。
    - 接口与方法要补充职责说明。

### 3. DAO 实现 (DAOImpl)
- **路径**: `.../aidatasetdao/aiDataSetDaoImpl/i_<name>_dao_impl.go`
- **结构**:
    - 包含 `type I<Name>DaoImpl struct { db *gorm.DB, table string }`
    - 包含构造函数 `func NewI<Name>DaoImpl(db *gorm.DB) aidatasetdao.I<Name>Dao`
    - 实现基础的 CRUD 逻辑，使用 `baizeContext` 获取用户和部门信息。
    - DAO 实现类型、构造函数、CRUD 方法都要带注释。

## 示例模板

### ServiceImpl 模板
```go
package aidatasetserviceimpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
)

// I<Name>ServiceImpl 提供 <Name> 相关的业务实现。
type I<Name>ServiceImpl struct {
	dao aidatasetdao.I<Name>Dao
}

// NewI<Name>ServiceImpl 创建 <Name> 服务实现。
func NewI<Name>ServiceImpl(dao aidatasetdao.I<Name>Dao) aidatasetservice.I<Name>Service {
	return &I<Name>ServiceImpl{
		dao: dao,
	}
}

// Create 创建 <Name> 数据。
// 其他方法同样补齐注释后再生成。
```

### DAOImpl 模板
```go
package aiDataSetDaoImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/utils/baizeContext"
)

// I<Name>DaoImpl 提供 <Name> 的数据库访问能力。
type I<Name>DaoImpl struct {
	db    *gorm.DB
	table string
}

// NewI<Name>DaoImpl 创建 <Name> DAO 实现。
func NewI<Name>DaoImpl(db *gorm.DB) aidatasetdao.I<Name>Dao {
	return &I<Name>DaoImpl{
		db:    db,
		table: "ai_<table_name>",
	}
}

// Create 新增 <Name> 记录。
// Update、GetByID、List、DeleteByIDs 同样补齐注释后再生成。
```

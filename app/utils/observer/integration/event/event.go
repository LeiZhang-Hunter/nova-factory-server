// 定义事件体系的基础接口。
// Base 为所有数据载体（如订单、商品、库存）提供通用能力，
// Event 为事件对象提供与配置、缓存及回调交互的统一入口。
// 所有具体事件类型（ProductEvent、StockEvent、OrderEvent）均组合了这两个基础接口。
package event

import (
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Base 事件数据基础接口，所有业务数据载体（ProductData、StockData、OrderData）均需实现。
// 提供元数据访问与指针获取能力，便于观察者在处理过程中传递额外上下文。
type Base interface {
	// Metadata 返回业务数据关联的扩展元数据键值对
	Metadata() map[string]any
	// Ptr 返回当前数据对象的指针，供序列化或深拷贝使用
	Ptr() any
}

// Event 事件基础接口，表示一次业务变更事件。
// 携带集成配置、事件类型、缓存实例及完成回调，观察者根据这些信息执行同步操作。
type Event interface {
	// Config 返回本次事件关联的集成配置，可能为 nil（表示未配置集成）
	Config() config.Config
	// Action 返回事件类型（创建、更新、删除等），观察者据此决定同步策略
	Action() EventType
	// Cache 返回缓存实例，用于在同步过程中读写临时数据（如 OAuth Token）
	GetCache() cache.Cache
	// GetCallback 返回处理完成后的回调接口，用于通知上游同步结果
	GetCallback() Callback
	// GetDB 读取DB
	GetDB() *gorm.DB
	// GetTransaction 是否打开事物
	GetTransaction() bool
	// GetCtx 获取gin.Context
	GetCtx() *gin.Context
}

// TransactionEvent 事务事件。
// T 用于指定转换后的业务事件类型，例如 OrderEvent、ProductEvent 或 StockEvent。
type TransactionEvent[T Event] interface {
	// GetDB 读取DB
	GetDB() *gorm.DB
	// WithDB 设置DB
	WithDB(tx *gorm.DB)
	// ToEvent 转换为具体业务事件
	ToEvent() T
}

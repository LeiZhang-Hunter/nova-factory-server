// Notifier 事件分发器，负责管理所有已注册的 Observer（观察者），
// 在业务事件（商品变更、库存变更、订单变更）发生时依次通知每个观察者。
// 采用全局单例模式，确保全应用共享同一组观察者列表。
package observer

import (
	"errors"
	"nova-factory-server/app/utils/observer/integration/event"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Notifier 事件分发器，管理所有 Observer 观察者并分发业务事件。
// 内部使用读写锁保证并发安全，在通知时先复制观察者列表再迭代，避免锁持有时间过长。
// db 由 Wire 在启动时通过 InitNotifier 注入，是事件分发开启事务的唯一来源。
type Notifier struct {
	mu        sync.RWMutex
	observers []Observer
	db        *gorm.DB
}

var (
	notifierOnce sync.Once
	notifierIns  *Notifier
)

// GetNotifier 获取全局单例事件分发器。
// 首次调用时初始化，后续调用返回同一实例。
func GetNotifier() *Notifier {
	notifierOnce.Do(func() {
		notifierIns = &Notifier{
			observers: make([]Observer, 0),
		}
	})
	return notifierIns
}

// InitNotifier 初始化全局 Notifier 的数据库连接，供 Wire 注入调用。
// 返回初始化后的全局单例，使其能作为 Wire provider 被依赖图驱动调用。
func InitNotifier(db *gorm.DB) *Notifier {
	n := GetNotifier()
	n.mu.Lock()
	n.db = db
	n.mu.Unlock()
	return n
}

// Register 注册一个观察者到分发器。
// 同一类型的 Observer 可注册多次，各适配器通常在 init() 中调用。
func (n *Notifier) Register(obs Observer) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.observers = append(n.observers, obs)
}

// notify 内部事件分发逻辑，通过回调函数决定对每个 Observer 执行的具体操作。
// 先加读锁复制观察者列表，释放锁后再迭代，避免在通知过程中阻塞注册操作。
// 一旦任一观察者返回错误，立即终止后续分发并返回该错误。
func (n *Notifier) notify(fn func(obs Observer) error) error {
	n.mu.RLock()
	observers := make([]Observer, len(n.observers))
	copy(observers, n.observers)
	n.mu.RUnlock()

	for _, ob := range observers {
		err := fn(ob)
		if err != nil {
			return err
		}
	}
	return nil
}

// withTx 统一事务入口：在 Notifier 持有的 db 上开启事务，把 tx 交给 fn。
// db 未初始化时显式报错，绝不静默退化为无事务分发。
func (n *Notifier) withTx(fn func(tx *gorm.DB) error) error {
	n.mu.RLock()
	db := n.db
	n.mu.RUnlock()
	if db == nil {
		return errors.New("notifier: db 未初始化（InitNotifier 未接入 Wire）")
	}
	return db.Transaction(fn)
}

// OnProductChanged 向所有观察者分发商品变更事件。
// 由 Notifier 统一开启事务，并将 tx 显式传给每个观察者，保证全员原子一致。
func (n *Notifier) OnProductChanged(ev event.ProductEvent) error {
	return n.withTx(func(tx *gorm.DB) error {
		return n.notify(func(ob Observer) error {
			_, err := ob.OnProductChanged(tx, ev)
			if err != nil {
				zap.L().Error("Observer OnProductChanged", zap.Error(err))
				return err
			}
			return nil
		})
	})
}

// OnStockChanged 向所有观察者分发库存变更事件。
// 由 Notifier 统一开启事务，并将 tx 显式传给每个观察者，保证全员原子一致。
func (n *Notifier) OnStockChanged(ev event.StockEvent) error {
	return n.withTx(func(tx *gorm.DB) error {
		return n.notify(func(ob Observer) error {
			err := ob.OnStockChanged(tx, ev)
			if err != nil {
				zap.L().Error("Observer OnStockChanged", zap.Error(err))
				return err
			}
			return nil
		})
	})
}

// OnOrderChanged 向所有观察者分发订单变更事件。
// 由 Notifier 统一开启事务，并将 tx 显式传给每个观察者，保证全员原子一致。
func (n *Notifier) OnOrderChanged(ev event.OrderEvent) error {
	return n.withTx(func(tx *gorm.DB) error {
		return n.notify(func(ob Observer) error {
			err := ob.OnOrderChanged(tx, ev)
			if err != nil {
				zap.L().Error("Observer OnOrderChanged", zap.Error(err))
				return err
			}
			return nil
		})
	})
}

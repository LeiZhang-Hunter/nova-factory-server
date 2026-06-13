// Notifier 事件分发器，负责管理所有已注册的 Observer（观察者），
// 在业务事件（商品变更、库存变更、订单变更）发生时依次通知每个观察者。
// 采用全局单例模式，确保全应用共享同一组观察者列表。
package observer

import (
	"nova-factory-server/app/utils/observer/integration/event"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Notifier 事件分发器，管理所有 Observer 观察者并分发业务事件。
// 内部使用读写锁保证并发安全，在通知时先复制观察者列表再迭代，避免锁持有时间过长。
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

// Register 注册一个观察者到分发器。
// 同一类型的 Observer 可注册多次，各适配器通常在 init() 中调用。
func (n *Notifier) Register(obs Observer) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.observers = append(n.observers, obs)
}

func (n *Notifier) SetDB(db *gorm.DB) *Notifier {
	n.db = db
	return n
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

// notifyWithTx 带事务的分发：如果 db 不为 nil，所有 observer 在同一事务中执行；
// 任一 observer 返回 error 则整体回滚。db 为 nil 时退化为普通分发。
func (n *Notifier) notifyWithTx(setDB func(*gorm.DB), fn func(obs Observer) error) error {
	if n.db == nil {
		return n.notify(fn)
	}

	return n.db.Transaction(func(tx *gorm.DB) error {
		setDB(tx)
		return n.notify(fn)
	})
}

// OnProductChanged 向所有观察者分发商品变更事件。
// 如果事件携带 DB，则在事务中执行，保证所有观察者原子一致。
func (n *Notifier) OnProductChanged(ev event.ProductEvent) error {
	return n.notifyWithTx(ev.SetDB, func(ob Observer) error {
		_, err := ob.OnProductChanged(ev)
		if err != nil {
			zap.L().Error("Observer OnProductChanged", zap.Error(err))
			return err
		}
		return nil
	})
}

// OnStockChanged 向所有观察者分发库存变更事件。
// 如果事件携带 DB，则在事务中执行，保证所有观察者原子一致。
func (n *Notifier) OnStockChanged(ev event.StockStoreTransactionEvent) error {
	return n.notifyWithTx(ev.SetDB, func(ob Observer) error {
		err := ob.OnStockChanged(ev.ToEvent())
		if err != nil {
			zap.L().Error("Observer OnStockChanged", zap.Error(err))
			return err
		}
		return nil
	})
}

// OnOrderChanged 向所有观察者分发订单变更事件。
// 如果事件携带 DB，则在事务中执行，保证所有观察者原子一致。
func (n *Notifier) OnOrderChanged(ev event.OrderEvent) error {
	return n.notifyWithTx(ev.SetDB, func(ob Observer) error {
		err := ob.OnOrderChanged(ev)
		if err != nil {
			zap.L().Error("Observer OnOrderChanged", zap.Error(err))
			return err
		}
		return nil
	})
}

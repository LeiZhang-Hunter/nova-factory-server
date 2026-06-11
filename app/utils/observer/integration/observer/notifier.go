package observer

import (
	"go.uber.org/zap"
	"nova-factory-server/app/utils/observer/integration/event"
	"sync"
)

// Notifier 事件分发器，管理所有观察者并异步分发事件
type Notifier struct {
	mu        sync.RWMutex
	observers []Observer
}

var (
	notifierOnce sync.Once
	notifierIns  *Notifier
)

// GetNotifier 获取全局单例事件分发器
func GetNotifier() *Notifier {
	notifierOnce.Do(func() {
		notifierIns = &Notifier{
			observers: make([]Observer, 0),
		}
	})
	return notifierIns
}

// Register 注册观察者
func (n *Notifier) Register(obs Observer) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.observers = append(n.observers, obs)
}

// Notify 异步分发事件，通过回调函数决定调用哪个 Observer 方法
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

// OnProductChanged 商品变更回调
func (n *Notifier) OnProductChanged(event event.ProductEvent) error {
	return n.notify(func(ob Observer) error {
		err := ob.OnProductChanged(event)
		if err != nil {
			zap.L().Error("Observer OnProductChanged", zap.Error(err))
			return err
		}
		return nil
	})
}

// OnStockChanged 库存变更回调
func (n *Notifier) OnStockChanged(event event.StockEvent) error {
	return n.notify(func(ob Observer) error {
		err := ob.OnStockChanged(event)
		if err != nil {
			zap.L().Error("Observer OnProductChanged", zap.Error(err))
			return err
		}
		return nil
	})
}

// OnOrderChanged 订单变更回调
func (n *Notifier) OnOrderChanged(event event.OrderEvent) error {
	return n.notify(func(ob Observer) error {
		err := ob.OnOrderChanged(event)
		if err != nil {
			zap.L().Error("Observer OnProductChanged", zap.Error(err))
			return err
		}
		return nil
	})
}

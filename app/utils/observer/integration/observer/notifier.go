package observer

import (
	"sync"

	"go.uber.org/zap"
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
func (n *Notifier) Notify(fn func(obs Observer)) {
	n.mu.RLock()
	observers := make([]Observer, len(n.observers))
	copy(observers, n.observers)
	n.mu.RUnlock()

	for _, obs := range observers {
		go func(o Observer) {
			defer func() {
				if err := recover(); err != nil {
					zap.L().Error("observer panic", zap.String("name", string(o.Name())), zap.Any("error", err))
				}
			}()
			fn(o)
		}(obs)
	}
}

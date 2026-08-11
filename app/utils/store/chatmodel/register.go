package chatmodel

import "sync"

var (
	once  sync.Once
	lock  sync.Mutex
	store Store
)

// RegisterStore 注册模型连接查询实现，仅首次注册生效。
func RegisterStore(i Store) {
	lock.Lock()
	defer lock.Unlock()
	once.Do(func() {
		store = i
	})
}

// GetStore 返回已注册的实现，未注册时返回空实现。
func GetStore() Store {
	lock.Lock()
	defer lock.Unlock()
	if store == nil {
		return NewEmptyStore()
	}
	return store
}

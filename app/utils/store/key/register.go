package key

import "sync"

var once sync.Once
var lock sync.Mutex
var store keys

func RegisterStore(i keys) {
	lock.Lock()
	defer lock.Unlock()
	once.Do(func() {
		store = i
	})
	return
}

func GetStore() keys {
	lock.Lock()
	defer lock.Unlock()
	if store == nil {
		return newEmptyPermissions()
	}
	return store
}

package requestlog

import "sync"

var store Store
var once sync.Once
var lock sync.Mutex

func RegisterStore(s Store) {
	lock.Lock()
	defer lock.Unlock()
	once.Do(func() {
		store = s
	})
}

func GetStore() Store {
	lock.Lock()
	defer lock.Unlock()
	if store == nil {
		return NewEmptyStore()
	}
	return store
}

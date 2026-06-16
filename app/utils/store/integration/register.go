package integration

import "sync"

var store Integration
var once sync.Once
var lock sync.Mutex

func RegisterStore(i Integration) {
	lock.Lock()
	defer lock.Unlock()
	once.Do(func() {
		store = i
	})
	return
}

func GetStore() Integration {
	lock.Lock()
	defer lock.Unlock()
	if store == nil {
		return NewEmptyIntegrationStore()
	}
	return store
}

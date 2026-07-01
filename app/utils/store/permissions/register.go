package permissions

import "sync"

var once sync.Once
var lock sync.Mutex
var store Permissions

func RegisterStore(i Permissions) {
	lock.Lock()
	defer lock.Unlock()
	once.Do(func() {
		store = i
	})
	return
}

func GetStore() Permissions {
	lock.Lock()
	defer lock.Unlock()
	if store == nil {
		return newEmptyPermissions()
	}
	return store
}

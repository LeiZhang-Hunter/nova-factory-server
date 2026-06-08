package store

import "sync"

var store map[string]IShopCategoryStore
var once sync.Once
var lock sync.Mutex

func init() {
	once.Do(func() {
		store = make(map[string]IShopCategoryStore)
	})
}

func RegisterStore(storeType string, categoryStore IShopCategoryStore) {
	lock.Lock()
	defer lock.Unlock()
	store[storeType] = categoryStore
	return
}

func GetStore(storeType string) IShopCategoryStore {
	lock.Lock()
	defer lock.Unlock()
	storeData, ok := store[storeType]
	if !ok {
		return nil
	}
	return storeData
}

package goods

import "sync"

var store IShopGoodsStore
var once sync.Once
var lock sync.Mutex

func RegisterStore(shoppGoodsStore IShopGoodsStore) {
	lock.Lock()
	defer lock.Unlock()
	once.Do(func() {
		store = shoppGoodsStore
	})
	return
}

func GetStore() IShopGoodsStore {
	lock.Lock()
	defer lock.Unlock()
	if store == nil {
		return NewEmptyIShopGoodsStore()
	}
	return store
}

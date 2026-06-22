package embedder

import "sync"

var once sync.Once
var lock sync.Mutex
var store Embedder

func RegisterStore(i Embedder) {
	lock.Lock()
	defer lock.Unlock()
	once.Do(func() {
		store = i
	})
	return
}

func GetStore() Embedder {
	lock.Lock()
	defer lock.Unlock()
	if store == nil {
		return NewEmptyEmbedder()
	}
	return store
}

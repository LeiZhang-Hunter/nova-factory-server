package callback

import (
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/result"
)

type GoodsCallback struct {
}

func NewGoodsCallback() *GoodsCallback {
	return &GoodsCallback{}
}

func (callback *GoodsCallback) OnSuccess(T event.Event, response result.SyncProductResponse) {
	return
}

func (callback *GoodsCallback) OnError(T event.Event, response result.SyncProductResponse, err error) {
	//TODO implement me
	return
}

func (callback *GoodsCallback) OnFinish(T event.Event) {
	//TODO implement me
	return
}

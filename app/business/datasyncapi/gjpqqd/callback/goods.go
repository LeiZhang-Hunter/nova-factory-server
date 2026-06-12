package callback

import "nova-factory-server/app/utils/observer/integration/event"

type GoodsCallback struct {
}

func NewGoodsCallback() event.Callback {
	return &GoodsCallback{}
}

func (callback *GoodsCallback) OnError() {
	//TODO implement me
	panic("implement me")
}

func (callback *GoodsCallback) OnSuccess() {

}

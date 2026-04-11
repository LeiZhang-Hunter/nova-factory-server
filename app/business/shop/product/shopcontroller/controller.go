package shopcontroller

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewCategory, NewGoods, NewSku, wire.Struct(new(Controller), "*"))

type Controller struct {
	Category *Category
	Goods    *Goods
	Sku      *Sku
}

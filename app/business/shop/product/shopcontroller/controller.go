package shopcontroller

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewAddress, NewCategory, NewGoods, NewSku, wire.Struct(new(Controller), "*"))

type Controller struct {
	Address  *Address
	Category *Category
	Goods    *Goods
	Sku      *Sku
}

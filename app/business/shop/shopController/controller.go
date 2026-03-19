package shopController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewCategory, NewGoods, NewSku, NewUser, wire.Struct(new(Controller), "*"))

type Controller struct {
	Category *Category
	Goods    *Goods
	Sku      *Sku
	User     *User
}

package product

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewProduct, NewCategory, NewCart, NewHome, wire.Struct(new(Controller), "*"))

type Controller struct {
	Product  *Product
	Category *Category
	Home     *Home
	Cart     *Cart
}

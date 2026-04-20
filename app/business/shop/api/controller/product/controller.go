package product

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewProduct, NewCategory, wire.Struct(new(Controller), "*"))

type Controller struct {
	Product  *Product
	Category *Category
}

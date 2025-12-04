package productController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewLaboratory, wire.Struct(new(Product), "*"))

type Product struct {
	Laboratory *Laboratory
}

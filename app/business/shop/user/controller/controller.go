package shopcontroller

import "github.com/google/wire"

type Controller struct {
	Address      *Address
	Cart         *Cart
	User         *User
	DiscountRule *DiscountRule
}

var ProviderSet = wire.NewSet(NewAddress, NewCart, NewUser, NewDiscountRule, wire.Struct(new(Controller), "*"))

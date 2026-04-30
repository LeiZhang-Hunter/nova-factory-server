package controller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewCombination, NewPink, NewSeckill, NewSeckillActivity, wire.Struct(new(Controller), "*"))

type Controller struct {
	Combination     *Combination
	Pink            *Pink
	Seckill         *Seckill
	SeckillActivity *SeckillActivity
}

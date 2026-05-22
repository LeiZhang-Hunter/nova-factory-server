package activity

import "github.com/google/wire"

var ProviderSet = wire.NewSet(wire.Struct(new(SeckillController), "*"), wire.Struct(new(CombinationController), "*"), wire.Struct(new(PinkController), "*"), wire.Struct(new(Controller), "*"))

type Controller struct {
	Seckill     *SeckillController
	Combination *CombinationController
	Pink        *PinkController
}

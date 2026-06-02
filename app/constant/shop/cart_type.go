package shop

const (
	CartProductTypeNormal      int32 = 0
	CartProductTypeSeckill     int32 = 1
	CartProductTypeCombination int32 = 3
)

const (
	CartModeCart   = "cart"
	CartModeBuyNow = "buyNow"
)

const (
	CartStateNormal  int32 = 0
	CartStateDeleted int32 = -1
	CartStateBuyNow  int32 = 10
	CartStateOrdered int32 = 20
)

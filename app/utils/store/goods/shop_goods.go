package goods

type IShopGoodsStore interface {
	GetProductList(request Request) DataResult
}

type EmptyIShopGoodsStore struct{}

func NewEmptyIShopGoodsStore() IShopGoodsStore {
	return &EmptyIShopGoodsStore{}
}

func (*EmptyIShopGoodsStore) GetProductList(request Request) DataResult {
	return &EmptyDataResult{}
}

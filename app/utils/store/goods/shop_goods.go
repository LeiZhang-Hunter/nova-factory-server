package goods

type IShopGoodsStore interface {
	GetProductList(request Request) DataResult
	GetProductCategory(DataCategoryRequest) DataCategoryResult
}

type EmptyIShopGoodsStore struct{}

func NewEmptyIShopGoodsStore() IShopGoodsStore {
	return &EmptyIShopGoodsStore{}
}

func (*EmptyIShopGoodsStore) GetProductList(request Request) DataResult {
	return &EmptyDataResult{}
}

func (*EmptyIShopGoodsStore) GetProductCategory(DataCategoryRequest) DataCategoryResult {
	return &EmptyCategoryDataResult{}
}

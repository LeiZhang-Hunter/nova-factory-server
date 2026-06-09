package store

type ShopCategoryData interface {
	CategoryID() int64
	ChildrenData() []ShopCategoryData
	SetChildren([]ShopCategoryData) error
	Name() string
}

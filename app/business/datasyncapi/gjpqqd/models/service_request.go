package models

// ProductListRequest 商品列表查询请求参数
type ProductListRequest struct {
	PageNo   int
	PageSize int
}

// ProductStockUpdateRequest 库存更新请求参数
// 支持按 skuid:qty 形式批量更新多个 SKU 的库存
type ProductStockUpdateRequest struct {
	ProductID  string
	ProductQty string
	Skus       string
}

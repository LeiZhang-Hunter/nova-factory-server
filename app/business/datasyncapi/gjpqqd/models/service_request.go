package models

import (
	"fmt"
	"strconv"
	"strings"
)

// ProductListRequest 商品列表查询请求参数
type ProductListRequest struct {
	PageNo   int
	PageSize int
}

// ProductStockUpdateRequest 库存更新请求参数
// productid: 商品ID（必填）
// productqty: 商品库存数量（无规格时必填，有规格时为空）
// skus: 规格库存，格式 "skuID:qty,skuID:qty"（无规格时为空）
type ProductStockUpdateRequest struct {
	ProductID  string `form:"productid" json:"productid"`
	ProductQty any    `form:"productqty" json:"productqty"`
	Skus       string `form:"skus" json:"skus"`
}

// ToStockSyncReq 将请求参数转换为 StockSyncReq 事件载体
func (r *ProductStockUpdateRequest) ToStockSyncReq() *StockSyncReq {
	productID, _ := strconv.ParseInt(r.ProductID, 10, 64)

	var items []StockItem

	if r.Skus != "" {
		pairs := strings.Split(r.Skus, ",")
		for _, pair := range pairs {
			pair = strings.TrimSpace(pair)
			if pair == "" {
				continue
			}
			parts := strings.SplitN(pair, ":", 2)
			if len(parts) != 2 {
				continue
			}
			skuID, _ := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
			qty, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			items = append(items, StockItem{
				ProductId:     productID,
				SkuId:         skuID,
				AfterQuantity: qty,
			})
		}
	} else {
		qty, _ := strconv.ParseFloat(fmt.Sprintf("%v", r.ProductQty), 64)
		items = append(items, StockItem{
			ProductId:     productID,
			AfterQuantity: qty,
		})
	}

	return &StockSyncReq{Stocks: items}
}

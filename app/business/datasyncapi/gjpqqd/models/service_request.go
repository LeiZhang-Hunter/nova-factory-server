package models

import (
	"fmt"
	"strconv"
	"strings"
)

// ProductListRequest 商品列表查询请求，实现 Request 接口。
type ProductListRequest struct {
	SellerCIds    string `form:"sellercids"`
	StartModified string `form:"startmodified"`
	EndModified   string `form:"endmodified"`
	Name          string `form:"name"`
	OuterId       string `form:"outerid"`
	Status        string `form:"status"`
	PageNo        int    `form:"pageno"`
	PageSize      int    `form:"pagesize"`
}

func (r *ProductListRequest) GetSellerCIds() string    { return r.SellerCIds }
func (r *ProductListRequest) GetStartModified() string { return r.StartModified }
func (r *ProductListRequest) GetEndModified() string   { return r.EndModified }
func (r *ProductListRequest) GetName() string          { return r.Name }
func (r *ProductListRequest) GetOuterId() string       { return r.OuterId }
func (r *ProductListRequest) GetStatus() string        { return r.Status }
func (r *ProductListRequest) GetPageNo() int           { return r.PageNo }
func (r *ProductListRequest) GetPageSize() int         { return r.PageSize }

// ProductStockUpdateRequest 库存更新请求参数
// productid: 商品ID（必填）
// productqty: 商品库存数量（无规格时必填，有规格时为空）
// skus: 规格库存，格式 "skuID:qty,skuID:qty"（无规格时为空）
type ProductStockUpdateRequest struct {
	ProductID  string `form:"productid" json:"productid" binding:"required"`
	ProductQty int64  `form:"productqty" json:"productqty"`
	Skus       string `form:"skus" json:"skus"`
}

// ToStockSyncReq 将请求参数转换为 StockSyncReq 事件载体
func (r *ProductStockUpdateRequest) ToStockSyncReq() *StockSyncReq {

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
				ProductId:     r.ProductID,
				SkuId:         skuID,
				AfterQuantity: qty,
			})
		}
	} else {
		qty, _ := strconv.ParseFloat(fmt.Sprintf("%v", r.ProductQty), 64)
		items = append(items, StockItem{
			ProductId:     r.ProductID,
			AfterQuantity: qty,
		})
	}

	return &StockSyncReq{Stocks: items}
}

type CategorySearchRequest struct {
	Parentcid string `json:"parentcid"`
}

func (r *CategorySearchRequest) GetParentcid() string { return r.Parentcid }

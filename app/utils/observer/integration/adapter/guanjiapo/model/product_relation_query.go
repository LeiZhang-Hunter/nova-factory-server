package model

import (
	"nova-factory-server/app/utils/observer/integration/result"
	"strconv"
)

// ProductRelationQueryData 商品对应关系数据，实现 result.ProductRelationQueryData。
type ProductRelationQueryData struct {
	EshopGoodsID string `json:"eshopgoodsid"`
	EshopSkuID   string `json:"eshopskuid"`
	GoodsID      string `json:"goodsid"`
	GoodsCode    string `json:"goodscode"`
	GoodsName    string `json:"goodsname"`
	SkuID        string `json:"skuid"`
	SkuCode      string `json:"skucode"`
	SkuName      string `json:"skuname"`
	UnitID       int64  `json:"unitid"`
	UnitName     string `json:"unitname"`
}

func (d *ProductRelationQueryData) GetEshopGoodsID() string { return d.EshopGoodsID }
func (d *ProductRelationQueryData) GetEshopSkuID() string   { return d.EshopSkuID }
func (d *ProductRelationQueryData) GetGoodsID() string      { return d.GoodsID }
func (d *ProductRelationQueryData) GetGoodsCode() string    { return d.GoodsCode }
func (d *ProductRelationQueryData) GetGoodsName() string    { return d.GoodsName }
func (d *ProductRelationQueryData) GetSkuID() int64 {
	parseInt, err := strconv.ParseInt(d.SkuID, 10, 64)
	if err != nil {
		return 0
	}
	return parseInt
}
func (d *ProductRelationQueryData) GetSkuCode() string  { return d.SkuCode }
func (d *ProductRelationQueryData) GetSkuName() string  { return d.SkuName }
func (d *ProductRelationQueryData) GetUnitID() int64    { return d.UnitID }
func (d *ProductRelationQueryData) GetUnitName() string { return d.UnitName }

// ProductRelationQueryResponse 商品对应关系查询响应，实现 result.ProductRelationQueryResponse。
type ProductRelationQueryResponse struct {
	Code          int64                       `json:"code"`
	Message       string                      `json:"message"`
	Total         int64                       `json:"total"`
	GoodsRelation []*ProductRelationQueryData `json:"goodsrelation"`
}

func (r *ProductRelationQueryResponse) GetCode() int64     { return r.Code }
func (r *ProductRelationQueryResponse) GetMessage() string { return r.Message }
func (r *ProductRelationQueryResponse) GetTotal() int64    { return r.Total }
func (r *ProductRelationQueryResponse) GetGoodsRelation() []result.ProductRelationQueryData {
	out := make([]result.ProductRelationQueryData, len(r.GoodsRelation))
	for i, v := range r.GoodsRelation {
		out[i] = v
	}
	return out
}

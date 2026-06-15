package model

import (
	"nova-factory-server/app/utils/observer/integration/result"
)

// ProductRemarkUpdateResult 单条商品备注更新结果，实现 result.ProductRemarkUpdateResult。
type ProductRemarkUpdateResult struct {
	GoodsID int64  `json:"goodsid"`
	Message string `json:"message"`
}

func (r *ProductRemarkUpdateResult) GetGoodsID() int64  { return r.GoodsID }
func (r *ProductRemarkUpdateResult) GetMessage() string { return r.Message }

// ProductRemarkUpdateResponse 商品备注更新响应，实现 result.ProductRemarkUpdateResponse。
type ProductRemarkUpdateResponse struct {
	Code    int64                        `json:"code"`
	Message string                       `json:"message"`
	Items   []*ProductRemarkUpdateResult `json:"items"`
}

func (p *ProductRemarkUpdateResponse) GetCode() int64     { return p.Code }
func (p *ProductRemarkUpdateResponse) GetMessage() string { return p.Message }
func (p *ProductRemarkUpdateResponse) GetItems() []result.ProductRemarkUpdateResult {
	items := make([]result.ProductRemarkUpdateResult, 0, len(p.Items))
	for _, item := range p.Items {
		items = append(items, item)
	}
	return items
}

// ProductSearchResponseSku SKU数据，实现 result.GoodsGetResponseDataSku。
type ProductSearchResponseSku struct {
	SkuID    int    `json:"skuid"`
	SkuCode  string `json:"skucode"`
	SkuName  string `json:"skuname"`
	Barcode  string `json:"barcode"`
	LCMCCode string `json:"lcmccode"`
	Weight   int    `json:"weight"`
	Size     int    `json:"size"`
	Price    int    `json:"price"`
	Price2   int    `json:"price2"`
	Price3   int    `json:"price3"`
	Price4   int    `json:"price4"`
	Price5   int    `json:"price5"`
}

func (s *ProductSearchResponseSku) GetSkuid() int       { return s.SkuID }
func (s *ProductSearchResponseSku) GetSkucode() string  { return s.SkuCode }
func (s *ProductSearchResponseSku) GetSkuname() string  { return s.SkuName }
func (s *ProductSearchResponseSku) GetBarcode() string  { return s.Barcode }
func (s *ProductSearchResponseSku) GetLcmccode() string { return s.LCMCCode }
func (s *ProductSearchResponseSku) GetWeight() int      { return s.Weight }
func (s *ProductSearchResponseSku) GetSize() int        { return s.Size }
func (s *ProductSearchResponseSku) GetPrice() int       { return s.Price }
func (s *ProductSearchResponseSku) GetPrice2() int      { return s.Price2 }
func (s *ProductSearchResponseSku) GetPrice3() int      { return s.Price3 }
func (s *ProductSearchResponseSku) GetPrice4() int      { return s.Price4 }
func (s *ProductSearchResponseSku) GetPrice5() int      { return s.Price5 }

// ProductSearchResponseUnit 单位数据，实现 result.GoodsGetResponseDataUnit。
type ProductSearchResponseUnit struct {
	UnitName string `json:"unitname"`
	Barcode  string `json:"barcode"`
	Rate     int    `json:"rate"`
	Price    int    `json:"price"`
	Price2   int    `json:"price2"`
	Price3   int    `json:"price3"`
	Price4   int    `json:"price4"`
	Price5   int    `json:"price5"`
}

func (u *ProductSearchResponseUnit) GetUnitname() string { return u.UnitName }
func (u *ProductSearchResponseUnit) GetBarcode() string  { return u.Barcode }
func (u *ProductSearchResponseUnit) GetRate() int        { return u.Rate }
func (u *ProductSearchResponseUnit) GetPrice() int       { return u.Price }
func (u *ProductSearchResponseUnit) GetPrice2() int      { return u.Price2 }
func (u *ProductSearchResponseUnit) GetPrice3() int      { return u.Price3 }
func (u *ProductSearchResponseUnit) GetPrice4() int      { return u.Price4 }
func (u *ProductSearchResponseUnit) GetPrice5() int      { return u.Price5 }

// ProductSearchResponseData 商品数据，实现 result.GoodsGetResponseData。
type ProductSearchResponseData struct {
	GoodsID      int                          `json:"goodsid"`
	EshopGoodsID string                       `json:"eshopgoodsid"`
	GoodsCode    string                       `json:"goodscode"`
	GoodsName    string                       `json:"goodsname"`
	Unit         string                       `json:"unit"`
	Remark       string                       `json:"remark"`
	Skus         []*ProductSearchResponseSku  `json:"skus"`
	Units        []*ProductSearchResponseUnit `json:"units"`
}

func (d *ProductSearchResponseData) GetGoodsid() int         { return d.GoodsID }
func (d *ProductSearchResponseData) GetEshopgoodsid() string { return d.EshopGoodsID }
func (d *ProductSearchResponseData) GetGoodscode() string    { return d.GoodsCode }
func (d *ProductSearchResponseData) GetGoodsname() string    { return d.GoodsName }
func (d *ProductSearchResponseData) GetUnit() string         { return d.Unit }
func (d *ProductSearchResponseData) GetRemark() string       { return d.Remark }
func (d *ProductSearchResponseData) GetSkus() []result.GoodsGetResponseDataSku {
	skus := make([]result.GoodsGetResponseDataSku, 0, len(d.Skus))
	for _, s := range d.Skus {
		skus = append(skus, s)
	}
	return skus
}
func (d *ProductSearchResponseData) GetUnits() []result.GoodsGetResponseDataUnit {
	units := make([]result.GoodsGetResponseDataUnit, 0, len(d.Units))
	for _, u := range d.Units {
		units = append(units, u)
	}
	return units
}

// ProductSearchResponse 商品查询响应，实现 result.GoodsGetResponse。
type ProductSearchResponse struct {
	RespCode    int                          `json:"code"`
	RespMessage string                       `json:"message"`
	RespTotal   int                          `json:"total"`
	RespGoods   []*ProductSearchResponseData `json:"goods"`
}

func (r *ProductSearchResponse) GetCode() int       { return r.RespCode }
func (r *ProductSearchResponse) GetMessage() string { return r.RespMessage }
func (r *ProductSearchResponse) GetTotal() int      { return r.RespTotal }
func (r *ProductSearchResponse) GetGoods() []result.GoodsGetResponseData {
	goods := make([]result.GoodsGetResponseData, 0, len(r.RespGoods))
	for _, g := range r.RespGoods {
		goods = append(goods, g)
	}
	return goods
}

package model

import "nova-factory-server/app/utils/observer/integration/result"

// StockGetData 单条库存数据，实现 result.StockGetResponseData。
type StockGetData struct {
	GoodsID   int64   `json:"goodsid"`
	GoodsCode string  `json:"goodscode"`
	GoodsName string  `json:"goodsname"`
	SkuID     string  `json:"skuid"`
	SkuCode   string  `json:"skucode"`
	SkuName   string  `json:"skuname"`
	Qty       float64 `json:"qty"`
	EnableNum float64 `json:"enablenum"`
	EnSaleNum float64 `json:"ensalenum"`
	WhsCode   string  `json:"whscode"`
	WhsName   string  `json:"whsname"`
}

func (d *StockGetData) GetGoodsid() int64     { return d.GoodsID }
func (d *StockGetData) GetGoodscode() string  { return d.GoodsCode }
func (d *StockGetData) GetGoodsname() string  { return d.GoodsName }
func (d *StockGetData) GetSkuid() string      { return d.SkuID }
func (d *StockGetData) GetSkucode() string    { return d.SkuCode }
func (d *StockGetData) GetSkuname() string    { return d.SkuName }
func (d *StockGetData) GetQty() float64       { return d.Qty }
func (d *StockGetData) GetEnableNum() float64 { return d.EnableNum }
func (d *StockGetData) GetEnSaleNum() float64 { return d.EnSaleNum }
func (d *StockGetData) GetWhsCode() string    { return d.WhsCode }
func (d *StockGetData) GetWhsName() string    { return d.WhsName }

// StockGetResponse 库存查询响应，实现 result.StockGetResponse。
type StockGetResponse struct {
	Code    int64           `json:"code"`
	Message string          `json:"message"`
	Total   int64           `json:"total"`
	Stocks  []*StockGetData `json:"stocks"`
}

func (r *StockGetResponse) GetCode() int64     { return r.Code }
func (r *StockGetResponse) GetMessage() string { return r.Message }
func (r *StockGetResponse) GetTotal() int64    { return r.Total }
func (r *StockGetResponse) GetStocks() []result.StockGetResponseData {
	out := make([]result.StockGetResponseData, len(r.Stocks))
	for i, v := range r.Stocks {
		out[i] = v
	}
	return out
}

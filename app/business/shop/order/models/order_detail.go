package models

import (
	"fmt"
	"nova-factory-server/app/baize"
	searchutil "nova-factory-server/app/utils/vectorsearch"
	"strings"
)

// OrderDetail ERP订单明细
type OrderDetail struct {
	ID                uint64  `json:"id,string" gorm:"column:id"`
	OrderID           uint64  `json:"order_id,string" gorm:"column:order_id"`
	Tid               string  `json:"tid" gorm:"column:tid"`
	OID               string  `json:"oid" gorm:"column:oid"`
	Barcode           string  `json:"barcode" gorm:"column:barcode"`
	OldBarcode        string  `json:"old_barcode" gorm:"-"`
	EShopGoodsID      string  `json:"eshop_goods_id" gorm:"column:eshop_goods_id"`
	OuterIID          string  `json:"outer_iid" gorm:"column:outer_iid"`
	EShopGoodsName    string  `json:"eshop_goods_name" gorm:"column:eshop_goods_name"`
	OldEShopGoodsName string  `json:"old_eshop_goods_name" gorm:"-"`
	EShopSkuID        string  `json:"eshop_sku_id" gorm:"column:eshop_sku_id"`
	OldEShopSkuID     string  `json:"old_eshop_sku_id" gorm:"-"`
	EShopSkuName      string  `json:"eshop_sku_name" gorm:"column:eshop_sku_name"`
	OldEShopSkuName   string  `json:"old_eshop_sku_name" gorm:"-"`
	NumIID            int64   `json:"num_iid" gorm:"column:num_iid"`
	SkuID             int64   `json:"sku_id" gorm:"column:sku_id"`
	Num               float64 `json:"num" gorm:"column:num"`
	Payment           float64 `json:"payment" gorm:"column:payment"`
	OldPayment        float64 `json:"old_payment" gorm:"-"`
	PicPath           string  `json:"pic_path" gorm:"column:pic_path"`
	Weight            float64 `json:"weight" gorm:"column:weight"`
	Size              float64 `json:"size" gorm:"column:size"`
	UnitID            int64   `json:"unit_id" gorm:"column:unit_id"`
	UnitQty           float64 `json:"unit_qty" gorm:"column:unit_qty"`
	DeptID            int64   `json:"dept_id" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// VectorSearchLabeledValues 返回订单明细可参与向量检索的结构化文本字段。
func (d *OrderDetail) VectorSearchLabeledValues() []searchutil.LabeledValue {
	if d == nil {
		return nil
	}
	values := make([]searchutil.LabeledValue, 0, 12)

	productCode := strings.TrimSpace(d.OuterIID)
	if productCode == "" {
		productCode = strings.TrimSpace(d.EShopGoodsID)
	}

	unitParts := make([]string, 0, 2)
	if d.UnitID > 0 {
		unitParts = append(unitParts, fmt.Sprintf("单位ID:%d", d.UnitID))
	}
	if d.UnitQty > 0 {
		unitParts = append(unitParts, fmt.Sprintf("单位数量:%.3f", d.UnitQty))
	}

	remarkParts := make([]string, 0, 6)
	if goodsID := strings.TrimSpace(d.EShopGoodsID); goodsID != "" {
		remarkParts = append(remarkParts, "商品ID:"+goodsID)
	}
	if skuID := strings.TrimSpace(d.EShopSkuID); skuID != "" {
		remarkParts = append(remarkParts, "SKU ID:"+skuID)
	}
	if d.NumIID > 0 {
		remarkParts = append(remarkParts, fmt.Sprintf("商品数字ID:%d", d.NumIID))
	}
	if d.SkuID > 0 {
		remarkParts = append(remarkParts, fmt.Sprintf("SKU数字ID:%d", d.SkuID))
	}
	if d.Size > 0 {
		remarkParts = append(remarkParts, fmt.Sprintf("体积:%.3f", d.Size))
	}

	values = append(values,
		searchutil.LabeledValue{Label: "产品名称", Value: d.EShopGoodsName},
		searchutil.LabeledValue{Label: "产品编码", Value: productCode},
		searchutil.LabeledValue{Label: "产品分类", Value: ""},
		searchutil.LabeledValue{Label: "单位", Value: strings.Join(unitParts, " ")},
		searchutil.LabeledValue{Label: "条码", Value: d.Barcode},
		searchutil.LabeledValue{Label: "规格", Value: d.EShopSkuName},
		searchutil.LabeledValue{Label: "备注", Value: strings.Join(remarkParts, " ")},
	)
	if d.Weight > 0 {
		values = append(values, searchutil.LabeledValue{Label: "重量", Value: fmt.Sprintf("%.3fkg", d.Weight)})
	}
	if d.Payment > 0 {
		values = append(values, searchutil.LabeledValue{Label: "销售价", Value: fmt.Sprintf("%.2f", d.Payment)})
	}
	return values
}

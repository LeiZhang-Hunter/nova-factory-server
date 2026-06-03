package mastermodels

import (
	"fmt"
	"nova-factory-server/app/baize"
	searchutil "nova-factory-server/app/utils/vectorsearch"
)

// Product ERP 产品
type Product struct {
	ID            int64   `json:"id,string" gorm:"column:id"`
	Name          string  `json:"name" gorm:"column:name"`
	ProductCode   string  `json:"product_code" gorm:"column:product_code"`
	BarCode       string  `json:"barCode" gorm:"column:bar_code"`
	CategoryId    int64   `json:"categoryId" gorm:"column:category_id"`
	CategoryName  string  `json:"categoryName" gorm:"-"`
	UnitId        int64   `json:"unitId,string" gorm:"column:unit_id"`
	UnitName      string  `json:"unitName" gorm:"-"`
	Status        int32   `json:"status" gorm:"column:status"`
	Standard      string  `json:"standard" gorm:"column:standard"`
	Remark        string  `json:"remark" gorm:"column:remark"`
	ExpiryDay     int32   `json:"expiryDay" gorm:"column:expiry_day"`
	Weight        float64 `json:"weight" gorm:"column:weight"`
	PurchasePrice float64 `json:"purchasePrice" gorm:"column:purchase_price"`
	SalePrice     float64 `json:"salePrice" gorm:"column:sale_price"`
	MinPrice      float64 `json:"minPrice" gorm:"column:min_price"`
	DeptID        int64   `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// VectorSearchLabeledValues 返回产品可参与向量生成的结构化文本字段。
func (p *Product) VectorSearchLabeledValues() []searchutil.LabeledValue {
	if p == nil {
		return nil
	}
	values := make([]searchutil.LabeledValue, 0, 12)
	values = append(values,
		searchutil.LabeledValue{Label: "产品名称", Value: p.Name},
		searchutil.LabeledValue{Label: "产品编码", Value: p.ProductCode},
		searchutil.LabeledValue{Label: "产品分类", Value: p.CategoryName},
		searchutil.LabeledValue{Label: "单位", Value: p.UnitName},
		searchutil.LabeledValue{Label: "条码", Value: p.BarCode},
		searchutil.LabeledValue{Label: "规格", Value: p.Standard},
		searchutil.LabeledValue{Label: "备注", Value: p.Remark},
	)
	if p.ExpiryDay > 0 {
		values = append(values, searchutil.LabeledValue{Label: "保质期", Value: fmt.Sprintf("%d天", p.ExpiryDay)})
	}
	if p.Weight > 0 {
		values = append(values, searchutil.LabeledValue{Label: "重量", Value: fmt.Sprintf("%.3fkg", p.Weight)})
	}
	if p.PurchasePrice > 0 {
		values = append(values, searchutil.LabeledValue{Label: "采购价", Value: fmt.Sprintf("%.2f", p.PurchasePrice)})
	}
	if p.SalePrice > 0 {
		values = append(values, searchutil.LabeledValue{Label: "销售价", Value: fmt.Sprintf("%.2f", p.SalePrice)})
	}
	if p.MinPrice > 0 {
		values = append(values, searchutil.LabeledValue{Label: "最低价", Value: fmt.Sprintf("%.2f", p.MinPrice)})
	}
	return values
}

// ProductUpsert ERP 产品新增修改参数
type ProductUpsert struct {
	ID            int64   `json:"id,string"`
	Name          string  `json:"name" binding:"required" label:"产品名称"`
	BarCode       string  `json:"barCode"`
	CategoryId    int64   `json:"categoryId"`
	UnitId        int64   `json:"unitId,string"`
	Status        int32   `json:"status"`
	Standard      string  `json:"standard"`
	Remark        string  `json:"remark"`
	ExpiryDay     int32   `json:"expiryDay"`
	Weight        float64 `json:"weight"`
	PurchasePrice float64 `json:"purchasePrice"`
	SalePrice     float64 `json:"salePrice"`
	MinPrice      float64 `json:"minPrice"`
}

// ProductQuery ERP 产品查询参数
type ProductQuery struct {
	Name   string `form:"name" filter:"like,name"`
	Code   string `form:"code"`
	Status *int32 `form:"status" filter:"eq,status"`
	Page   int64  `form:"page"`
	Size   int64  `form:"size"`
}

// ProductListData ERP 产品分页数据
type ProductListData struct {
	Rows  []*Product `json:"rows"`
	Total int64      `json:"total"`
}

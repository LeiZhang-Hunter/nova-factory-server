package craftRouteModels

import (
	"nova-factory-server/app/baize"
)

// SysProRouteProduct 产品制程
type SysProRouteProduct struct {
	RecordID       int64   `gorm:"column:record_id;primaryKey;autoIncrement:true;comment:记录ID" json:"record_id"` // 记录ID
	RouteID        int64   `gorm:"column:route_id;not null;comment:工艺路线ID" json:"route_id"`                      // 工艺路线ID
	ItemID         int64   `gorm:"column:item_id;not null;comment:产品物料ID" json:"item_id"`                        // 产品物料ID
	ItemCode       string  `gorm:"column:item_code;not null;comment:产品物料编码" json:"item_code"`                    // 产品物料编码
	ItemName       string  `gorm:"column:item_name;not null;comment:产品物料名称" json:"item_name"`                    // 产品物料名称
	Specification  string  `gorm:"column:specification;comment:规格型号" json:"specification"`                       // 规格型号
	UnitOfMeasure  string  `gorm:"column:unit_of_measure;not null;comment:单位" json:"unit_of_measure"`            // 单位
	UnitName       string  `gorm:"column:unit_name;comment:单位名称" json:"unit_name"`                               // 单位名称
	Quantity       int32   `gorm:"column:quantity;default:1;comment:生产数量" json:"quantity"`                       // 生产数量
	ProductionTime float64 `gorm:"column:production_time;default:1.00;comment:生产用时" json:"production_time"`      // 生产用时
	TimeUnitType   string  `gorm:"column:time_unit_type;default:MINUTE;comment:时间单位" json:"time_unit_type"`      // 时间单位
	Remark         string  `gorm:"column:remark;comment:备注" json:"remark"`                                       // 备注
	Attr1          string  `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                      // 预留字段1
	Attr2          string  `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                      // 预留字段2
	Attr3          int32   `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                      // 预留字段3
	Attr4          int32   `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                      // 预留字段4
	DeptID         int64   `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                   // 部门ID
	State          bool    `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                             // 操作状态（0正常 -1删除）
	baize.BaseEntity
}

func NewSysProRouteProduct(req *SysProRouteSetProduct) *SysProRouteProduct {
	return &SysProRouteProduct{
		RecordID:       req.RecordID,
		RouteID:        req.RouteID,
		ItemID:         req.ItemID,
		ItemCode:       req.ItemCode,
		ItemName:       req.ItemName,
		Specification:  req.Specification,
		UnitOfMeasure:  req.UnitOfMeasure,
		UnitName:       req.UnitName,
		Quantity:       req.Quantity,
		ProductionTime: req.ProductionTime,
		TimeUnitType:   req.TimeUnitType,
		Remark:         req.Remark,
		Attr1:          req.Attr1,
		Attr2:          req.Attr2,
		Attr3:          req.Attr3,
		Attr4:          req.Attr4,
	}
}

type SysProRouteProductReq struct {
	ItemCode string `gorm:"column:item_code;not null;comment:产品物料编码" form:"item_code" json:"item_code"` // 产品物料编码
	ItemName string `gorm:"column:item_name;not null;comment:产品物料名称" form:"item_name" json:"item_name"` // 产品物料名称
	baize.BaseEntityDQL
}

type SysProRouteProductList struct {
	Rows  []*SysProRouteProduct `json:"rows"`
	Total int64                 `json:"total"`
}

type SysProRouteSetProduct struct {
	RecordID       int64   `gorm:"column:record_id;primaryKey;autoIncrement:true;comment:记录ID" json:"record_id"` // 记录ID
	RouteID        int64   `gorm:"column:route_id;not null;comment:工艺路线ID" json:"route_id"`                      // 工艺路线ID
	ItemID         int64   `gorm:"column:item_id;not null;comment:产品物料ID" json:"item_id"`                        // 产品物料ID
	ItemCode       string  `gorm:"column:item_code;not null;comment:产品物料编码" json:"item_code"`                    // 产品物料编码
	ItemName       string  `gorm:"column:item_name;not null;comment:产品物料名称" json:"item_name"`                    // 产品物料名称
	Specification  string  `gorm:"column:specification;comment:规格型号" json:"specification"`                       // 规格型号
	UnitOfMeasure  string  `gorm:"column:unit_of_measure;not null;comment:单位" json:"unit_of_measure"`            // 单位
	UnitName       string  `gorm:"column:unit_name;comment:单位名称" json:"unit_name"`                               // 单位名称
	Quantity       int32   `gorm:"column:quantity;default:1;comment:生产数量" json:"quantity"`                       // 生产数量
	ProductionTime float64 `gorm:"column:production_time;default:1.00;comment:生产用时" json:"production_time"`      // 生产用时
	TimeUnitType   string  `gorm:"column:time_unit_type;default:MINUTE;comment:时间单位" json:"time_unit_type"`      // 时间单位
	Remark         string  `gorm:"column:remark;comment:备注" json:"remark"`                                       // 备注
	Attr1          string  `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                      // 预留字段1
	Attr2          string  `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                      // 预留字段2
	Attr3          int32   `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                      // 预留字段3
	Attr4          int32   `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                      // 预留字段4
}

package craftRouteModels

import "time"

// SysProRouteProductBom 产品制程物料BOM表
type SysProRouteProductBom struct {
	RecordID      int64     `gorm:"column:record_id;primaryKey;autoIncrement:true;comment:记录ID" json:"record_id"` // 记录ID
	RouteID       int64     `gorm:"column:route_id;not null;comment:工艺路线ID" json:"route_id"`                      // 工艺路线ID
	ProcessID     int64     `gorm:"column:process_id;not null;comment:工序ID" json:"process_id"`                    // 工序ID
	ProductID     int64     `gorm:"column:product_id;not null;comment:产品BOM中的唯一ID" json:"product_id"`             // 产品BOM中的唯一ID
	ItemID        int64     `gorm:"column:item_id;not null;comment:产品物料ID" json:"item_id"`                        // 产品物料ID
	ItemCode      string    `gorm:"column:item_code;not null;comment:产品物料编码" json:"item_code"`                    // 产品物料编码
	ItemName      string    `gorm:"column:item_name;not null;comment:产品物料名称" json:"item_name"`                    // 产品物料名称
	Specification string    `gorm:"column:specification;comment:规格型号" json:"specification"`                       // 规格型号
	UnitOfMeasure string    `gorm:"column:unit_of_measure;not null;comment:单位" json:"unit_of_measure"`            // 单位
	UnitName      string    `gorm:"column:unit_name;comment:单位名称" json:"unit_name"`                               // 单位名称
	Quantity      float64   `gorm:"column:quantity;default:1.00;comment:用料比例" json:"quantity"`                    // 用料比例
	Remark        string    `gorm:"column:remark;comment:备注" json:"remark"`                                       // 备注
	Attr1         string    `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                      // 预留字段1
	Attr2         string    `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                      // 预留字段2
	Attr3         int32     `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                      // 预留字段3
	Attr4         int32     `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                      // 预留字段4
	CreateBy      string    `gorm:"column:create_by;comment:创建者" json:"create_by"`                                // 创建者
	CreateTime    time.Time `gorm:"column:create_time;comment:创建时间" json:"create_time"`                           // 创建时间
	UpdateBy      string    `gorm:"column:update_by;comment:更新者" json:"update_by"`                                // 更新者
	UpdateTime    time.Time `gorm:"column:update_time;comment:更新时间" json:"update_time"`                           // 更新时间
}

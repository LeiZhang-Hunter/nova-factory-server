package craftRouteModels

import (
	"nova-factory-server/app/baize"
	"time"
)

// SysProWorkorder 生产工单表
type SysProWorkorder struct {
	WorkorderID       int64     `gorm:"column:workorder_id;primaryKey;autoIncrement:true;comment:工单ID" json:"workorder_id"` // 工单ID
	WorkorderCode     string    `gorm:"column:workorder_code;not null;comment:工单编码" json:"workorder_code"`                  // 工单编码
	WorkorderName     string    `gorm:"column:workorder_name;not null;comment:工单名称" json:"workorder_name"`                  // 工单名称
	WorkorderType     string    `gorm:"column:workorder_type;default:SELF;comment:工单类型" json:"workorder_type"`              // 工单类型
	OrderSource       string    `gorm:"column:order_source;not null;comment:来源类型" json:"order_source"`                      // 来源类型
	SourceCode        string    `gorm:"column:source_code;comment:来源单据" json:"source_code"`                                 // 来源单据
	RouterID          int64     `gorm:"column:router_id;not null;comment:工艺路线id" json:"router_id"`                          // 工艺路线id
	ProductID         int64     `gorm:"column:product_id;not null;comment:产品ID" json:"product_id"`                          // 产品ID
	ProductCode       string    `gorm:"column:product_code;not null;comment:产品编号" json:"product_code"`                      // 产品编号
	ProductName       string    `gorm:"column:product_name;not null;comment:产品名称" json:"product_name"`                      // 产品名称
	ProductSpc        string    `gorm:"column:product_spc;comment:规格型号" json:"product_spc"`                                 // 规格型号
	UnitOfMeasure     string    `gorm:"column:unit_of_measure;not null;comment:单位" json:"unit_of_measure"`                  // 单位
	Quantity          float64   `gorm:"column:quantity;not null;default:0.00;comment:生产数量" json:"quantity"`                 // 生产数量
	QuantityProduced  float64   `gorm:"column:quantity_produced;default:0.00;comment:已生产数量" json:"quantity_produced"`       // 已生产数量
	QuantityChanged   float64   `gorm:"column:quantity_changed;default:0.00;comment:调整数量" json:"quantity_changed"`          // 调整数量
	QuantityScheduled float64   `gorm:"column:quantity_scheduled;default:0.00;comment:已排产数量" json:"quantity_scheduled"`     // 已排产数量
	ClientID          int64     `gorm:"column:client_id;comment:客户ID" json:"client_id"`                                     // 客户ID
	ClientCode        string    `gorm:"column:client_code;comment:客户编码" json:"client_code"`                                 // 客户编码
	ClientName        string    `gorm:"column:client_name;comment:客户名称" json:"client_name"`                                 // 客户名称
	VendorID          int64     `gorm:"column:vendor_id;comment:供应商ID" json:"vendor_id"`                                    // 供应商ID
	VendorCode        string    `gorm:"column:vendor_code;comment:供应商编号" json:"vendor_code"`                                // 供应商编号
	VendorName        string    `gorm:"column:vendor_name;comment:供应商名称" json:"vendor_name"`                                // 供应商名称
	BatchCode         string    `gorm:"column:batch_code;comment:批次号" json:"batch_code"`                                    // 批次号
	RequestDate       time.Time `gorm:"column:request_date;not null;comment:需求日期" json:"request_date"`                      // 需求日期
	ParentID          int64     `gorm:"column:parent_id;not null;comment:父工单" json:"parent_id"`                             // 父工单
	Ancestors         string    `gorm:"column:ancestors;not null;comment:所有父节点ID" json:"ancestors"`                         // 所有父节点ID
	FinishDate        time.Time `gorm:"column:finish_date;comment:完成时间" json:"finish_date"`                                 // 完成时间
	Status            string    `gorm:"column:status;default:PREPARE;comment:单据状态" json:"status"`                           // 单据状态
	Remark            string    `gorm:"column:remark;comment:备注" json:"remark"`                                             // 备注
	Attr1             string    `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                            // 预留字段1
	Attr2             string    `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                            // 预留字段2
	Attr3             int32     `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                            // 预留字段3
	Attr4             int32     `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                            // 预留字段4
	DeptID            int64     `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                         // 部门ID
	State             bool      `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                                   // 操作状态（0正常 -1删除）
	baize.BaseEntity
}

func OfSysSetProWorkorder(order *SysSetProWorkorder) *SysProWorkorder {
	return &SysProWorkorder{
		WorkorderID:       order.WorkorderID,
		WorkorderCode:     order.WorkorderCode,
		WorkorderName:     order.WorkorderName,
		WorkorderType:     order.WorkorderType,
		OrderSource:       order.OrderSource,
		SourceCode:        order.SourceCode,
		RouterID:          order.RouterID,
		ProductID:         order.ProductID,
		ProductCode:       order.ProductCode,
		ProductName:       order.ProductName,
		ProductSpc:        order.ProductSpc,
		UnitOfMeasure:     order.UnitOfMeasure,
		Quantity:          order.Quantity,
		QuantityProduced:  order.QuantityProduced,
		QuantityChanged:   order.QuantityChanged,
		QuantityScheduled: order.QuantityScheduled,
		ClientID:          order.ClientID,
		ClientCode:        order.ClientCode,
		ClientName:        order.ClientName,
		VendorID:          order.VendorID,
		VendorCode:        order.VendorCode,
		VendorName:        order.VendorName,
		BatchCode:         order.BatchCode,
		RequestDate:       order.RequestDate,
		ParentID:          order.ParentID,
		Ancestors:         order.Ancestors,
		FinishDate:        order.FinishDate,
		Status:            order.Status,
		Remark:            order.Remark,
		Attr1:             order.Attr1,
		Attr2:             order.Attr2,
		Attr3:             order.Attr3,
		Attr4:             order.Attr4,
		DeptID:            order.DeptID,
		State:             order.State,
	}
}

type SysSetProWorkorder struct {
	WorkorderID       int64     `gorm:"column:workorder_id;primaryKey;autoIncrement:true;comment:工单ID" json:"workorder_id"` // 工单ID
	WorkorderCode     string    `gorm:"column:workorder_code;not null;comment:工单编码" json:"workorder_code"`                  // 工单编码
	WorkorderName     string    `gorm:"column:workorder_name;not null;comment:工单名称" json:"workorder_name"`                  // 工单名称
	WorkorderType     string    `gorm:"column:workorder_type;default:SELF;comment:工单类型" json:"workorder_type"`              // 工单类型
	OrderSource       string    `gorm:"column:order_source;not null;comment:来源类型" json:"order_source"`                      // 来源类型
	SourceCode        string    `gorm:"column:source_code;comment:来源单据" json:"source_code"`                                 // 来源单据
	RouterID          int64     `gorm:"column:router_id;not null;comment:工艺路线id" json:"router_id"`                          // 工艺路线id
	ProductID         int64     `gorm:"column:product_id;not null;comment:产品ID" json:"product_id"`                          // 产品ID
	ProductCode       string    `gorm:"column:product_code;not null;comment:产品编号" json:"product_code"`                      // 产品编号
	ProductName       string    `gorm:"column:product_name;not null;comment:产品名称" json:"product_name"`                      // 产品名称
	ProductSpc        string    `gorm:"column:product_spc;comment:规格型号" json:"product_spc"`                                 // 规格型号
	UnitOfMeasure     string    `gorm:"column:unit_of_measure;not null;comment:单位" json:"unit_of_measure"`                  // 单位
	Quantity          float64   `gorm:"column:quantity;not null;default:0.00;comment:生产数量" json:"quantity"`                 // 生产数量
	QuantityProduced  float64   `gorm:"column:quantity_produced;default:0.00;comment:已生产数量" json:"quantity_produced"`       // 已生产数量
	QuantityChanged   float64   `gorm:"column:quantity_changed;default:0.00;comment:调整数量" json:"quantity_changed"`          // 调整数量
	QuantityScheduled float64   `gorm:"column:quantity_scheduled;default:0.00;comment:已排产数量" json:"quantity_scheduled"`     // 已排产数量
	ClientID          int64     `gorm:"column:client_id;comment:客户ID" json:"client_id"`                                     // 客户ID
	ClientCode        string    `gorm:"column:client_code;comment:客户编码" json:"client_code"`                                 // 客户编码
	ClientName        string    `gorm:"column:client_name;comment:客户名称" json:"client_name"`                                 // 客户名称
	VendorID          int64     `gorm:"column:vendor_id;comment:供应商ID" json:"vendor_id"`                                    // 供应商ID
	VendorCode        string    `gorm:"column:vendor_code;comment:供应商编号" json:"vendor_code"`                                // 供应商编号
	VendorName        string    `gorm:"column:vendor_name;comment:供应商名称" json:"vendor_name"`                                // 供应商名称
	BatchCode         string    `gorm:"column:batch_code;comment:批次号" json:"batch_code"`                                    // 批次号
	RequestDate       time.Time `gorm:"column:request_date;not null;comment:需求日期" json:"request_date"`                      // 需求日期
	ParentID          int64     `gorm:"column:parent_id;not null;comment:父工单" json:"parent_id"`                             // 父工单
	Ancestors         string    `gorm:"column:ancestors;not null;comment:所有父节点ID" json:"ancestors"`                         // 所有父节点ID
	FinishDate        time.Time `gorm:"column:finish_date;comment:完成时间" json:"finish_date"`                                 // 完成时间
	Status            string    `gorm:"column:status;default:PREPARE;comment:单据状态" json:"status"`                           // 单据状态
	Remark            string    `gorm:"column:remark;comment:备注" json:"remark"`                                             // 备注
	Attr1             string    `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                            // 预留字段1
	Attr2             string    `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                            // 预留字段2
	Attr3             int32     `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                            // 预留字段3
	Attr4             int32     `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                            // 预留字段4
	DeptID            int64     `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                         // 部门ID
	State             bool      `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                                   // 操作状态（0正常 -1删除）
}

type SysProWorkorderReq struct {
	WorkorderCode string `form:"workorder_code"` // 工单编码
	WorkorderName string `form:"workorder_name"` // 工单名称
	SourceCode    string `form:"source_code"`    // 来源单据
	ProductCode   string `form:"product_code"`   // 产品编号
	ProductName   string `form:"product_name"`   // 产品名称
	ClientCode    string `form:"client_code"`    // 客户编码
	ClientName    string `form:"client_name"`    // 客户名称
	WorkorderType string `form:"workorder_type"` // 工单类型
	RequestDate   string `form:"request_date"`   // 需求日期
	baize.BaseEntityDQL
}

type SysProWorkorderList struct {
	Rows  []*SysProWorkorder `json:"rows"`
	Total int64              `json:"total"`
}

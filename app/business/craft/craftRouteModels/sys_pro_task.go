package craftRouteModels

import (
	"nova-factory-server/app/baize"
	"time"
)

// SysProTask 生产任务表
type SysProTask struct {
	TaskID             int64     `gorm:"column:task_id;primaryKey;autoIncrement:true;comment:任务ID" json:"task_id"`         // 任务ID
	TaskCode           string    `gorm:"column:task_code;not null;comment:任务编号" json:"task_code"`                          // 任务编号
	TaskName           string    `gorm:"column:task_name;not null;comment:任务名称" json:"task_name"`                          // 任务名称
	WorkorderID        int64     `gorm:"column:workorder_id;not null;comment:生产工单ID" json:"workorder_id"`                  // 生产工单ID
	WorkorderCode      string    `gorm:"column:workorder_code;not null;comment:生产工单编号" json:"workorder_code"`              // 生产工单编号
	WorkorderName      string    `gorm:"column:workorder_name;not null;comment:工单名称" json:"workorder_name"`                // 工单名称
	WorkstationID      int64     `gorm:"column:workstation_id;not null;comment:工作站ID" json:"workstation_id"`               // 工作站ID
	WorkstationCode    string    `gorm:"column:workstation_code;not null;comment:工作站编号" json:"workstation_code"`           // 工作站编号
	WorkstationName    string    `gorm:"column:workstation_name;not null;comment:工作站名称" json:"workstation_name"`           // 工作站名称
	RouteID            int64     `gorm:"column:route_id;not null;comment:工艺ID" json:"route_id"`                            // 工艺ID
	RouteCode          string    `gorm:"column:route_code;comment:工艺编号" json:"route_code"`                                 // 工艺编号
	ProcessID          int64     `gorm:"column:process_id;not null;comment:工序ID" json:"process_id"`                        // 工序ID
	ProcessCode        string    `gorm:"column:process_code;comment:工序编码" json:"process_code"`                             // 工序编码
	ProcessName        string    `gorm:"column:process_name;comment:工序名称" json:"process_name"`                             // 工序名称
	ItemID             int64     `gorm:"column:item_id;not null;comment:产品物料ID" json:"item_id"`                            // 产品物料ID
	ItemCode           string    `gorm:"column:item_code;not null;comment:产品物料编码" json:"item_code"`                        // 产品物料编码
	ItemName           string    `gorm:"column:item_name;not null;comment:产品物料名称" json:"item_name"`                        // 产品物料名称
	Specification      string    `gorm:"column:specification;comment:规格型号" json:"specification"`                           // 规格型号
	UnitOfMeasure      string    `gorm:"column:unit_of_measure;not null;comment:单位" json:"unit_of_measure"`                // 单位
	UnitName           string    `gorm:"column:unit_name;comment:单位名称" json:"unit_name"`                                   // 单位名称
	Quantity           float64   `gorm:"column:quantity;not null;default:1.00;comment:排产数量" json:"quantity"`               // 排产数量
	QuantityProduced   float64   `gorm:"column:quantity_produced;default:0.00;comment:已生产数量" json:"quantity_produced"`     // 已生产数量
	QuantityQuanlify   float64   `gorm:"column:quantity_quanlify;default:0.00;comment:合格品数量" json:"quantity_quanlify"`     // 合格品数量
	QuantityUnquanlify float64   `gorm:"column:quantity_unquanlify;default:0.00;comment:不良品数量" json:"quantity_unquanlify"` // 不良品数量
	QuantityChanged    float64   `gorm:"column:quantity_changed;default:0.00;comment:调整数量" json:"quantity_changed"`        // 调整数量
	ClientID           int64     `gorm:"column:client_id;comment:客户ID" json:"client_id"`                                   // 客户ID
	ClientCode         string    `gorm:"column:client_code;comment:客户编码" json:"client_code"`                               // 客户编码
	ClientName         string    `gorm:"column:client_name;comment:客户名称" json:"client_name"`                               // 客户名称
	ClientNick         string    `gorm:"column:client_nick;comment:客户简称" json:"client_nick"`                               // 客户简称
	StartTime          time.Time `gorm:"column:start_time;default:CURRENT_TIMESTAMP;comment:开始生产时间" json:"start_time"`     // 开始生产时间
	Duration           int32     `gorm:"column:duration;default:1;comment:生产时长" json:"duration"`                           // 生产时长
	EndTime            time.Time `gorm:"column:end_time;comment:完成生产时间" json:"end_time"`                                   // 完成生产时间
	ColorCode          string    `gorm:"column:color_code;default:#00AEF3;comment:甘特图显示颜色" json:"color_code"`              // 甘特图显示颜色
	RequestDate        time.Time `gorm:"column:request_date;comment:需求日期" json:"request_date"`                             // 需求日期
	Status             string    `gorm:"column:status;default:NORMAL;comment:生产状态" json:"status"`                          // 生产状态
	Remark             string    `gorm:"column:remark;comment:备注" json:"remark"`                                           // 备注
	Attr1              string    `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                          // 预留字段1
	Attr2              string    `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                          // 预留字段2
	Attr3              int32     `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                          // 预留字段3
	Attr4              int32     `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                          // 预留字段4
	DeptID             int64     `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                       // 部门ID
	State              bool      `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                                 // 操作状态（0正常 -1删除）
	baize.BaseEntity
}

func OfSysProTask(order *SysSetProTask) *SysProTask {
	return &SysProTask{
		TaskID:             order.TaskID,
		TaskCode:           order.TaskCode,
		TaskName:           order.TaskName,
		WorkorderID:        order.WorkorderID,
		WorkorderCode:      order.WorkorderCode,
		WorkorderName:      order.WorkorderName,
		WorkstationID:      order.WorkstationID,
		WorkstationCode:    order.WorkstationCode,
		WorkstationName:    order.WorkstationName,
		RouteID:            order.RouteID,
		RouteCode:          order.RouteCode,
		ProcessID:          order.ProcessID,
		ProcessCode:        order.ProcessCode,
		ProcessName:        order.ProcessName,
		ItemID:             order.ItemID,
		ItemCode:           order.ItemCode,
		ItemName:           order.ItemName,
		Specification:      order.Specification,
		UnitOfMeasure:      order.UnitOfMeasure,
		UnitName:           order.UnitName,
		Quantity:           order.Quantity,
		QuantityProduced:   order.QuantityProduced,
		QuantityQuanlify:   order.QuantityQuanlify,
		QuantityUnquanlify: order.QuantityUnquanlify,
		QuantityChanged:    order.QuantityChanged,
		ClientID:           order.ClientID,
		ClientCode:         order.ClientCode,
		ClientName:         order.ClientName,
		ClientNick:         order.ClientNick,
		StartTime:          order.StartTime,
		Duration:           order.Duration,
		EndTime:            order.EndTime,
		ColorCode:          order.ColorCode,
		RequestDate:        order.RequestDate,
		Status:             order.Status,
		Remark:             order.Remark,
		Attr1:              order.Attr1,
		Attr2:              order.Attr2,
		Attr3:              order.Attr3,
		Attr4:              order.Attr4,
	}
}

type SysSetProTask struct {
	TaskID             int64     `gorm:"column:task_id;primaryKey;autoIncrement:true;comment:任务ID" json:"task_id"`         // 任务ID
	TaskCode           string    `gorm:"column:task_code;not null;comment:任务编号" json:"task_code"`                          // 任务编号
	TaskName           string    `gorm:"column:task_name;not null;comment:任务名称" json:"task_name"`                          // 任务名称
	WorkorderID        int64     `gorm:"column:workorder_id;not null;comment:生产工单ID" json:"workorder_id"`                  // 生产工单ID
	WorkorderCode      string    `gorm:"column:workorder_code;not null;comment:生产工单编号" json:"workorder_code"`              // 生产工单编号
	WorkorderName      string    `gorm:"column:workorder_name;not null;comment:工单名称" json:"workorder_name"`                // 工单名称
	WorkstationID      int64     `gorm:"column:workstation_id;not null;comment:工作站ID" json:"workstation_id"`               // 工作站ID
	WorkstationCode    string    `gorm:"column:workstation_code;not null;comment:工作站编号" json:"workstation_code"`           // 工作站编号
	WorkstationName    string    `gorm:"column:workstation_name;not null;comment:工作站名称" json:"workstation_name"`           // 工作站名称
	RouteID            int64     `gorm:"column:route_id;not null;comment:工艺ID" json:"route_id"`                            // 工艺ID
	RouteCode          string    `gorm:"column:route_code;comment:工艺编号" json:"route_code"`                                 // 工艺编号
	ProcessID          int64     `gorm:"column:process_id;not null;comment:工序ID" json:"process_id"`                        // 工序ID
	ProcessCode        string    `gorm:"column:process_code;comment:工序编码" json:"process_code"`                             // 工序编码
	ProcessName        string    `gorm:"column:process_name;comment:工序名称" json:"process_name"`                             // 工序名称
	ItemID             int64     `gorm:"column:item_id;not null;comment:产品物料ID" json:"item_id"`                            // 产品物料ID
	ItemCode           string    `gorm:"column:item_code;not null;comment:产品物料编码" json:"item_code"`                        // 产品物料编码
	ItemName           string    `gorm:"column:item_name;not null;comment:产品物料名称" json:"item_name"`                        // 产品物料名称
	Specification      string    `gorm:"column:specification;comment:规格型号" json:"specification"`                           // 规格型号
	UnitOfMeasure      string    `gorm:"column:unit_of_measure;not null;comment:单位" json:"unit_of_measure"`                // 单位
	UnitName           string    `gorm:"column:unit_name;comment:单位名称" json:"unit_name"`                                   // 单位名称
	Quantity           float64   `gorm:"column:quantity;not null;default:1.00;comment:排产数量" json:"quantity"`               // 排产数量
	QuantityProduced   float64   `gorm:"column:quantity_produced;default:0.00;comment:已生产数量" json:"quantity_produced"`     // 已生产数量
	QuantityQuanlify   float64   `gorm:"column:quantity_quanlify;default:0.00;comment:合格品数量" json:"quantity_quanlify"`     // 合格品数量
	QuantityUnquanlify float64   `gorm:"column:quantity_unquanlify;default:0.00;comment:不良品数量" json:"quantity_unquanlify"` // 不良品数量
	QuantityChanged    float64   `gorm:"column:quantity_changed;default:0.00;comment:调整数量" json:"quantity_changed"`        // 调整数量
	ClientID           int64     `gorm:"column:client_id;comment:客户ID" json:"client_id"`                                   // 客户ID
	ClientCode         string    `gorm:"column:client_code;comment:客户编码" json:"client_code"`                               // 客户编码
	ClientName         string    `gorm:"column:client_name;comment:客户名称" json:"client_name"`                               // 客户名称
	ClientNick         string    `gorm:"column:client_nick;comment:客户简称" json:"client_nick"`                               // 客户简称
	StartTime          time.Time `gorm:"column:start_time;default:CURRENT_TIMESTAMP;comment:开始生产时间" json:"start_time"`     // 开始生产时间
	Duration           int32     `gorm:"column:duration;default:1;comment:生产时长" json:"duration"`                           // 生产时长
	EndTime            time.Time `gorm:"column:end_time;comment:完成生产时间" json:"end_time"`                                   // 完成生产时间
	ColorCode          string    `gorm:"column:color_code;default:#00AEF3;comment:甘特图显示颜色" json:"color_code"`              // 甘特图显示颜色
	RequestDate        time.Time `gorm:"column:request_date;comment:需求日期" json:"request_date"`                             // 需求日期
	Status             string    `gorm:"column:status;default:NORMAL;comment:生产状态" json:"status"`                          // 生产状态
	Remark             string    `gorm:"column:remark;comment:备注" json:"remark"`                                           // 备注
	Attr1              string    `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                          // 预留字段1
	Attr2              string    `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                          // 预留字段2
	Attr3              int32     `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                          // 预留字段3
	Attr4              int32     `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                          // 预留字段4
}

// SysProTaskReq 生产任务请求
type SysProTaskReq struct {
	WorkorderID int64 `gorm:"column:workorder_id;not null;comment:生产工单ID" form:"workorder_id"` // 生产工单ID
	baize.BaseEntityDQL
}

type SysProTaskList struct {
	Rows  []*SysProTask `json:"rows"`
	Total int64         `json:"total"`
}

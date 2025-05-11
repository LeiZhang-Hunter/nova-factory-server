package craftRouteModels

import (
	"nova-factory-server/app/baize"
)

type SysCraftRouteRequest struct {
	RouteID   int64  `gorm:"column:route_id;primaryKey;comment:工艺路线ID" json:"datasetId,string"`              // 出库id
	RouteCode string `gorm:"column:route_code;not null;comment:工艺路线编号" json:"route_code" binding:"required"` // 工艺路线编号
	RouteName string `gorm:"column:route_name;not null;comment:工艺路线名称" json:"route_name" binding:"required"` // 工艺路线名称
	RouteDesc string `gorm:"column:route_desc;comment:工艺路线说明" json:"route_desc"`                             // 工艺路线说明
	Remark    string `gorm:"column:remark;comment:备注" json:"remark"`                                         // 备注
	Status    bool   `gorm:"column:status;comment:操作状态（0正常 1异常）" json:"status"`                              // 操作状态（0正常 1异常）
}

type SysCraftRoute struct {
	RouteID        int64  `gorm:"column:route_id;primaryKey;autoIncrement:true;comment:工艺路线ID" json:"route_id"` // 工艺路线ID
	RouteCode      string `gorm:"column:route_code;not null;comment:工艺路线编号" json:"route_code"`                  // 工艺路线编号
	RouteName      string `gorm:"column:route_name;not null;comment:工艺路线名称" json:"route_name"`                  // 工艺路线名称
	RouteDesc      string `gorm:"column:route_desc;comment:工艺路线说明" json:"route_desc"`                           // 工艺路线说明
	Remark         string `gorm:"column:remark;comment:备注" json:"remark"`                                       // 备注
	Attr1          string `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                      // 预留字段1
	Attr2          string `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                      // 预留字段2
	Attr3          int32  `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                      // 预留字段3
	Attr4          int32  `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                      // 预留字段4
	DeptID         int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                   // 部门ID
	Status         bool   `gorm:"column:status;comment:操作状态（0正常 1异常）" json:"status"`                            // 操作状态（0正常 1异常）
	State          bool   `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                             // 操作状态（0正常 -1删除）
	CreateUserName string `json:"createUserName" gorm:"-"`
	UpdateUserName string `json:"updateUserName" gorm:"-"`
	baize.BaseEntity
}

// SysCraftRouteListReq 读取列表
type SysCraftRouteListReq struct {
	RouteCode string `gorm:"column:route_code;not null;comment:工艺路线编号" json:"route_code"` // 工艺路线编号
	RouteName string `gorm:"column:route_name;not null;comment:工艺路线名称" json:"route_name"` // 工艺路线名称
	Status    *bool  `gorm:"column:status;comment:操作状态（0正常 1异常）" json:"status"`           // 操作状态（0正常 1异常）
	baize.BaseEntityDQL
}

type SysCraftRouteListData struct {
	Rows  []*SysCraftRoute `json:"rows"`
	Total int64            `json:"total"`
}

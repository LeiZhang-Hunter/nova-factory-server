package craftRouteModels

import (
	"nova-factory-server/app/baize"
)

type SysCraftRouteRequest struct {
	RouteID   int64  `gorm:"column:route_id;primaryKey;comment:工艺路线ID" json:"route_id,string"`               // 出库id
	RouteCode string `gorm:"column:route_code;not null;comment:工艺路线编号" json:"route_code" binding:"required"` // 工艺路线编号
	RouteName string `gorm:"column:route_name;not null;comment:工艺路线名称" json:"route_name" binding:"required"` // 工艺路线名称
	RouteDesc string `gorm:"column:route_desc;comment:工艺路线说明" json:"route_desc"`                             // 工艺路线说明
	LoopTime  int32  `gorm:"column:loop_time;default:30;comment:循环执行时间" json:"loop_time"`                    // 循环执行时间
	Remark    string `gorm:"column:remark;comment:备注" json:"remark"`                                         // 备注
	Status    bool   `gorm:"column:status;comment:操作状态（0正常 1异常）" json:"status"`                              // 操作状态（0正常 1异常）
}

type SysCraftRoute struct {
	RouteID        int64  `gorm:"column:route_id;primaryKey;autoIncrement:true;comment:工艺路线ID" json:"route_id,string"` // 工艺路线ID
	RouteCode      string `gorm:"column:route_code;not null;comment:工艺路线编号" json:"route_code"`                         // 工艺路线编号
	RouteName      string `gorm:"column:route_name;not null;comment:工艺路线名称" json:"route_name"`                         // 工艺路线名称
	RouteDesc      string `gorm:"column:route_desc;comment:工艺路线说明" json:"route_desc"`                                  // 工艺路线说明
	LoopTime       int32  `gorm:"column:loop_time;default:30;comment:循环执行时间" json:"loop_time"`                         // 循环执行时间
	Remark         string `gorm:"column:remark;comment:备注" json:"remark"`                                              // 备注
	Attr1          string `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                             // 预留字段1
	Attr2          string `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                             // 预留字段2
	Attr3          int32  `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                             // 预留字段3
	Attr4          int32  `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                             // 预留字段4
	DeptID         int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id,string"`                                   // 部门ID
	Status         bool   `gorm:"column:status;comment:操作状态（0正常 1异常）" json:"status"`                                   // 操作状态（0正常 1异常）
	State          bool   `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                                    // 操作状态（0正常 -1删除）
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

type SysCraftRouteDetailRequest struct {
	RouteID int64 `gorm:"column:route_id;primaryKey;comment:工艺路线ID" form:"route_id" json:"route_id,string"` // 出库id
}

type Dimensions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type ComputedPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type HandleBoundsSource struct {
	Id       string  `json:"id"`
	Type     string  `json:"type"`
	NodeId   string  `json:"nodeId"`
	Position string  `json:"position"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
}

type HandleBoundsTarget struct {
	Id       string  `json:"id"`
	Type     string  `json:"type"`
	NodeId   string  `json:"nodeId"`
	Position string  `json:"position"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
}

type ProcessTopoNodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type ProcessTopoData struct {
	Label       string      `json:"label"`
	Description string      `json:"description"`
	Config      interface{} `json:"config"`
}

type ProcessTopoNode struct {
	Id               string            `json:"id"`
	Type             string            `json:"type"`
	Dimensions       *Dimensions       `json:"dimensions"`
	ComputedPosition *ComputedPosition `json:"computedPosition"`
	HandleBounds     struct {
		Source []HandleBoundsSource `json:"source"`
		Target []HandleBoundsTarget `json:"target"`
	} `json:"handleBounds"`
	Selected    bool                     `json:"selected"`
	Dragging    bool                     `json:"dragging"`
	Resizing    bool                     `json:"resizing"`
	Initialized bool                     `json:"initialized"`
	IsParent    bool                     `json:"isParent"`
	Position    *ProcessTopoNodePosition `json:"position"`
	Data        *ProcessTopoData         `json:"data"`
	Events      struct {
	} `json:"events"`
}

type ProcessData struct {
	ProcessCode string `json:"process_code" mapstructure:"process_code"`
	ProcessName string `json:"process_name" mapstructure:"process_name"`
	ProcessId   string `json:"process_id,string" mapstructure:"process_id,string"`
}

type ProcessTopoEdge struct {
	Id           string `json:"id"`
	Type         string `json:"type"`
	Source       string `json:"source"`
	Target       string `json:"target"`
	SourceHandle string `json:"sourceHandle"`
	TargetHandle string `json:"targetHandle"`
	Data         struct {
		Label       string      `json:"label"`
		Description string      `json:"description"`
		Config      interface{} `json:"config"`
	} `json:"data"`
	Events struct {
	} `json:"events"`
	Label     string `json:"label"`
	Animated  bool   `json:"animated"`
	MarkerEnd struct {
		Type   string `json:"type"`
		Color  string `json:"color"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"markerEnd"`
	Style struct {
		Stroke      string `json:"stroke"`
		StrokeWidth int    `json:"strokeWidth"`
	} `json:"style"`
	SourceX float64 `json:"sourceX"`
	SourceY float64 `json:"sourceY"`
	TargetX float64 `json:"targetX"`
	TargetY float64 `json:"targetY"`
}

// ProcessTopo 工艺流程
type ProcessTopo struct {
	Route    *SysCraftRoute           `json:"route,omitempty"`
	Nodes    []*ProcessTopoNode       `json:"nodes,omitempty"`
	Edges    []*ProcessTopoEdge       `json:"edges,omitempty"`
	Products []*SysProRouteSetProduct `json:"products,omitempty"`
}

// SysCraftRouteConfig 工艺配置详情
type SysCraftRouteConfig struct {
	RouteConfigID int64        `gorm:"column:route_config_id;primaryKey;comment:工艺路线配置id" json:"route_config_id"` // 工艺路线配置id
	RouteID       int64        `gorm:"column:route_id;not null;comment:工艺路线ID" json:"route_id"`                   // 工艺路线ID
	Context       string       `gorm:"column:context;comment:配置" json:"-"`                                        // 配置
	Config        string       `gorm:"column:config;comment:配置" json:"-"`                                         // 配置
	DeptID        int64        `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                // 部门ID
	Topo          *ProcessTopo `gorm:"-" json:"topo"`
	baize.BaseEntity
	State bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

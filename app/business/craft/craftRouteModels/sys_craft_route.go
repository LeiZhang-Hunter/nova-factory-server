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
	RouteID        int64  `gorm:"column:route_id;primaryKey;autoIncrement:true;comment:工艺路线ID" json:"route_id,string"` // 工艺路线ID
	RouteCode      string `gorm:"column:route_code;not null;comment:工艺路线编号" json:"route_code"`                         // 工艺路线编号
	RouteName      string `gorm:"column:route_name;not null;comment:工艺路线名称" json:"route_name"`                         // 工艺路线名称
	RouteDesc      string `gorm:"column:route_desc;comment:工艺路线说明" json:"route_desc"`                                  // 工艺路线说明
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

type ProcessTopoNode struct {
	ProcessID       int64                   `gorm:"column:process_id;primaryKey;autoIncrement:true;comment:工序ID" json:"process_id,string"` // 工序ID
	ProcessCode     string                  `gorm:"column:process_code;not null;comment:工序编码" json:"process_code" binding:"required"`      // 工序编码
	ProcessName     string                  `gorm:"column:process_name;not null;comment:工序名称" json:"process_name" binding:"required"`      // 工序名称
	Attention       string                  `gorm:"column:attention;comment:工艺要求" json:"attention"`                                        // 工艺要求
	Remark          string                  `gorm:"column:remark;comment:备注" json:"remark"`                                                // 备注
	Attr1           string                  `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                               // 预留字段1
	Attr2           string                  `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                               // 预留字段2
	Attr3           int32                   `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                               // 预留字段3
	Attr4           int32                   `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                               // 预留字段4
	Status          bool                    `gorm:"column:status;comment:是否启用（0禁用 1启用）" json:"status"`                                     // 是否启用（0禁用 1启用）
	ProcessContexts []*SysProProcessContent `json:"context"`
	Boms            []*SysSetProRouteProductBom
}

func NewProcessTopoNode(process *SysProProcess) *ProcessTopoNode {
	return &ProcessTopoNode{
		ProcessID:       process.ProcessID,
		ProcessCode:     process.ProcessCode,
		ProcessName:     process.ProcessName,
		Attention:       process.Attention,
		Remark:          process.Remark,
		Attr1:           process.Attr1,
		Attr2:           process.Attr2,
		Attr3:           process.Attr3,
		Attr4:           process.Attr4,
		Status:          process.Status,
		ProcessContexts: make([]*SysProProcessContent, 0),
	}
}

type ProcessTopoEdge struct {
	Source          uint64 `json:"source,string"`
	Target          uint64 `json:"target,string"`
	Esize           string `json:"esize"`
	RouteID         int64  `gorm:"column:route_id;not null;comment:工艺路线ID" json:"route_id"`             // 工艺路线ID
	ProcessCode     string `gorm:"column:process_code;comment:工序编码" json:"process_code"`                // 工序编码
	ProcessName     string `gorm:"column:process_name;comment:工序名称" json:"process_name"`                // 工序名称
	OrderNum        int32  `gorm:"column:order_num;default:1;comment:序号" json:"order_num"`              // 序号
	NextProcessCode string `gorm:"column:next_process_code;comment:工序编码" json:"next_process_code"`      // 工序编码
	NextProcessName string `gorm:"column:next_process_name;comment:工序名称" json:"next_process_name"`      // 工序名称
	LinkType        string `gorm:"column:link_type;default:SS;comment:与下一道工序关系" json:"link_type"`       // 与下一道工序关系
	DefaultPreTime  int64  `gorm:"column:default_pre_time;comment:准备时间" json:"default_pre_time"`        // 准备时间
	DefaultSufTime  int64  `gorm:"column:default_suf_time;comment:等待时间" json:"default_suf_time"`        // 等待时间
	ColorCode       string `gorm:"column:color_code;default:#00AEF3;comment:甘特图显示颜色" json:"color_code"` // 甘特图显示颜色
	KeyFlag         string `gorm:"column:key_flag;default:N;comment:关键工序" json:"key_flag"`              // 关键工序
	IsCheck         string `gorm:"column:is_check;default:N;comment:是否检验" json:"is_check"`              // 是否检验
	Remark          string `gorm:"column:remark;comment:备注" json:"remark"`                              // 备注
	Attr1           string `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                             // 预留字段1
	Attr2           string `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                             // 预留字段2
	Attr3           int32  `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                             // 预留字段3
}

func NewProcessTopoEdge(edge *SysProRouteProcess) *ProcessTopoEdge {
	return &ProcessTopoEdge{
		Source:          uint64(edge.ProcessID),
		Target:          uint64(edge.NextProcessID),
		RouteID:         edge.RouteID,
		ProcessCode:     edge.ProcessCode,
		ProcessName:     edge.ProcessName,
		OrderNum:        edge.OrderNum,
		NextProcessCode: edge.NextProcessCode,
		NextProcessName: edge.NextProcessName,
		LinkType:        edge.LinkType,
		DefaultPreTime:  edge.DefaultPreTime,
		DefaultSufTime:  edge.DefaultSufTime,
		ColorCode:       edge.ColorCode,
		KeyFlag:         edge.KeyFlag,
		IsCheck:         edge.IsCheck,
		Remark:          edge.Remark,
		Attr1:           edge.Attr1,
		Attr2:           edge.Attr2,
		Attr3:           edge.Attr3,
	}
}

// ProcessTopo 工艺流程
type ProcessTopo struct {
	Route    *SysCraftRoute           `json:"route,omitempty"`
	Nodes    []*ProcessTopoNode       `json:"nodes,omitempty"`
	Edges    []*ProcessTopoEdge       `json:"edges,omitempty"`
	Products []*SysProRouteSetProduct `json:"products,omitempty"`
}

func NewProcessTopo() *ProcessTopo {
	return &ProcessTopo{
		Nodes:    make([]*ProcessTopoNode, 0),
		Edges:    make([]*ProcessTopoEdge, 0),
		Products: make([]*SysProRouteSetProduct, 0),
	}
}

func (p *ProcessTopo) OfProcess(processes []*SysProProcess, processContexts []*SysProProcessContent,
	boms []*SysProRouteProductBom, products []*SysProRouteSetProduct) {
	if len(processes) == 0 {
		return
	}

	// 格式化工序内容
	processContextMap := make(map[uint64][]*SysProProcessContent)
	for _, processContext := range processContexts {
		_, ok := processContextMap[uint64(processContext.ProcessID)]
		if !ok {
			processContextMap[uint64(processContext.ProcessID)] = make([]*SysProProcessContent, 0)
		}
		processContextMap[uint64(processContext.ProcessID)] = append(processContextMap[uint64(processContext.ProcessID)], processContext)
	}

	// 投放bom
	productBomMap := make(map[uint64][]*SysSetProRouteProductBom)
	for _, bom := range boms {
		_, ok := productBomMap[uint64(bom.ProcessID)]
		if !ok {
			productBomMap[uint64(bom.ProcessID)] = make([]*SysSetProRouteProductBom, 0)
		}
		productBomMap[uint64(bom.ProcessID)] = append(productBomMap[uint64(bom.ProcessID)], OfSysSetProRouteProductBom(bom))
	}

	for _, process := range processes {
		node := NewProcessTopoNode(process)
		value, ok := processContextMap[uint64(process.ProcessID)]
		if ok {
			node.ProcessContexts = value
		}

		bomValue, ok := productBomMap[uint64(process.ProcessID)]
		if ok {
			node.Boms = bomValue
		}
		p.Nodes = append(p.Nodes, node)

	}

	p.Products = products

}

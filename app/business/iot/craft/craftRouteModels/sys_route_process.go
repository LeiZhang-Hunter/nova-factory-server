package craftRouteModels

import (
	"nova-factory-server/app/baize"
)

// SysProRouteProcessSetRequest 工艺组成表
type SysProRouteProcessSetRequest struct {
	RecordID        int64  `gorm:"column:record_id;primaryKey;autoIncrement:true;comment:记录ID" json:"record_id"` // 记录ID
	RouteID         int64  `gorm:"column:route_id;not null;comment:工艺路线ID" json:"route_id"`                      // 工艺路线ID
	ProcessID       int64  `gorm:"column:process_id;not null;comment:工序ID" json:"process_id"`                    // 工序ID
	ProcessCode     string `gorm:"column:process_code;comment:工序编码" json:"process_code"`                         // 工序编码
	ProcessName     string `gorm:"column:process_name;comment:工序名称" json:"process_name"`                         // 工序名称
	OrderNum        int32  `gorm:"column:order_num;default:1;comment:序号" json:"order_num"`                       // 序号
	NextProcessID   int64  `gorm:"column:next_process_id;not null;comment:工序ID" json:"next_process_id"`          // 工序ID
	NextProcessCode string `gorm:"column:next_process_code;comment:工序编码" json:"next_process_code"`               // 工序编码
	NextProcessName string `gorm:"column:next_process_name;comment:工序名称" json:"next_process_name"`               // 工序名称
	LinkType        string `gorm:"column:link_type;default:SS;comment:与下一道工序关系" json:"link_type"`                // 与下一道工序关系
	DefaultPreTime  int64  `gorm:"column:default_pre_time;comment:准备时间" json:"default_pre_time"`                 // 准备时间
	DefaultSufTime  int64  `gorm:"column:default_suf_time;comment:等待时间" json:"default_suf_time"`                 // 等待时间
	ColorCode       string `gorm:"column:color_code;default:#00AEF3;comment:甘特图显示颜色" json:"color_code"`          // 甘特图显示颜色
	KeyFlag         string `gorm:"column:key_flag;default:N;comment:关键工序" json:"key_flag"`                       // 关键工序
	IsCheck         string `gorm:"column:is_check;default:N;comment:是否检验" json:"is_check"`                       // 是否检验
	Remark          string `gorm:"column:remark;comment:备注" json:"remark"`                                       // 备注
	Attr1           string `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                      // 预留字段1
	Attr2           string `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                      // 预留字段2
	Attr3           int32  `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                      // 预留字段3
	Attr4           int32  `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                      // 预留字段4
	DeptID          int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                   // 部门ID
}

// SysProRouteProcess 工艺组成表
type SysProRouteProcess struct {
	RecordID        int64  `gorm:"column:record_id;primaryKey;autoIncrement:true;comment:记录ID" json:"record_id"` // 记录ID
	RouteID         int64  `gorm:"column:route_id;not null;comment:工艺路线ID" json:"route_id"`                      // 工艺路线ID
	ProcessID       int64  `gorm:"column:process_id;not null;comment:工序ID" json:"process_id"`                    // 工序ID
	ProcessCode     string `gorm:"column:process_code;comment:工序编码" json:"process_code"`                         // 工序编码
	ProcessName     string `gorm:"column:process_name;comment:工序名称" json:"process_name"`                         // 工序名称
	OrderNum        int32  `gorm:"column:order_num;default:1;comment:序号" json:"order_num"`                       // 序号
	NextProcessID   int64  `gorm:"column:next_process_id;not null;comment:工序ID" json:"next_process_id"`          // 工序ID
	NextProcessCode string `gorm:"column:next_process_code;comment:工序编码" json:"next_process_code"`               // 工序编码
	NextProcessName string `gorm:"column:next_process_name;comment:工序名称" json:"next_process_name"`               // 工序名称
	LinkType        string `gorm:"column:link_type;default:SS;comment:与下一道工序关系" json:"link_type"`                // 与下一道工序关系
	DefaultPreTime  int64  `gorm:"column:default_pre_time;comment:准备时间" json:"default_pre_time"`                 // 准备时间
	DefaultSufTime  int64  `gorm:"column:default_suf_time;comment:等待时间" json:"default_suf_time"`                 // 等待时间
	ColorCode       string `gorm:"column:color_code;default:#00AEF3;comment:甘特图显示颜色" json:"color_code"`          // 甘特图显示颜色
	KeyFlag         string `gorm:"column:key_flag;default:N;comment:关键工序" json:"key_flag"`                       // 关键工序
	IsCheck         string `gorm:"column:is_check;default:N;comment:是否检验" json:"is_check"`                       // 是否检验
	Remark          string `gorm:"column:remark;comment:备注" json:"remark"`                                       // 备注
	Attr1           string `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                      // 预留字段1
	Attr2           string `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                      // 预留字段2
	Attr3           int32  `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                      // 预留字段3
	Attr4           int32  `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                      // 预留字段4
	DeptID          int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                   // 部门ID
	State           bool   `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                             // 操作状态（0正常 -1删除）
	baize.BaseEntity
}

func (s *SysProRouteProcess) Of(req *SysProRouteProcessSetRequest) {
	s.RecordID = req.RecordID
	s.RouteID = req.RouteID
	s.ProcessID = req.ProcessID
	s.ProcessCode = req.ProcessCode
	s.ProcessName = req.ProcessName
	s.OrderNum = req.OrderNum
	s.NextProcessID = req.NextProcessID
	s.NextProcessCode = req.NextProcessCode
	s.NextProcessName = req.NextProcessName
	s.LinkType = req.LinkType
	s.DefaultPreTime = req.DefaultPreTime
	s.DefaultSufTime = req.DefaultSufTime
	s.ColorCode = req.ColorCode
	s.KeyFlag = req.KeyFlag
	s.IsCheck = req.IsCheck
	s.Remark = req.Remark
	s.Attr1 = req.Attr1
	s.Attr2 = req.Attr2
	s.Attr3 = req.Attr3
	s.Attr4 = req.Attr4
	return
}

type SysProRouteProcessListReq struct {
	baize.BaseEntityDQL
}

type SysProRouteProcessList struct {
	Rows  []*SysProRouteProcess `json:"rows"`
	Total int64                 `json:"total"`
}

package craftRouteModels

import (
	"nova-factory-server/app/baize"
)

// SysProProcess 生产工序表
type SysProProcess struct {
	ProcessID      int64  `gorm:"column:process_id;primaryKey;autoIncrement:true;comment:工序ID" json:"process_id"`   // 工序ID
	ProcessCode    string `gorm:"column:process_code;not null;comment:工序编码" json:"process_code" binding:"required"` // 工序编码
	ProcessName    string `gorm:"column:process_name;not null;comment:工序名称" json:"process_name" binding:"required"` // 工序名称
	Attention      string `gorm:"column:attention;comment:工艺要求" json:"attention"`                                   // 工艺要求
	Remark         string `gorm:"column:remark;comment:备注" json:"remark"`                                           // 备注
	Attr1          string `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                          // 预留字段1
	Attr2          string `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                          // 预留字段2
	Attr3          int32  `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                          // 预留字段3
	Attr4          int32  `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                          // 预留字段4
	Status         bool   `gorm:"column:status;comment:是否启用（0禁用 1启用）" json:"status"`                                // 是否启用（0禁用 1启用）
	State          bool   `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                                 // 操作状态（0正常 -1删除）
	CreateUserName string `json:"createUserName" gorm:"-"`
	UpdateUserName string `json:"updateUserName" gorm:"-"`
	baize.BaseEntity
}

func NewSysProProcess(req *SysProSetProcessReq) *SysProProcess {
	return &SysProProcess{
		ProcessID:   req.ProcessID,
		ProcessCode: req.ProcessCode,
		ProcessName: req.ProcessName,
		Attention:   req.Attention,
		Remark:      req.Remark,
		Attr1:       req.Attr1,
		Attr2:       req.Attr2,
		Attr3:       req.Attr3,
		Attr4:       req.Attr4,
		Status:      req.Status,
	}
}

type SysProSetProcessReq struct {
	ProcessID   int64  `gorm:"column:process_id;primaryKey;autoIncrement:true;comment:工序ID" json:"process_id"`   // 工序ID
	ProcessCode string `gorm:"column:process_code;not null;comment:工序编码" json:"process_code" binding:"required"` // 工序编码
	ProcessName string `gorm:"column:process_name;not null;comment:工序名称" json:"process_name" binding:"required"` // 工序名称
	Attention   string `gorm:"column:attention;comment:工艺要求" json:"attention"`                                   // 工艺要求
	Remark      string `gorm:"column:remark;comment:备注" json:"remark"`                                           // 备注
	Attr1       string `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                          // 预留字段1
	Attr2       string `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                          // 预留字段2
	Attr3       int32  `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                          // 预留字段3
	Attr4       int32  `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                          // 预留字段4
	Status      bool   `gorm:"column:status;comment:是否启用（0禁用 1启用）" json:"status"`                                // 是否启用（0禁用 1启用）
}

type SysProProcessListReq struct {
	ProcessCode string `gorm:"column:process_code;not null;comment:工序编码" json:"process_code"` // 工序编码
	ProcessName string `gorm:"column:process_name;not null;comment:工序名称" json:"process_name"` // 工序名称
	Status      *bool  `gorm:"column:status;comment:是否启用（0禁用 1启用）" json:"status"`             // 是否启用（0禁用 1启用）
	baize.BaseEntityDQL
}

// SysProProcessListData 工序列表
type SysProProcessListData struct {
	Rows  []*SysProProcess `json:"rows"`
	Total int64            `json:"total"`
}

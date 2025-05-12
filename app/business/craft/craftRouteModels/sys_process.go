package craftRouteModels

import (
	"nova-factory-server/app/baize"
	"time"
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

type SysProProcessListReq struct {
	ProcessCode string `gorm:"column:process_code;not null;comment:工序编码" json:"process_code" binding:"required"` // 工序编码
	ProcessName string `gorm:"column:process_name;not null;comment:工序名称" json:"process_name" binding:"required"` // 工序名称
	Status      *bool  `gorm:"column:status;comment:是否启用（0禁用 1启用）" json:"status"`                                // 是否启用（0禁用 1启用）
	baize.BaseEntityDQL
}

// SysProProcessListData 工序列表
type SysProProcessListData struct {
	Rows  []*SysProProcess `json:"rows"`
	Total int64            `json:"total"`
}

// SysProProcessContent 生产工序内容表
type SysProProcessContent struct {
	ContentID   int64     `gorm:"column:content_id;primaryKey;autoIncrement:true;comment:内容ID" json:"content_id"` // 内容ID
	ProcessID   int64     `gorm:"column:process_id;not null;comment:工序ID" json:"process_id"`                      // 工序ID
	OrderNum    int32     `gorm:"column:order_num;comment:顺序编号" json:"order_num"`                                 // 顺序编号
	ContentText string    `gorm:"column:content_text;comment:内容说明" json:"content_text"`                           // 内容说明
	Device      string    `gorm:"column:device;comment:辅助设备" json:"device"`                                       // 辅助设备
	Material    string    `gorm:"column:material;comment:辅助材料" json:"material"`                                   // 辅助材料
	DocURL      string    `gorm:"column:doc_url;comment:材料URL" json:"doc_url"`                                    // 材料URL
	Remark      string    `gorm:"column:remark;comment:备注" json:"remark"`                                         // 备注
	Attr1       string    `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                        // 预留字段1
	Attr2       string    `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                        // 预留字段2
	Attr3       int32     `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                        // 预留字段3
	Attr4       int32     `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                        // 预留字段4
	CreateBy    string    `gorm:"column:create_by;comment:创建者" json:"create_by"`                                  // 创建者
	CreateTime  time.Time `gorm:"column:create_time;comment:创建时间" json:"create_time"`                             // 创建时间
	UpdateBy    string    `gorm:"column:update_by;comment:更新者" json:"update_by"`                                  // 更新者
	UpdateTime  time.Time `gorm:"column:update_time;comment:更新时间" json:"update_time"`                             // 更新时间
}

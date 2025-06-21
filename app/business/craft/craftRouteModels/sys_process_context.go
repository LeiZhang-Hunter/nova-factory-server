package craftRouteModels

import (
	"nova-factory-server/app/baize"
)

// SysProProcessContent 生产工序内容表
type SysProProcessContent struct {
	ContentID      int64  `gorm:"column:content_id;primaryKey;autoIncrement:true;comment:内容ID" json:"content_id"` // 内容ID
	ProcessID      int64  `gorm:"column:process_id;not null;comment:工序ID" json:"process_id" binding:"required"`   // 工序ID
	OrderNum       int32  `gorm:"column:order_num;comment:顺序编号" json:"order_num"`                                 // 顺序编号
	ContentText    string `gorm:"column:content_text;comment:内容说明" json:"content_text" binding:"required"`        // 内容说明
	Device         string `gorm:"column:device;comment:辅助设备" json:"device"`                                       // 辅助设备
	Material       string `gorm:"column:material;comment:辅助材料" json:"material"`                                   // 辅助材料
	DocURL         string `gorm:"column:doc_url;comment:材料URL" json:"doc_url"`                                    // 材料URL
	Remark         string `gorm:"column:remark;comment:备注" json:"remark"`                                         // 备注
	Attr1          string `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                        // 预留字段1
	Attr2          string `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                        // 预留字段2
	Attr3          int32  `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                        // 预留字段3
	Attr4          int32  `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                        // 预留字段4
	CreateUserName string `json:"createUserName" gorm:"-"`
	UpdateUserName string `json:"updateUserName" gorm:"-"`
	baize.BaseEntity
}

func NewSysProProcessContent(context *SysProSetProcessContent) *SysProProcessContent {
	return &SysProProcessContent{
		ContentID:   context.ContentID,
		ProcessID:   context.ProcessID,
		OrderNum:    context.OrderNum,
		ContentText: context.ContentText,
		Device:      context.Device,
		Material:    context.Material,
		DocURL:      context.DocURL,
		Remark:      context.Remark,
		Attr1:       context.Attr1,
		Attr2:       context.Attr2,
		Attr3:       context.Attr3,
		Attr4:       context.Attr4,
	}
}

type SysProSetProcessContent struct {
	ContentID   int64  `gorm:"column:content_id;primaryKey;autoIncrement:true;comment:内容ID" json:"content_id"`      // 内容ID
	ProcessID   int64  `gorm:"column:process_id;not null;comment:工序ID" json:"process_id,string" binding:"required"` // 工序ID
	OrderNum    int32  `gorm:"column:order_num;comment:顺序编号" json:"order_num"`                                      // 顺序编号
	ContentText string `gorm:"column:content_text;comment:内容说明" json:"content_text" binding:"required"`             // 内容说明
	Device      string `gorm:"column:device;comment:辅助设备" json:"device"`                                            // 辅助设备
	Material    string `gorm:"column:material;comment:辅助材料" json:"material"`                                        // 辅助材料
	DocURL      string `gorm:"column:doc_url;comment:材料URL" json:"doc_url"`                                         // 材料URL
	Remark      string `gorm:"column:remark;comment:备注" json:"remark"`                                              // 备注
	Attr1       string `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                                             // 预留字段1
	Attr2       string `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                                             // 预留字段2
	Attr3       int32  `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                                             // 预留字段3
	Attr4       int32  `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                                             // 预留字段4
}

func OfSysProSetProcessContent(ctx *SysProProcessContent) *SysProSetProcessContent {
	return &SysProSetProcessContent{
		ContentID:   ctx.ContentID,
		ProcessID:   ctx.ProcessID,
		OrderNum:    ctx.OrderNum,
		ContentText: ctx.ContentText,
		Device:      ctx.Device,
		Material:    ctx.Material,
		DocURL:      ctx.DocURL,
		Remark:      ctx.Remark,
		Attr1:       ctx.Attr1,
		Attr2:       ctx.Attr2,
		Attr3:       ctx.Attr3,
		Attr4:       ctx.Attr4,
	}
}

type SysProProcessContextListReq struct {
	ProcessID int64 `gorm:"column:process_id;not null;comment:工序ID" json:"process_id"` // 工序ID
	baize.BaseEntityDQL
}

type SysProProcessContextListData struct {
	Rows  []*SysProProcessContent `json:"rows"`
	Total int64                   `json:"total"`
}

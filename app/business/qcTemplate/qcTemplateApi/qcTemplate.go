package qcTemplateApi

import (
	"time"
)

// QcTemplate 检测模板表
type QcTemplate struct {
	TemplateId   int64      `json:"templateId" db:"template_id"`     // 检测模板ID
	TemplateCode string     `json:"templateCode" db:"template_code"` // 检测模板编号
	TemplateName string     `json:"templateName" db:"template_name"` // 检测模板名称
	QcTypes      string     `json:"qcTypes" db:"qc_types"`           // 检测种类
	EnableFlag   *string    `json:"enableFlag" db:"enable_flag"`     // 是否启用
	Remark       *string    `json:"remark" db:"remark"`              // 备注
	Attr1        *string    `json:"attr1" db:"attr1"`                // 预留字段1
	Attr2        *string    `json:"attr2" db:"attr2"`                // 预留字段2
	Attr3        *int       `json:"attr3" db:"attr3"`                // 预留字段3
	Attr4        *int       `json:"attr4" db:"attr4"`                // 预留字段4
	CreateBy     string     `json:"createBy" db:"create_by"`         // 创建者
	CreateTime   time.Time  `json:"createTime" db:"create_time"`     // 创建时间
	UpdateBy     string     `json:"updateBy" db:"update_by"`         // 更新者
	UpdateTime   *time.Time `json:"updateTime" db:"update_time"`     // 更新时间
}

// QcTemplateQueryReq 查询请求
type QcTemplateQueryReq struct {
	TemplateCode *string `json:"templateCode" form:"templateCode"`
	TemplateName *string `json:"templateName" form:"templateName"`
	QcTypes      *string `json:"qcTypes" form:"qcTypes"`
	EnableFlag   *string `json:"enableFlag" form:"enableFlag"`
	CreateBy     *string `json:"createBy" form:"createBy"`
	PageNum      int     `json:"pageNum" form:"pageNum"`
	PageSize     int     `json:"pageSize" form:"pageSize"`
}

// QcTemplateCreateReq 创建请求
type QcTemplateCreateReq struct {
	TemplateCode string  `json:"templateCode" binding:"required"` // 检测模板编号
	TemplateName string  `json:"templateName" binding:"required"` // 检测模板名称
	QcTypes      string  `json:"qcTypes" binding:"required"`      // 检测种类
	EnableFlag   *string `json:"enableFlag"`                      // 是否启用
	Remark       *string `json:"remark"`                          // 备注
	Attr1        *string `json:"attr1"`                           // 预留字段1
	Attr2        *string `json:"attr2"`                           // 预留字段2
	Attr3        *int    `json:"attr3"`                           // 预留字段3
	Attr4        *int    `json:"attr4"`                           // 预留字段4
}

// QcTemplateUpdateReq 更新请求
type QcTemplateUpdateReq struct {
	TemplateId   int64   `json:"templateId" binding:"required"` // 检测模板ID
	TemplateCode *string `json:"templateCode"`                  // 检测模板编号
	TemplateName *string `json:"templateName"`                  // 检测模板名称
	QcTypes      *string `json:"qcTypes"`                       // 检测种类
	EnableFlag   *string `json:"enableFlag"`                    // 是否启用
	Remark       *string `json:"remark"`                        // 备注
	Attr1        *string `json:"attr1"`                         // 预留字段1
	Attr2        *string `json:"attr2"`                         // 预留字段2
	Attr3        *int    `json:"attr3"`                         // 预留字段3
	Attr4        *int    `json:"attr4"`                         // 预留字段4
}

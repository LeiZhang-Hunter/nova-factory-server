package qcTemplateModels

import (
	"gorm.io/gorm"
)

// QcTemplate 检测模板表
type QcTemplate struct {
	gorm.Model
	TemplateCode string  `json:"templateCode" db:"template_code" gorm:"not null"` // 检测模板编号
	TemplateName string  `json:"templateName" db:"template_name" gorm:"not null"` // 检测模板名称
	QcTypes      string  `json:"qcTypes" db:"qc_types" gorm:"not null"`           // 检测种类
	EnableFlag   *string `json:"enableFlag" db:"enable_flag"`                     // 是否启用
	Remark       *string `json:"remark" db:"remark"`                              // 备注
	Attr1        *string `json:"attr1" db:"attr1"`                                // 预留字段1
	Attr2        *string `json:"attr2" db:"attr2"`                                // 预留字段2
	Attr3        *int    `json:"attr3" db:"attr3"`                                // 预留字段3
	Attr4        *int    `json:"attr4" db:"attr4"`                                // 预留字段4
	CreateBy     *string `json:"createBy" db:"create_by"`                         // 创建者
	UpdateBy     *string `json:"updateBy" db:"update_by"`                         // 更新者
	CreateById   int64   `json:"createById" db:"create_by_id" gorm:"not null"`    // 创建者ID
	UpdateById   int64   `json:"updateById" db:"update_by_id" gorm:"not null"`    // 更新者ID
}

// TableName 指定表名
func (QcTemplate) TableName() string {
	return "qc_template"
}

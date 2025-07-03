package qcOqcModels

import (
	"time"

	"gorm.io/gorm"
)

// QcOqc 出货检验单表
type QcOqc struct {
	gorm.Model
	OqcCode                string     `json:"oqcCode" db:"oqc_code" gorm:"not null"`                // 出货检验单编号
	OqcName                *string    `json:"oqcName" db:"oqc_name"`                                // 出货检验单名称
	TemplateId             int64      `json:"templateId" db:"template_id" gorm:"not null"`          // 检验模板ID
	SourceDocId            *int64     `json:"sourceDocId" db:"source_doc_id"`                       // 来源单据ID
	SourceDocType          *string    `json:"sourceDocType" db:"source_doc_type"`                   // 来源单据类型
	SourceDocCode          *string    `json:"sourceDocCode" db:"source_doc_code"`                   // 来源单据编号
	SourceLineId           *int64     `json:"sourceLineId" db:"source_line_id"`                     // 来源单据行ID
	ClientId               int64      `json:"clientId" db:"client_id" gorm:"not null"`              // 客户ID
	ClientCode             string     `json:"clientCode" db:"client_code" gorm:"not null"`          // 客户编码
	ClientName             string     `json:"clientName" db:"client_name" gorm:"not null"`          // 客户名称
	BatchCode              *string    `json:"batchCode" db:"batch_code"`                            // 批次号
	ItemId                 int64      `json:"itemId" db:"item_id" gorm:"not null"`                  // 产品物料ID
	ItemCode               *string    `json:"itemCode" db:"item_code"`                              // 产品物料编码
	ItemName               *string    `json:"itemName" db:"item_name"`                              // 产品物料名称
	Specification          *string    `json:"specification" db:"specification"`                     // 规格型号
	UnitOfMeasure          *string    `json:"unitOfMeasure" db:"unit_of_measure"`                   // 单位
	QuantityMinCheck       *float64   `json:"quantityMinCheck" db:"quantity_min_check"`             // 最低检测数
	QuantityMaxUnqualified *float64   `json:"quantityMaxUnqualified" db:"quantity_max_unqualified"` // 最大不合格数
	QuantityOut            float64    `json:"quantityOut" db:"quantity_out" gorm:"not null"`        // 发货数量
	QuantityCheck          float64    `json:"quantityCheck" db:"quantity_check" gorm:"not null"`    // 本次检测数量
	QuantityUnqualified    *float64   `json:"quantityUnqualified" db:"quantity_unqualified"`        // 不合格数
	QuantityQualified      *float64   `json:"quantityQualified" db:"quantity_qualified"`            // 合格数量
	CrRate                 *float64   `json:"crRate" db:"cr_rate"`                                  // 致命缺陷率
	MajRate                *float64   `json:"majRate" db:"maj_rate"`                                // 严重缺陷率
	MinRate                *float64   `json:"minRate" db:"min_rate"`                                // 轻微缺陷率
	CrQuantity             *float64   `json:"crQuantity" db:"cr_quantity"`                          // 致命缺陷数量
	MajQuantity            *float64   `json:"majQuantity" db:"maj_quantity"`                        // 严重缺陷数量
	MinQuantity            *float64   `json:"minQuantity" db:"min_quantity"`                        // 轻微缺陷数量
	CheckResult            *string    `json:"checkResult" db:"check_result"`                        // 检测结果
	OutDate                *time.Time `json:"outDate" db:"out_date"`                                // 出货日期
	InspectDate            *time.Time `json:"inspectDate" db:"inspect_date"`                        // 检测日期
	Inspector              *string    `json:"inspector" db:"inspector"`                             // 检测人员
	Status                 *string    `json:"status" db:"status"`                                   // 单据状态
	Remark                 *string    `json:"remark" db:"remark"`                                   // 备注
	Attr1                  *string    `json:"attr1" db:"attr1"`                                     // 预留字段1
	Attr2                  *string    `json:"attr2" db:"attr2"`                                     // 预留字段2
	Attr3                  *int       `json:"attr3" db:"attr3"`                                     // 预留字段3
	Attr4                  *int       `json:"attr4" db:"attr4"`                                     // 预留字段4
	CreateBy               *string    `json:"createBy" db:"create_by"`                              // 创建者
	UpdateBy               *string    `json:"updateBy" db:"update_by"`                              // 更新者
	CreateById             int64      `json:"createById" db:"create_by_id" gorm:"not null"`         // 创建者ID
	UpdateById             int64      `json:"updateById" db:"update_by_id" gorm:"not null"`         // 更新者ID
}

// TableName 指定表名
func (QcOqc) TableName() string {
	return "qc_oqc"
}

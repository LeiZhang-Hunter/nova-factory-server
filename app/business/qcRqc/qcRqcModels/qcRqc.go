package qcRqcModels

import (
	"time"

	"gorm.io/gorm"
)

// QcRqc 退料检验单表
type QcRqc struct {
	gorm.Model
	RqcCode             string     `json:"rqcCode" db:"rqc_code" gorm:"not null"`         // 检验单编号
	RqcName             *string    `json:"rqcName" db:"rqc_name"`                         // 检验单名称
	TemplateId          int64      `json:"templateId" db:"template_id" gorm:"not null"`   // 检验模板ID
	SourceDocId         *int64     `json:"sourceDocId" db:"source_doc_id"`                // 来源单据ID
	SourceDocType       *string    `json:"sourceDocType" db:"source_doc_type"`            // 来源单据类型
	SourceDocCode       *string    `json:"sourceDocCode" db:"source_doc_code"`            // 来源单据编号
	SourceLineId        *int64     `json:"sourceLineId" db:"source_line_id"`              // 来源单据行ID
	ItemId              int64      `json:"itemId" db:"item_id" gorm:"not null"`           // 产品物料ID
	ItemCode            *string    `json:"itemCode" db:"item_code"`                       // 产品物料编码
	ItemName            *string    `json:"itemName" db:"item_name"`                       // 产品物料名称
	Specification       *string    `json:"specification" db:"specification"`              // 规格型号
	UnitOfMeasure       *string    `json:"unitOfMeasure" db:"unit_of_measure"`            // 单位
	UnitName            *string    `json:"unitName" db:"unit_name"`                       // 单位名称
	BatchId             *int64     `json:"batchId" db:"batch_id"`                         // 批次ID
	BatchCode           *string    `json:"batchCode" db:"batch_code"`                     // 批次号
	QuantityCheck       *float64   `json:"quantityCheck" db:"quantity_check"`             // 检测数量
	QuantityUnqualified *float64   `json:"quantityUnqualified" db:"quantity_unqualified"` // 不合格数
	QuantityQualified   *float64   `json:"quantityQualified" db:"quantity_qualified"`     // 合格品数量
	CheckResult         *string    `json:"checkResult" db:"check_result"`                 // 检测结果
	InspectDate         *time.Time `json:"inspectDate" db:"inspect_date"`                 // 检测日期
	UserId              *int64     `json:"userId" db:"user_id"`                           // 检测人员ID
	UserName            *string    `json:"userName" db:"user_name"`                       // 检测人员名称
	NickName            *string    `json:"nickName" db:"nick_name"`                       // 检测人员
	Status              *string    `json:"status" db:"status"`                            // 单据状态
	Remark              *string    `json:"remark" db:"remark"`                            // 备注
	Attr1               *string    `json:"attr1" db:"attr1"`                              // 预留字段1
	Attr2               *string    `json:"attr2" db:"attr2"`                              // 预留字段2
	Attr3               *int       `json:"attr3" db:"attr3"`                              // 预留字段3
	Attr4               *int       `json:"attr4" db:"attr4"`                              // 预留字段4
	CreateBy            *string    `json:"createBy" db:"create_by"`                       // 创建者
	UpdateBy            *string    `json:"updateBy" db:"update_by"`                       // 更新者
	CreateById          int64      `json:"createById" db:"create_by_id" gorm:"not null"`  // 创建者ID
	UpdateById          int64      `json:"updateById" db:"update_by_id" gorm:"not null"`  // 更新者ID
}

// TableName 指定表名
func (QcRqc) TableName() string {
	return "qc_rqc"
}

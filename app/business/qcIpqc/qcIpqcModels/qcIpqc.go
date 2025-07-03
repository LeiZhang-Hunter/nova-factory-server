package qcIpqcModels

import (
	"time"

	"gorm.io/gorm"
)

// QcIpqc 过程检验单表
type QcIpqc struct {
	gorm.Model
	IpqcCode              string     `json:"ipqcCode" db:"ipqc_code" gorm:"not null"`
	IpqcName              *string    `json:"ipqcName" db:"ipqc_name"`
	IpqcType              *string    `json:"ipqcType" db:"ipqc_type"`
	TemplateId            *int64     `json:"templateId" db:"template_id"`
	SourceDocId           *int64     `json:"sourceDocId" db:"source_doc_id"`
	SourceDocType         *string    `json:"sourceDocType" db:"source_doc_type"`
	SourceDocCode         *string    `json:"sourceDocCode" db:"source_doc_code"`
	SourceLineId          *int64     `json:"sourceLineId" db:"source_line_id"`
	WorkorderId           *int64     `json:"workorderId" db:"workorder_id"`
	WorkorderCode         *string    `json:"workorderCode" db:"workorder_code"`
	WorkorderName         *string    `json:"workorderName" db:"workorder_name"`
	TaskId                *int64     `json:"taskId" db:"task_id"`
	TaskCode              *string    `json:"taskCode" db:"task_code"`
	TaskName              *string    `json:"taskName" db:"task_name"`
	WorkstationId         *int64     `json:"workstationId" db:"workstation_id"`
	WorkstationCode       *string    `json:"workstationCode" db:"workstation_code"`
	WorkstationName       *string    `json:"workstationName" db:"workstation_name"`
	ProcessId             *int64     `json:"processId" db:"process_id"`
	ProcessCode           *string    `json:"processCode" db:"process_code"`
	ProcessName           *string    `json:"processName" db:"process_name"`
	ItemId                *int64     `json:"itemId" db:"item_id"`
	ItemCode              *string    `json:"itemCode" db:"item_code"`
	ItemName              *string    `json:"itemName" db:"item_name"`
	Specification         *string    `json:"specification" db:"specification"`
	UnitOfMeasure         *string    `json:"unitOfMeasure" db:"unit_of_measure"`
	UnitName              *string    `json:"unitName" db:"unit_name"`
	QuantityCheck         *float64   `json:"quantityCheck" db:"quantity_check"`
	QuantityUnqualified   *float64   `json:"quantityUnqualified" db:"quantity_unqualified"`
	QuantityQualified     *float64   `json:"quantityQualified" db:"quantity_qualified"`
	QuantityLaborScrap    *float64   `json:"quantityLaborScrap" db:"quantity_labor_scrap"`
	QuantityMaterialScrap *float64   `json:"quantityMaterialScrap" db:"quantity_material_scrap"`
	QuantityOtherScrap    *float64   `json:"quantityOtherScrap" db:"quantity_other_scrap"`
	CrRate                *float64   `json:"crRate" db:"cr_rate"`
	MajRate               *float64   `json:"majRate" db:"maj_rate"`
	MinRate               *float64   `json:"minRate" db:"min_rate"`
	CrQuantity            *float64   `json:"crQuantity" db:"cr_quantity"`
	MajQuantity           *float64   `json:"majQuantity" db:"maj_quantity"`
	MinQuantity           *float64   `json:"minQuantity" db:"min_quantity"`
	CheckResult           *string    `json:"checkResult" db:"check_result"`
	InspectDate           *time.Time `json:"inspectDate" db:"inspect_date"`
	Inspector             *string    `json:"inspector" db:"inspector"`
	Status                *string    `json:"status" db:"status"`
	Remark                *string    `json:"remark" db:"remark"`
	Attr1                 *string    `json:"attr1" db:"attr1"`
	Attr2                 *string    `json:"attr2" db:"attr2"`
	Attr3                 *int       `json:"attr3" db:"attr3"`
	Attr4                 *int       `json:"attr4" db:"attr4"`
	CreateBy              *string    `json:"createBy" db:"create_by"`
	UpdateBy              *string    `json:"updateBy" db:"update_by"`
	CreateById            *int64     `json:"createById" db:"create_by_id"`
	UpdateById            *int64     `json:"updateById" db:"update_by_id"`
}

func (QcIpqc) TableName() string {
	return "qc_ipqc"
}

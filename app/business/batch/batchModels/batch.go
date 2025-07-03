package batchModels

import (
	"gorm.io/gorm"
)

type Batch struct {
	gorm.Model      `gorm:"table:wm_batch"`
	BatchCode       string  `json:"batchCode" db:"batch_code" gorm:"not null"`
	ItemId          int64   `json:"itemId" db:"item_id" gorm:"not null"`
	ItemCode        *string `json:"itemCode" db:"item_code"`
	ItemName        *string `json:"itemName" db:"item_name"`
	Specification   *string `json:"specification" db:"specification"`
	UnitOfMeasure   *string `json:"unitOfMeasure" db:"unit_of_measure"`
	VendorId        *int64  `json:"vendorId" db:"vendor_id"`
	VendorCode      *string `json:"vendorCode" db:"vendor_code"`
	VendorName      *string `json:"vendorName" db:"vendor_name"`
	VendorNick      *string `json:"vendorNick" db:"vendor_nick"`
	ClientId        *int64  `json:"clientId" db:"client_id"`
	ClientCode      *string `json:"clientCode" db:"client_code"`
	ClientName      *string `json:"clientName" db:"client_name"`
	ClientNick      *string `json:"clientNick" db:"client_nick"`
	CoCode          *string `json:"coCode" db:"co_code"`
	PoCode          *string `json:"poCode" db:"po_code"`
	WorkorderId     *int64  `json:"workorderId" db:"workorder_id"`
	WorkorderCode   *string `json:"workorderCode" db:"workorder_code"`
	TaskId          *int64  `json:"taskId" db:"task_id"`
	TaskCode        *string `json:"taskCode" db:"task_code"`
	WorkstationId   *int64  `json:"workstationId" db:"workstation_id"`
	WorkstationCode *string `json:"workstationCode" db:"workstation_code"`
	ToolId          *int64  `json:"toolId" db:"tool_id"`
	ToolCode        *string `json:"toolCode" db:"tool_code"`
	MoldId          *int64  `json:"moldId" db:"mold_id"`
	MoldCode        *string `json:"moldCode" db:"mold_code"`
	LotNumber       *string `json:"lotNumber" db:"lot_number"`
	QualityStatus   *string `json:"qualityStatus" db:"quality_status"`
	Remark          *string `json:"remark" db:"remark"`
	Attr1           *string `json:"attr1" db:"attr1"`
	Attr2           *string `json:"attr2" db:"attr2"`
	Attr3           *int    `json:"attr3" db:"attr3"`
	Attr4           *int    `json:"attr4" db:"attr4"`
	CreateBy        *string `json:"createBy" db:"create_by"`
	UpdateBy        *string `json:"updateBy" db:"update_by"`
	CreateById      *int64  `json:"createById" db:"create_by_id"`
	UpdateById      *int64  `json:"updateById" db:"update_by_id"`
}

func (Batch) TableName() string {
	return "wm_batch"
}

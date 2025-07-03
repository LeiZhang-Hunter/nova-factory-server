package batchApi

import (
	"time"
)

// Batch 批次主表 API 结构体
type Batch struct {
	BatchId         int64      `json:"batchId"`
	BatchCode       string     `json:"batchCode"`
	ItemId          int64      `json:"itemId"`
	ItemCode        *string    `json:"itemCode"`
	ItemName        *string    `json:"itemName"`
	Specification   *string    `json:"specification"`
	UnitOfMeasure   *string    `json:"unitOfMeasure"`
	VendorId        *int64     `json:"vendorId"`
	VendorCode      *string    `json:"vendorCode"`
	VendorName      *string    `json:"vendorName"`
	VendorNick      *string    `json:"vendorNick"`
	ClientId        *int64     `json:"clientId"`
	ClientCode      *string    `json:"clientCode"`
	ClientName      *string    `json:"clientName"`
	ClientNick      *string    `json:"clientNick"`
	CoCode          *string    `json:"coCode"`
	PoCode          *string    `json:"poCode"`
	WorkorderId     *int64     `json:"workorderId"`
	WorkorderCode   *string    `json:"workorderCode"`
	TaskId          *int64     `json:"taskId"`
	TaskCode        *string    `json:"taskCode"`
	WorkstationId   *int64     `json:"workstationId"`
	WorkstationCode *string    `json:"workstationCode"`
	ToolId          *int64     `json:"toolId"`
	ToolCode        *string    `json:"toolCode"`
	MoldId          *int64     `json:"moldId"`
	MoldCode        *string    `json:"moldCode"`
	LotNumber       *string    `json:"lotNumber"`
	QualityStatus   *string    `json:"qualityStatus"`
	Remark          *string    `json:"remark"`
	Attr1           *string    `json:"attr1"`
	Attr2           *string    `json:"attr2"`
	Attr3           *int       `json:"attr3"`
	Attr4           *int       `json:"attr4"`
	CreateBy        *string    `json:"createBy"`
	UpdateBy        *string    `json:"updateBy"`
	CreateTime      *time.Time `json:"createTime"`
	UpdateTime      *time.Time `json:"updateTime"`
}

type BatchQueryReq struct {
	BatchCode     *string `json:"batchCode" form:"batchCode"`
	ItemId        *int64  `json:"itemId" form:"itemId"`
	ItemCode      *string `json:"itemCode" form:"itemCode"`
	ItemName      *string `json:"itemName" form:"itemName"`
	Specification *string `json:"specification" form:"specification"`
	PageNum       int     `json:"pageNum" form:"pageNum"`
	PageSize      int     `json:"pageSize" form:"pageSize"`
}

type BatchCreateReq struct {
	BatchCode       string  `json:"batchCode" binding:"required"`
	ItemId          int64   `json:"itemId" binding:"required"`
	ItemCode        *string `json:"itemCode"`
	ItemName        *string `json:"itemName"`
	Specification   *string `json:"specification"`
	UnitOfMeasure   *string `json:"unitOfMeasure"`
	VendorId        *int64  `json:"vendorId"`
	VendorCode      *string `json:"vendorCode"`
	VendorName      *string `json:"vendorName"`
	VendorNick      *string `json:"vendorNick"`
	ClientId        *int64  `json:"clientId"`
	ClientCode      *string `json:"clientCode"`
	ClientName      *string `json:"clientName"`
	ClientNick      *string `json:"clientNick"`
	CoCode          *string `json:"coCode"`
	PoCode          *string `json:"poCode"`
	WorkorderId     *int64  `json:"workorderId"`
	WorkorderCode   *string `json:"workorderCode"`
	TaskId          *int64  `json:"taskId"`
	TaskCode        *string `json:"taskCode"`
	WorkstationId   *int64  `json:"workstationId"`
	WorkstationCode *string `json:"workstationCode"`
	ToolId          *int64  `json:"toolId"`
	ToolCode        *string `json:"toolCode"`
	MoldId          *int64  `json:"moldId"`
	MoldCode        *string `json:"moldCode"`
	LotNumber       *string `json:"lotNumber"`
	QualityStatus   *string `json:"qualityStatus"`
	Remark          *string `json:"remark"`
	Attr1           *string `json:"attr1"`
	Attr2           *string `json:"attr2"`
	Attr3           *int    `json:"attr3"`
	Attr4           *int    `json:"attr4"`
}

type BatchUpdateReq struct {
	BatchId         int64   `json:"batchId" binding:"required"`
	BatchCode       string  `json:"batchCode" binding:"required"`
	ItemId          int64   `json:"itemId" binding:"required"`
	ItemCode        *string `json:"itemCode"`
	ItemName        *string `json:"itemName"`
	Specification   *string `json:"specification"`
	UnitOfMeasure   *string `json:"unitOfMeasure"`
	VendorId        *int64  `json:"vendorId"`
	VendorCode      *string `json:"vendorCode"`
	VendorName      *string `json:"vendorName"`
	VendorNick      *string `json:"vendorNick"`
	ClientId        *int64  `json:"clientId"`
	ClientCode      *string `json:"clientCode"`
	ClientName      *string `json:"clientName"`
	ClientNick      *string `json:"clientNick"`
	CoCode          *string `json:"coCode"`
	PoCode          *string `json:"poCode"`
	WorkorderId     *int64  `json:"workorderId"`
	WorkorderCode   *string `json:"workorderCode"`
	TaskId          *int64  `json:"taskId"`
	TaskCode        *string `json:"taskCode"`
	WorkstationId   *int64  `json:"workstationId"`
	WorkstationCode *string `json:"workstationCode"`
	ToolId          *int64  `json:"toolId"`
	ToolCode        *string `json:"toolCode"`
	MoldId          *int64  `json:"moldId"`
	MoldCode        *string `json:"moldCode"`
	LotNumber       *string `json:"lotNumber"`
	QualityStatus   *string `json:"qualityStatus"`
	Remark          *string `json:"remark"`
	Attr1           *string `json:"attr1"`
	Attr2           *string `json:"attr2"`
	Attr3           *int    `json:"attr3"`
	Attr4           *int    `json:"attr4"`
}

type BatchDeleteReq struct {
	BatchIds []int64 `json:"batchIds" binding:"required"`
}

type BatchListRes struct {
	Rows  []*Batch `json:"rows"`
	Total int64    `json:"total"`
}

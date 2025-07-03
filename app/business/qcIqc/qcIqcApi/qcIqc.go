package qcIqcApi

import (
	"time"
)

// QcIqcQueryReq 查询请求
type QcIqcQueryReq struct {
	IqcCode                *string    `json:"iqcCode" form:"iqcCode"`
	IqcName                *string    `json:"iqcName" form:"iqcName"`
	TemplateId             *int64     `json:"templateId" form:"templateId"`
	SourceDocId            *int64     `json:"sourceDocId" form:"sourceDocId"`
	SourceDocType          *string    `json:"sourceDocType" form:"sourceDocType"`
	SourceDocCode          *string    `json:"sourceDocCode" form:"sourceDocCode"`
	SourceLineId           *int64     `json:"sourceLineId" form:"sourceLineId"`
	VendorId               *int64     `json:"vendorId" form:"vendorId"`
	VendorCode             *string    `json:"vendorCode" form:"vendorCode"`
	VendorName             *string    `json:"vendorName" form:"vendorName"`
	VendorNick             *string    `json:"vendorNick" form:"vendorNick"`
	VendorBatch            *string    `json:"vendorBatch" form:"vendorBatch"`
	ItemId                 *int64     `json:"itemId" form:"itemId"`
	ItemCode               *string    `json:"itemCode" form:"itemCode"`
	ItemName               *string    `json:"itemName" form:"itemName"`
	Specification          *string    `json:"specification" form:"specification"`
	UnitOfMeasure          *string    `json:"unitOfMeasure" form:"unitOfMeasure"`
	UnitName               *string    `json:"unitName" form:"unitName"`
	QuantityMinCheck       *int       `json:"quantityMinCheck" form:"quantityMinCheck"`
	QuantityMaxUnqualified *int       `json:"quantityMaxUnqualified" form:"quantityMaxUnqualified"`
	QuantityRecived        *float64   `json:"quantityRecived" form:"quantityRecived"`
	QuantityCheck          *int       `json:"quantityCheck" form:"quantityCheck"`
	QuantityQualified      *int       `json:"quantityQualified" form:"quantityQualified"`
	QuantityUnqualified    *int       `json:"quantityUnqualified" form:"quantityUnqualified"`
	CrRate                 *float64   `json:"crRate" form:"crRate"`
	MajRate                *float64   `json:"majRate" form:"majRate"`
	MinRate                *float64   `json:"minRate" form:"minRate"`
	CrQuantity             *int       `json:"crQuantity" form:"crQuantity"`
	MajQuantity            *int       `json:"majQuantity" form:"majQuantity"`
	MinQuantity            *int       `json:"minQuantity" form:"minQuantity"`
	CheckResult            *string    `json:"checkResult" form:"checkResult"`
	ReciveDate             *time.Time `json:"reciveDate" form:"reciveDate"`
	InspectDate            *time.Time `json:"inspectDate" form:"inspectDate"`
	Inspector              *string    `json:"inspector" form:"inspector"`
	Status                 *string    `json:"status" form:"status"`
	CreateBy               *string    `json:"createBy" form:"createBy"`
	PageNum                int        `json:"pageNum" form:"pageNum"`
	PageSize               int        `json:"pageSize" form:"pageSize"`
}

// QcIqcCreateReq 创建请求
type QcIqcCreateReq struct {
	IqcName                string     `json:"iqcName" binding:"required"`         // 来料检验单名称
	TemplateId             int64      `json:"templateId" binding:"required"`      // 检验模板ID
	SourceDocId            *int64     `json:"sourceDocId"`                        // 来源单据ID
	SourceDocType          *string    `json:"sourceDocType"`                      // 来源单据类型
	SourceDocCode          *string    `json:"sourceDocCode"`                      // 来源单据编号
	SourceLineId           *int64     `json:"sourceLineId"`                       // 来源单据行ID
	VendorId               int64      `json:"vendorId" binding:"required"`        // 供应商ID
	VendorCode             string     `json:"vendorCode" binding:"required"`      // 供应商编码
	VendorName             string     `json:"vendorName" binding:"required"`      // 供应商名称
	VendorNick             *string    `json:"vendorNick"`                         // 供应商简称
	VendorBatch            *string    `json:"vendorBatch"`                        // 供应商批次号
	ItemId                 int64      `json:"itemId" binding:"required"`          // 产品物料ID
	ItemCode               *string    `json:"itemCode"`                           // 产品物料编码
	ItemName               *string    `json:"itemName"`                           // 产品物料名称
	Specification          *string    `json:"specification"`                      // 规格型号
	UnitOfMeasure          *string    `json:"unitOfMeasure"`                      // 单位
	UnitName               *string    `json:"unitName"`                           // 单位名称
	QuantityMinCheck       *int       `json:"quantityMinCheck"`                   // 最低检测数
	QuantityMaxUnqualified *int       `json:"quantityMaxUnqualified"`             // 最大不合格数
	QuantityRecived        float64    `json:"quantityRecived" binding:"required"` // 本次接收数量
	QuantityCheck          *int       `json:"quantityCheck"`                      // 本次检测数量
	QuantityQualified      *int       `json:"quantityQualified"`                  // 合格数
	QuantityUnqualified    *int       `json:"quantityUnqualified"`                // 不合格数
	CrRate                 *float64   `json:"crRate"`                             // 致命缺陷率
	MajRate                *float64   `json:"majRate"`                            // 严重缺陷率
	MinRate                *float64   `json:"minRate"`                            // 轻微缺陷率
	CrQuantity             *int       `json:"crQuantity"`                         // 致命缺陷数量
	MajQuantity            *int       `json:"majQuantity"`                        // 严重缺陷数量
	MinQuantity            *int       `json:"minQuantity"`                        // 轻微缺陷数量
	CheckResult            *string    `json:"checkResult"`                        // 检测结果
	ReciveDate             *time.Time `json:"reciveDate"`                         // 来料日期
	InspectDate            *time.Time `json:"inspectDate"`                        // 检测日期
	Inspector              *string    `json:"inspector"`                          // 检测人员
	Status                 *string    `json:"status"`                             // 单据状态
	Remark                 *string    `json:"remark"`                             // 备注
	Attr1                  *string    `json:"attr1"`                              // 预留字段1
	Attr2                  *string    `json:"attr2"`                              // 预留字段2
	Attr3                  *int       `json:"attr3"`                              // 预留字段3
	Attr4                  *int       `json:"attr4"`                              // 预留字段4
}

// QcIqcUpdateReq 更新请求
type QcIqcUpdateReq struct {
	IqcId                  int64      `json:"iqcId" binding:"required"` // 来料检验单ID
	IqcName                *string    `json:"iqcName"`                  // 来料检验单名称
	TemplateId             *int64     `json:"templateId"`               // 检验模板ID
	SourceDocId            *int64     `json:"sourceDocId"`              // 来源单据ID
	SourceDocType          *string    `json:"sourceDocType"`            // 来源单据类型
	SourceDocCode          *string    `json:"sourceDocCode"`            // 来源单据编号
	SourceLineId           *int64     `json:"sourceLineId"`             // 来源单据行ID
	VendorId               *int64     `json:"vendorId"`                 // 供应商ID
	VendorCode             *string    `json:"vendorCode"`               // 供应商编码
	VendorName             *string    `json:"vendorName"`               // 供应商名称
	VendorNick             *string    `json:"vendorNick"`               // 供应商简称
	VendorBatch            *string    `json:"vendorBatch"`              // 供应商批次号
	ItemId                 *int64     `json:"itemId"`                   // 产品物料ID
	ItemCode               *string    `json:"itemCode"`                 // 产品物料编码
	ItemName               *string    `json:"itemName"`                 // 产品物料名称
	Specification          *string    `json:"specification"`            // 规格型号
	UnitOfMeasure          *string    `json:"unitOfMeasure"`            // 单位
	UnitName               *string    `json:"unitName"`                 // 单位名称
	QuantityMinCheck       *int       `json:"quantityMinCheck"`         // 最低检测数
	QuantityMaxUnqualified *int       `json:"quantityMaxUnqualified"`   // 最大不合格数
	QuantityRecived        *float64   `json:"quantityRecived"`          // 本次接收数量
	QuantityCheck          *int       `json:"quantityCheck"`            // 本次检测数量
	QuantityQualified      *int       `json:"quantityQualified"`        // 合格数
	QuantityUnqualified    *int       `json:"quantityUnqualified"`      // 不合格数
	CrRate                 *float64   `json:"crRate"`                   // 致命缺陷率
	MajRate                *float64   `json:"majRate"`                  // 严重缺陷率
	MinRate                *float64   `json:"minRate"`                  // 轻微缺陷率
	CrQuantity             *int       `json:"crQuantity"`               // 致命缺陷数量
	MajQuantity            *int       `json:"majQuantity"`              // 严重缺陷数量
	MinQuantity            *int       `json:"minQuantity"`              // 轻微缺陷数量
	CheckResult            *string    `json:"checkResult"`              // 检测结果
	ReciveDate             *time.Time `json:"reciveDate"`               // 来料日期
	InspectDate            *time.Time `json:"inspectDate"`              // 检测日期
	Inspector              *string    `json:"inspector"`                // 检测人员
	Status                 *string    `json:"status"`                   // 单据状态
	Remark                 *string    `json:"remark"`                   // 备注
	Attr1                  *string    `json:"attr1"`                    // 预留字段1
	Attr2                  *string    `json:"attr2"`                    // 预留字段2
	Attr3                  *int       `json:"attr3"`                    // 预留字段3
	Attr4                  *int       `json:"attr4"`                    // 预留字段4
}

// QcIqcData 来料检验单数据
type QcIqcData struct {
	IqcId                  int64      `json:"iqcId"`
	IqcCode                string     `json:"iqcCode"`
	IqcName                string     `json:"iqcName"`
	TemplateId             int64      `json:"templateId"`
	SourceDocId            *int64     `json:"sourceDocId"`
	SourceDocType          *string    `json:"sourceDocType"`
	SourceDocCode          *string    `json:"sourceDocCode"`
	SourceLineId           *int64     `json:"sourceLineId"`
	VendorId               int64      `json:"vendorId"`
	VendorCode             string     `json:"vendorCode"`
	VendorName             string     `json:"vendorName"`
	VendorNick             *string    `json:"vendorNick"`
	VendorBatch            *string    `json:"vendorBatch"`
	ItemId                 int64      `json:"itemId"`
	ItemCode               *string    `json:"itemCode"`
	ItemName               *string    `json:"itemName"`
	Specification          *string    `json:"specification"`
	UnitOfMeasure          *string    `json:"unitOfMeasure"`
	UnitName               *string    `json:"unitName"`
	QuantityMinCheck       int        `json:"quantityMinCheck"`
	QuantityMaxUnqualified int        `json:"quantityMaxUnqualified"`
	QuantityRecived        float64    `json:"quantityRecived"`
	QuantityCheck          int        `json:"quantityCheck"`
	QuantityQualified      int        `json:"quantityQualified"`
	QuantityUnqualified    int        `json:"quantityUnqualified"`
	CrRate                 float64    `json:"crRate"`
	MajRate                float64    `json:"majRate"`
	MinRate                float64    `json:"minRate"`
	CrQuantity             int        `json:"crQuantity"`
	MajQuantity            int        `json:"majQuantity"`
	MinQuantity            int        `json:"minQuantity"`
	CheckResult            *string    `json:"checkResult"`
	ReciveDate             *time.Time `json:"reciveDate"`
	InspectDate            *time.Time `json:"inspectDate"`
	Inspector              *string    `json:"inspector"`
	Status                 string     `json:"status"`
	Remark                 *string    `json:"remark"`
	Attr1                  *string    `json:"attr1"`
	Attr2                  *string    `json:"attr2"`
	Attr3                  int        `json:"attr3"`
	Attr4                  int        `json:"attr4"`
	CreateBy               string     `json:"createBy"`
	UpdateBy               string     `json:"updateBy"`
	CreateById             int64      `json:"createById"`
	UpdateById             int64      `json:"updateById"`
	CreateTime             time.Time  `json:"createTime"`
	UpdateTime             time.Time  `json:"updateTime"`
}

// QcIqcListRes 来料检验单列表响应
type QcIqcListRes struct {
	Rows  []*QcIqcData `json:"rows"`
	Total int64        `json:"total"`
}

package qcOqcApi

import (
	"time"
)

// QcOqc 出货检验单表
type QcOqc struct {
	OqcId                  int64      `json:"oqcId" db:"oqc_id"`                                    // 出货检验单ID
	OqcCode                string     `json:"oqcCode" db:"oqc_code"`                                // 出货检验单编号
	OqcName                string     `json:"oqcName" db:"oqc_name"`                                // 出货检验单名称
	TemplateId             int64      `json:"templateId" db:"template_id"`                          // 检验模板ID
	SourceDocId            *int64     `json:"sourceDocId" db:"source_doc_id"`                       // 来源单据ID
	SourceDocType          *string    `json:"sourceDocType" db:"source_doc_type"`                   // 来源单据类型
	SourceDocCode          *string    `json:"sourceDocCode" db:"source_doc_code"`                   // 来源单据编号
	SourceLineId           *int64     `json:"sourceLineId" db:"source_line_id"`                     // 来源单据行ID
	ClientId               int64      `json:"clientId" db:"client_id"`                              // 客户ID
	ClientCode             string     `json:"clientCode" db:"client_code"`                          // 客户编码
	ClientName             string     `json:"clientName" db:"client_name"`                          // 客户名称
	BatchCode              *string    `json:"batchCode" db:"batch_code"`                            // 批次号
	ItemId                 int64      `json:"itemId" db:"item_id"`                                  // 产品物料ID
	ItemCode               *string    `json:"itemCode" db:"item_code"`                              // 产品物料编码
	ItemName               *string    `json:"itemName" db:"item_name"`                              // 产品物料名称
	Specification          *string    `json:"specification" db:"specification"`                     // 规格型号
	UnitOfMeasure          *string    `json:"unitOfMeasure" db:"unit_of_measure"`                   // 单位
	QuantityMinCheck       *float64   `json:"quantityMinCheck" db:"quantity_min_check"`             // 最低检测数
	QuantityMaxUnqualified *float64   `json:"quantityMaxUnqualified" db:"quantity_max_unqualified"` // 最大不合格数
	QuantityOut            float64    `json:"quantityOut" db:"quantity_out"`                        // 发货数量
	QuantityCheck          float64    `json:"quantityCheck" db:"quantity_check"`                    // 本次检测数量
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
	Remark                 string     `json:"remark" db:"remark"`                                   // 备注
	Attr1                  *string    `json:"attr1" db:"attr1"`                                     // 预留字段1
	Attr2                  *string    `json:"attr2" db:"attr2"`                                     // 预留字段2
	Attr3                  *int       `json:"attr3" db:"attr3"`                                     // 预留字段3
	Attr4                  *int       `json:"attr4" db:"attr4"`                                     // 预留字段4
	CreateBy               string     `json:"createBy" db:"create_by"`                              // 创建者
	CreateTime             *time.Time `json:"createTime" db:"create_time"`                          // 创建时间
	UpdateBy               string     `json:"updateBy" db:"update_by"`                              // 更新者
	UpdateTime             *time.Time `json:"updateTime" db:"update_time"`                          // 更新时间
}

// QcOqcQueryReq 查询请求
type QcOqcQueryReq struct {
	OqcCode                *string    `json:"oqcCode" form:"oqcCode"`
	OqcName                *string    `json:"oqcName" form:"oqcName"`
	TemplateId             *int64     `json:"templateId" form:"templateId"`
	SourceDocId            *int64     `json:"sourceDocId" form:"sourceDocId"`
	SourceDocType          *string    `json:"sourceDocType" form:"sourceDocType"`
	SourceDocCode          *string    `json:"sourceDocCode" form:"sourceDocCode"`
	SourceLineId           *int64     `json:"sourceLineId" form:"sourceLineId"`
	ClientId               *int64     `json:"clientId" form:"clientId"`
	ClientCode             *string    `json:"clientCode" form:"clientCode"`
	ClientName             *string    `json:"clientName" form:"clientName"`
	BatchCode              *string    `json:"batchCode" form:"batchCode"`
	ItemId                 *int64     `json:"itemId" form:"itemId"`
	ItemCode               *string    `json:"itemCode" form:"itemCode"`
	ItemName               *string    `json:"itemName" form:"itemName"`
	Specification          *string    `json:"specification" form:"specification"`
	UnitOfMeasure          *string    `json:"unitOfMeasure" form:"unitOfMeasure"`
	QuantityMinCheck       *float64   `json:"quantityMinCheck" form:"quantityMinCheck"`
	QuantityMaxUnqualified *float64   `json:"quantityMaxUnqualified" form:"quantityMaxUnqualified"`
	QuantityOut            *float64   `json:"quantityOut" form:"quantityOut"`
	QuantityCheck          *float64   `json:"quantityCheck" form:"quantityCheck"`
	QuantityUnqualified    *float64   `json:"quantityUnqualified" form:"quantityUnqualified"`
	QuantityQualified      *float64   `json:"quantityQualified" form:"quantityQualified"`
	CrRate                 *float64   `json:"crRate" form:"crRate"`
	MajRate                *float64   `json:"majRate" form:"majRate"`
	MinRate                *float64   `json:"minRate" form:"minRate"`
	CrQuantity             *float64   `json:"crQuantity" form:"crQuantity"`
	MajQuantity            *float64   `json:"majQuantity" form:"majQuantity"`
	MinQuantity            *float64   `json:"minQuantity" form:"minQuantity"`
	CheckResult            *string    `json:"checkResult" form:"checkResult"`
	OutDate                *time.Time `json:"outDate" form:"outDate"`
	InspectDate            *time.Time `json:"inspectDate" form:"inspectDate"`
	Inspector              *string    `json:"inspector" form:"inspector"`
	Status                 *string    `json:"status" form:"status"`
	CreateBy               *string    `json:"createBy" form:"createBy"`
	PageNum                int        `json:"pageNum" form:"pageNum"`
	PageSize               int        `json:"pageSize" form:"pageSize"`
}

// QcOqcCreateReq 创建请求
type QcOqcCreateReq struct {
	OqcCode                string     `json:"oqcCode" binding:"required"`       // 出货检验单编号
	OqcName                string     `json:"oqcName" binding:"required"`       // 出货检验单名称
	TemplateId             int64      `json:"templateId" binding:"required"`    // 检验模板ID
	SourceDocId            *int64     `json:"sourceDocId"`                      // 来源单据ID
	SourceDocType          *string    `json:"sourceDocType"`                    // 来源单据类型
	SourceDocCode          *string    `json:"sourceDocCode"`                    // 来源单据编号
	SourceLineId           *int64     `json:"sourceLineId"`                     // 来源单据行ID
	ClientId               int64      `json:"clientId" binding:"required"`      // 客户ID
	ClientCode             string     `json:"clientCode" binding:"required"`    // 客户编码
	ClientName             string     `json:"clientName" binding:"required"`    // 客户名称
	BatchCode              *string    `json:"batchCode"`                        // 批次号
	ItemId                 int64      `json:"itemId" binding:"required"`        // 产品物料ID
	ItemCode               *string    `json:"itemCode"`                         // 产品物料编码
	ItemName               *string    `json:"itemName"`                         // 产品物料名称
	Specification          *string    `json:"specification"`                    // 规格型号
	UnitOfMeasure          *string    `json:"unitOfMeasure"`                    // 单位
	QuantityMinCheck       *float64   `json:"quantityMinCheck"`                 // 最低检测数
	QuantityMaxUnqualified *float64   `json:"quantityMaxUnqualified"`           // 最大不合格数
	QuantityOut            float64    `json:"quantityOut" binding:"required"`   // 发货数量
	QuantityCheck          float64    `json:"quantityCheck" binding:"required"` // 本次检测数量
	QuantityUnqualified    *float64   `json:"quantityUnqualified"`              // 不合格数
	QuantityQualified      *float64   `json:"quantityQualified"`                // 合格数量
	CrRate                 *float64   `json:"crRate"`                           // 致命缺陷率
	MajRate                *float64   `json:"majRate"`                          // 严重缺陷率
	MinRate                *float64   `json:"minRate"`                          // 轻微缺陷率
	CrQuantity             *float64   `json:"crQuantity"`                       // 致命缺陷数量
	MajQuantity            *float64   `json:"majQuantity"`                      // 严重缺陷数量
	MinQuantity            *float64   `json:"minQuantity"`                      // 轻微缺陷数量
	CheckResult            *string    `json:"checkResult"`                      // 检测结果
	OutDate                *time.Time `json:"outDate"`                          // 出货日期
	InspectDate            *time.Time `json:"inspectDate"`                      // 检测日期
	Inspector              *string    `json:"inspector"`                        // 检测人员
	Status                 *string    `json:"status"`                           // 单据状态
	Remark                 string     `json:"remark"`                           // 备注
	Attr1                  *string    `json:"attr1"`                            // 预留字段1
	Attr2                  *string    `json:"attr2"`                            // 预留字段2
	Attr3                  *int       `json:"attr3"`                            // 预留字段3
	Attr4                  *int       `json:"attr4"`                            // 预留字段4
}

// QcOqcUpdateReq 更新请求
type QcOqcUpdateReq struct {
	OqcId                  int64      `json:"oqcId" binding:"required"`         // 出货检验单ID
	OqcCode                string     `json:"oqcCode" binding:"required"`       // 出货检验单编号
	OqcName                string     `json:"oqcName" binding:"required"`       // 出货检验单名称
	TemplateId             int64      `json:"templateId" binding:"required"`    // 检验模板ID
	SourceDocId            *int64     `json:"sourceDocId"`                      // 来源单据ID
	SourceDocType          *string    `json:"sourceDocType"`                    // 来源单据类型
	SourceDocCode          *string    `json:"sourceDocCode"`                    // 来源单据编号
	SourceLineId           *int64     `json:"sourceLineId"`                     // 来源单据行ID
	ClientId               int64      `json:"clientId" binding:"required"`      // 客户ID
	ClientCode             string     `json:"clientCode" binding:"required"`    // 客户编码
	ClientName             string     `json:"clientName" binding:"required"`    // 客户名称
	BatchCode              *string    `json:"batchCode"`                        // 批次号
	ItemId                 int64      `json:"itemId" binding:"required"`        // 产品物料ID
	ItemCode               *string    `json:"itemCode"`                         // 产品物料编码
	ItemName               *string    `json:"itemName"`                         // 产品物料名称
	Specification          *string    `json:"specification"`                    // 规格型号
	UnitOfMeasure          *string    `json:"unitOfMeasure"`                    // 单位
	QuantityMinCheck       *float64   `json:"quantityMinCheck"`                 // 最低检测数
	QuantityMaxUnqualified *float64   `json:"quantityMaxUnqualified"`           // 最大不合格数
	QuantityOut            float64    `json:"quantityOut" binding:"required"`   // 发货数量
	QuantityCheck          float64    `json:"quantityCheck" binding:"required"` // 本次检测数量
	QuantityUnqualified    *float64   `json:"quantityUnqualified"`              // 不合格数
	QuantityQualified      *float64   `json:"quantityQualified"`                // 合格数量
	CrRate                 *float64   `json:"crRate"`                           // 致命缺陷率
	MajRate                *float64   `json:"majRate"`                          // 严重缺陷率
	MinRate                *float64   `json:"minRate"`                          // 轻微缺陷率
	CrQuantity             *float64   `json:"crQuantity"`                       // 致命缺陷数量
	MajQuantity            *float64   `json:"majQuantity"`                      // 严重缺陷数量
	MinQuantity            *float64   `json:"minQuantity"`                      // 轻微缺陷数量
	CheckResult            *string    `json:"checkResult"`                      // 检测结果
	OutDate                *time.Time `json:"outDate"`                          // 出货日期
	InspectDate            *time.Time `json:"inspectDate"`                      // 检测日期
	Inspector              *string    `json:"inspector"`                        // 检测人员
	Status                 *string    `json:"status"`                           // 单据状态
	Remark                 string     `json:"remark"`                           // 备注
	Attr1                  *string    `json:"attr1"`                            // 预留字段1
	Attr2                  *string    `json:"attr2"`                            // 预留字段2
	Attr3                  *int       `json:"attr3"`                            // 预留字段3
	Attr4                  *int       `json:"attr4"`                            // 预留字段4
}

// QcOqcDeleteReq 删除请求
type QcOqcDeleteReq struct {
	OqcIds []int64 `json:"oqcIds" binding:"required"` // 出货检验单ID列表
}

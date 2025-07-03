package qcRqcApi

import (
	"time"
)

// QcRqc 退料检验单表
type QcRqc struct {
	RqcId               int64      `json:"rqcId" db:"rqc_id"`                             // 检验单ID
	RqcCode             string     `json:"rqcCode" db:"rqc_code"`                         // 检验单编号
	RqcName             *string    `json:"rqcName" db:"rqc_name"`                         // 检验单名称
	TemplateId          int64      `json:"templateId" db:"template_id"`                   // 检验模板ID
	SourceDocId         *int64     `json:"sourceDocId" db:"source_doc_id"`                // 来源单据ID
	SourceDocType       *string    `json:"sourceDocType" db:"source_doc_type"`            // 来源单据类型
	SourceDocCode       *string    `json:"sourceDocCode" db:"source_doc_code"`            // 来源单据编号
	SourceLineId        *int64     `json:"sourceLineId" db:"source_line_id"`              // 来源单据行ID
	ItemId              int64      `json:"itemId" db:"item_id"`                           // 产品物料ID
	ItemCode            *string    `json:"itemCode" db:"item_code"`                       // 产品物料编码
	ItemName            *string    `json:"itemName" db:"item_name"`                       // 产品物料名称
	Specification       *string    `json:"specification" db:"specification"`              // 规格型号
	UnitOfMeasure       *string    `json:"unitOfMeasure" db:"unit_of_measure"`            // 单位
	UnitName            *string    `json:"unitName" db:"unit_name"`                       // 单位名称
	BatchId             *int64     `json:"batchId" db:"batch_id"`                         // 批次ID
	BatchCode           *string    `json:"batchCode" db:"batch_code"`                     // 批次号
	QuantityCheck       float64    `json:"quantityCheck" db:"quantity_check"`             // 检测数量
	QuantityUnqualified float64    `json:"quantityUnqualified" db:"quantity_unqualified"` // 不合格数
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
	CreateBy            string     `json:"createBy" db:"create_by"`                       // 创建者
	CreateTime          time.Time  `json:"createTime" db:"create_time"`                   // 创建时间
	UpdateBy            string     `json:"updateBy" db:"update_by"`                       // 更新者
	UpdateTime          *time.Time `json:"updateTime" db:"update_time"`                   // 更新时间
}

// QcRqcQueryReq 查询请求
type QcRqcQueryReq struct {
	RqcCode             *string    `json:"rqcCode" form:"rqcCode"`
	RqcName             *string    `json:"rqcName" form:"rqcName"`
	TemplateId          *int64     `json:"templateId" form:"templateId"`
	SourceDocId         *int64     `json:"sourceDocId" form:"sourceDocId"`
	SourceDocType       *string    `json:"sourceDocType" form:"sourceDocType"`
	SourceDocCode       *string    `json:"sourceDocCode" form:"sourceDocCode"`
	SourceLineId        *int64     `json:"sourceLineId" form:"sourceLineId"`
	ItemId              *int64     `json:"itemId" form:"itemId"`
	ItemCode            *string    `json:"itemCode" form:"itemCode"`
	ItemName            *string    `json:"itemName" form:"itemName"`
	Specification       *string    `json:"specification" form:"specification"`
	UnitOfMeasure       *string    `json:"unitOfMeasure" form:"unitOfMeasure"`
	UnitName            *string    `json:"unitName" form:"unitName"`
	BatchId             *int64     `json:"batchId" form:"batchId"`
	BatchCode           *string    `json:"batchCode" form:"batchCode"`
	QuantityCheck       *float64   `json:"quantityCheck" form:"quantityCheck"`
	QuantityUnqualified *float64   `json:"quantityUnqualified" form:"quantityUnqualified"`
	QuantityQualified   *float64   `json:"quantityQualified" form:"quantityQualified"`
	CheckResult         *string    `json:"checkResult" form:"checkResult"`
	InspectDate         *time.Time `json:"inspectDate" form:"inspectDate"`
	UserId              *int64     `json:"userId" form:"userId"`
	UserName            *string    `json:"userName" form:"userName"`
	NickName            *string    `json:"nickName" form:"nickName"`
	Status              *string    `json:"status" form:"status"`
	CreateBy            *string    `json:"createBy" form:"createBy"`
	PageNum             int        `json:"pageNum" form:"pageNum"`
	PageSize            int        `json:"pageSize" form:"pageSize"`
}

// QcRqcCreateReq 创建请求
type QcRqcCreateReq struct {
	RqcName             *string    `json:"rqcName"`                          // 检验单名称
	TemplateId          int64      `json:"templateId" binding:"required"`    // 检验模板ID
	SourceDocId         *int64     `json:"sourceDocId"`                      // 来源单据ID
	SourceDocType       *string    `json:"sourceDocType"`                    // 来源单据类型
	SourceDocCode       *string    `json:"sourceDocCode"`                    // 来源单据编号
	SourceLineId        *int64     `json:"sourceLineId"`                     // 来源单据行ID
	ItemId              int64      `json:"itemId" binding:"required"`        // 产品物料ID
	ItemCode            *string    `json:"itemCode"`                         // 产品物料编码
	ItemName            *string    `json:"itemName"`                         // 产品物料名称
	Specification       *string    `json:"specification"`                    // 规格型号
	UnitOfMeasure       *string    `json:"unitOfMeasure"`                    // 单位
	UnitName            *string    `json:"unitName"`                         // 单位名称
	BatchId             *int64     `json:"batchId"`                          // 批次ID
	BatchCode           *string    `json:"batchCode"`                        // 批次号
	QuantityCheck       float64    `json:"quantityCheck" binding:"required"` // 检测数量
	QuantityUnqualified float64    `json:"quantityUnqualified"`              // 不合格数
	QuantityQualified   *float64   `json:"quantityQualified"`                // 合格品数量
	CheckResult         *string    `json:"checkResult"`                      // 检测结果
	InspectDate         *time.Time `json:"inspectDate"`                      // 检测日期
	UserId              *int64     `json:"userId"`                           // 检测人员ID
	UserName            *string    `json:"userName"`                         // 检测人员名称
	NickName            *string    `json:"nickName"`                         // 检测人员
	Status              *string    `json:"status"`                           // 单据状态
	Remark              *string    `json:"remark"`                           // 备注
	Attr1               *string    `json:"attr1"`                            // 预留字段1
	Attr2               *string    `json:"attr2"`                            // 预留字段2
	Attr3               *int       `json:"attr3"`                            // 预留字段3
	Attr4               *int       `json:"attr4"`                            // 预留字段4
}

// QcRqcUpdateReq 更新请求
type QcRqcUpdateReq struct {
	RqcId               int64      `json:"rqcId" binding:"required"` // 检验单ID
	RqcName             *string    `json:"rqcName"`                  // 检验单名称
	TemplateId          *int64     `json:"templateId"`               // 检验模板ID
	SourceDocId         *int64     `json:"sourceDocId"`              // 来源单据ID
	SourceDocType       *string    `json:"sourceDocType"`            // 来源单据类型
	SourceDocCode       *string    `json:"sourceDocCode"`            // 来源单据编号
	SourceLineId        *int64     `json:"sourceLineId"`             // 来源单据行ID
	ItemId              *int64     `json:"itemId"`                   // 产品物料ID
	ItemCode            *string    `json:"itemCode"`                 // 产品物料编码
	ItemName            *string    `json:"itemName"`                 // 产品物料名称
	Specification       *string    `json:"specification"`            // 规格型号
	UnitOfMeasure       *string    `json:"unitOfMeasure"`            // 单位
	UnitName            *string    `json:"unitName"`                 // 单位名称
	BatchId             *int64     `json:"batchId"`                  // 批次ID
	BatchCode           *string    `json:"batchCode"`                // 批次号
	QuantityCheck       *float64   `json:"quantityCheck"`            // 检测数量
	QuantityUnqualified *float64   `json:"quantityUnqualified"`      // 不合格数
	QuantityQualified   *float64   `json:"quantityQualified"`        // 合格品数量
	CheckResult         *string    `json:"checkResult"`              // 检测结果
	InspectDate         *time.Time `json:"inspectDate"`              // 检测日期
	UserId              *int64     `json:"userId"`                   // 检测人员ID
	UserName            *string    `json:"userName"`                 // 检测人员名称
	NickName            *string    `json:"nickName"`                 // 检测人员
	Status              *string    `json:"status"`                   // 单据状态
	Remark              *string    `json:"remark"`                   // 备注
	Attr1               *string    `json:"attr1"`                    // 预留字段1
	Attr2               *string    `json:"attr2"`                    // 预留字段2
	Attr3               *int       `json:"attr3"`                    // 预留字段3
	Attr4               *int       `json:"attr4"`                    // 预留字段4
}

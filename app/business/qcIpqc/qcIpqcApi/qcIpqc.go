package qcIpqcApi

import (
	"time"
)

// QcIpqc 过程检验单表 API 结构体
type QcIpqc struct {
	IpqcId                int64      `json:"ipqcId"`                // 主键ID
	IpqcCode              string     `json:"ipqcCode"`              // 检验单编号
	IpqcName              *string    `json:"ipqcName"`              // 检验单名称
	IpqcType              *string    `json:"ipqcType"`              // 检验类型
	TemplateId            *int64     `json:"templateId"`            // 检验模板ID
	SourceDocId           *int64     `json:"sourceDocId"`           // 来源单据ID
	SourceDocType         *string    `json:"sourceDocType"`         // 来源单据类型
	SourceDocCode         *string    `json:"sourceDocCode"`         // 来源单据编号
	SourceLineId          *int64     `json:"sourceLineId"`          // 来源单据行ID
	WorkorderId           *int64     `json:"workorderId"`           // 工单ID
	WorkorderCode         *string    `json:"workorderCode"`         // 工单编号
	WorkorderName         *string    `json:"workorderName"`         // 工单名称
	TaskId                *int64     `json:"taskId"`                // 任务ID
	TaskCode              *string    `json:"taskCode"`              // 任务编号
	TaskName              *string    `json:"taskName"`              // 任务名称
	WorkstationId         *int64     `json:"workstationId"`         // 工位ID
	WorkstationCode       *string    `json:"workstationCode"`       // 工位编号
	WorkstationName       *string    `json:"workstationName"`       // 工位名称
	ProcessId             *int64     `json:"processId"`             // 工序ID
	ProcessCode           *string    `json:"processCode"`           // 工序编号
	ProcessName           *string    `json:"processName"`           // 工序名称
	ItemId                *int64     `json:"itemId"`                // 物料ID
	ItemCode              *string    `json:"itemCode"`              // 物料编码
	ItemName              *string    `json:"itemName"`              // 物料名称
	Specification         *string    `json:"specification"`         // 规格型号
	UnitOfMeasure         *string    `json:"unitOfMeasure"`         // 单位
	UnitName              *string    `json:"unitName"`              // 单位名称
	QuantityCheck         *float64   `json:"quantityCheck"`         // 检测数量
	QuantityUnqualified   *float64   `json:"quantityUnqualified"`   // 不合格数
	QuantityQualified     *float64   `json:"quantityQualified"`     // 合格品数量
	QuantityLaborScrap    *float64   `json:"quantityLaborScrap"`    // 工废数量
	QuantityMaterialScrap *float64   `json:"quantityMaterialScrap"` // 料废数量
	QuantityOtherScrap    *float64   `json:"quantityOtherScrap"`    // 其他废品数量
	CrRate                *float64   `json:"crRate"`                // 严重不合格率
	MajRate               *float64   `json:"majRate"`               // 主要不合格率
	MinRate               *float64   `json:"minRate"`               // 次要不合格率
	CrQuantity            *float64   `json:"crQuantity"`            // 严重不合格数
	MajQuantity           *float64   `json:"majQuantity"`           // 主要不合格数
	MinQuantity           *float64   `json:"minQuantity"`           // 次要不合格数
	CheckResult           *string    `json:"checkResult"`           // 检测结果
	InspectDate           *time.Time `json:"inspectDate"`           // 检测日期
	Inspector             *string    `json:"inspector"`             // 检验员
	Status                *string    `json:"status"`                // 单据状态
	Remark                *string    `json:"remark"`                // 备注
	Attr1                 *string    `json:"attr1"`                 // 预留字段1
	Attr2                 *string    `json:"attr2"`                 // 预留字段2
	Attr3                 *int       `json:"attr3"`                 // 预留字段3
	Attr4                 *int       `json:"attr4"`                 // 预留字段4
	CreateBy              string     `json:"createBy"`              // 创建者
	CreateTime            time.Time  `json:"createTime"`            // 创建时间
	UpdateBy              string     `json:"updateBy"`              // 更新者
	UpdateTime            *time.Time `json:"updateTime"`            // 更新时间
}

type QcIpqcQueryReq struct {
	IpqcCode   *string `json:"ipqcCode" form:"ipqcCode"`
	IpqcName   *string `json:"ipqcName" form:"ipqcName"`
	IpqcType   *string `json:"ipqcType" form:"ipqcType"`
	TemplateId *int64  `json:"templateId" form:"templateId"`
	Status     *string `json:"status" form:"status"`
	PageNum    int     `json:"pageNum" form:"pageNum"`
	PageSize   int     `json:"pageSize" form:"pageSize"`
}

type QcIpqcCreateReq struct {
	IpqcCode              string     `json:"ipqcCode"`
	IpqcName              *string    `json:"ipqcName"`
	IpqcType              *string    `json:"ipqcType"`
	TemplateId            *int64     `json:"templateId"`
	SourceDocId           *int64     `json:"sourceDocId"`
	SourceDocType         *string    `json:"sourceDocType"`
	SourceDocCode         *string    `json:"sourceDocCode"`
	SourceLineId          *int64     `json:"sourceLineId"`
	WorkorderId           *int64     `json:"workorderId"`
	WorkorderCode         *string    `json:"workorderCode"`
	WorkorderName         *string    `json:"workorderName"`
	TaskId                *int64     `json:"taskId"`
	TaskCode              *string    `json:"taskCode"`
	TaskName              *string    `json:"taskName"`
	WorkstationId         *int64     `json:"workstationId"`
	WorkstationCode       *string    `json:"workstationCode"`
	WorkstationName       *string    `json:"workstationName"`
	ProcessId             *int64     `json:"processId"`
	ProcessCode           *string    `json:"processCode"`
	ProcessName           *string    `json:"processName"`
	ItemId                *int64     `json:"itemId"`
	ItemCode              *string    `json:"itemCode"`
	ItemName              *string    `json:"itemName"`
	Specification         *string    `json:"specification"`
	UnitOfMeasure         *string    `json:"unitOfMeasure"`
	UnitName              *string    `json:"unitName"`
	QuantityCheck         *float64   `json:"quantityCheck"`
	QuantityUnqualified   *float64   `json:"quantityUnqualified"`
	QuantityQualified     *float64   `json:"quantityQualified"`
	QuantityLaborScrap    *float64   `json:"quantityLaborScrap"`
	QuantityMaterialScrap *float64   `json:"quantityMaterialScrap"`
	QuantityOtherScrap    *float64   `json:"quantityOtherScrap"`
	CrRate                *float64   `json:"crRate"`
	MajRate               *float64   `json:"majRate"`
	MinRate               *float64   `json:"minRate"`
	CrQuantity            *float64   `json:"crQuantity"`
	MajQuantity           *float64   `json:"majQuantity"`
	MinQuantity           *float64   `json:"minQuantity"`
	CheckResult           *string    `json:"checkResult"`
	InspectDate           *time.Time `json:"inspectDate"`
	Inspector             *string    `json:"inspector"`
	Status                *string    `json:"status"`
	Remark                *string    `json:"remark"`
	Attr1                 *string    `json:"attr1"`
	Attr2                 *string    `json:"attr2"`
	Attr3                 *int       `json:"attr3"`
	Attr4                 *int       `json:"attr4"`
}

type QcIpqcUpdateReq struct {
	IpqcId     int64   `json:"ipqcId"`
	IpqcName   *string `json:"ipqcName"`
	IpqcType   *string `json:"ipqcType"`
	TemplateId *int64  `json:"templateId"`
	Status     *string `json:"status"`
	Remark     *string `json:"remark"`
	Attr1      *string `json:"attr1"`
	Attr2      *string `json:"attr2"`
	Attr3      *int    `json:"attr3"`
	Attr4      *int    `json:"attr4"`
}

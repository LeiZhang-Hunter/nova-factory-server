package qcIndexApi

import (
	"nova-factory-server/app/baize"
	"time"
)

// QcIndexListReq 检测项列表请求
type QcIndexListReq struct {
	IndexCode    string `form:"indexCode"`    // 检测项编码
	IndexName    string `form:"indexName"`    // 检测项名称
	IndexType    string `form:"indexType"`    // 检测项类型
	QcResultType string `form:"qcResultType"` // 质检值类型
	baize.BaseEntityDQL
}

// QcIndexCreateReq 检测项创建请求
type QcIndexCreateReq struct {
	IndexCode    string `json:"index_code" binding:"required"`     // 检测项编码
	IndexName    string `json:"index_name" binding:"required"`     // 检测项名称
	IndexType    string `json:"index_type" binding:"required"`     // 检测项类型
	QcTool       string `json:"qc_tool"`                           // 检测工具
	QcResultType string `json:"qc_result_type" binding:"required"` // 质检值类型
	QcResultSpc  string `json:"qc_result_spc"`                     // 值属性
	Remark       string `json:"remark"`                            // 备注
	Attr1        string `json:"attr1"`                             // 预留字段1
	Attr2        string `json:"attr2"`                             // 预留字段2
	Attr3        int    `json:"attr3"`                             // 预留字段3
	Attr4        int    `json:"attr4"`                             // 预留字段4
}

// QcIndexUpdateReq 检测项修改请求
type QcIndexUpdateReq struct {
	IndexId      int64  `json:"index_id" binding:"required"`       // 检测项ID
	IndexCode    string `json:"index_code" binding:"required"`     // 检测项编码
	IndexName    string `json:"index_name" binding:"required"`     // 检测项名称
	IndexType    string `json:"index_type" binding:"required"`     // 检测项类型
	QcTool       string `json:"qc_tool"`                           // 检测工具
	QcResultType string `json:"qc_result_type" binding:"required"` // 质检值类型
	QcResultSpc  string `json:"qc_result_spc"`                     // 值属性
	Remark       string `json:"remark"`                            // 备注
	Attr1        string `json:"attr1"`                             // 预留字段1
	Attr2        string `json:"attr2"`                             // 预留字段2
	Attr3        int    `json:"attr3"`                             // 预留字段3
	Attr4        int    `json:"attr4"`                             // 预留字段4
}

// QcIndexData 检测项数据
type QcIndexData struct {
	IndexId      int64     `json:"index_id"`       // 检测项ID
	IndexCode    string    `json:"index_code"`     // 检测项编码
	IndexName    string    `json:"index_name"`     // 检测项名称
	IndexType    string    `json:"index_type"`     // 检测项类型
	QcTool       string    `json:"qc_tool"`        // 检测工具
	QcResultType string    `json:"qc_result_type"` // 质检值类型
	QcResultSpc  string    `json:"qc_result_spc"`  // 值属性
	Remark       string    `json:"remark"`         // 备注
	Attr1        string    `json:"attr1"`          // 预留字段1
	Attr2        string    `json:"attr2"`          // 预留字段2
	Attr3        int       `json:"attr3"`          // 预留字段3
	Attr4        int       `json:"attr4"`          // 预留字段4
	CreateBy     string    `json:"create_by"`      // 创建者
	CreateTime   time.Time `json:"create_time"`    // 创建时间
	UpdateBy     string    `json:"update_by"`      // 更新者
	UpdateTime   time.Time `json:"update_time"`    // 更新时间
}

// QcIndexListRes 检测项列表响应
type QcIndexListRes struct {
	Rows  []*QcIndexData `json:"rows"`
	Total int64          `json:"total"`
}

// QcIndexCreateRes 检测项创建响应
type QcIndexCreateRes struct {
	IndexId      int64     `json:"index_id"`       // 检测项ID
	IndexCode    string    `json:"index_code"`     // 检测项编码
	IndexName    string    `json:"index_name"`     // 检测项名称
	IndexType    string    `json:"index_type"`     // 检测项类型
	QcTool       string    `json:"qc_tool"`        // 检测工具
	QcResultType string    `json:"qc_result_type"` // 质检值类型
	QcResultSpc  string    `json:"qc_result_spc"`  // 值属性
	Remark       string    `json:"remark"`         // 备注
	Attr1        string    `json:"attr1"`          // 预留字段1
	Attr2        string    `json:"attr2"`          // 预留字段2
	Attr3        int       `json:"attr3"`          // 预留字段3
	Attr4        int       `json:"attr4"`          // 预留字段4
	CreateBy     string    `json:"create_by"`      // 创建者
	CreateTime   time.Time `json:"create_time"`    // 创建时间
	UpdateBy     string    `json:"update_by"`      // 更新者
	UpdateTime   time.Time `json:"update_time"`    // 更新时间
}

// QcIndexUpdateRes 检测项修改响应
type QcIndexUpdateRes struct {
	IndexId      int64     `json:"index_id"`       // 检测项ID
	IndexCode    string    `json:"index_code"`     // 检测项编码
	IndexName    string    `json:"index_name"`     // 检测项名称
	IndexType    string    `json:"index_type"`     // 检测项类型
	QcTool       string    `json:"qc_tool"`        // 检测工具
	QcResultType string    `json:"qc_result_type"` // 质检值类型
	QcResultSpc  string    `json:"qc_result_spc"`  // 值属性
	Remark       string    `json:"remark"`         // 备注
	Attr1        string    `json:"attr1"`          // 预留字段1
	Attr2        string    `json:"attr2"`          // 预留字段2
	Attr3        int       `json:"attr3"`          // 预留字段3
	Attr4        int       `json:"attr4"`          // 预留字段4
	CreateBy     string    `json:"create_by"`      // 创建者
	CreateTime   time.Time `json:"create_time"`    // 创建时间
	UpdateBy     string    `json:"update_by"`      // 更新者
	UpdateTime   time.Time `json:"update_time"`    // 更新时间
}

package defectApi

import (
	"nova-factory-server/app/baize"
)

// DefectListReq 缺陷列表请求
type DefectListReq struct {
	DefectCode  string `form:"defectCode"`
	DefectName  string `form:"defectName"`
	IndexType   string `form:"indexType"`
	DefectLevel string `form:"defectLevel"`
	baize.BaseEntityDQL
}

// DefectCreateReq 缺陷创建请求
type DefectCreateReq struct {
	DefectCode  string `json:"defect_code" binding:"required"`  // 缺陷编码
	DefectName  string `json:"defect_name" binding:"required"`  // 缺陷名称
	IndexType   string `json:"index_type" binding:"required"`   // 指标类型
	DefectLevel string `json:"defect_level" binding:"required"` // 缺陷等级
	Remark      string `json:"remark"`                          // 备注
	Attr1       string `json:"attr_1"`                          // 属性1
	Attr2       string `json:"attr_2"`                          // 属性2
	Attr3       string `json:"attr_3"`                          // 属性3
	Attr4       string `json:"attr_4"`                          // 属性4
}

// DefectUpdateReq 缺陷修改请求
type DefectUpdateReq struct {
	DefectId    int64  `json:"defect_id" binding:"required"`    // 缺陷ID
	DefectCode  string `json:"defect_code" binding:"required"`  // 缺陷编码
	DefectName  string `json:"defect_name" binding:"required"`  // 缺陷名称
	IndexType   string `json:"index_type" binding:"required"`   // 指标类型
	DefectLevel string `json:"defect_level" binding:"required"` // 缺陷等级
	State       string `json:"state"`                           // 状态
	Remark      string `json:"remark"`                          // 备注
	Attr1       string `json:"attr_1"`                          // 属性1
	Attr2       string `json:"attr_2"`                          // 属性2
	Attr3       string `json:"attr_3"`                          // 属性3
	Attr4       string `json:"attr_4"`                          // 属性4
}

// DefectCreateRes 缺陷创建响应
type DefectCreateRes struct {
}

// DefectUpdateRes 缺陷修改响应
type DefectUpdateRes struct {
}

type DefectData struct {
	Id          int64  `json:"id"`
	DefectCode  string `json:"defect_code"`
	DefectName  string `json:"defect_name"`
	IndexType   string `json:"index_type"`
	DefectLevel string `json:"defect_level"`
	DeptId      int64  `json:"dept_id"`
	State       string `json:"state"`
	CreateBy    string `json:"create_by"`
	UpdateBy    string `json:"update_by"`
	Remark      string `json:"remark"`
	Attr1       string `json:"attr_1"`
	Attr2       string `json:"attr_2"`
	Attr3       string `json:"attr_3"`
	Attr4       string `json:"attr_4"`
}
type DefectListRes struct {
	Rows  []*DefectData `json:"rows"`
	Total int64         `json:"total"`
}

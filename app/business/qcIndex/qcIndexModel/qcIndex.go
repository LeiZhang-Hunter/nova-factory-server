package qcIndexModel

import "gorm.io/gorm"

// QcIndex 检测项表
type QcIndex struct {
	gorm.Model
	IndexCode    string `json:"indexCode" gorm:"column:index_code;size:64;not null;comment:检测项编码"`
	IndexName    string `json:"indexName" gorm:"column:index_name;size:255;not null;comment:检测项名称"`
	IndexType    string `json:"indexType" gorm:"column:index_type;size:64;not null;comment:检测项类型"`
	QcTool       string `json:"qcTool" gorm:"column:qc_tool;size:255;comment:检测工具"`
	QcResultType string `json:"qcResultType" gorm:"column:qc_result_type;size:64;not null;comment:质检值类型"`
	QcResultSpc  string `json:"qcResultSpc" gorm:"column:qc_result_spc;size:255;comment:值属性"`
	Remark       string `json:"remark" gorm:"column:remark;size:500;comment:备注"`
	Attr1        string `json:"attr1" gorm:"column:attr1;size:64;comment:预留字段1"`
	Attr2        string `json:"attr2" gorm:"column:attr2;size:255;comment:预留字段2"`
	Attr3        int    `json:"attr3" gorm:"column:attr3;default:0;comment:预留字段3"`
	Attr4        int    `json:"attr4" gorm:"column:attr4;default:0;comment:预留字段4"`
	CreateBy     string `json:"createBy" gorm:"column:create_by;size:64;comment:创建者"`
	UpdateBy     string `json:"updateBy" gorm:"column:update_by;size:64;comment:更新者"`
	CreateById   int64  `json:"createById" gorm:"column:create_by_id;comment:创建者ID"`
	UpdateById   int64  `json:"updateById" gorm:"column:update_by_id;comment:更新者ID"`
}

func (QcIndex) TableName() string {
	return "qc_index"
}

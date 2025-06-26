package defectModel

import "gorm.io/gorm"

// Defect 缺陷信息表
type Defect struct {
	gorm.Model
	DefectCode  string `json:"defectCode" gorm:"column:defect_code;size:64;not null;comment:缺陷编码"`
	DefectName  string `json:"defectName" gorm:"column:defect_name;size:128;not null;comment:缺陷名称"`
	IndexType   string `json:"indexType" gorm:"column:index_type;size:32;comment:指标类型"`
	DefectLevel string `json:"defectLevel" gorm:"column:defect_level;size:32;comment:缺陷等级"`
	DeptId      int64  `json:"deptId" gorm:"column:dept_id;comment:部门ID"`
	State       string `json:"state" gorm:"column:state;size:1;default:0;comment:状态（0正常 1停用）"`
	CreateById  int64  `json:"createById" gorm:"column:create_by_id;comment:创建者ID"`
	CreateBy    string `json:"createBy" gorm:"column:create_by;size:64;comment:创建者"`
	UpdateById  int64  `json:"updateById" gorm:"column:update_by_id;comment:更新者ID"`
	UpdateBy    string `json:"updateBy" gorm:"column:update_by;size:64;comment:更新者"`
	Remark      string `json:"remark" gorm:"column:remark;size:500;comment:备注"`
	Attr1       string `json:"attr1" gorm:"column:attr_1;size:64;comment:预留字段1"`
	Attr2       string `json:"attr2" gorm:"column:attr_2;size:64;comment:预留字段2"`
	Attr3       string `json:"attr3" gorm:"column:attr_3;size:64;comment:预留字段3"`
	Attr4       string `json:"attr4" gorm:"column:attr_4;size:64;comment:预留字段4"`
}

func (Defect) TableName() string {
	return "sys_defect"
}

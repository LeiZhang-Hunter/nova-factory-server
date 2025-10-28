package deviceModels

import (
	"nova-factory-server/app/baize"
)

// SysDeviceSubject 设备点检保养项目表
type SysDeviceSubject struct {
	ID              int64  `gorm:"column:id;primaryKey;autoIncrement:true;comment:项目ID" json:"id"`      // 项目ID
	SubjectCode     string `gorm:"column:subject_code;not null;comment:项目编码" json:"subject_code"`       // 项目编码
	SubjectName     string `gorm:"column:subject_name;comment:项目名称" json:"subject_name"`                // 项目名称
	SubjectType     string `gorm:"column:subject_type;comment:项目类型" json:"subject_type"`                // 项目类型
	SubjectContent  string `gorm:"column:subject_content;not null;comment:项目内容" json:"subject_content"` // 项目内容
	SubjectStandard string `gorm:"column:subject_standard;comment:标准" json:"subject_standard"`          // 标准
	Remark          string `gorm:"column:remark;comment:备注" json:"remark"`                              // 备注
	Attr1           string `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                             // 预留字段1
	Attr2           string `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                             // 预留字段2
	Attr3           int32  `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                             // 预留字段3
	Attr4           int32  `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                             // 预留字段4
	Status          bool   `gorm:"column:status;comment:操作状态（0禁用 1启用）" json:"status"`                   // 操作状态（0禁用 1启用）
	DeptID          int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                          // 部门ID
	baize.BaseEntity
	State bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

func ToSysDeviceSubject(vo *SysDeviceSubjectVO) *SysDeviceSubject {
	return &SysDeviceSubject{
		ID:              vo.ID,
		SubjectCode:     vo.SubjectCode,
		SubjectName:     vo.SubjectName,
		SubjectType:     vo.SubjectType,
		SubjectContent:  vo.SubjectContent,
		SubjectStandard: vo.SubjectStandard,
		Remark:          vo.Remark,
		Attr1:           vo.Attr1,
		Attr2:           vo.Attr2,
		Attr3:           vo.Attr3,
		Attr4:           vo.Attr4,
		Status:          vo.Status,
	}
}

// SysDeviceSubjectVO 设备点检保养项目表
type SysDeviceSubjectVO struct {
	ID              int64  `gorm:"column:id;primaryKey;autoIncrement:true;comment:项目ID" json:"id"`      // 项目ID
	SubjectCode     string `gorm:"column:subject_code;not null;comment:项目编码" json:"subject_code"`       // 项目编码
	SubjectName     string `gorm:"column:subject_name;comment:项目名称" json:"subject_name"`                // 项目名称
	SubjectType     string `gorm:"column:subject_type;comment:项目类型" json:"subject_type"`                // 项目类型
	SubjectContent  string `gorm:"column:subject_content;not null;comment:项目内容" json:"subject_content"` // 项目内容
	SubjectStandard string `gorm:"column:subject_standard;comment:标准" json:"subject_standard"`          // 标准
	Remark          string `gorm:"column:remark;comment:备注" json:"remark"`                              // 备注
	Attr1           string `gorm:"column:attr1;comment:预留字段1" json:"attr1"`                             // 预留字段1
	Attr2           string `gorm:"column:attr2;comment:预留字段2" json:"attr2"`                             // 预留字段2
	Attr3           int32  `gorm:"column:attr3;comment:预留字段3" json:"attr3"`                             // 预留字段3
	Attr4           int32  `gorm:"column:attr4;comment:预留字段4" json:"attr4"`                             // 预留字段4
	Status          bool   `gorm:"column:status;comment:操作状态（0禁用 1启用）" json:"status"`                   // 操作状态（0禁用 1启用）
}

type SysDeviceSubjectReq struct {
	Name string `form:"name"`
	baize.BaseEntityDQL
}

package gatewaymodels

import "nova-factory-server/app/baize"

type InstalledSkill struct {
	ID          int64  `json:"id,string" gorm:"column:id"`
	Name        string `json:"name" gorm:"column:name"`
	Slug        string `json:"slug" gorm:"column:slug"`
	Version     string `json:"version" gorm:"column:version"`
	Source      string `json:"source" gorm:"column:source"`
	Description string `json:"description" gorm:"column:description"`
	Enabled     *bool  `json:"enabled" gorm:"column:enabled"`
	DeptID      int64  `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

type InstalledSkillUpsert struct {
	ID          int64  `json:"id,string"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Version     string `json:"version"`
	Source      string `json:"source"`
	Description string `json:"description"`
	Enabled     *bool  `json:"enabled"`
}

type InstalledSkillQuery struct {
	Name    string `form:"name"`
	Slug    string `form:"slug"`
	Enabled *bool  `form:"enabled"`
	Page    int64  `form:"page"`
	Size    int64  `form:"size"`
}

type InstalledSkillListData struct {
	Rows  []*InstalledSkill `json:"rows"`
	Total int64             `json:"total"`
}

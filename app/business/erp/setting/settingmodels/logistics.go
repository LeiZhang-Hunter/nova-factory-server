package settingmodels

import "nova-factory-server/app/baize"

// LogisticsCompany ERP物流公司
type LogisticsCompany struct {
	ID           uint64 `json:"id,string" db:"id"`
	Code         string `json:"code" db:"code"`
	Name         string `json:"name" db:"name"`
	ShortName    string `json:"shortName" db:"short_name"`
	ContactName  string `json:"contactName" db:"contact_name"`
	ContactPhone string `json:"contactPhone" db:"contact_phone"`
	Address      string `json:"address" db:"address"`
	Website      string `json:"website" db:"website"`
	Remark       string `json:"remark" db:"remark"`
	Sort         int32  `json:"sort" db:"sort"`
	Status       *bool  `json:"status" db:"status"`
	DeptID       int64  `json:"deptId" db:"dept_id"`
	baize.BaseEntity
	State int32 `json:"state" db:"state"`
}

// LogisticsCompanyUpsert ERP物流公司新增修改参数
type LogisticsCompanyUpsert struct {
	ID           uint64 `json:"id,string"`
	Code         string `json:"code" binding:"required"`
	Name         string `json:"name" binding:"required"`
	ShortName    string `json:"shortName"`
	ContactName  string `json:"contactName"`
	ContactPhone string `json:"contactPhone"`
	Address      string `json:"address"`
	Website      string `json:"website"`
	Remark       string `json:"remark"`
	Sort         int32  `json:"sort"`
	Status       *bool  `json:"status"`
}

// LogisticsCompanyQuery ERP物流公司查询参数
type LogisticsCompanyQuery struct {
	Code   string `form:"code"`
	Name   string `form:"name"`
	Status *bool  `form:"status"`
	Page   int64  `form:"page"`
	Size   int64  `form:"size"`
}

// LogisticsCompanyListData ERP物流公司分页数据
type LogisticsCompanyListData struct {
	Rows  []*LogisticsCompany `json:"rows"`
	Total int64               `json:"total"`
}

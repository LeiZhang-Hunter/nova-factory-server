package models

import "nova-factory-server/app/baize"

// LogisticsCompany ERP物流公司
type LogisticsCompany struct {
	ID      int64  `json:"id,string" gorm:"id"`
	Code    string `json:"code" gorm:"code"`
	Name    string `json:"name" gorm:"name"`
	Company string `json:"company" gorm:"company"`

	ProvinceName string `json:"provinceName" gorm:"province_name"`
	ProvinceCode string `json:"provinceCode" gorm:"province_code"`

	CityName string `json:"cityName" gorm:"city_name"`
	CityCode string `json:"cityCode" gorm:"city_code"`

	DistrictName string `json:"districtName" gorm:"district_name"`
	DistrictCode string `json:"districtCode" gorm:"district_code"`

	StreetName string `json:"streetName" gorm:"street_name"`
	StreetCode string `json:"streetCode" gorm:"street_code"`

	ShortName    string `json:"shortName" gorm:"short_name"`
	ContactName  string `json:"contactName" gorm:"contact_name"`
	ContactPhone string `json:"contactPhone" gorm:"contact_phone"`
	Address      string `json:"address" gorm:"address"`
	Remark       string `json:"remark" gorm:"remark"`
	Sort         int32  `json:"sort" gorm:"sort"`
	Status       *bool  `json:"status" gorm:"status"`
	DeptID       int64  `json:"deptId" gorm:"dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"state"`
}

// LogisticsCompanyUpsert ERP物流公司新增修改参数
type LogisticsCompanyUpsert struct {
	ID      int64  `json:"id,string"`
	Code    string `json:"code" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Company string `json:"company"`

	ProvinceName string `json:"provinceName"`
	ProvinceCode string `json:"provinceCode"`

	CityName string `json:"cityName"`
	CityCode string `json:"cityCode"`

	DistrictName string `json:"districtName"`
	DistrictCode string `json:"districtCode"`

	StreetName string `json:"streetName" `
	StreetCode string `json:"streetCode" `

	ShortName    string `json:"shortName"`
	ContactName  string `json:"contactName"`
	ContactPhone string `json:"contactPhone"`
	Address      string `json:"address"`
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

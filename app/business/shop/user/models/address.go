package models

import "nova-factory-server/app/baize"

// Address 商城用户地址
type Address struct {
	ID             int64  `json:"id,string" db:"id"`                         // 主键ID
	UserID         string `json:"userId" db:"user_id"`                       // 用户业务ID
	ReceiverName   string `json:"receiverName" db:"receiver_name"`           // 收货人姓名
	ReceiverMobile string `json:"receiverMobile" db:"receiver_mobile"`       // 收货人手机号
	ProvinceCode   string `json:"provinceCode" db:"province_code"`           // 省编码
	ProvinceName   string `json:"provinceName" db:"province_name"`           // 省名称
	CityCode       string `json:"cityCode" db:"city_code"`                   // 市编码
	CityName       string `json:"cityName" db:"city_name"`                   // 市名称
	DistrictCode   string `json:"districtCode" db:"district_code"`           // 区编码
	DistrictName   string `json:"districtName" db:"district_name"`           // 区名称
	StreetCode     string `json:"streetCode" db:"street_code"`               // 街道编码
	StreetName     string `json:"streetName" db:"street_name"`               // 街道名称
	DetailAddress  string `json:"detailAddress" db:"detail_address"`         // 详细地址
	PostalCode     string `json:"postalCode" db:"postal_code"`               // 邮政编码
	AddressLabel   string `json:"addressLabel" db:"address_label"`           // 地址标签
	IsDefault      *bool  `json:"isDefault" db:"is_default"`                 // 是否默认地址
	Status         int32  `json:"status" db:"status"`                        // 状态
	DeptID         int64  `json:"deptId" gorm:"column:dept_id" db:"dept_id"` // 部门ID
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state" db:"state"` // 操作状态
}

// AddressSetReq 地址新增修改参数
type AddressSetReq struct {
	ID             int64  `json:"id,string"`                         // 主键ID
	UserID         string `json:"userId"`                            // 用户业务ID
	ReceiverName   string `json:"receiverName"`                      // 收货人姓名
	ReceiverMobile string `json:"receiverMobile" binding:"required"` // 收货人手机号
	ProvinceCode   string `json:"provinceCode"`                      // 省编码
	ProvinceName   string `json:"provinceName"`                      // 省名称
	CityCode       string `json:"cityCode"`                          // 市编码
	CityName       string `json:"cityName"`                          // 市名称
	DistrictCode   string `json:"districtCode"`                      // 区编码
	DistrictName   string `json:"districtName"`                      // 区名称
	StreetCode     string `json:"streetCode"`                        // 街道编码
	StreetName     string `json:"streetName"`                        // 街道名称
	DetailAddress  string `json:"detailAddress" binding:"required"`  // 详细地址
	PostalCode     string `json:"postalCode"`                        // 邮政编码
	AddressLabel   string `json:"addressLabel"`                      // 地址标签
	IsDefault      *bool  `json:"isDefault"`                         // 是否默认地址
	Status         *int32 `json:"status"`                            // 状态
}

// AddressQuery 地址查询参数
type AddressQuery struct {
	UserID         string `form:"userId"`         // 用户业务ID
	ReceiverName   string `form:"receiverName"`   // 收货人姓名
	ReceiverMobile string `form:"receiverMobile"` // 收货人手机号
	Status         *int32 `form:"status"`         // 状态
	Page           int64  `form:"page"`           // 页码
	Size           int64  `form:"size"`           // 每页数量
}

// AddressListData 地址列表结果
type AddressListData struct {
	Rows  []*Address `json:"rows"`  // 数据列表
	Total int64      `json:"total"` // 总数
}

// AddressRegionQuery 行政区查询参数
type AddressRegionQuery struct {
	ParentCode string `form:"parentCode"` // 父级行政区编码，为空时查询省级
	Type       string `form:"type"`       // 行政区层级：province/city/district/street
}

// AddressRegionItem 行政区节点
type AddressRegionItem struct {
	Code       string `json:"code"`       // 行政区编码
	Name       string `json:"name"`       // 行政区名称
	Level      string `json:"level"`      // 层级：province/city/district
	ParentCode string `json:"parentCode"` // 父级行政区编码
}

type UserAddressInfoQuery struct {
	Username string `form:"username"`
}

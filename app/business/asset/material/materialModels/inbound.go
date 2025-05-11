package materialModels

import "nova-factory-server/app/baize"

type InboundListReq struct {
	Keyword string `json:"keyword,string" db:"keyword" form:"keyword"`
	baize.BaseEntityDQL
}

// InboundInfo 入库
type InboundInfo struct {
	MaterialId uint64 `json:"materialId,string" db:"material_id"`
	Number     uint64 `json:"number" db:"number"`
}

type InboundVO struct {
	Name           string  `json:"name" db:"name" binding:"required" gorm:"-"`
	Code           string  `json:"code" db:"code" binding:"required" gorm:"-"`
	Model          string  `json:"model" db:"model" binding:"required" gorm:"-"`
	Unit           string  `json:"unit" db:"device_building_id" binding:"required" gorm:"-"`
	Factory        string  `json:"factory" db:"device_building_id" binding:"required" gorm:"-"`
	Address        string  `json:"address" db:"device_building_id" binding:"required" gorm:"-"`
	Price          float64 `json:"price" db:"device_building_id" binding:"required" gorm:"-"`
	Total          uint64  `json:"total" db:"device_building_id" gorm:"-"`
	Outbound       uint64  `json:"outbound" db:"device_building_id" gorm:"-"`
	CurrentTotal   uint64  `json:"currentTotal" db:"current_total" gorm:"-"`
	InboundId      uint64  `json:"inboundId,string" db:"inbound_id"`
	MaterialId     uint64  `json:"materialId,string" db:"material_id"`
	DeptId         int64   `json:"-" db:"dept_id"`
	Number         uint64  `json:"number" db:"number"`
	CreateUserName string  `json:"createUserName" gorm:"-"`
	UpdateUserName string  `json:"updateUserName" gorm:"-"`
	baize.BaseEntity
}

type InboundValue struct {
	InboundId      uint64 `json:"inboundId,string" db:"inbound_id"`
	MaterialId     uint64 `json:"materialId,string" db:"material_id"`
	DeptId         int64  `json:"-" db:"dept_id"`
	Number         uint64 `json:"number" db:"number"`
	CreateUserName string `json:"createUserName"`
	UpdateUserName string `json:"updateUserName"`
	baize.BaseEntity
}

type InboundListData struct {
	Rows  []*InboundData `json:"rows"`
	Total int64          `json:"total"`
}

type InboundListValue struct {
	Rows  []*InboundValue `json:"rows"`
	Total int64           `json:"total"`
}

type InboundData struct {
	Name           string  `json:"name" db:"name" binding:"required"`
	Code           string  `json:"code" db:"code" binding:"required"`
	Model          string  `json:"model" db:"model" binding:"required"`
	Unit           string  `json:"unit" db:"device_building_id" binding:"required"`
	Factory        string  `json:"factory" db:"device_building_id" binding:"required"`
	Address        string  `json:"address" db:"device_building_id" binding:"required"`
	Price          float64 `json:"price" db:"device_building_id" binding:"required"`
	Total          uint64  `json:"total" db:"device_building_id"`
	Outbound       uint64  `json:"outbound" db:"device_building_id"`
	CurrentTotal   uint64  `json:"currentTotal" db:"current_total"`
	InboundId      uint64  `json:"inboundId,string" db:"inbound_id"`
	MaterialId     uint64  `json:"materialId,string" db:"material_id"`
	DeptId         int64   `json:"-" db:"dept_id"`
	Number         uint64  `json:"number" db:"number"`
	CreateUserName string  `json:"createUserName"`
	UpdateUserName string  `json:"updateUserName"`
	baize.BaseEntity
}

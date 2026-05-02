package models

import "nova-factory-server/app/baize"

type Combination struct {
	ID            int64   `json:"id,string" db:"id"`
	ProductID     string  `json:"productId" db:"product_id"`
	MerID         int64   `json:"merId,string" db:"mer_id"`
	Image         string  `json:"image" db:"image"`
	Images        string  `json:"images" db:"images"`
	Title         string  `json:"title" db:"title"`
	Attr          string  `json:"attr" db:"attr"`
	People        int32   `json:"people" db:"people"`
	Info          string  `json:"info" db:"info"`
	Price         float64 `json:"price" db:"price"`
	Sort          int32   `json:"sort" db:"sort"`
	Sales         int64   `json:"sales" db:"sales"`
	Stock         int64   `json:"stock" db:"stock"`
	IsHost        int32   `json:"isHost" db:"is_host"`
	IsShow        int32   `json:"isShow" db:"is_show"`
	IsPostage     int32   `json:"isPostage" db:"is_postage"`
	Postage       float64 `json:"postage" db:"postage"`
	StartTime     int64   `json:"startTime" db:"start_time"`
	StopTime      int64   `json:"stopTime" db:"stop_time"`
	EffectiveTime int64   `json:"effectiveTime" db:"effective_time"`
	Browse        int64   `json:"browse" db:"browse"`
	UnitName      string  `json:"unitName" db:"unit_name"`
	Weight        float64 `json:"weight" db:"weight"`
	Volume        float64 `json:"volume" db:"volume"`
	Num           int64   `json:"num" db:"num"`
	OnceNum       int64   `json:"onceNum" db:"once_num"`
	Quota         int64   `json:"quota" db:"quota"`
	QuotaShow     int64   `json:"quotaShow" db:"quota_show"`
	Virtual       int64   `json:"virtual" db:"virtual"`
	HomeModuleIDs string  `json:"homeModuleIds" db:"home_module_ids"`
	DeptID        int64   `json:"deptId" db:"dept_id"`
	baize.BaseEntity
	State int32 `json:"state" db:"state"`
}

type CombinationSet struct {
	ID            int64    `json:"id,string"`
	ProductID     string   `json:"productId" binding:"required"`
	MerID         int64    `json:"merId,string"`
	Image         string   `json:"image" binding:"required"`
	Images        string   `json:"images"`
	Title         string   `json:"title" binding:"required"`
	Attr          string   `json:"attr"`
	People        int32    `json:"people" binding:"required"`
	Info          string   `json:"info"`
	Price         float64  `json:"price"`
	Sort          int32    `json:"sort"`
	Sales         int64    `json:"sales"`
	Stock         int64    `json:"stock"`
	IsHost        int32    `json:"isHost"`
	IsShow        int32    `json:"isShow"`
	IsPostage     int32    `json:"isPostage"`
	Postage       float64  `json:"postage"`
	StartTime     int64    `json:"startTime"`
	StopTime      int64    `json:"stopTime"`
	EffectiveTime int64    `json:"effectiveTime"`
	Browse        int64    `json:"browse"`
	UnitName      string   `json:"unitName"`
	Weight        float64  `json:"weight"`
	Volume        float64  `json:"volume"`
	Num           int64    `json:"num"`
	OnceNum       int64    `json:"onceNum"`
	Quota         int64    `json:"quota"`
	QuotaShow     int64    `json:"quotaShow"`
	Virtual       int64    `json:"virtual"`
	HomeModuleIDs []string `json:"homeModuleIds"`
}

type CombinationQuery struct {
	Title     string `form:"title"`
	ProductID int64  `form:"productId"`
	IsShow    *int32 `form:"isShow"`
	IsHost    *int32 `form:"isHost"`
	Page      int64  `form:"page"`
	Size      int64  `form:"size"`
}

type CombinationListData struct {
	Rows  []*Combination `json:"rows"`
	Total int64          `json:"total"`
}

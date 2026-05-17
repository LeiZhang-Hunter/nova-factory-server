package models

import "nova-factory-server/app/baize"

type CombinationMainInfo struct {
	ID            int64   `json:"id,string" gorm:"id"`
	ProductID     string  `json:"productId" gorm:"product_id"`
	MerID         int64   `json:"merId,string" gorm:"mer_id"`
	Attr          string  `json:"attr" gorm:"attr"`
	Price         float64 `json:"price" gorm:"price"`
	Sort          int32   `json:"sort" gorm:"sort"`
	Sales         int64   `json:"sales" gorm:"sales"`
	Stock         int64   `json:"stock" gorm:"stock"`
	IsHost        int32   `json:"isHost" gorm:"is_host"`
	IsShow        int32   `json:"isShow" gorm:"is_show"`
	IsPostage     int32   `json:"isPostage" gorm:"is_postage"`
	Postage       float64 `json:"postage" gorm:"postage"`
	StartTime     int64   `json:"startTime" gorm:"start_time"`
	StopTime      int64   `json:"stopTime" gorm:"stop_time"`
	EffectiveTime int64   `json:"effectiveTime" gorm:"effective_time"`
	Browse        int64   `json:"browse" gorm:"browse"`
	UnitName      string  `json:"unitName" gorm:"unit_name"`
	Weight        float64 `json:"weight" gorm:"weight"`
	Volume        float64 `json:"volume" gorm:"volume"`
	Num           int64   `json:"num" gorm:"num"`
	OnceNum       int64   `json:"onceNum" gorm:"once_num"`
	Quota         int64   `json:"quota" gorm:"quota"`
	QuotaShow     int64   `json:"quotaShow" gorm:"quota_show"`
	Virtual       int64   `json:"virtual" gorm:"virtual"`
	HomeModuleIDs string  `json:"homeModuleIds" gorm:"home_module_ids"`
}

// FormatCombinationMainInfo 提取拼团商品主信息。
func FormatCombinationMainInfo(sec *Combination) *CombinationMainInfo {
	if sec == nil {
		return nil
	}
	return &CombinationMainInfo{
		ID:            sec.ID,
		ProductID:     sec.ProductID,
		MerID:         sec.MerID,
		Attr:          sec.Attr,
		Price:         sec.Price,
		Sort:          sec.Sort,
		Sales:         sec.Sales,
		Stock:         sec.Stock,
		IsHost:        sec.IsHost,
		IsShow:        sec.IsShow,
		IsPostage:     sec.IsPostage,
		Postage:       sec.Postage,
		StartTime:     sec.StartTime,
		StopTime:      sec.StopTime,
		EffectiveTime: sec.EffectiveTime,
		Browse:        sec.Browse,
		UnitName:      sec.UnitName,
		Weight:        sec.Weight,
		Volume:        sec.Volume,
		Num:           sec.Num,
		OnceNum:       sec.OnceNum,
		Quota:         sec.Quota,
		QuotaShow:     sec.QuotaShow,
		Virtual:       sec.Virtual,
		HomeModuleIDs: sec.HomeModuleIDs,
	}
}

type Combination struct {
	ID            int64   `json:"id,string" gorm:"id"`
	ProductID     string  `json:"productId" gorm:"product_id"`
	MerID         int64   `json:"merId,string" gorm:"mer_id"`
	Image         string  `json:"image" gorm:"image"`
	Images        string  `json:"images" gorm:"images"`
	Title         string  `json:"title" gorm:"title"`
	Attr          string  `json:"attr" gorm:"attr"`
	People        int32   `json:"people" gorm:"people"`
	Info          string  `json:"info" gorm:"info"`
	Price         float64 `json:"price" gorm:"price"`
	Sort          int32   `json:"sort" gorm:"sort"`
	Sales         int64   `json:"sales" gorm:"sales"`
	Stock         int64   `json:"stock" gorm:"stock"`
	IsHost        int32   `json:"isHost" gorm:"is_host"`
	IsShow        int32   `json:"isShow" gorm:"is_show"`
	IsPostage     int32   `json:"isPostage" gorm:"is_postage"`
	Postage       float64 `json:"postage" gorm:"postage"`
	StartTime     int64   `json:"startTime" gorm:"start_time"`
	StopTime      int64   `json:"stopTime" gorm:"stop_time"`
	EffectiveTime int64   `json:"effectiveTime" gorm:"effective_time"`
	Browse        int64   `json:"browse" gorm:"browse"`
	UnitName      string  `json:"unitName" gorm:"unit_name"`
	Weight        float64 `json:"weight" gorm:"weight"`
	Volume        float64 `json:"volume" gorm:"volume"`
	Num           int64   `json:"num" gorm:"num"`
	OnceNum       int64   `json:"onceNum" gorm:"once_num"`
	Quota         int64   `json:"quota" gorm:"quota"`
	QuotaShow     int64   `json:"quotaShow" gorm:"quota_show"`
	Virtual       int64   `json:"virtual" gorm:"virtual"`
	HomeModuleIDs string  `json:"homeModuleIds" gorm:"home_module_ids"`
	DeptID        int64   `json:"deptId" gorm:"dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"state"`
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

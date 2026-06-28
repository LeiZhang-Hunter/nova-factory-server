package models

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/shop/logistics/client/api"
	factory "nova-factory-server/app/business/shop/logistics/client/factory"
)

// ShopLogisticsConfig 物流配置
type ShopLogisticsConfig struct {
	ID     int64  `json:"id,string" gorm:"id"`
	Type   string `json:"type" gorm:"type"`
	Data   string `json:"data" gorm:"data"`
	Status *bool  `json:"status" gorm:"status"`
	DeptID int64  `json:"deptId" gorm:"dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"state"`
}

func (i *ShopLogisticsConfig) GetData() string {
	return i.Data
}

// GetType 返回集成类型，实现 observer.ShopErpIntegrationConfig 接口
func (i *ShopLogisticsConfig) GetType() string {
	return i.Type
}

func (i *ShopLogisticsConfig) Service() (api.ExpressClient, error) {
	return factory.GetRegistry().Create(i)
}

// ShopLogisticsConfigSet 物流配置新增修改参数
type ShopLogisticsConfigSet struct {
	ID     int64  `json:"id,string"`
	Type   string `json:"type" binding:"required"`
	Data   string `json:"data"`
	Status *bool  `json:"status"`
}

// ShopLogisticsConfigQuery 物流配置查询参数
type ShopLogisticsConfigQuery struct {
	Type   string `form:"type"`
	Status *bool  `form:"status"`
	Page   int64  `form:"page"`
	Size   int64  `form:"size"`
}

// ShopLogisticsConfigListData 物流配置分页数据
type ShopLogisticsConfigListData struct {
	Rows  []*ShopLogisticsConfig `json:"rows"`
	Total int64                  `json:"total"`
}

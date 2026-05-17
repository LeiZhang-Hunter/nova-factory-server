package mastermodels

import (
	"nova-factory-server/app/baize"
)

// Account ERP 结算账户
type Account struct {
	ID            int64  `json:"id,string" gorm:"column:id"`
	Name          string `json:"name" gorm:"column:name"`
	No            string `json:"no" gorm:"column:no"`
	Remark        string `json:"remark" gorm:"column:remark"`
	Status        int32  `json:"status" gorm:"column:status"`
	Sort          int32  `json:"sort" gorm:"column:sort"`
	DefaultStatus *bool  `json:"defaultStatus" gorm:"column:default_status"`
	DeptID        int64  `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// AccountUpsert ERP 结算账户新增修改参数
type AccountUpsert struct {
	ID            int64  `json:"id,string"`
	Name          string `json:"name" binding:"required" label:"账户名称"`
	No            string `json:"no" binding:"required" label:"账户编码，对接 erp_order_account.finance_code"`
	Remark        string `json:"remark"`
	Status        int32  `json:"status"`
	Sort          int32  `json:"sort"`
	DefaultStatus *bool  `json:"defaultStatus"`
}

// AccountQuery ERP 结算账户查询参数
type AccountQuery struct {
	No            string `form:"no" filter:"like,no"`
	Name          string `form:"name" filter:"like,name"`
	Status        *int32 `form:"status" filter:"eq,status"`
	DefaultStatus *bool  `form:"defaultStatus" filter:"eq,default_status"`
	Page          int64  `form:"page"`
	Size          int64  `form:"size"`
}

// AccountListData ERP 结算账户分页数据
type AccountListData struct {
	Rows  []*Account `json:"rows"`
	Total int64      `json:"total"`
}

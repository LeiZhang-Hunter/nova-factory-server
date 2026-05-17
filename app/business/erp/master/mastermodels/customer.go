package mastermodels

import (
	"nova-factory-server/app/baize"
)

// Customer ERP 客户
type Customer struct {
	ID          int64   `json:"id,string" gorm:"column:id"`
	Name        string  `json:"name" gorm:"column:name"`
	Code        string  `json:"code" gorm:"column:code"`
	Contact     string  `json:"contact" gorm:"column:contact"`
	Mobile      string  `json:"mobile" gorm:"column:mobile"`
	Telephone   string  `json:"telephone" gorm:"column:telephone"`
	Email       string  `json:"email" gorm:"column:email"`
	Fax         string  `json:"fax" gorm:"column:fax"`
	Remark      string  `json:"remark" gorm:"column:remark"`
	Status      int32   `json:"status" gorm:"column:status"`
	Sort        int32   `json:"sort" gorm:"column:sort"`
	TaxNo       string  `json:"taxNo" gorm:"column:tax_no"`
	TaxPercent  float64 `json:"taxPercent" gorm:"column:tax_percent"`
	BankName    string  `json:"bankName" gorm:"column:bank_name"`
	BankAccount string  `json:"bankAccount" gorm:"column:bank_account"`
	BankAddress string  `json:"bankAddress" gorm:"column:bank_address"`
	DeptID      int64   `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// CustomerUpsert ERP 客户新增修改参数
type CustomerUpsert struct {
	ID          int64   `json:"id,string"`
	Name        string  `json:"name" binding:"required" label:"客户名称"`
	Code        string  `json:"code" binding:"required" label:"客户编码，对接 erp_order.b_type_code"`
	Contact     string  `json:"contact"`
	Mobile      string  `json:"mobile"`
	Telephone   string  `json:"telephone"`
	Email       string  `json:"email"`
	Fax         string  `json:"fax"`
	Remark      string  `json:"remark"`
	Status      int32   `json:"status"`
	Sort        int32   `json:"sort"`
	TaxNo       string  `json:"taxNo"`
	TaxPercent  float64 `json:"taxPercent"`
	BankName    string  `json:"bankName"`
	BankAccount string  `json:"bankAccount"`
	BankAddress string  `json:"bankAddress"`
}

// CustomerQuery ERP 客户查询参数
type CustomerQuery struct {
	Name   string `form:"name" filter:"like,name"`
	Code   string `form:"code" filter:"like,code"`
	Status *int32 `form:"status" filter:"eq,status"`
	Page   int64  `form:"page"`
	Size   int64  `form:"size"`
}

// CustomerListData ERP 客户分页数据
type CustomerListData struct {
	Rows  []*Customer `json:"rows"`
	Total int64       `json:"total"`
}

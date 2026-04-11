package models

import "nova-factory-server/app/baize"

// User 商城用户
type User struct {
	ID           int64  `json:"id,string" db:"id"`               // 主键ID
	UserID       string `json:"userId" db:"user_id"`             // 用户业务ID
	Username     string `json:"username" db:"username"`          // 用户名
	Nickname     string `json:"nickname" db:"nickname"`          // 用户昵称
	Mobile       string `json:"mobile" db:"mobile"`              // 手机号
	Email        string `json:"email" db:"email"`                // 邮箱
	Password     string `json:"password" db:"password"`          // 密码摘要
	UserType     int32  `json:"userType" db:"user_type"`         // 用户类型
	Avatar       string `json:"avatar" db:"avatar"`              // 头像
	CompanyName  string `json:"companyName" db:"company_name"`   // 企业名称
	ContactName  string `json:"contactName" db:"contact_name"`   // 联系人
	ContactPhone string `json:"contactPhone" db:"contact_phone"` // 联系人手机号
	Status       *bool  `json:"status" db:"status"`              // 状态
	DeptID       int64  `json:"deptId" gorm:"column:dept_id"`    // 部门ID
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"` // 操作状态
}

// UserUpsert 商城用户新增修改参数
type UserUpsert struct {
	ID           int64  `json:"id,string"`                   // 主键ID
	UserID       string `json:"userId"`                      // 用户业务ID
	Username     string `json:"username" binding:"required"` // 用户名
	Nickname     string `json:"nickname"`                    // 用户昵称
	Mobile       string `json:"mobile"`                      // 手机号
	Email        string `json:"email"`                       // 邮箱
	Password     string `json:"password"`                    // 密码摘要
	UserType     int32  `json:"userType" binding:"required"` // 用户类型
	Avatar       string `json:"avatar"`                      // 头像
	CompanyName  string `json:"companyName"`                 // 企业名称
	ContactName  string `json:"contactName"`                 // 联系人
	ContactPhone string `json:"contactPhone"`                // 联系人手机号
	Status       *bool  `json:"status"`                      // 状态
}

// UserQuery 商城用户查询参数
type UserQuery struct {
	Username string `form:"username"` // 用户名
	Mobile   string `form:"mobile"`   // 手机号
	UserType int32  `form:"userType"` // 用户类型
	Status   *bool  `form:"status"`   // 状态
	Page     int64  `form:"page"`     // 页码
	Size     int64  `form:"size"`     // 每页数量
}

// UserListData 商城用户列表结果
type UserListData struct {
	Rows  []*User `json:"rows"`  // 数据列表
	Total int64   `json:"total"` // 总数
}

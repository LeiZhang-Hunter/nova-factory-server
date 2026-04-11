package shopmodels

import "time"

// User 商城用户
type User struct {
	ID           uint64    `json:"id" db:"id"`                      // 主键ID
	UserID       string    `json:"userId" db:"user_id"`             // 用户业务ID
	Username     string    `json:"username" db:"username"`          // 用户名
	Nickname     string    `json:"nickname" db:"nickname"`          // 昵称
	Mobile       string    `json:"mobile" db:"mobile"`              // 手机号
	Email        string    `json:"email" db:"email"`                // 邮箱
	Password     string    `json:"password" db:"password"`          // 密码
	UserType     int32     `json:"userType" db:"user_type"`         // 用户类型
	CompanyName  string    `json:"companyName" db:"company_name"`   // 公司名称
	ContactName  string    `json:"contactName" db:"contact_name"`   // 联系人
	ContactPhone string    `json:"contactPhone" db:"contact_phone"` // 联系电话
	Status       int32     `json:"status" db:"status"`              // 状态
	IsDeleted    int32     `json:"isDeleted" db:"is_deleted"`       // 是否删除
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`       // 创建时间
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`       // 更新时间
}

// UserUpsert 商城用户新增修改参数
type UserUpsert struct {
	ID           uint64 `json:"id"`                          // 主键ID
	UserID       string `json:"userId" binding:"required"`   // 用户业务ID
	Username     string `json:"username" binding:"required"` // 用户名
	Nickname     string `json:"nickname"`                    // 昵称
	Mobile       string `json:"mobile"`                      // 手机号
	Email        string `json:"email"`                       // 邮箱
	Password     string `json:"password"`                    // 密码
	UserType     int32  `json:"userType" binding:"required"` // 用户类型
	CompanyName  string `json:"companyName"`                 // 公司名称
	ContactName  string `json:"contactName"`                 // 联系人
	ContactPhone string `json:"contactPhone"`                // 联系电话
	Status       int32  `json:"status"`                      // 状态
}

// UserQuery 商城用户查询参数
type UserQuery struct {
	Username string `form:"username"` // 用户名
	Mobile   string `form:"mobile"`   // 手机号
	UserType int32  `form:"userType"` // 用户类型
	Status   int32  `form:"status"`   // 状态
	Page     int64  `form:"page"`     // 页码
	Size     int64  `form:"size"`     // 每页数量
}

// UserListData 商城用户列表结果
type UserListData struct {
	Rows  []*User `json:"rows"`  // 数据列表
	Total int64   `json:"total"` // 总数
}

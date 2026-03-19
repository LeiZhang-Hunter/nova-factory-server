package shopModels

import "time"

type User struct {
	ID           uint64    `json:"id" db:"id"`
	UserID       string    `json:"userId" db:"user_id"`
	Username     string    `json:"username" db:"username"`
	Nickname     string    `json:"nickname" db:"nickname"`
	Mobile       string    `json:"mobile" db:"mobile"`
	Email        string    `json:"email" db:"email"`
	Password     string    `json:"password" db:"password"`
	UserType     int32     `json:"userType" db:"user_type"`
	CompanyName  string    `json:"companyName" db:"company_name"`
	ContactName  string    `json:"contactName" db:"contact_name"`
	ContactPhone string    `json:"contactPhone" db:"contact_phone"`
	Status       int32     `json:"status" db:"status"`
	IsDeleted    int32     `json:"isDeleted" db:"is_deleted"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type UserUpsert struct {
	ID           uint64 `json:"id"`
	UserID       string `json:"userId" binding:"required"`
	Username     string `json:"username" binding:"required"`
	Nickname     string `json:"nickname"`
	Mobile       string `json:"mobile"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	UserType     int32  `json:"userType" binding:"required"`
	CompanyName  string `json:"companyName"`
	ContactName  string `json:"contactName"`
	ContactPhone string `json:"contactPhone"`
	Status       int32  `json:"status"`
}

type UserQuery struct {
	Username string `form:"username"`
	Mobile   string `form:"mobile"`
	UserType int32  `form:"userType"`
	Status   int32  `form:"status"`
	Page     int64  `form:"page"`
	Size     int64  `form:"size"`
}

type UserListData struct {
	Rows  []*User `json:"rows"`
	Total int64   `json:"total"`
}

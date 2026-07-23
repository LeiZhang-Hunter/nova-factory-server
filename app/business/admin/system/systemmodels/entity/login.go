package entity

type User struct {
	UserId    int64  `json:"userId,string" db:"user_id"`
	DeptId    int64  `json:"-" db:"dept_id"`
	UserName  string `json:"userName" db:"user_name"`
	Avatar    string `json:"avatar" db:"avatar" `
	DataScope string `json:"dataScope" db:"data_scope"`
	Password  string `json:"-" db:"password"`
	Status    string `json:"-" db:"status"`
	DelFlag   string `json:"-" db:"del_flag"`
	Os        string `json:"-"`
	Browser   string `json:"-"`
}

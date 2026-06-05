package systemmodels

type LoginBody struct {
	Username string `json:"username" binding:"required"` //用户名
	Password string `json:"password" binding:"required"` //密码
	Code     string `json:"code"`                        //验证码
	Uuid     string `json:"uuid"`                        //uuid
}

// LoginResp 登录返回
type LoginResp struct {
	Token      string `json:"token"`      // 登录 token
	ExpireTime int64  `json:"expireTime"` // 有效时间，单位秒
}

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

type GetInfo struct {
	User        *User    `json:"user"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

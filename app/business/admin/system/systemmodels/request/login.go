package request

type LoginBody struct {
	Username string `json:"username" binding:"required"` //用户名
	Password string `json:"password" binding:"required"` //密码
	Code     string `json:"code"`                        //验证码
	Uuid     string `json:"uuid"`                        //uuid
}

package request

import "nova-factory-server/app/baize"

type SysUserDML struct {
	UserId      int64    `json:"userId,string" db:"user_id"swaggerignore:"true"` //用户ID
	DeptId      int64    `json:"deptId,string" db:"dept_id" binding:"required"`  //部门ID
	UserName    string   `json:"userName" db:"user_name" binding:"required"`     //用户名
	NickName    string   `json:"nickName" db:"nick_name" binding:"required"`     //用户昵称
	Email       string   `json:"email" db:"email"`                               //邮箱
	Avatar      string   `json:"avatar" db:"avatar"`                             //头像
	Phonenumber string   `json:"phonenumber" db:"phonenumber"`                   //手机号
	Sex         string   `json:"sex" db:"sex"  binding:"required"`               //性别
	Password    string   `json:"password" db:"password" binding:"required"`      //密码
	DataScope   string   `json:"dataScope" db:"data_scope"`                      //权限范围
	Status      string   `json:"status" db:"status"`                             //状态
	Remark      string   `json:"remark" db:"remark"`                             //备注
	PostIds     []string `json:"postIds"`                                        //岗位IDS
	RoleIds     []string `json:"roleIds"`                                        //角色IDS
	baize.BaseEntity
}

type SysUserDataScope struct {
	UserId    int64    `json:"userId,string"  binding:"required"` //用户ID
	DataScope string   `json:"dataScope"  binding:"required"`     //数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限,无任何）权限
	DeptIds   []string `json:"deptIds"`                           //如果是自定义就是部门ID 其他不填
}

type ResetPwd struct {
	UserId   int64  `json:"userId,string" db:"user_id" binding:"required"` //用户ID
	Password string `json:"password" db:"password" binding:"required"`     //新密码
}

type EditUserStatus struct {
	UserId int64  `json:"userId,string" binding:"required"` //用户id
	Status string `json:"status" binding:"required"`        //状态
	baize.BaseEntity
}

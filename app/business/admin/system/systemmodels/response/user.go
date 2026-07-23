package response

import (
	"nova-factory-server/app/baize"
)

type SysUserVo struct {
	UserId      int64   `json:"userId,string" db:"user_id"`
	UserName    string  `json:"userName" db:"user_name" bze:"1,用户名"`
	NickName    string  `json:"nickName" db:"nick_name" bze:"2,用户昵称"`
	Sex         string  `json:"sex" db:"sex" bze:"3,性别"`
	Status      string  `json:"status" db:"status"`
	DelFlag     string  `json:"delFlag" db:"del_flag"`
	DeptId      int64   `json:"deptId,string" db:"dept_id"`
	DeptName    *string `json:"deptName" db:"dept_name" bze:"4,部门名称"`
	Leader      string  `json:"leader" db:"leader"`
	Email       string  `json:"email" db:"email"`
	Phonenumber string  `json:"phonenumber"db:"phonenumber" bze:"5,电话"`
	Avatar      string  `json:"avatar" db:"avatar"`
	DataScope   string  `json:"dataScope" db:"data_scope"`
	Remark      string  `json:"remark" db:"remark"`
	baize.BaseEntity
}

type Accredit struct {
	Posts []*SysPostVo        `json:"posts"` //岗位
	Roles []*SysRoleIdAndName `json:"roles"` //角色
}

type UserAndAccredit struct {
	User    *SysUserVo          `json:"user"`    //user
	Roles   []*SysRoleIdAndName `json:"roles"`   //角色
	RoleIds []string            `json:"roleIds"` //选择的角色Id
	Posts   []*SysPostVo        `json:"posts"`   //岗位
	PostIds []string            `json:"postIds"` //选择的岗位Id
}

type UserAndRoles struct {
	User    *SysUserVo          `json:"user"`    //user
	Roles   []*SysRoleIdAndName `json:"roles"`   //角色
	RoleIds []string            `json:"roleIds"` //选择的角色Id
}

type UserProfile struct {
	User      *SysUserVo `json:"user"`      //user
	RoleGroup string     `json:"roleGroup"` //角色
	PostGroup string     `json:"postGroup"` //选择的角色Id
}

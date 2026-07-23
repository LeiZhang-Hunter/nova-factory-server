package query

import "nova-factory-server/app/baize"

type SysRoleDQL struct {
	RoleName  string `form:"roleName" db:"role_name"`
	Status    string `form:"status" db:"status"`
	BeginTime string `form:"beginTime" db:"begin_time"`
	EndTime   string `form:"endTime" db:"end_time"`
	CreateBy  int64  `db:"create_by" swaggerignore:"true"` //创建人
	baize.BaseEntityDQL
}

type SysRoleAndUserDQL struct {
	RoleId      string `form:"roleId" db:"role_id" binding:"required"`
	UserName    string `form:"userName" db:"user_name"`
	Phonenumber string `form:"phonenumber" db:"phonenumber"`
	baize.BaseEntityDQL
}

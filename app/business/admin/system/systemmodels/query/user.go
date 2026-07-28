package query

import "nova-factory-server/app/baize"

type SysUserDQL struct {
	UserName    string `form:"userName" db:"user_name"`      //用户名
	Status      string `form:"status" db:"status"`           //状态
	Phonenumber string `form:"phonenumber" db:"phonenumber"` //电话
	BeginTime   string `form:"beginTime" db:"begin_time"`    //注册开始时间
	EndTime     string `form:"endTime" db:"end_time"`        //注册结束时间
	DeptId      int64  `form:"deptId" db:"dept_id"`          //部门ID
	baize.BaseEntityDQL
}

package query

import "nova-factory-server/app/baize"

type SysDeptDQL struct {
	ParentId int64  `form:"parentId,string" db:"parent_id"` //上级id
	DeptName string `form:"deptName" db:"dept_name"`        //部门名称
	Status   string `form:"status" db:"status"`             //状态
	baize.BaseEntityDQL
}

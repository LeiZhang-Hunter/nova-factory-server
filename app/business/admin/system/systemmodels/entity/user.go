package entity

type SysUserDeptScope struct {
	UserId int64 `db:"user_id"`
	DeptId int64 `db:"dept_id"`
}

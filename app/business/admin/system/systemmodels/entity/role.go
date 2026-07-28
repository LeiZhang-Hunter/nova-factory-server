package entity

type SysRole struct {
	RoleId   int64  `db:"role_id"`
	RoleName string `db:"role_name"`
}

type SysRolePermission struct {
	RoleId       int64 `db:"role_id"`
	PermissionId int64 `db:"permission_id"`
}

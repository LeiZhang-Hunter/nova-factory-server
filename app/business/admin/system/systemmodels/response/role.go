package response

import "nova-factory-server/app/baize"

type SysRoleVo struct {
	RoleId        int64    `json:"roleId,string" db:"role_id"`
	RoleName      string   `json:"roleName" db:"role_name" bze:"1,角色名称"`
	RoleSort      int      `json:"roleSort" db:"role_sort"`
	Status        string   `json:"status"  db:"status"`
	Remake        string   `json:"remark" db:"remark"`
	PermissionIds []string `json:"permissionIds"`
	baize.BaseEntity
}

type SysRoleIdAndName struct {
	RoleId   int64  `json:"roleId,string" db:"role_id"`
	RoleName string `json:"roleName" db:"role_name" `
}

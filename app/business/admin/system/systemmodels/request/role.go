package request

import "nova-factory-server/app/baize"

type SysRoleDML struct {
	RoleId        int64    `json:"RoleId,string" db:"role_id"`
	RoleName      string   `json:"roleName" db:"role_name"`
	RoleSort      int      `json:"roleSort" db:"role_sort"`
	Status        string   `json:"status" db:"status"`
	DelFlag       string   `json:"delFlag" db:"del_flag"`
	Remake        string   `json:"remark" db:"remark"`
	PermissionIds []string `json:"permissionIds"`
	baize.BaseEntity
}

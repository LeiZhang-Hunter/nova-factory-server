package query

import "nova-factory-server/app/baize"

type SysPermissionDQL struct {
	Status string `form:"status" db:"status"` // 状态
	baize.BaseEntity
}

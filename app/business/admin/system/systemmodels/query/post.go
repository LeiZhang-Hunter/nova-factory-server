package query

import "nova-factory-server/app/baize"

type SysPostDQL struct {
	PostCode string `form:"postCode" db:"post_code"`
	Status   string `form:"status" db:"status"`
	PostName string `form:"postName" db:"post_name"`
	baize.BaseEntityDQL
}

package query

import "nova-factory-server/app/baize"

type SysDictTypeDQL struct {
	DictName string `form:"dictName" db:"dict_name"`
	Status   string `form:"status" db:"status"`
	DictType string `form:"dictType" db:"dict_type"`
	baize.BaseEntityDQL
}

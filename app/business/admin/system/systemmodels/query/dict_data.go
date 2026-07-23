package query

import "nova-factory-server/app/baize"

type SysDictDataDQL struct {
	DictType  string `form:"dictType" db:"dict_type"`
	DictLabel string `form:"dictLabel" db:"dict_label"`
	Status    string `form:"status" db:"status"`
	baize.BaseEntityDQL
}

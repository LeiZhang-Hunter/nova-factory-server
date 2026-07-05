package metricdaoimpl

import "nova-factory-server/app/datasource/iotdb"

type iPrometheusDbExecutor struct {
	iotDb *iotdb.TSDBStorage
}

func newPrometheusDbExecutor() iDaoExport {
	return &iotDbExport{}
}

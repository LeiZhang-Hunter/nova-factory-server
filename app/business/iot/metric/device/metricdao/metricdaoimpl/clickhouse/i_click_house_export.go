package clickhouse

import (
	"nova-factory-server/app/business/iot/metric/device/metricdao"
	"nova-factory-server/app/datasource/clickhouse"

	"github.com/gin-gonic/gin"
)

type iClickHouseExport struct {
	tableName  string
	adder      *adder
	query      *query
	clickhouse *clickhouse.ClickHouse
}

func (i *iClickHouseExport) Questioner() metricdao.IMetricQueryDao {
	return i.query
}

func NewIClickHouseExport(clickhouse *clickhouse.ClickHouse) metricdao.IMetricStorageDao {
	i := &iClickHouseExport{
		clickhouse: clickhouse,
		tableName:  "nova_metrics_device",
	}
	i.adder = newAdder(i.tableName, clickhouse)
	i.query = newQuery(i.tableName, clickhouse)
	return i
}

func (i *iClickHouseExport) InstallDevice(c *gin.Context, deviceId int64, templateId int64) error {
	return nil
}
func (i *iClickHouseExport) UnInStallDevice(c *gin.Context, deviceId int64, templateId int64) error {
	return nil
}

// InstallRunStatusDevice 运行状态设备模板
func (i *iClickHouseExport) InstallRunStatusDevice(c *gin.Context, deviceId int64) error {
	return nil
}

// UnInStallRunStatusDevice 卸载设备运行状态模板
func (i *iClickHouseExport) UnInStallRunStatusDevice(c *gin.Context, deviceId int64) error {
	return nil
}

// Template 模板
func (i *iClickHouseExport) Template() metricdao.IIotStorageTemplateDao {
	return nil
}

// Adder 添加
func (i *iClickHouseExport) Adder() metricdao.IMetricAdderDao {
	return i.adder
}

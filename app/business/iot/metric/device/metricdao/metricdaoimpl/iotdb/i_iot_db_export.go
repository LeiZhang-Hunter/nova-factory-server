package iotdb

import (
	"fmt"
	"github.com/apache/iotdb-client-go/client"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/iot/metric/device/metricdao"
	iotdb2 "nova-factory-server/app/constant/iotdb"
	"nova-factory-server/app/datasource/iotdb"
	stdtime "time"
)

type iotDbExport struct {
	iotDb *iotdb.IotDb
	// template 设备模板
	template *deviceTemplate

	adder *adder
	query *query
}

func (i *iotDbExport) Questioner() metricdao.IMetricQueryDao {
	//TODO implement me
	return i.query
}

func NewIotDbExport(iotDb *iotdb.IotDb) metricdao.IMetricStorageDao {
	i := &iotDbExport{
		iotDb:    iotDb,
		template: newDeviceTemplate(iotDb),
		adder:    newAdder(iotDb),
		query:    newQuery(iotDb),
	}
	i.init()
	return i
}

func (i *iotDbExport) init() {
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Fatal("iotdb.GetSession()", zap.Error(err))
		return
	}
	defer i.iotDb.PutSession(session)

	var initErr error
	for attempt := 1; attempt <= 10; attempt++ {
		initErr = i.ensureDatabase(session, "root.device")
		if initErr == nil {
			initErr = i.ensureDatabase(session, "root.run_status_device")
		}
		if initErr == nil {
			break
		}

		zap.L().Error("initialize iotdb database failed",
			zap.Int("attempt", attempt),
			zap.Error(initErr),
		)
		if attempt < 10 {
			stdtime.Sleep(stdtime.Second)
		}
	}
	if initErr != nil {
		zap.L().Fatal("initialize iotdb database failed after retries", zap.Error(initErr))
		return
	}

	// 创建设备数据采集模板
	session.ExecuteStatement(fmt.Sprintf("create device template %s ALIGNED (value DOUBLE)", iotdb2.NOVA_DEVICE_TEMPLATE))
	// 创建设备运行时间统计模板
	session.ExecuteStatement(fmt.Sprintf("create device template %s ALIGNED (duration INT64, status INT64)", iotdb2.NOVA_DEVICE_RUN_TEMPLATE))
}

func (i *iotDbExport) ensureDatabase(session client.Session, database string) error {
	statement, err := session.ExecuteStatement(fmt.Sprintf("count databases %s", database))
	if err != nil {
		return err
	}

	hasDatabase, err := statement.Next()
	if err != nil {
		return err
	}
	if !hasDatabase {
		return fmt.Errorf("unable to read database count for %s", database)
	}
	if statement.GetInt32("count") > 0 {
		return nil
	}

	_, err = session.ExecuteStatement(fmt.Sprintf("create database %s", database))
	return err
}

// InstallDevice 安装设备模板
func (i *iotDbExport) InstallDevice(c *gin.Context, deviceId int64, templateId int64) error {
	if templateId <= 0 {
		return fmt.Errorf("invalid device template id: %d", templateId)
	}

	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return err
	}
	defer i.iotDb.PutSession(session)

	name := iotdb2.MakeDeviceDataName(deviceId)
	templateName := iotdb2.MakeDeviceTemplateName(templateId)
	// 创建设备模板
	group, err := session.SetStorageGroup(name)
	if err != nil {
		zap.L().Error("创建设备数据库失败, ", zap.Error(err), zap.Any("code", group.GetCode()))
		return err
	}

	// 挂载设备模板
	_, err = session.ExecuteStatement(fmt.Sprintf("set device template %s to %s", templateName, name))
	if err != nil {
		zap.L().Error("绑定设备数据库失败, ", zap.String("template", templateName), zap.Error(err))
		return err
	}

	// 激活设备模板
	_, err = session.ExecuteStatement(fmt.Sprintf("create timeseries using device template on %s", name))
	if err != nil {
		zap.L().Error("激活设备模板失败, ", zap.String("template", templateName), zap.Error(err))
		return err
	}
	return nil
}

// InstallRunStatusDevice 运行状态设备模板
func (i *iotDbExport) InstallRunStatusDevice(c *gin.Context, deviceId int64) error {
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return err
	}
	defer i.iotDb.PutSession(session)

	name := iotdb2.MakeRunDeviceTemplateName(deviceId)
	// 创建设备模板
	group, err := session.SetStorageGroup(name)
	if err != nil {
		zap.L().Error("创建设备数据库失败, ", zap.Error(err), zap.Any("code", group.GetCode()))
		return err
	}

	// 挂载设备模板
	_, err = session.ExecuteStatement(fmt.Sprintf("set device template %s to %s", iotdb2.NOVA_DEVICE_RUN_TEMPLATE, name))
	if err != nil {
		zap.L().Error("绑定设备数据库失败, ", zap.Error(err))
		return err
	}

	// 激活设备模板
	_, err = session.ExecuteStatement(fmt.Sprintf("create timeseries using device template on %s", name))
	if err != nil {
		zap.L().Error("激活设备模板失败, ", zap.Error(err))
		return err
	}
	return nil
}

// UnInStallRunStatusDevice 卸载设备运行状态模板
func (i *iotDbExport) UnInStallRunStatusDevice(c *gin.Context, deviceId int64) error {
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return err
	}
	defer i.iotDb.PutSession(session)

	name := iotdb2.MakeRunDeviceTemplateName(deviceId)

	// 删除模板表示的某一组时间序列
	_, err = session.ExecuteStatement(fmt.Sprintf("deactivate device template %s from %s", iotdb2.NOVA_DEVICE_RUN_TEMPLATE, name))
	if err != nil {
		zap.L().Error("deactivate  device template", zap.Error(err))
		return err
	}

	_, err = session.ExecuteStatement(fmt.Sprintf("unset device template %s from %s", iotdb2.NOVA_DEVICE_RUN_TEMPLATE, name))
	if err != nil {
		zap.L().Error("unset  device template", zap.Error(err))
		return err
	}

	_, err = session.ExecuteStatement(fmt.Sprintf("drop database %s", name))
	if err != nil {
		zap.L().Error("unset  device template", zap.Error(err))
		return err
	}
	return nil
}

// UnInStallDevice 卸载设备模板
func (i *iotDbExport) UnInStallDevice(c *gin.Context, deviceId int64, templateId int64) error {
	if templateId <= 0 {
		return fmt.Errorf("invalid device template id: %d", templateId)
	}

	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return err
	}
	defer i.iotDb.PutSession(session)

	name := iotdb2.MakeDeviceDataName(deviceId)
	templateName := iotdb2.MakeDeviceTemplateName(templateId)

	// 删除模板表示的某一组时间序列
	_, err = session.ExecuteStatement(fmt.Sprintf("deactivate device template %s from %s", templateName, name))
	if err != nil {
		zap.L().Error("deactivate  device template", zap.String("template", templateName), zap.Error(err))
		return err
	}

	_, err = session.ExecuteStatement(fmt.Sprintf("unset device template %s from %s", templateName, name))
	if err != nil {
		zap.L().Error("unset  device template", zap.String("template", templateName), zap.Error(err))
		return err
	}

	_, err = session.ExecuteStatement(fmt.Sprintf("drop database %s", name))
	if err != nil {
		zap.L().Error("unset  device template", zap.Error(err))
		return err
	}
	return nil
}

// ExportTimeData 导入时序数据
func metricEndTime(end uint64) stdtime.Time {
	if end == 0 {
		return stdtime.Now()
	}
	return stdtime.UnixMilli(int64(end))
}

// Template returns the template DAO for the active IoTDB storage backend.
func (i *iotDbExport) Template() metricdao.IIotStorageTemplateDao {
	return i.template
}

// Adder returns the metric adder for the active IoTDB storage backend.
func (i *iotDbExport) Adder() metricdao.IMetricAdderDao {
	return i.adder
}

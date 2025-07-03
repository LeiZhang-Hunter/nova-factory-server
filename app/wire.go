//go:build wireinject
// +build wireinject

package main

import (
	"nova-factory-server/app/business/ai/aiDataSetController"
	"nova-factory-server/app/business/ai/aiDataSetDao/aiDataSetDaoImpl"
	"nova-factory-server/app/business/ai/aiDataSetService/aiDataSetServiceImpl"
	"nova-factory-server/app/business/asset/device/deviceController"
	"nova-factory-server/app/business/asset/device/deviceDao/deviceDaoImpl"
	"nova-factory-server/app/business/asset/device/deviceService/deviceServiceImpl"
	"nova-factory-server/app/business/asset/material/materialController"
	"nova-factory-server/app/business/asset/material/materialDao/materialDaoImpl"
	"nova-factory-server/app/business/asset/material/materialService/materialServiceImpl"
	"nova-factory-server/app/business/batch/batchController"
	"nova-factory-server/app/business/batch/batchDao/batchDaoImpl"
	"nova-factory-server/app/business/batch/batchService/batchServiceImpl"
	"nova-factory-server/app/business/craft/craftRouteController"
	"nova-factory-server/app/business/craft/craftRouteDao/craftRouteDaoImpl"
	"nova-factory-server/app/business/craft/craftRouteService/craftRouteServiceImpl"
	"nova-factory-server/app/business/daemonize/daemonizeController"
	"nova-factory-server/app/business/daemonize/daemonizeDao/daemonizeDaoImpl"
	"nova-factory-server/app/business/daemonize/daemonizeService/daemonizeServiceImpl"
	"nova-factory-server/app/business/defect/defectController"
	"nova-factory-server/app/business/defect/defectDao/serviceImpl"
	"nova-factory-server/app/business/defect/defectService/defectServiceImpl"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorController"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService/deviceMonitorServiceImpl"
	"nova-factory-server/app/business/metric/device/metricController"
	"nova-factory-server/app/business/metric/device/metricDao/metricDaoIMpl"
	"nova-factory-server/app/business/metric/device/metricService/metricServiceImpl"
	"nova-factory-server/app/business/monitor/monitorController"
	"nova-factory-server/app/business/monitor/monitorDao/monitorDaoImpl"
	"nova-factory-server/app/business/monitor/monitorService/monitorServiceImpl"
	"nova-factory-server/app/business/qcIndex/qcIndexController"
	"nova-factory-server/app/business/qcIndex/qcIndexDao/qcIndexDaoImpl"
	"nova-factory-server/app/business/qcIndex/qcIndexService/qcIndexServiceImpl"
	"nova-factory-server/app/business/qcIpqc/qcIpqcController"
	"nova-factory-server/app/business/qcIpqc/qcIpqcDao/qcIpqcDaoImpl"
	"nova-factory-server/app/business/qcIpqc/qcIpqcService/qcIpqcServiceImpl"
	"nova-factory-server/app/business/qcIqc/qcIqcController"
	"nova-factory-server/app/business/qcIqc/qcIqcDao/qcIqcDaoImpl"
	"nova-factory-server/app/business/qcIqc/qcIqcService/qcIqcServiceImpl"
	"nova-factory-server/app/business/qcOqc/qcOqcController"
	"nova-factory-server/app/business/qcOqc/qcOqcDao/qcOqcDaoImpl"
	"nova-factory-server/app/business/qcOqc/qcOqcService/qcOqcServiceImpl"
	"nova-factory-server/app/business/qcRqc/qcRqcController"
	"nova-factory-server/app/business/qcRqc/qcRqcDao/qcRqcDaoImpl"
	"nova-factory-server/app/business/qcRqc/qcRqcService/qcRqcServiceImpl"
	"nova-factory-server/app/business/qcTemplate/qcTemplateController"
	"nova-factory-server/app/business/qcTemplate/qcTemplateDao/qcTemplateDaoImpl"
	"nova-factory-server/app/business/qcTemplate/qcTemplateService/qcTemplateServiceImpl"
	"nova-factory-server/app/business/system/systemController"
	"nova-factory-server/app/business/system/systemDao/systemDaoImpl"
	"nova-factory-server/app/business/system/systemService/systemServiceImpl"
	"nova-factory-server/app/business/tool/toolController"
	"nova-factory-server/app/business/tool/toolDao/toolDaoImpl"
	"nova-factory-server/app/business/tool/toolService/toolServiceImpl"
	"nova-factory-server/app/datasource"
	"nova-factory-server/app/routes"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func wireApp() (*gin.Engine, func(), error) {
	panic(wire.Build(
		toolDaoImpl.ProviderSet,
		toolServiceImpl.ProviderSet,
		toolController.ProviderSet,

		systemDaoImpl.ProviderSet,
		systemServiceImpl.ProviderSet,
		systemController.ProviderSet,

		monitorDaoImpl.ProviderSet,
		monitorServiceImpl.ProviderSet,
		monitorController.ProviderSet,

		deviceDaoImpl.ProviderSet,
		deviceServiceImpl.ProviderSet,
		deviceController.ProviderSet,

		materialDaoImpl.ProviderSet,
		materialServiceImpl.ProviderSet,
		materialController.ProviderSet,

		aiDataSetDaoImpl.ProviderSet,
		aiDataSetServiceImpl.ProviderSet,
		aiDataSetController.ProviderSet,

		craftRouteDaoImpl.ProviderSet,
		craftRouteServiceImpl.ProviderSet,
		craftRouteController.ProviderSet,

		metricDaoIMpl.ProviderSet,
		metricServiceImpl.ProviderSet,
		metricController.ProviderSet,

		daemonizeDaoImpl.ProviderSet,
		daemonizeServiceImpl.ProviderSet,
		daemonizeController.ProviderSet,

		deviceMonitorServiceImpl.ProviderSet,
		deviceMonitorController.ProviderSet,

		serviceImpl.ProviderSet,
		defectServiceImpl.ProviderSet,
		defectController.ProviderSet,

		batchDaoImpl.ProviderSet,
		batchServiceImpl.ProviderSet,
		batchController.ProviderSet,

		qcIndexDaoImpl.ProviderSet,
		qcIndexServiceImpl.ProviderSet,
		qcIndexController.ProviderSet,

		qcIqcDaoImpl.ProviderSet,
		qcIqcServiceImpl.ProviderSet,
		qcIqcController.ProviderSet,

		qcOqcDaoImpl.ProviderSet,
		qcOqcServiceImpl.ProviderSet,
		qcOqcController.ProviderSet,

		qcRqcDaoImpl.ProviderSet,
		qcRqcServiceImpl.ProviderSet,
		qcRqcController.ProviderSet,

		qcTemplateDaoImpl.ProviderSet,
		qcTemplateServiceImpl.ProviderSet,
		qcTemplateController.ProviderSet,

		qcIpqcDaoImpl.ProviderSet,
		qcIpqcServiceImpl.ProviderSet,
		qcIpqcController.ProviderSet,

		datasource.ProviderSet,
		routes.ProviderSet,
	))
}

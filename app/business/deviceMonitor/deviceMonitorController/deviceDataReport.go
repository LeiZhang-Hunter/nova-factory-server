package deviceMonitorController

import (
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"net/url"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
	"nova-factory-server/app/business/metric/device/metricService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"strconv"
)

type DeviceReport struct {
	service             metricService.IMetricService
	deviceReportService deviceMonitorService.IDeviceDataReportService
}

func NewDeviceReport(service metricService.IMetricService, deviceReportService deviceMonitorService.IDeviceDataReportService) *DeviceReport {
	return &DeviceReport{
		service:             service,
		deviceReportService: deviceReportService,
	}
}

func (d *DeviceReport) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/device/monitor/data")
	group.GET("/dev/list", middlewares.HasPermission("device:monitor:data:dev:list"), d.DevList)
	group.POST("/list", middlewares.HasPermission("device:monitor:data:list"), d.List)
	group.POST("/export", middlewares.HasPermission("device:monitor:data:export"), d.Export)
}

// DevList 设备测点列表
// @Summary 设备测点列表
// @Description 设备测点列表
// @Tags 设备监控/设备监控
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /device/monitor/data/dev/list [get]
func (d *DeviceReport) DevList(c *gin.Context) {
	list, err := d.deviceReportService.DevList(c)
	if err != nil {
		return
	}
	baizeContext.SuccessData(c, list)
}

// List 实时数据
// @Summary 实时数据
// @Description 实时数据
// @Tags 设备监控/设备监控
// @Param  object body deviceMonitorModel.DevDataReq true "实时数据表请求参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /device/monitor/data/list [post]
func (d *DeviceReport) List(c *gin.Context) {
	req := new(deviceMonitorModel.DevDataReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := d.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// Export 导出实时数据
// @Summary 导出实时数据
// @Description 导出实时数据
// @Tags 设备监控/设备监控
// @Param  object body deviceMonitorModel.DevDataReq true "实时数据表请求参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /device/monitor/data/export [post]
func (d *DeviceReport) Export(c *gin.Context) {
	req := new(deviceMonitorModel.DevDataReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	// 创建Excel文件
	f := excelize.NewFile()
	defer f.Close()
	// 创建工作表
	index, _ := f.NewSheet("Sheet1")
	// 设置表头
	_ = f.SetCellValue("Sheet1", "A1", "测点名称")
	_ = f.SetCellValue("Sheet1", "B1", "测点单位")
	_ = f.SetCellValue("Sheet1", "C1", "设备数值")
	_ = f.SetCellValue("Sheet1", "D1", "时间")

	// 设置活动工作表
	f.SetActiveSheet(index)
	req.Size = 1500
	req.Page = 1
	data, err := d.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	if data == nil {
		baizeContext.Waring(c, "导出失败")
		return
	}
	for i, row := range data.Rows {
		_ = f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), row.Name)
		_ = f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), row.Unit)
		_ = f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), row.Value)
		_ = f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), row.Time)
	}

	// 设置响应头 - 解决中文文件名问题
	filename := url.QueryEscape("实时数据导出报表.xlsx")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+filename)

	// 写入响应
	_ = f.Write(c.Writer)
}

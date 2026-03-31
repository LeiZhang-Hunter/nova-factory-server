package devicemonitorcontroller

import (
	"net/url"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitorservice"
	metricService2 "nova-factory-server/app/business/iot/metric/device/metricservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

// DeviceReport 设备接入接口
type DeviceReport struct {
	service             metricService2.IMetricService
	deviceReportService devicemonitorservice.IDeviceDataReportService
	devMapService       metricService2.IDevMapService
}

func NewDeviceReport(service metricService2.IMetricService,
	deviceReportService devicemonitorservice.IDeviceDataReportService,
	devMapService metricService2.IDevMapService) *DeviceReport {
	return &DeviceReport{
		service:             service,
		deviceReportService: deviceReportService,
		devMapService:       devMapService,
	}
}

func (d *DeviceReport) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/device/monitor/data")
	group.GET("/dev/list", middlewares.HasPermission("device:monitor:data:dev:list"), d.DevList)
	group.POST("/list", middlewares.HasPermission("device:monitor:data:list"), d.List)
	group.POST("/export", middlewares.HasPermission("device:monitor:data:export"), d.Export)

	metric := router.Group("/metric/time/seq")
	metric.GET("/list", middlewares.HasPermission("metric:time:seq:list"), d.SearchTimeSeqList)
	metric.POST("/set", middlewares.HasPermission("metric:time:seq:set"), d.SetDevMap)
	metric.DELETE("/remove", middlewares.HasPermission("metric:time:seq:remove"), d.RemoveDevMap)
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
// @Param  object body devicemonitormodel.DevDataReq true "实时数据表请求参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /device/monitor/data/list [post]
func (d *DeviceReport) List(c *gin.Context) {
	req := new(devicemonitormodel.DevDataReq)
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

	if list == nil {
		baizeContext.SuccessData(c, list)
		return
	}

	devList, err := d.devMapService.GetDevList(c, req.Dev)
	if err != nil {
		zap.L().Error("get dev list error", zap.Error(err))
		return
	}
	var devMap map[string]*devicemonitormodel.SysIotDbDevMapData = make(map[string]*devicemonitormodel.SysIotDbDevMapData)

	for _, dev := range devList {
		devMap[dev.Device] = &dev
	}

	for k, v := range list.Rows {
		value, ok := devMap[v.Dev]
		if !ok {
			continue
		}
		list.Rows[k].Name = value.DataName
		list.Rows[k].Unit = value.Unit
		list.Rows[k].DeviceID = value.DeviceID
		list.Rows[k].TemplateID = value.TemplateID
		list.Rows[k].DataID = value.DataID
		list.Rows[k].DevName = value.DevName
	}
	baizeContext.SuccessData(c, list)
}

// Export 导出实时数据
// @Summary 导出实时数据
// @Description 导出实时数据
// @Tags 设备监控/设备监控
// @Param  object body devicemonitormodel.DevDataReq true "实时数据表请求参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /device/monitor/data/export [post]
func (d *DeviceReport) Export(c *gin.Context) {
	req := new(devicemonitormodel.DevDataReq)
	err := c.ShouldBind(req)
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

	devList, err := d.devMapService.GetDevList(c, req.Dev)
	if err != nil {
		zap.L().Error("get dev list error", zap.Error(err))
		return
	}
	var devMap map[string]*devicemonitormodel.SysIotDbDevMapData = make(map[string]*devicemonitormodel.SysIotDbDevMapData)

	for _, dev := range devList {
		devMap[dev.Device] = &dev
	}

	for k, v := range data.Rows {
		value, ok := devMap[v.Dev]
		if !ok {
			continue
		}
		data.Rows[k].Name = value.DataName
		data.Rows[k].Unit = value.Unit
		data.Rows[k].DeviceID = value.DeviceID
		data.Rows[k].TemplateID = value.TemplateID
		data.Rows[k].DataID = value.DataID
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

// SearchTimeSeqList 时序数据测点
// @Summary 时序数据测点
// @Description 时序数据测点
// @Param  object query devicemodels.DeviceListReq true "时序数据测点"
// @Tags 设备监控/设备监控
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /metric/time/seq/list [get]
func (d *DeviceReport) SearchTimeSeqList(c *gin.Context) {
	req := new(devicemonitormodel.DevListReq)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := d.deviceReportService.GetDevList(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}

// SetDevMap 设置测点映射
// @Summary 设置测点映射
// @Description 设置测点映射
// @Param  object body devicemonitormodel.SetDevMapInfo true "时序数据测点"
// @Tags 设备监控/设置测点映射
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /metric/time/seq/set [post]
func (d *DeviceReport) SetDevMap(c *gin.Context) {
	info := new(devicemonitormodel.SetDevMapInfo)
	err := c.ShouldBindJSON(info)
	if err != nil {
		zap.L().Error("SetDevMap <UNK>", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	err = d.deviceReportService.SetDevMap(c, info)
	if err != nil {
		zap.L().Error("SetDevMap <UNK>", zap.Error(err))
		return
	}
	baizeContext.Success(c)
}

// RemoveDevMap 删除测点映射
// @Summary 删除测点映射
// @Description 删除测点映射
// @Param  object body devicemonitormodel.RemoveDevMapInfo true "时序数据测点"
// @Tags 设备监控/删除测点映射
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /metric/time/seq/remove [delete]
func (d *DeviceReport) RemoveDevMap(c *gin.Context) {
	info := new(devicemonitormodel.RemoveDevMapInfo)
	err := c.ShouldBindJSON(info)
	if err != nil {
		zap.L().Error("SetDevMap <UNK>", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	err = d.deviceReportService.RemoveDevMap(c, info.Device)
	if err != nil {
		zap.L().Error("SetDevMap <UNK>", zap.Error(err))
		return
	}
	baizeContext.Success(c)
}

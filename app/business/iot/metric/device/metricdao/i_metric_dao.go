package metricdao

import (
	"github.com/gin-gonic/gin"
)

type IMetricDao interface {
	IIotStorageMetricDao
	IMetricAdderDao
	IMetricQueryDao
}

type IMetricStorageDao interface {
	InstallDevice(c *gin.Context, deviceId int64, templateId int64) error
	UnInStallDevice(c *gin.Context, deviceId int64, templateId int64) error
	InstallRunStatusDevice(c *gin.Context, deviceId int64) error
	UnInStallRunStatusDevice(c *gin.Context, deviceId int64) error
	Template() IIotStorageTemplateDao
	Adder() IMetricAdderDao
	Questioner() IMetricQueryDao
}

type IIotStorageMetricDao interface {
	InstallDevice(c *gin.Context, deviceId int64, templateId int64) error
	UnInStallDevice(c *gin.Context, deviceId int64, templateId int64) error
	// InstallRunStatusDevice 运行状态设备模板
	InstallRunStatusDevice(c *gin.Context, deviceId int64) error
	// UnInStallRunStatusDevice 卸载设备运行状态模板
	UnInStallRunStatusDevice(c *gin.Context, deviceId int64) error

	// Template 模板
	Template() IIotStorageTemplateDao

	// Adder 添加
	Adder() IMetricAdderDao
}

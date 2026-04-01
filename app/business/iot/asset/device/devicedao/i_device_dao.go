package devicedao

import (
	"nova-factory-server/app/business/iot/asset/device/devicemodels"

	"github.com/gin-gonic/gin"
)

type IDeviceDao interface {
	InsertDevice(c *gin.Context, device *devicemodels.DeviceInfo) (*devicemodels.DeviceVO, error)
	UpdateDevice(c *gin.Context, device *devicemodels.DeviceInfo) (*devicemodels.DeviceVO, error)
	GetDeviceGroupByName(c *gin.Context, name string) (*devicemodels.DeviceVO, error)
	GetNoExitIdDeviceGroupByName(c *gin.Context, name string, id uint64) (*devicemodels.DeviceVO, error)
	SelectDeviceList(c *gin.Context, req *devicemodels.DeviceListReq) (*devicemodels.DeviceInfoListData, error)
	DeleteByDeviceIds(c *gin.Context, ids []int64) error
	GetLocalByGateWayId(c *gin.Context, id int64) ([]*devicemodels.DeviceVO, error)
	GetByIds(c *gin.Context, ids []int64) ([]*devicemodels.DeviceVO, error)
	GetByIdString(c *gin.Context, id string) (*devicemodels.DeviceVO, error)
	GetById(c *gin.Context, id int64) (*devicemodels.DeviceVO, error)
	GetByTag(c *gin.Context, number string) (*devicemodels.DeviceVO, error)
	// SelectPublicDeviceList 非登录情况下请求的接口
	SelectPublicDeviceList(c *gin.Context, req *devicemodels.DeviceListReq) (*devicemodels.DeviceInfoListData, error)
	// All 读取全部信息
	All(c *gin.Context) ([]*devicemodels.DeviceVO, error)
	// Count 设备数量统计
	Count(c *gin.Context) (int64, error)
}

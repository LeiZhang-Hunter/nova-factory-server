package deviceMonitorModel

type DeviceLayoutRequest struct {
	FloorId int64 `form:"floorId,string" binding:"required"`
}

package deviceMonitorDaoImpl

import "nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"

type deviceStatusSort struct {
	runSortUtilization *[]*deviceMonitorModel.DeviceUtilizationDataV2
	status             int
}

func newDeviceStatusSort(status int, runSortUtilization *[]*deviceMonitorModel.DeviceUtilizationDataV2) *deviceStatusSort {
	return &deviceStatusSort{
		status:             status,
		runSortUtilization: runSortUtilization,
	}
}

func (p deviceStatusSort) Less(i, j int) bool {
	return (*p.runSortUtilization)[i].StatusMap[p.status].Time < (*p.runSortUtilization)[i].StatusMap[p.status].Time
	// 返回 p[i] > p[j] // 降序排序
}

func (p deviceStatusSort) Swap(i, j int) {
	(*p.runSortUtilization)[i], (*p.runSortUtilization)[j] = (*p.runSortUtilization)[j], (*p.runSortUtilization)[i]
}

func (p deviceStatusSort) Len() int {
	return len(*p.runSortUtilization)
}

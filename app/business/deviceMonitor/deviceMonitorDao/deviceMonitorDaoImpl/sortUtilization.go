package deviceMonitorDaoImpl

import "nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"

type stopSortUtilization []*deviceMonitorModel.DeviceUtilizationData

func (p stopSortUtilization) Less(i, j int) bool {
	return p[i].StopTime < p[j].StopTime // 升序排序
	// 返回 p[i] > p[j] // 降序排序
}

func (p stopSortUtilization) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p stopSortUtilization) Len() int {
	return len(p)
}

type runSortUtilization []*deviceMonitorModel.DeviceUtilizationData

func (p runSortUtilization) Less(i, j int) bool {
	return p[i].RunTime < p[j].RunTime // 升序排序
	// 返回 p[i] > p[j] // 降序排序
}

func (p runSortUtilization) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p runSortUtilization) Len() int {
	return len(p)
}

type waitSortUtilization []*deviceMonitorModel.DeviceUtilizationData

func (p waitSortUtilization) Less(i, j int) bool {
	return p[i].WaitTime < p[j].WaitTime // 升序排序
}

func (p waitSortUtilization) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p waitSortUtilization) Len() int {
	return len(p)
}

package region

import (
	"github.com/itmisx/go_regions"
	"sync"
)

var mtx sync.Mutex

type Region struct {
	ID    int64  `json:"id" gorm:"column:id"`
	Name  string `json:"name" gorm:"column:name"`
	Level int    `json:"level" gorm:"column:level"` // 0-省 1-市 2-区 3-街道
}

func GetRegionInfo(id int) *Region {
	mtx.Lock()
	defer mtx.Unlock()
	info := go_regions.RegionInfo(id)
	return &Region{
		ID:    info.ID,
		Name:  info.Name,
		Level: info.Level,
	}
}

func GetRegionList(pid int) []*Region {
	var records []*Region = make([]*Region, 0)
	list := go_regions.RegionList(pid)
	for _, item := range list {
		records = append(records, &Region{
			ID:    item.ID,
			Name:  item.Name,
			Level: item.Level,
		})
	}
	return records
}

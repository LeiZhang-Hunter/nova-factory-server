package impl

import (
	"encoding/json"
	"fmt"

	"nova-factory-server/app/business/data/dao"
	"nova-factory-server/app/business/data/models/dto"
	"nova-factory-server/app/business/data/models/entity"
	"nova-factory-server/app/business/data/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

// ICollectorServiceImpl 采集器服务实现
type ICollectorServiceImpl struct {
	collectorDao dao.ICollectorDAO
}

// NewICollectorServiceImpl 创建服务
func NewICollectorServiceImpl(collectorDao dao.ICollectorDAO) service.ICollectorService {
	return &ICollectorServiceImpl{collectorDao: collectorDao}
}

// CreateCollector 创建采集器
func (s *ICollectorServiceImpl) CreateCollector(c *gin.Context, req *dto.CreateCollectorRequest) (*entity.Collector, error) {
	// 设备 ID 唯一性检查
	if req.DeviceID != "" {
		existing, err := s.collectorDao.GetByDeviceID(c, req.DeviceID)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, fmt.Errorf("设备ID '%s' 已存在", req.DeviceID)
		}
	}

	col := &entity.Collector{
		ID:       uuid.New().String(),
		Name:     req.Name,
		DeviceID: req.DeviceID,
		Token:    req.Token,
		Status:   entity.StatusOffline,
		Tags:     marshalTags(req.Tags),
	}

	if err := s.collectorDao.Create(c, col); err != nil {
		return nil, err
	}

	return col, nil
}

// GetCollector 获取采集器
func (s *ICollectorServiceImpl) GetCollector(c *gin.Context, id string) (*entity.Collector, error) {
	col, err := s.collectorDao.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	if col == nil {
		return nil, fmt.Errorf("采集器不存在")
	}
	return col, nil
}

// ListCollectors 列表
func (s *ICollectorServiceImpl) ListCollectors(c *gin.Context, offset, limit int, name, deviceID, status string) ([]entity.Collector, int64, error) {
	return s.collectorDao.List(c, offset, limit, name, deviceID, status)
}

// ListAllOnline 获取全部在线采集器（用于下发弹窗）
func (s *ICollectorServiceImpl) ListAllOnline(c *gin.Context) ([]entity.Collector, error) {
	all, err := s.collectorDao.ListAll(c)
	if err != nil {
		return nil, err
	}
	var online []entity.Collector
	for _, col := range all {
		if col.Status == entity.StatusOnline {
			online = append(online, col)
		}
	}
	return online, nil
}

// UpdateCollector 更新采集器
func (s *ICollectorServiceImpl) UpdateCollector(c *gin.Context, id string, req *dto.UpdateCollectorRequest) (*entity.Collector, error) {
	col, err := s.collectorDao.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	if col == nil {
		return nil, fmt.Errorf("采集器不存在")
	}

	if req.Name != "" {
		col.Name = req.Name
	}
	if req.DeviceID != "" && req.DeviceID != col.DeviceID {
		existing, err := s.collectorDao.GetByDeviceID(c, req.DeviceID)
		if err == nil && existing != nil && existing.ID != id {
			return nil, fmt.Errorf("设备ID '%s' 已存在", req.DeviceID)
		}
		col.DeviceID = req.DeviceID
	}
	// Token 允许清空
	col.Token = req.Token
	if req.Tags != nil {
		col.Tags = marshalTags(req.Tags)
	}

	if err := s.collectorDao.Update(c, col); err != nil {
		return nil, err
	}

	return col, nil
}

// DeleteCollector 删除
func (s *ICollectorServiceImpl) DeleteCollector(c *gin.Context, id string) error {
	col, err := s.collectorDao.GetByID(c, id)
	if err != nil {
		return err
	}
	if col == nil {
		return fmt.Errorf("采集器不存在")
	}
	return s.collectorDao.Delete(c, id)
}

func marshalTags(tags []string) string {
	if len(tags) == 0 {
		return "[]"
	}
	b, _ := json.Marshal(tags)
	return cast.ToString(b)
}

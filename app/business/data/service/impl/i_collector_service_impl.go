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

// Create 创建采集器
func (s *ICollectorServiceImpl) Create(c *gin.Context, req *dto.CreateCollectorRequest) (*dto.CollectorResponse, error) {
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

	return collectorResponse(col), nil
}

// Get 获取采集器
func (s *ICollectorServiceImpl) Get(c *gin.Context, id string) (*dto.CollectorResponse, error) {
	col, err := s.collectorDao.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	if col == nil {
		return nil, fmt.Errorf("采集器不存在")
	}
	return collectorResponse(col), nil
}

// List 列表
func (s *ICollectorServiceImpl) List(c *gin.Context, offset, limit int, name, deviceID, status string) ([]dto.CollectorResponse, int64, error) {
	rows, total, err := s.collectorDao.List(c, offset, limit, name, deviceID, status)
	if err != nil {
		return nil, 0, err
	}
	out := make([]dto.CollectorResponse, 0, len(rows))
	for i := range rows {
		out = append(out, *collectorResponse(&rows[i]))
	}
	return out, total, nil
}

// ListOnline 获取全部在线采集器（用于下发弹窗）
func (s *ICollectorServiceImpl) ListOnline(c *gin.Context) ([]dto.CollectorResponse, error) {
	all, err := s.collectorDao.ListAll(c)
	if err != nil {
		return nil, err
	}
	online := make([]dto.CollectorResponse, 0)
	for i := range all {
		if all[i].Status == entity.StatusOnline {
			online = append(online, *collectorResponse(&all[i]))
		}
	}
	return online, nil
}

// Update 更新采集器
func (s *ICollectorServiceImpl) Update(c *gin.Context, id string, req *dto.UpdateCollectorRequest) (*dto.CollectorResponse, error) {
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

	return collectorResponse(col), nil
}

// Delete 删除
func (s *ICollectorServiceImpl) Delete(c *gin.Context, id string) error {
	col, err := s.collectorDao.GetByID(c, id)
	if err != nil {
		return err
	}
	if col == nil {
		return fmt.Errorf("采集器不存在")
	}
	return s.collectorDao.Delete(c, id)
}

func collectorResponse(col *entity.Collector) *dto.CollectorResponse {
	return &dto.CollectorResponse{
		ID:              col.ID,
		Name:            col.Name,
		DeviceID:        col.DeviceID,
		Token:           col.Token,
		Status:          col.Status,
		LastHeartbeat:   col.LastHeartbeat,
		LastConnectedAt: col.LastConnectedAt,
		Version:         col.Version,
		Tags:            unmarshalTags(col.Tags),
		CreatedAt:       col.CreatedAt,
		UpdatedAt:       col.UpdatedAt,
	}
}

func unmarshalTags(tags string) []string {
	if tags == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(tags), &out); err != nil {
		return nil
	}
	return out
}

func marshalTags(tags []string) string {
	if len(tags) == 0 {
		return "[]"
	}
	b, _ := json.Marshal(tags)
	return cast.ToString(b)
}

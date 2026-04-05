package gatewaydaoimpl

import (
	"strconv"
	"strings"
	"time"

	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MCPServerDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewMCPServerDao(db *gorm.DB) gatewaydao.IMCPServerDao {
	return &MCPServerDaoImpl{
		db:    db,
		table: "mcp_servers",
	}
}

func (m *MCPServerDaoImpl) Create(c *gin.Context, req *gatewaymodels.MCPServerUpsert) (*gatewaymodels.MCPServer, error) {
	item := &gatewaymodels.MCPServer{
		ID:          strconv.FormatInt(snowflake.GenID(), 10),
		Name:        req.Name,
		Description: req.Description,
		Transport:   req.Transport,
		Command:     req.Command,
		Args:        req.Args,
		Env:         req.Env,
		URL:         req.URL,
		Headers:     req.Headers,
		Timeout:     req.Timeout,
		IsCommon:    req.IsCommon,
		Enabled:     req.Enabled,
		DeptID:      baizeContext.GetDeptId(c),
		State:       commonStatus.NORMAL,
	}
	item.SetCreateBy(baizeContext.GetUserId(c))
	if err := m.db.WithContext(c).Table(m.table).Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (m *MCPServerDaoImpl) Update(c *gin.Context, req *gatewaymodels.MCPServerUpsert) (*gatewaymodels.MCPServer, error) {
	item := &gatewaymodels.MCPServer{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Transport:   req.Transport,
		Command:     req.Command,
		Args:        req.Args,
		Env:         req.Env,
		URL:         req.URL,
		Headers:     req.Headers,
		Timeout:     req.Timeout,
		IsCommon:    req.IsCommon,
		Enabled:     req.Enabled,
	}
	item.SetUpdateBy(baizeContext.GetUserId(c))
	if err := m.db.WithContext(c).Table(m.table).
		Where("id = ?", item.ID).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		Select("name", "description", "transport", "command", "args", "env", "url", "headers", "timeout", "is_common", "enabled", "update_by", "update_time").
		Updates(item).Error; err != nil {
		return nil, err
	}
	return m.getByID(c, item.ID)
}

func (m *MCPServerDaoImpl) DeleteByIDs(c *gin.Context, ids []string) error {
	now := time.Now()
	return m.db.WithContext(c).Table(m.table).
		Where("id IN ?", ids).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Updates(map[string]interface{}{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": now,
		}).Error
}

func (m *MCPServerDaoImpl) List(c *gin.Context, req *gatewaymodels.MCPServerQuery) (*gatewaymodels.MCPServerListData, error) {
	db := m.db.WithContext(c).Table(m.table).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL)
	if name := strings.TrimSpace(req.Name); name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if transport := strings.TrimSpace(req.Transport); transport != "" {
		db = db.Where("transport = ?", transport)
	}
	if req.Enabled != nil {
		db = db.Where("enabled = ?", req.Enabled)
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*gatewaymodels.MCPServer, 0)
	if err := db.Order("create_time DESC, id DESC").
		Offset(int((req.Page - 1) * req.Size)).
		Limit(int(req.Size)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return &gatewaymodels.MCPServerListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func (m *MCPServerDaoImpl) getByID(c *gin.Context, id string) (*gatewaymodels.MCPServer, error) {
	var item gatewaymodels.MCPServer
	if err := m.db.WithContext(c).Table(m.table).
		Where("id = ?", id).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

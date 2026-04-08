package aiDataSetDaoImpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IDatasetRolePermissionDaoImpl struct {
	db    *gorm.DB
	table string
}

func (i *IDatasetRolePermissionDaoImpl) translateDatasetIDsToUUIDs(c *gin.Context, datasetIDs []string) ([]string, error) {
	if len(datasetIDs) == 0 {
		return nil, errors.New("datasetId不能为空")
	}
	idList := make([]int64, 0, len(datasetIDs))
	for _, idStr := range datasetIDs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, errors.New("datasetId格式不合法")
		}
		idList = append(idList, id)
	}
	if len(idList) == 0 {
		return nil, errors.New("datasetId不能为空")
	}

	rows := make([]*aidatasetmodels.SysDataset, 0)
	if err := i.db.WithContext(c).Table("sys_dataset").
		Select("dataset_id", "dataset_uuid").
		Where("dataset_id IN ?", idList).
		Where("state = ?", commonStatus.NORMAL).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("datasetId不存在")
	}
	m := make(map[int64]string, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		m[row.DatasetID] = strings.TrimSpace(row.DatasetUUID)
	}

	uuidList := make([]string, 0, len(idList))
	missing := make([]string, 0)
	for _, id := range idList {
		uuid := strings.TrimSpace(m[id])
		if uuid == "" {
			missing = append(missing, strconv.FormatInt(id, 10))
			continue
		}
		uuidList = append(uuidList, uuid)
	}
	if len(missing) > 0 {
		return nil, errors.New(fmt.Sprintf("datasetId不存在: %s", strings.Join(missing, ",")))
	}
	return uuidList, nil
}

func (i *IDatasetRolePermissionDaoImpl) translateDocumentIDsToUUIDs(c *gin.Context, documentIDs []string) ([]string, error) {
	if len(documentIDs) == 0 {
		return make([]string, 0), nil
	}
	idList := make([]int64, 0, len(documentIDs))
	for _, idStr := range documentIDs {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, errors.New("documentId格式不合法")
		}
		idList = append(idList, id)
	}
	if len(idList) == 0 {
		return make([]string, 0), nil
	}

	rows := make([]*aidatasetmodels.SysDatasetDocument, 0)
	if err := i.db.WithContext(c).Table("sys_dataset_document").
		Select("document_id", "dataset_document_uuid").
		Where("document_id IN ?", idList).
		Where("state = ?", commonStatus.NORMAL).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("documentId不存在")
	}
	m := make(map[int64]string, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		m[row.DocumentID] = strings.TrimSpace(row.DatasetDocumentUUID)
	}

	uuidList := make([]string, 0, len(idList))
	missing := make([]string, 0)
	for _, id := range idList {
		uuid := strings.TrimSpace(m[id])
		if uuid == "" {
			missing = append(missing, strconv.FormatInt(id, 10))
			continue
		}
		uuidList = append(uuidList, uuid)
	}
	if len(missing) > 0 {
		return nil, errors.New(fmt.Sprintf("documentId不存在: %s", strings.Join(missing, ",")))
	}
	return uuidList, nil
}

// NewIDatasetRolePermissionDaoImpl 知识库/文档-角色权限DAO构造函数。
func NewIDatasetRolePermissionDaoImpl(db *gorm.DB) aidatasetdao.IDatasetRolePermissionDao {
	return &IDatasetRolePermissionDaoImpl{
		db:    db,
		table: "sys_dataset_role_permission",
	}
}

// Create 创建知识库/文档-角色权限记录。
func (i *IDatasetRolePermissionDaoImpl) Create(c *gin.Context, req *aidatasetmodels.SetDatasetRolePermission) (*aidatasetmodels.DatasetRolePermission, error) {
	datasetUUIDs, err := i.translateDatasetIDsToUUIDs(c, req.DatasetID)
	if err != nil {
		return nil, err
	}
	datasetUuIdStr, err := json.Marshal(datasetUUIDs)
	if err != nil {
		zap.L().Error("json marshal datasetUUIDs failed", zap.Error(err))
		return nil, err
	}

	documentUUIDs, err := i.translateDocumentIDsToUUIDs(c, req.DocumentID)
	if err != nil {
		return nil, err
	}
	documentUuIDStr, err := json.Marshal(documentUUIDs)
	if err != nil {
		zap.L().Error("json marshal documentUUIDs failed", zap.Error(err))
		return nil, err
	}

	datasetIDStr, err := json.Marshal(req.DatasetID)
	if err != nil {
		zap.L().Error("json marshal datasetId failed", zap.Error(err))
		return nil, err
	}

	documentIDStr, err := json.Marshal(req.DocumentID)
	if err != nil {
		zap.L().Error("json marshal document failed", zap.Error(err))
		return nil, err
	}

	data := &aidatasetmodels.DatasetRolePermission{
		ID:            snowflake.GenID(),
		RoleID:        req.RoleID,
		DatasetIDs:    string(datasetIDStr),
		DatasetUuIDs:  string(datasetUuIdStr),
		DocumentIDs:   string(documentIDStr),
		DocumentUuIDs: string(documentUuIDStr),
		Permission:    req.Permission,
		DeptID:        baizeContext.GetDeptId(c),
		State:         commonStatus.NORMAL,
	}
	data.SetCreateBy(baizeContext.GetUserId(c))
	if err := i.db.WithContext(c).Table(i.table).Create(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Update 更新知识库/文档-角色权限记录。
func (i *IDatasetRolePermissionDaoImpl) Update(c *gin.Context, req *aidatasetmodels.SetDatasetRolePermission) (*aidatasetmodels.DatasetRolePermission, error) {
	if req.ID == 0 {
		return nil, errors.New("id不能为空")
	}
	data := &aidatasetmodels.DatasetRolePermission{}
	if err := i.db.WithContext(c).Table(i.table).
		Where("id = ?", req.ID).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		First(data).Error; err != nil {
		return nil, err
	}
	data.RoleID = req.RoleID

	datasetUUIDs, err := i.translateDatasetIDsToUUIDs(c, req.DatasetID)
	if err != nil {
		return nil, err
	}
	datasetUuIdStr, err := json.Marshal(datasetUUIDs)
	if err != nil {
		zap.L().Error("json marshal datasetUUIDs failed", zap.Error(err))
		return nil, err
	}

	documentUUIDs, err := i.translateDocumentIDsToUUIDs(c, req.DocumentID)
	if err != nil {
		return nil, err
	}
	documentUuIDStr, err := json.Marshal(documentUUIDs)
	if err != nil {
		zap.L().Error("json marshal documentUUIDs failed", zap.Error(err))
		return nil, err
	}

	datasetIdStr, err := json.Marshal(req.DatasetID)
	if err != nil {
		zap.L().Error("json marshal datasetId failed", zap.Error(err))
		return nil, err
	}

	documentIDStr, err := json.Marshal(req.DocumentID)
	if err != nil {
		zap.L().Error("json marshal DocumentID failed", zap.Error(err))
		return nil, err
	}

	data.DatasetIDs = string(datasetIdStr)
	data.DocumentIDs = string(documentIDStr)
	data.DocumentUuIDs = string(documentUuIDStr)
	data.DatasetUuIDs = string(datasetUuIdStr)
	data.Permission = req.Permission
	data.SetUpdateBy(baizeContext.GetUserId(c))
	if err := i.db.WithContext(c).Table(i.table).
		Where("id = ?", req.ID).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Updates(map[string]interface{}{
			"role_id":      data.RoleID,
			"dataset_ids":  data.DatasetIDs,
			"document_ids": data.DocumentIDs,
			"permission":   data.Permission,
			"status":       req.Status,
			"update_by":    data.UpdateBy,
			"update_time":  data.UpdateTime,
		}).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// List 查询知识库/文档-角色权限列表。
func (i *IDatasetRolePermissionDaoImpl) List(c *gin.Context, req *aidatasetmodels.DatasetRolePermissionQuery) (*aidatasetmodels.DatasetRolePermissionListData, error) {
	db := i.db.WithContext(c).Table(i.table)
	if req.RoleID > 0 {
		db = db.Where("role_id = ?", req.RoleID)
	}
	if req.DatasetID > 0 {
		db = db.Where("dataset_id = ?", req.DatasetID)
	}
	if req.DocumentID > 0 {
		db = db.Where("document_id = ?", req.DocumentID)
	}
	if strings.TrimSpace(req.Permission) != "" {
		db = db.Where("permission = ?", strings.TrimSpace(req.Permission))
	}
	db = db.Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL)
	page := req.GetPage()
	size := req.GetSize()
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*aidatasetmodels.DatasetRolePermission, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}

	for k, row := range rows {
		documentIDArray := make([]string, 0)
		datasetIDArray := make([]string, 0)

		err := json.Unmarshal([]byte(row.DocumentIDs), &documentIDArray)
		if err != nil {
			zap.L().Error("json unmarshal DocumentIDs failed", zap.Error(err))
		}

		err = json.Unmarshal([]byte(row.DatasetIDs), &datasetIDArray)
		if err != nil {
			zap.L().Error("json unmarshal DocumentIDs failed", zap.Error(err))
		}

		rows[k].DatasetIDArray = datasetIDArray
		rows[k].DocumentIDsArray = documentIDArray
	}
	return &aidatasetmodels.DatasetRolePermissionListData{
		Rows:  rows,
		Total: total,
	}, nil
}

// Remove 删除知识库/文档-角色权限记录（软删除）。
func (i *IDatasetRolePermissionDaoImpl) Remove(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return i.db.WithContext(c).Table(i.table).
		Where("id IN ?", ids).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Updates(map[string]interface{}{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

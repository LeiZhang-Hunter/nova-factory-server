package daoimpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"nova-factory-server/app/constant/shop"
	"time"

	"nova-factory-server/app/business/shop/finance/dao"
	"nova-factory-server/app/business/shop/finance/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountDaoImpl struct {
	db    *gorm.DB
	table string
	cache cache.Cache
}

func NewAccountDao(db *gorm.DB, cache cache.Cache) dao.IAccountDao {
	return &AccountDaoImpl{
		db:    db,
		table: "shop_account",
		cache: cache,
	}
}

func (d *AccountDaoImpl) Create(c *gin.Context, req *models.AccountUpsert) (*models.Account, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	model := models.AccountUpsertToEntity(req)
	if model == nil {
		return nil, errors.New("参数不能为空")
	}
	model.ID = snowflake.GenID()
	model.DeptID = baizeContext.GetDeptId(c)
	model.State = commonStatus.NORMAL
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := d.db.WithContext(c).Table(d.table).Create(model).Error; err != nil {
		return nil, err
	}
	d.cacheDefaultIfNeeded(c, model)
	return model, nil
}

func (d *AccountDaoImpl) Update(c *gin.Context, req *models.AccountUpsert) (*models.Account, error) {
	if req == nil || req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	updates := make(map[string]any)
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.No != "" {
		updates["no"] = req.No
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	if req.DefaultStatus != nil {
		updates["default_status"] = req.DefaultStatus
	}
	updates["status"] = req.Status
	updates["sort"] = req.Sort
	updates["update_by"] = baizeContext.GetUserId(c)
	updates["update_time"] = time.Now()
	db := d.db.WithContext(c).Table(d.table).Where("id = ?", req.ID)
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	model, err := d.GetByID(c, req.ID)
	if err != nil {
		return nil, err
	}
	d.cacheDefaultIfNeeded(c, model)
	return model, nil
}

func (d *AccountDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table(d.table).
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": time.Now(),
		}).Error
}

func (d *AccountDaoImpl) GetByID(c *gin.Context, id int64) (*models.Account, error) {
	item := new(models.Account)
	if err := d.db.WithContext(c).Table(d.table).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *AccountDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*models.Account, error) {
	if column == "" {
		return nil, nil
	}
	item := new(models.Account)
	if err := d.db.WithContext(c).Table(d.table).
		Where(fmt.Sprintf("%s = ?", column), value).
		Where("state = ?", commonStatus.NORMAL).
		First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *AccountDaoImpl) ListPage(c *gin.Context, req *models.AccountQuery) (*models.AccountListData, error) {
	if req == nil {
		req = new(models.AccountQuery)
	}
	db := d.db.WithContext(c).Table(d.table).Where("state = ?", commonStatus.NORMAL)
	db = applyAccountFilters(db, req)
	page, size := getPageSize(req.Page, req.Size)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]models.Account, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*models.Account, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &models.AccountListData{Rows: result, Total: total}, nil
}

func (d *AccountDaoImpl) ResetDefaultStatus(c *gin.Context, excludeID int64) error {
	err := d.db.WithContext(c).Table(d.table).
		Where("state = ?", commonStatus.NORMAL).
		Where("id != ?", excludeID).
		Where("default_status = ?", 1).
		Updates(map[string]any{
			"default_status": 0,
			"update_time":    time.Now(),
			"update_by":      baizeContext.GetUserId(c),
		}).Error
	if err != nil {
		return err
	}
	if excludeID > 0 {
		model, err := d.GetByID(c, excludeID)
		if err != nil {
			return err
		}
		d.cacheDefaultIfNeeded(c, model)
	} else {
		d.cache.Del(c.Request.Context(), shop.ShopAccountDefaultCacheKey)
	}
	return nil
}

// cacheDefaultIfNeeded 如果账户是默认账户则写入缓存，否则清除缓存中的旧默认账户。
func (d *AccountDaoImpl) cacheDefaultIfNeeded(c *gin.Context, model *models.Account) {
	if model == nil {
		return
	}
	if model.DefaultStatus != nil && *model.DefaultStatus {
		data, err := json.Marshal(model)
		if err != nil {
			return
		}
		d.cache.Set(c.Request.Context(), shop.ShopAccountDefaultCacheKey, string(data), 0)
	} else {
		d.cache.Del(c.Request.Context(), shop.ShopAccountDefaultCacheKey)
	}
}

// GetDefaultFromCache 从缓存获取默认结算账户。
func (d *AccountDaoImpl) GetDefaultFromCache(c *gin.Context) *models.Account {
	val, err := d.cache.Get(c.Request.Context(), shop.ShopAccountDefaultCacheKey)
	if err != nil || val == "" {
		return nil
	}
	item := new(models.Account)
	if err := json.Unmarshal([]byte(val), item); err != nil {
		return nil
	}
	return item
}

func (d *AccountDaoImpl) List(c *gin.Context, req *models.AccountQuery) (*models.AccountListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &models.AccountListData{Rows: result.Rows, Total: result.Total}, nil
}

// GetDefaultAccountNo 从缓存读取默认结算账户编码，缓存未命中返回空字符串。
func (d *AccountDaoImpl) GetDefaultAccountNo(c *gin.Context) string {
	val, err := d.cache.Get(c, shop.ShopAccountDefaultCacheKey)
	if err != nil || val == "" {
		return ""
	}
	// 只解析 No 字段，避免引入 finance models 包
	var account struct {
		No string `json:"no"`
	}
	if err := json.Unmarshal([]byte(val), &account); err != nil {
		return ""
	}
	return account.No
}

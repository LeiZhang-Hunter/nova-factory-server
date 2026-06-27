package impl

import (
	"errors"
	"nova-factory-server/app/utils/snowflake"
	"strings"
	"time"

	"nova-factory-server/app/business/shop/order/dao"
	"nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// OrderDaoImpl shop订单数据访问实现。
type OrderDaoImpl struct {
	db         *gorm.DB
	table      string
	detailDao  *OrderDetailDaoImpl
	accountDao *OrderAccountDaoImpl
}

// NewOrderDao 创建 shop 订单 DAO。
func NewOrderDao(db *gorm.DB) dao.IOrderDao {
	return &OrderDaoImpl{
		db:         db,
		table:      "shop_order",
		detailDao:  newOrderDetailDaoImpl(db),
		accountDao: newOrderAccountDaoImpl(db),
	}
}

// Set 新增或修改 shop 订单。
func (o *OrderDaoImpl) Set(c *gin.Context, req *models.OrderSet) (*models.Order, error) {
	var resultID uint64
	err := o.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		id, setErr := o.doSet(c, tx, req)
		if setErr != nil {
			return setErr
		}
		resultID = id
		return nil
	})
	if err != nil {
		return nil, err
	}
	return o.GetByID(c, resultID)
}

// SetWithTx 新增或修改 shop 订单（带事务）。
func (o *OrderDaoImpl) SetWithTx(c *gin.Context, tx *gorm.DB, req *models.OrderSet) (*models.Order, error) {
	resultID, err := o.doSet(c, tx, req)
	if err != nil {
		return nil, err
	}
	return o.getByIDWithTx(c, tx, resultID)
}

// Create 在事务内创建 shop 订单主表记录。
func (o *OrderDaoImpl) Create(tx *gorm.DB, order *models.Order) error {
	if tx == nil {
		return errors.New("shop订单主表创建失败：事务不能为空")
	}
	if order == nil {
		return errors.New("shop订单主表创建失败：订单不能为空")
	}
	if order.ID == 0 {
		order.ID = uint64(snowflake.GenID())
	}
	now := time.Now()
	order.CreateTime = &now
	order.UpdateTime = &now
	return tx.Table(o.table).Create(order).Error
}

// UpdateByID 在事务内按 ID 更新 shop 订单主表记录。
func (o *OrderDaoImpl) UpdateByID(tx *gorm.DB, id uint64, updates map[string]any) error {
	if tx == nil {
		return errors.New("shop订单主表更新失败：事务不能为空")
	}
	if id == 0 {
		return errors.New("shop订单主表更新失败：订单ID不能为空")
	}
	if len(updates) == 0 {
		return nil
	}
	return tx.Table(o.table).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		Updates(updates).Error
}

// doSet 执行新增或修改 shop 订单的核心逻辑。
func (o *OrderDaoImpl) doSet(c *gin.Context, tx *gorm.DB, req *models.OrderSet) (uint64, error) {
	var resultID uint64
	exists, findErr := o.findExists(tx, c, req)
	if findErr != nil {
		return 0, findErr
	}
	data, parseErr := models.BuildOrderModel(req)
	if parseErr != nil {
		return 0, parseErr
	}
	if exists == nil {
		data.ID = uint64(snowflake.GenID())
		data.SetCreateBy(baizeContext.GetUserId(c))
		if err := tx.Table(o.table).Create(data).Error; err != nil {
			return 0, err
		}
		resultID = data.ID
	} else {
		data.ID = exists.ID
		data.SetUpdateBy(baizeContext.GetUserId(c))
		if err := tx.Table(o.table).
			Where("id = ?", exists.ID).
			Where("state = ?", commonStatus.NORMAL).
			Updates(buildOrderUpdateMap(data)).Error; err != nil {
			return 0, err
		}
		resultID = exists.ID
		if err := o.detailDao.DeleteByOrderID(tx, resultID); err != nil {
			return 0, err
		}
		if err := o.accountDao.DeleteByOrderID(tx, resultID); err != nil {
			return 0, err
		}
	}
	if err := o.detailDao.DeleteByTidAndOIDs(tx, data.Tid, req.Details); err != nil {
		return 0, err
	}
	if err := o.detailDao.BatchCreate(tx, c, resultID, data.Tid, req.Details); err != nil {
		return 0, err
	}
	if err := o.accountDao.BatchCreate(tx, c, resultID, data.Tid, req.Accounts); err != nil {
		return 0, err
	}
	return resultID, nil
}

// GetByID 查询 shop 订单详情。
func (o *OrderDaoImpl) GetByID(c *gin.Context, id uint64) (*models.Order, error) {
	var item models.Order
	if err := o.db.WithContext(c).Table(o.table).
		Where("id = ?", id).
		//Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}
func (o *OrderDaoImpl) GetByTid(c *gin.Context, tid string) (*models.Order, error) {
	var item models.Order
	if err := o.db.WithContext(c).Table(o.table).
		Where("tid = ?", strings.TrimSpace(tid)).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}

// getByIDWithTx 在事务内读取订单主表，避免未提交数据在事务外不可见。
func (o *OrderDaoImpl) getByIDWithTx(c *gin.Context, tx *gorm.DB, id uint64) (*models.Order, error) {
	if tx == nil {
		return o.GetByID(c, id)
	}
	var item models.Order
	if err := tx.WithContext(c).Table(o.table).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// List 分页查询 shop 订单。
func (o *OrderDaoImpl) List(c *gin.Context, req *models.OrderQuery) (*models.OrderListData, error) {
	db := o.db.WithContext(c).Table(o.table).
		//Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL)
	if req != nil {
		if strings.TrimSpace(req.Tid) != "" {
			db = db.Where("tid LIKE ?", "%"+strings.TrimSpace(req.Tid)+"%")
		}
		if strings.TrimSpace(req.Status) != "" {
			db = db.Where("status = ?", strings.TrimSpace(req.Status))
		}
		if req.SyncStatus != nil {
			db = db.Where("sync_status = ?", *req.SyncStatus)
		}
		if strings.TrimSpace(req.ReceiverName) != "" {
			db = db.Where("receivername LIKE ?", "%"+strings.TrimSpace(req.ReceiverName)+"%")
		}
	}
	page := int64(1)
	size := int64(20)
	if req != nil && req.Page > 0 {
		page = req.Page
	}
	if req != nil && req.Size > 0 {
		size = req.Size
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*models.Order, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}

	return &models.OrderListData{
		Rows:  rows,
		Total: total,
	}, nil
}

// DeleteByIDs 软删除 shop 订单及其子表。
func (o *OrderDaoImpl) DeleteByIDs(c *gin.Context, ids []uint64) error {
	now := time.Now()
	return o.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(o.table).
			Where("id IN ?", ids).
			//Where("dept_id = ?", baizeContext.GetDeptId(c)).
			Where("state = ?", commonStatus.NORMAL).
			Updates(map[string]interface{}{
				"state":       commonStatus.DELETE,
				"update_by":   baizeContext.GetUserId(c),
				"update_time": now,
			}).Error; err != nil {
			return err
		}
		orderIDs := make([]uint64, 0, len(ids))
		orderIDs = append(orderIDs, ids...)
		if err := o.detailDao.DeleteByOrderIDs(tx, orderIDs); err != nil {
			return err
		}
		return o.accountDao.DeleteByOrderIDs(tx, orderIDs)
	})
}

// findExists 根据订单ID或订单编号查询当前部门下的有效订单。
func (o *OrderDaoImpl) findExists(tx *gorm.DB, c *gin.Context, req *models.OrderSet) (*models.Order, error) {
	var item models.Order
	db := tx.Table(o.table)
	if req.ID > 0 {
		db = db.Where("id = ?", req.ID)
	} else {
		db = db.Where("tid = ?", req.Tid)
	}
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// buildOrderUpdateMap 构建 shop 订单更新字段映射。
func buildOrderUpdateMap(order *models.Order) map[string]interface{} {
	if order == nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"tid":                    order.Tid,
		"weight":                 order.Weight,
		"size":                   order.Size,
		"buyer_nick":             order.BuyerNick,
		"buyer_message":          order.BuyerMessage,
		"seller_memo":            order.SellerMemo,
		"total":                  order.Total,
		"privilege":              order.Privilege,
		"post_fee":               order.PostFee,
		"receiver_name":          order.ReceiverName,
		"receiver_province":      order.ReceiverProvince,
		"receiver_province_name": order.ReceiverProvinceName,
		"receiver_city":          order.ReceiverCity,
		"receiver_city_name":     order.ReceiverCityName,
		"receiver_district":      order.ReceiverDistrict,
		"receiver_district_name": order.ReceiverDistrictName,
		"receiver_street":        order.ReceiverStreet,
		"receiver_street_name":   order.ReceiverStreetName,
		"receiver_address":       order.ReceiverAddress,
		"receiver_phone":         order.ReceiverPhone,
		"receiver_mobile":        order.ReceiverMobile,
		"receiver_zip":           order.ReceiverZip,
		"status":                 order.Status,
		"order_type":             order.Type,
		"invoice_name":           order.InvoiceName,
		"seller_flag":            order.SellerFlag,
		"pay_time":               order.PayTime,
		"logist_b_type_code":     order.LogistBTypeCode,
		"logist_bill_code":       order.LogistBillCode,
		"b_type_code":            order.BTypeCode,
		"details_json":           order.DetailsJSON,
		"accounts_json":          order.AccountsJSON,
		"bill_code":              order.BillCode,
		"sync_message":           order.SyncMessage,
		"sync_status":            order.SyncStatus,
		"sync_time":              order.SyncTime,
		"update_by":              order.UpdateBy,
		"update_time":            order.UpdateTime,
	}
}

func (o *OrderDaoImpl) GetByTidTx(tx *gorm.DB, tid string) (*models.Order, error) {
	var item models.Order
	if err := tx.Table(o.table).
		Where("tid = ?", tid).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// GetByTidForUpdateTx 在事务内按订单编号查询并加行锁。
func (o *OrderDaoImpl) GetByTidForUpdateTx(tx *gorm.DB, tid string) (*models.Order, error) {
	var item models.Order
	if err := tx.Table(o.table).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("tid = ?", tid).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// Transaction 开启订单同步事务。
//
// service 层会在该事务中组合调用订单主表 DAO、订单明细 DAO、订单账户 DAO。
// 只要 fn 返回 error，GORM 会回滚整个事务；fn 返回 nil 时才提交。
func (o *OrderDaoImpl) Transaction(fn func(tx *gorm.DB) error) error {
	if o.db == nil {
		return errors.New("shop order dao db is nil")
	}
	return o.db.Transaction(fn)
}

package impl

import (
	"errors"
	"strings"
	"time"

	"nova-factory-server/app/business/shop/order/dao"
	"nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/constant/commonStatus"
	orderConstant "nova-factory-server/app/constant/order"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderRefundDaoImpl 售后单数据访问实现。
type OrderRefundDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewOrderRefundDao 创建售后单 DAO。
func NewOrderRefundDao(db *gorm.DB) dao.IOrderRefundDao {
	return &OrderRefundDaoImpl{
		db:        db,
		tableName: "shop_order_refund",
	}
}

func (d *OrderRefundDaoImpl) Create(c *gin.Context, aftersale *models.OrderRefund) error {
	if aftersale == nil {
		return errors.New("售后单不能为空")
	}
	aftersale.SetCreateBy(baizeContext.GetUserId(c))
	return d.CreateWithTx(d.db.WithContext(c), aftersale)
}

func (d *OrderRefundDaoImpl) CreateWithTx(tx *gorm.DB, aftersale *models.OrderRefund) error {
	if tx == nil {
		return errors.New("事务不能为空")
	}
	if aftersale == nil {
		return errors.New("售后单不能为空")
	}
	if aftersale.ID == 0 {
		aftersale.ID = snowflake.GenID()
	}
	now := time.Now()
	aftersale.CreateTime = &now
	aftersale.UpdateTime = &now
	aftersale.State = commonStatus.NORMAL
	return tx.Table(d.tableName).Create(aftersale).Error
}

func (d *OrderRefundDaoImpl) GetByID(c *gin.Context, id int64) (*models.OrderRefund, error) {
	var result models.OrderRefund
	if err := d.db.WithContext(c).Table(d.tableName).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (d *OrderRefundDaoImpl) GetByOutRefundNo(c *gin.Context, outRefundNo string) (*models.OrderRefund, error) {
	var result models.OrderRefund
	if err := d.db.WithContext(c).Table(d.tableName).
		Where("out_refund_no = ?", strings.TrimSpace(outRefundNo)).
		//Where("state = ?", commonStatus.NORMAL).
		First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (d *OrderRefundDaoImpl) GetByOrderId(c *gin.Context, orderId int64) (*models.OrderRefund, error) {
	var result models.OrderRefund
	if err := d.db.WithContext(c).Table(d.tableName).
		Where("order_id = ?", orderId).
		Where("state = ?", commonStatus.NORMAL).
		Where("status NOT IN ?", []int32{
			orderConstant.AftersaleStatusRefundSuccess,
			orderConstant.AftersaleStatusRefundFailed,
			orderConstant.AftersaleStatusRefundClosed,
			orderConstant.AftersaleStatusRejected,
		}).
		First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (d *OrderRefundDaoImpl) UpdateStatus(c *gin.Context, id int64, status int32, updates map[string]any) error {
	return d.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		return d.UpdateStatusWithTx(tx, id, status, updates)
	})
}
func (d *OrderRefundDaoImpl) UpdateByID(c *gin.Context, id int64, updates map[string]any) error {
	return d.db.WithContext(c).Table(d.tableName).
		Where("id = ?", id).
		Updates(updates).Error
}
func (d *OrderRefundDaoImpl) UpdateStatusWithTx(tx *gorm.DB, id int64, status int32, updates map[string]any) error {
	if tx == nil {
		return errors.New("事务不能为空")
	}
	if id == 0 {
		return errors.New("售后单ID不能为空")
	}
	if updates == nil {
		updates = make(map[string]any)
	}
	updates["status"] = status
	updates["update_time"] = time.Now()
	return tx.Table(d.tableName).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		Updates(updates).Error
}

func (d *OrderRefundDaoImpl) List(c *gin.Context, req *models.RefundQuery) (*models.RefundListData, error) {
	db := d.db.WithContext(c).Table(d.tableName).
		Where("state = ?", commonStatus.NORMAL)

	if req != nil {
		if strings.TrimSpace(req.Tid) != "" {
			db = db.Where("tid LIKE ?", "%"+strings.TrimSpace(req.Tid)+"%")
		}
		if req.Status != nil {
			db = db.Where("status = ?", *req.Status)
		}
		if req.UserID > 0 {
			db = db.Where("user_id = ?", req.UserID)
		}
		if req.SyncStatus != nil {
			db = db.Where("sync_status = ?", *req.SyncStatus)
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

	var rows []*models.OrderRefund
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}

	return &models.RefundListData{
		Rows:  rows,
		Total: total,
	}, nil
}

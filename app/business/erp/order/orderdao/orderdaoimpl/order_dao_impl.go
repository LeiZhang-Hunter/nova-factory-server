package orderdaoimpl

import (
	"errors"
	"strings"
	"time"

	"nova-factory-server/app/business/erp/order/orderdao"
	"nova-factory-server/app/business/erp/order/ordermodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderDaoImpl ERP订单数据访问实现。
type OrderDaoImpl struct {
	db           *gorm.DB
	table        string
	detailTable  string
	accountTable string
}

// NewOrderDao 创建 ERP 订单 DAO。
func NewOrderDao(db *gorm.DB) orderdao.IOrderDao {
	return &OrderDaoImpl{
		db:           db,
		table:        "erp_order",
		detailTable:  "erp_order_detail",
		accountTable: "erp_order_account",
	}
}

// Set 新增或修改 ERP 订单。
func (o *OrderDaoImpl) Set(c *gin.Context, req *ordermodels.OrderSet) (*ordermodels.Order, error) {
	var resultID uint64
	err := o.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		exists, findErr := o.findExists(tx, c, req)
		if findErr != nil {
			return findErr
		}
		data, parseErr := buildOrderModel(c, req)
		if parseErr != nil {
			return parseErr
		}
		if exists == nil {
			if err := tx.Table(o.table).Create(data).Error; err != nil {
				return err
			}
			resultID = data.ID
		} else {
			data.ID = exists.ID
			data.SetUpdateBy(baizeContext.GetUserId(c))
			if err := tx.Table(o.table).
				Where("id = ?", exists.ID).
				Where("dept_id = ?", baizeContext.GetDeptId(c)).
				Where("state = ?", commonStatus.NORMAL).
				Select("tid", "weight", "size", "buyernick", "buyermessage", "sellermemo", "total", "privilege",
					"postfee", "receivername", "receiverstate", "receivercity", "receiverdistrict", "receiveraddress",
					"receiverphone", "receivermobile", "receiverzip", "created", "status", "type", "invoicename",
					"sellerflag", "paytime", "logistbtypecode", "logistbillcode", "btypecode", "billcode",
					"sync_message", "sync_status", "sync_time", "update_by", "update_time").
				Updates(data).Error; err != nil {
				return err
			}
			resultID = exists.ID
			if err := o.softDeleteChildren(tx, c, resultID); err != nil {
				return err
			}
		}
		if err := o.createDetails(tx, c, resultID, data.Tid, req.Details); err != nil {
			return err
		}
		if err := o.createAccounts(tx, c, resultID, data.Tid, req.Accounts); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return o.GetByID(c, resultID)
}

// GetByID 查询 ERP 订单详情。
func (o *OrderDaoImpl) GetByID(c *gin.Context, id uint64) (*ordermodels.Order, error) {
	var item ordermodels.Order
	if err := o.db.WithContext(c).Table(o.table).
		Where("id = ?", id).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if err := o.attachChildren(c, []*ordermodels.Order{&item}); err != nil {
		return nil, err
	}
	return &item, nil
}

// List 分页查询 ERP 订单。
func (o *OrderDaoImpl) List(c *gin.Context, req *ordermodels.OrderQuery) (*ordermodels.OrderListData, error) {
	db := o.db.WithContext(c).Table(o.table).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
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
	rows := make([]*ordermodels.Order, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	if err := o.attachChildren(c, rows); err != nil {
		return nil, err
	}
	return &ordermodels.OrderListData{
		Rows:  rows,
		Total: total,
	}, nil
}

// DeleteByIDs 软删除 ERP 订单及其子表。
func (o *OrderDaoImpl) DeleteByIDs(c *gin.Context, ids []uint64) error {
	now := time.Now()
	return o.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(o.table).
			Where("id IN ?", ids).
			Where("dept_id = ?", baizeContext.GetDeptId(c)).
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
		if err := tx.Table(o.detailTable).
			Where("order_id IN ?", orderIDs).
			Delete(nil).Error; err != nil {
			return err
		}
		return tx.Table(o.accountTable).
			Where("order_id IN ?", orderIDs).
			Delete(nil).Error
	})
}

// findExists 根据订单ID或订单编号查询当前部门下的有效订单。
func (o *OrderDaoImpl) findExists(tx *gorm.DB, c *gin.Context, req *ordermodels.OrderSet) (*ordermodels.Order, error) {
	var item ordermodels.Order
	db := tx.Table(o.table).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL)
	if req.ID > 0 {
		db = db.Where("id = ?", req.ID)
	} else {
		db = db.Where("tid = ?", req.Tid)
	}
	if err := db.First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// buildOrderModel 将保存参数转换为订单主表模型，并完成时间字段解析。
func buildOrderModel(c *gin.Context, req *ordermodels.OrderSet) (*ordermodels.Order, error) {
	payTime, err := parseOrderTimePtr(req.PayTime)
	if err != nil {
		return nil, errors.New("paytime时间格式错误，要求: 2006-01-02 15:04:05")
	}
	syncTime, err := parseOrderTimePtr(req.SyncTime)
	if err != nil {
		return nil, errors.New("syncTime时间格式错误，要求: 2006-01-02 15:04:05")
	}
	data := &ordermodels.Order{
		Tid:              req.Tid,
		Weight:           req.Weight,
		Size:             req.Size,
		BuyerNick:        req.BuyerNick,
		BuyerMessage:     req.BuyerMessage,
		SellerMemo:       req.SellerMemo,
		Total:            req.Total,
		Privilege:        req.Privilege,
		PostFee:          req.PostFee,
		ReceiverName:     req.ReceiverName,
		ReceiverState:    req.ReceiverState,
		ReceiverCity:     req.ReceiverCity,
		ReceiverDistrict: req.ReceiverDistrict,
		ReceiverAddress:  req.ReceiverAddress,
		ReceiverPhone:    req.ReceiverPhone,
		ReceiverMobile:   req.ReceiverMobile,
		ReceiverZip:      req.ReceiverZip,
		Status:           req.Status,
		Type:             req.Type,
		InvoiceName:      req.InvoiceName,
		SellerFlag:       req.SellerFlag,
		PayTime:          payTime,
		LogistBTypeCode:  req.LogistBTypeCode,
		LogistBillCode:   req.LogistBillCode,
		BTypeCode:        req.BTypeCode,
		BillCode:         req.BillCode,
		SyncMessage:      req.SyncMessage,
		SyncStatus:       req.SyncStatus,
		SyncTime:         syncTime,
		DeptID:           baizeContext.GetDeptId(c),
		State:            commonStatus.NORMAL,
	}
	data.SetCreateBy(baizeContext.GetUserId(c))
	return data, nil
}

// parseOrderTime 解析必填的订单时间字段。
func parseOrderTime(value string) (*time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// parseOrderTimePtr 解析可为空的订单时间字段。
func parseOrderTimePtr(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	return parseOrderTime(value)
}

// softDeleteChildren 在订单更新前清理旧的明细和账户子表数据。
func (o *OrderDaoImpl) softDeleteChildren(tx *gorm.DB, c *gin.Context, orderID uint64) error {
	if err := tx.Table(o.detailTable).
		Where("order_id = ?", orderID).
		Delete(nil).Error; err != nil {
		return err
	}
	return tx.Table(o.accountTable).
		Where("order_id = ?", orderID).
		Delete(nil).Error
}

// createDetails 批量创建订单明细记录。
func (o *OrderDaoImpl) createDetails(tx *gorm.DB, c *gin.Context, orderID uint64, tid string, details []*ordermodels.OrderDetailSet) error {
	if len(details) == 0 {
		return nil
	}
	rows := make([]*ordermodels.OrderDetail, 0, len(details))
	for _, item := range details {
		if item == nil {
			continue
		}
		row := &ordermodels.OrderDetail{
			OrderID:        orderID,
			Tid:            tid,
			OID:            item.OID,
			Barcode:        item.Barcode,
			EShopGoodsID:   item.EShopGoodsID,
			OuterIID:       item.OuterIID,
			EShopGoodsName: item.EShopGoodsName,
			EShopSkuID:     item.EShopSkuID,
			EShopSkuName:   item.EShopSkuName,
			NumIID:         item.NumIID,
			SkuID:          item.SkuID,
			Num:            item.Num,
			Payment:        item.Payment,
			PicPath:        item.PicPath,
			Weight:         item.Weight,
			Size:           item.Size,
			UnitID:         item.UnitID,
			UnitQty:        item.UnitQty,
			DeptID:         baizeContext.GetDeptId(c),
			State:          commonStatus.NORMAL,
		}
		row.SetCreateBy(baizeContext.GetUserId(c))
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		return nil
	}
	return tx.Table(o.detailTable).Create(&rows).Error
}

// createAccounts 批量创建订单账户记录。
func (o *OrderDaoImpl) createAccounts(tx *gorm.DB, c *gin.Context, orderID uint64, tid string, accounts []*ordermodels.OrderAccountSet) error {
	if len(accounts) == 0 {
		return nil
	}
	rows := make([]*ordermodels.OrderAccount, 0, len(accounts))
	for _, item := range accounts {
		if item == nil {
			continue
		}
		row := &ordermodels.OrderAccount{
			OrderID:     orderID,
			Tid:         tid,
			FinanceCode: item.FinanceCode,
			Total:       item.Total,
			DeptID:      baizeContext.GetDeptId(c),
			State:       commonStatus.NORMAL,
		}
		row.SetCreateBy(baizeContext.GetUserId(c))
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		return nil
	}
	return tx.Table(o.accountTable).Create(&rows).Error
}

// attachChildren 为订单结果批量挂载明细与账户列表。
func (o *OrderDaoImpl) attachChildren(c *gin.Context, orders []*ordermodels.Order) error {
	if len(orders) == 0 {
		return nil
	}
	orderIDs := make([]uint64, 0, len(orders))
	for _, item := range orders {
		if item == nil {
			continue
		}
		orderIDs = append(orderIDs, item.ID)
	}
	if len(orderIDs) == 0 {
		return nil
	}
	details := make([]*ordermodels.OrderDetail, 0)
	if err := o.db.WithContext(c).Table(o.detailTable).
		Where("order_id IN ?", orderIDs).
		Where("state = ?", commonStatus.NORMAL).
		Order("id ASC").
		Find(&details).Error; err != nil {
		return err
	}
	accounts := make([]*ordermodels.OrderAccount, 0)
	if err := o.db.WithContext(c).Table(o.accountTable).
		Where("order_id IN ?", orderIDs).
		Where("state = ?", commonStatus.NORMAL).
		Order("id ASC").
		Find(&accounts).Error; err != nil {
		return err
	}
	detailMap := make(map[uint64][]*ordermodels.OrderDetail)
	for _, detail := range details {
		detailMap[detail.OrderID] = append(detailMap[detail.OrderID], detail)
	}
	accountMap := make(map[uint64][]*ordermodels.OrderAccount)
	for _, account := range accounts {
		accountMap[account.OrderID] = append(accountMap[account.OrderID], account)
	}
	for _, item := range orders {
		if item == nil {
			continue
		}
		item.Details = detailMap[item.ID]
		item.Accounts = accountMap[item.ID]
	}
	return nil
}

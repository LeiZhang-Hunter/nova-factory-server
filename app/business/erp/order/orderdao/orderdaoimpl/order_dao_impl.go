package orderdaoimpl

import (
	"errors"
	"nova-factory-server/app/utils/snowflake"
	"strings"
	"time"

	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/erp/order/orderdao"
	"nova-factory-server/app/business/erp/order/ordermodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderDaoImpl ERP订单数据访问实现。
type OrderDaoImpl struct {
	db         *gorm.DB
	table      string
	detailDao  *OrderDetailDaoImpl
	accountDao *OrderAccountDaoImpl
}

// erpOrderRow ERP 订单主表行模型，显式绑定真实表字段名。
type erpOrderRow struct {
	ID                   uint64     `gorm:"column:id"`
	Tid                  string     `gorm:"column:tid"`
	Weight               float64    `gorm:"column:weight"`
	Size                 float64    `gorm:"column:size"`
	BuyerNick            string     `gorm:"column:buyer_nick"`
	BuyerMessage         string     `gorm:"column:buyer_message"`
	SellerMemo           string     `gorm:"column:seller_memo"`
	Total                float64    `gorm:"column:total"`
	Privilege            float64    `gorm:"column:privilege"`
	PostFee              float64    `gorm:"column:post_fee"`
	ReceiverName         string     `gorm:"column:receiver_name"`
	ReceiverProvince     string     `gorm:"column:receiver_province"`
	ReceiverProvinceName string     `gorm:"column:receiver_province_name"`
	ReceiverCity         string     `gorm:"column:receiver_city"`
	ReceiverCityName     string     `gorm:"column:receiver_city_name"`
	ReceiverDistrict     string     `gorm:"column:receiver_district"`
	ReceiverDistrictName string     `gorm:"column:receiver_district_name"`
	ReceiverStreet       string     `gorm:"column:receiver_street"`
	ReceiverStreetName   string     `gorm:"column:receiver_street_name"`
	ReceiverAddress      string     `gorm:"column:receiver_address"`
	ReceiverPhone        string     `gorm:"column:receiver_phone"`
	ReceiverMobile       string     `gorm:"column:receiver_mobile"`
	ReceiverZip          string     `gorm:"column:receiver_zip"`
	Status               string     `gorm:"column:status"`
	Type                 string     `gorm:"column:order_type"`
	InvoiceName          string     `gorm:"column:invoice_name"`
	SellerFlag           string     `gorm:"column:seller_flag"`
	PayTime              *time.Time `gorm:"column:pay_time"`
	LogistBTypeCode      string     `gorm:"column:logist_b_type_code"`
	LogistBillCode       string     `gorm:"column:logist_bill_code"`
	BTypeCode            string     `gorm:"column:b_type_code"`
	DetailsJSON          string     `gorm:"column:details_json"`
	AccountsJSON         string     `gorm:"column:accounts_json"`
	BillCode             string     `gorm:"column:bill_code"`
	SyncMessage          string     `gorm:"column:sync_message"`
	SyncStatus           int32      `gorm:"column:sync_status"`
	SyncTime             *time.Time `gorm:"column:sync_time"`
	DeptID               int64      `gorm:"column:dept_id"`
	CreateBy             int64      `gorm:"column:create_by"`
	CreateTime           *time.Time `gorm:"column:create_time"`
	UpdateBy             int64      `gorm:"column:update_by"`
	UpdateTime           *time.Time `gorm:"column:update_time"`
	State                int32      `gorm:"column:state"`
}

// NewOrderDao 创建 ERP 订单 DAO。
func NewOrderDao(db *gorm.DB) orderdao.IOrderDao {
	return &OrderDaoImpl{
		db:         db,
		table:      "erp_order",
		detailDao:  NewOrderDetailDao(db),
		accountDao: NewOrderAccountDao(db),
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
		row := buildOrderRow(data)
		if exists == nil {
			row.ID = uint64(snowflake.GenID())
			if err := tx.Table(o.table).Create(row).Error; err != nil {
				return err
			}
			resultID = row.ID
		} else {
			data.ID = exists.ID
			row.ID = exists.ID
			data.SetUpdateBy(baizeContext.GetUserId(c))
			row.UpdateBy = data.UpdateBy
			row.UpdateTime = data.UpdateTime
			if err := tx.Table(o.table).
				Where("id = ?", exists.ID).
				//Where("dept_id = ?", baizeContext.GetDeptId(c)).
				Where("state = ?", commonStatus.NORMAL).
				Updates(buildOrderUpdateMap(row)).Error; err != nil {
				return err
			}
			resultID = exists.ID
			if err := o.detailDao.DeleteByOrderID(tx, resultID); err != nil {
				return err
			}
			if err := o.accountDao.DeleteByOrderID(tx, resultID); err != nil {
				return err
			}
		}
		if err := o.detailDao.BatchCreate(tx, c, resultID, data.Tid, req.Details); err != nil {
			return err
		}
		if err := o.accountDao.BatchCreate(tx, c, resultID, data.Tid, req.Accounts); err != nil {
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
	var row erpOrderRow
	if err := o.db.WithContext(c).Table(o.table).
		Where("id = ?", id).
		//Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	item := row.toModel()
	if err := o.attachChildren(c, []*ordermodels.Order{&item}); err != nil {
		return nil, err
	}
	return &item, nil
}

// List 分页查询 ERP 订单。
func (o *OrderDaoImpl) List(c *gin.Context, req *ordermodels.OrderQuery) (*ordermodels.OrderListData, error) {
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
	rowList := make([]*erpOrderRow, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rowList).Error; err != nil {
		return nil, err
	}
	rows := make([]*ordermodels.Order, 0, len(rowList))
	for _, row := range rowList {
		if row == nil {
			continue
		}
		item := row.toModel()
		rows = append(rows, &item)
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
func (o *OrderDaoImpl) findExists(tx *gorm.DB, c *gin.Context, req *ordermodels.OrderSet) (*ordermodels.Order, error) {
	var row erpOrderRow
	db := tx.Table(o.table)
	if req.ID > 0 {
		db = db.Where("id = ?", req.ID)
	} else {
		db = db.Where("tid = ?", req.Tid)
	}
	db = db.Where("state = ?", commonStatus.NORMAL)
	if err := db.First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	item := row.toModel()
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
		Tid:                  req.Tid,
		Weight:               req.Weight,
		Size:                 req.Size,
		BuyerNick:            req.BuyerNick,
		BuyerMessage:         req.BuyerMessage,
		SellerMemo:           req.SellerMemo,
		Total:                req.Total,
		Privilege:            req.Privilege,
		PostFee:              req.PostFee,
		ReceiverName:         req.ReceiverName,
		ReceiverProvince:     req.ReceiverProvince,
		ReceiverProvinceName: req.ReceiverProvinceName,
		ReceiverCity:         req.ReceiverCity,
		ReceiverCityName:     req.ReceiverCityName,
		ReceiverDistrict:     req.ReceiverDistrict,
		ReceiverDistrictName: req.ReceiverDistrictName,
		ReceiverStreet:       req.ReceiverStreet,
		ReceiverStreetName:   req.ReceiverStreetName,
		ReceiverAddress:      req.ReceiverAddress,
		ReceiverPhone:        req.ReceiverPhone,
		ReceiverMobile:       req.ReceiverMobile,
		ReceiverZip:          req.ReceiverZip,
		Status:               req.Status,
		Type:                 req.OrderType,
		InvoiceName:          req.InvoiceName,
		SellerFlag:           req.SellerFlag,
		PayTime:              payTime,
		LogistBTypeCode:      req.LogistBTypeCode,
		LogistBillCode:       req.LogistBillCode,
		BTypeCode:            req.BTypeCode,
		BillCode:             req.BillCode,
		SyncMessage:          req.SyncMessage,
		SyncStatus:           req.SyncStatus,
		SyncTime:             syncTime,
		//DeptID:               baizeContext.GetDeptId(c),
		State: commonStatus.NORMAL,
	}
	data.SetCreateBy(baizeContext.GetUserId(c))
	return data, nil
}

// buildOrderRow 将领域模型转换为真实表结构行模型。
func buildOrderRow(data *ordermodels.Order) *erpOrderRow {
	if data == nil {
		return nil
	}
	return &erpOrderRow{
		ID:                   data.ID,
		Tid:                  data.Tid,
		Weight:               data.Weight,
		Size:                 data.Size,
		BuyerNick:            data.BuyerNick,
		BuyerMessage:         data.BuyerMessage,
		SellerMemo:           data.SellerMemo,
		Total:                data.Total,
		Privilege:            data.Privilege,
		PostFee:              data.PostFee,
		ReceiverName:         data.ReceiverName,
		ReceiverProvince:     data.ReceiverProvince,
		ReceiverProvinceName: data.ReceiverProvinceName,
		ReceiverCity:         data.ReceiverCity,
		ReceiverCityName:     data.ReceiverCityName,
		ReceiverDistrict:     data.ReceiverDistrict,
		ReceiverDistrictName: data.ReceiverDistrictName,
		ReceiverStreet:       data.ReceiverStreet,
		ReceiverStreetName:   data.ReceiverStreetName,
		ReceiverAddress:      data.ReceiverAddress,
		ReceiverPhone:        data.ReceiverPhone,
		ReceiverMobile:       data.ReceiverMobile,
		ReceiverZip:          data.ReceiverZip,
		Status:               data.Status,
		Type:                 data.Type,
		InvoiceName:          data.InvoiceName,
		SellerFlag:           data.SellerFlag,
		PayTime:              data.PayTime,
		LogistBTypeCode:      data.LogistBTypeCode,
		LogistBillCode:       data.LogistBillCode,
		BTypeCode:            data.BTypeCode,
		BillCode:             data.BillCode,
		SyncMessage:          data.SyncMessage,
		SyncStatus:           data.SyncStatus,
		SyncTime:             data.SyncTime,
		DeptID:               data.DeptID,
		CreateBy:             data.CreateBy,
		CreateTime:           data.CreateTime,
		UpdateBy:             data.UpdateBy,
		UpdateTime:           data.UpdateTime,
		State:                data.State,
	}
}

// toModel 将真实表结构行模型转换为领域模型。
func (r *erpOrderRow) toModel() ordermodels.Order {
	if r == nil {
		return ordermodels.Order{}
	}
	return ordermodels.Order{
		ID:                   r.ID,
		Tid:                  r.Tid,
		Weight:               r.Weight,
		Size:                 r.Size,
		BuyerNick:            r.BuyerNick,
		BuyerMessage:         r.BuyerMessage,
		SellerMemo:           r.SellerMemo,
		Total:                r.Total,
		Privilege:            r.Privilege,
		PostFee:              r.PostFee,
		ReceiverName:         r.ReceiverName,
		ReceiverProvince:     r.ReceiverProvince,
		ReceiverProvinceName: r.ReceiverProvinceName,
		ReceiverCity:         r.ReceiverCity,
		ReceiverCityName:     r.ReceiverCityName,
		ReceiverDistrict:     r.ReceiverDistrict,
		ReceiverDistrictName: r.ReceiverDistrictName,
		ReceiverStreet:       r.ReceiverStreet,
		ReceiverStreetName:   r.ReceiverStreetName,
		ReceiverAddress:      r.ReceiverAddress,
		ReceiverPhone:        r.ReceiverPhone,
		ReceiverMobile:       r.ReceiverMobile,
		ReceiverZip:          r.ReceiverZip,
		Status:               r.Status,
		Type:                 r.Type,
		InvoiceName:          r.InvoiceName,
		SellerFlag:           r.SellerFlag,
		PayTime:              r.PayTime,
		LogistBTypeCode:      r.LogistBTypeCode,
		LogistBillCode:       r.LogistBillCode,
		BTypeCode:            r.BTypeCode,
		BillCode:             r.BillCode,
		SyncMessage:          r.SyncMessage,
		SyncStatus:           r.SyncStatus,
		SyncTime:             r.SyncTime,
		DeptID:               r.DeptID,
		BaseEntity: baize.BaseEntity{
			CreateBy:   r.CreateBy,
			CreateTime: r.CreateTime,
			UpdateBy:   r.UpdateBy,
			UpdateTime: r.UpdateTime,
		},
		State: r.State,
	}
}

// buildOrderUpdateMap 构建 ERP 订单更新字段映射。
func buildOrderUpdateMap(row *erpOrderRow) map[string]interface{} {
	if row == nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"tid":                    row.Tid,
		"weight":                 row.Weight,
		"size":                   row.Size,
		"buyer_nick":             row.BuyerNick,
		"buyer_message":          row.BuyerMessage,
		"seller_memo":            row.SellerMemo,
		"total":                  row.Total,
		"privilege":              row.Privilege,
		"post_fee":               row.PostFee,
		"receiver_name":          row.ReceiverName,
		"receiver_province":      row.ReceiverProvince,
		"receiver_province_name": row.ReceiverProvinceName,
		"receiver_city":          row.ReceiverCity,
		"receiver_city_name":     row.ReceiverCityName,
		"receiver_district":      row.ReceiverDistrict,
		"receiver_district_name": row.ReceiverDistrictName,
		"receiver_street":        row.ReceiverStreet,
		"receiver_street_name":   row.ReceiverStreetName,
		"receiver_address":       row.ReceiverAddress,
		"receiver_phone":         row.ReceiverPhone,
		"receiver_mobile":        row.ReceiverMobile,
		"receiver_zip":           row.ReceiverZip,
		"status":                 row.Status,
		"order_type":             row.Type,
		"invoice_name":           row.InvoiceName,
		"seller_flag":            row.SellerFlag,
		"pay_time":               row.PayTime,
		"logist_b_type_code":     row.LogistBTypeCode,
		"logist_bill_code":       row.LogistBillCode,
		"b_type_code":            row.BTypeCode,
		"bill_code":              row.BillCode,
		"sync_message":           row.SyncMessage,
		"sync_status":            row.SyncStatus,
		"sync_time":              row.SyncTime,
		"update_by":              row.UpdateBy,
		"update_time":            row.UpdateTime,
	}
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
	details, err := o.detailDao.ListByOrderIDs(c, orderIDs)
	if err != nil {
		return err
	}
	accounts, err := o.accountDao.ListByOrderIDs(c, orderIDs)
	if err != nil {
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

package saledaoimpl

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"nova-factory-server/app/business/erp/sale/saledao"
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderAuditDaoImpl ERP订单审核数据访问实现。
type OrderAuditDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewOrderAuditDao 创建 ERP 订单审核 DAO。
func NewOrderAuditDao(db *gorm.DB) saledao.IOrderAuditDao {
	return &OrderAuditDaoImpl{
		db:    db,
		table: "erp_order_audit",
	}
}

func (o *OrderAuditDaoImpl) Set(c *gin.Context, req *salemodels.OrderAuditSet) (*salemodels.OrderAudit, error) {
	var resultID uint64
	db := o.db.WithContext(c)
	err := db.Transaction(func(tx *gorm.DB) error {
		exists, err := o.findExists(tx, req)
		if err != nil {
			return err
		}
		data, err := buildOrderAuditModel(c, req)
		if err != nil {
			return err
		}
		if exists == nil {
			if err := tx.Table(o.table).Create(data).Error; err != nil {
				return err
			}
			resultID = data.ID
			return nil
		}
		data.ID = exists.ID
		data.SetUpdateBy(baizeContext.GetUserId(c))
		if err := tx.Table(o.table).
			Where("id = ?", exists.ID).
			Where("state = ?", commonStatus.NORMAL).
			Select("tid", "weight", "size", "buyer_nick", "buyer_message", "seller_memo", "total", "privilege",
				"post_fee", "receiver_name", "receiver_province", "receiver_province_name", "receiver_city",
				"receiver_city_name", "receiver_district", "receiver_district_name", "receiver_street",
				"receiver_street_name", "receiver_address", "receiver_phone", "receiver_mobile", "receiver_zip",
				"status", "order_type", "invoice_name", "seller_flag", "pay_time", "logist_b_type_code",
				"logist_bill_code", "b_type_code", "details_json", "accounts_json", "source_json",
				"audit_status", "audit_remark", "audit_by", "audit_time", "transfer_status", "transfer_message",
				"transfer_time", "erp_order_id", "update_by", "update_time").
			Updates(data).Error; err != nil {
			return err
		}
		resultID = exists.ID
		return nil
	})
	if err != nil {
		return nil, err
	}
	return o.GetByID(c, resultID)
}

func (o *OrderAuditDaoImpl) GetByID(c *gin.Context, id uint64) (*salemodels.OrderAudit, error) {
	var item salemodels.OrderAudit
	if err := o.db.WithContext(c).Table(o.table).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	decodeOrderAudit(&item)
	return &item, nil
}

func (o *OrderAuditDaoImpl) List(c *gin.Context, req *salemodels.OrderAuditQuery) (*salemodels.OrderAuditListData, error) {
	db := o.db.WithContext(c).Table(o.table)
	if req != nil {
		if strings.TrimSpace(req.Tid) != "" {
			db = db.Where("tid LIKE ?", "%"+strings.TrimSpace(req.Tid)+"%")
		}
		if req.AuditStatus != nil {
			db = db.Where("audit_status = ?", *req.AuditStatus)
		}
		if req.TransferStatus != nil {
			db = db.Where("transfer_status = ?", *req.TransferStatus)
		}
		if strings.TrimSpace(req.ReceiverName) != "" {
			db = db.Where("receiver_name LIKE ?", "%"+strings.TrimSpace(req.ReceiverName)+"%")
		}
	}

	userId := baizeContext.GetUserId(c)
	db = db.Where("create_by = ?", userId)
	db = db.Where("state = ?", commonStatus.NORMAL)
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
	rows := make([]*salemodels.OrderAudit, 0)
	if err := db.Order("id DESC").Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		decodeOrderAudit(row)
	}
	return &salemodels.OrderAuditListData{Rows: rows, Total: total}, nil
}

func (o *OrderAuditDaoImpl) DeleteByIDs(c *gin.Context, ids []uint64) error {

	return o.db.WithContext(c).Table(o.table).
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Delete(salemodels.OrderAudit{}).Error
}

func (o *OrderAuditDaoImpl) Approve(c *gin.Context, id uint64, remark string, erpOrderID uint64) error {
	return o.ApproveWithTx(c, o.db, id, remark, erpOrderID)
}

func (o *OrderAuditDaoImpl) ApproveWithTx(c *gin.Context, tx *gorm.DB, id uint64, remark string, erpOrderID uint64) error {
	now := time.Now()
	return tx.WithContext(c).Table(o.table).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"audit_status":     1,
			"audit_remark":     remark,
			"audit_by":         baizeContext.GetUserId(c),
			"audit_time":       now,
			"transfer_status":  1,
			"transfer_message": "已转入正式订单",
			"transfer_time":    now,
			"erp_order_id":     erpOrderID,
			"update_by":        baizeContext.GetUserId(c),
			"update_time":      now,
		}).Error
}

func (o *OrderAuditDaoImpl) Reject(c *gin.Context, id uint64, remark string) error {
	return o.RejectWithTx(c, o.db, id, remark)
}

func (o *OrderAuditDaoImpl) RejectWithTx(c *gin.Context, tx *gorm.DB, id uint64, remark string) error {
	now := time.Now()
	return tx.WithContext(c).Table(o.table).
		Where("id = ?", id).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"audit_status": 2,
			"audit_remark": remark,
			"audit_by":     baizeContext.GetUserId(c),
			"audit_time":   now,
			"update_by":    baizeContext.GetUserId(c),
			"update_time":  now,
		}).Error
}

func (o *OrderAuditDaoImpl) GetByIDWithTx(c *gin.Context, tx *gorm.DB, id uint64) (*salemodels.OrderAudit, error) {
	var item salemodels.OrderAudit
	if err := tx.WithContext(c).Table(o.table).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	decodeOrderAudit(&item)
	return &item, nil
}

func (o *OrderAuditDaoImpl) findExists(tx *gorm.DB, req *salemodels.OrderAuditSet) (*salemodels.OrderAudit, error) {
	var item salemodels.OrderAudit
	db := tx.Table(o.table)
	if req.ID > 0 {
		db = db.Where("id = ?", req.ID)
	} else {
		db = db.Where("tid = ?", req.Tid)
	}
	if err := db.Where("state = ?", commonStatus.NORMAL).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func buildOrderAuditModel(c *gin.Context, req *salemodels.OrderAuditSet) (*salemodels.OrderAudit, error) {
	payTime, err := parseOrderTimePtr(req.PayTime)
	if err != nil {
		return nil, errors.New("paytime时间格式错误，要求: 2006-01-02 15:04:05")
	}
	detailsJSON, err := json.Marshal(req.Details)
	if err != nil {
		return nil, err
	}
	accountsJSON, err := json.Marshal(req.Accounts)
	if err != nil {
		return nil, err
	}
	sourceJSON := strings.TrimSpace(req.SourceJSON)
	if sourceJSON == "" {
		raw, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}
		sourceJSON = string(raw)
	}
	data := &salemodels.OrderAudit{
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
		DetailsJSON:          string(detailsJSON),
		AccountsJSON:         string(accountsJSON),
		SourceJSON:           sourceJSON,
		DeptID:               baizeContext.GetDeptId(c),
		State:                commonStatus.NORMAL,
		Details:              convertDetailSets(req.Details),
		Accounts:             convertAccountSets(req.Accounts),
	}
	data.SetCreateBy(baizeContext.GetUserId(c))
	return data, nil
}

func decodeOrderAudit(item *salemodels.OrderAudit) {
	if item == nil {
		return
	}
	if strings.TrimSpace(item.DetailsJSON) != "" {
		_ = json.Unmarshal([]byte(item.DetailsJSON), &item.Details)
	}
	if strings.TrimSpace(item.AccountsJSON) != "" {
		_ = json.Unmarshal([]byte(item.AccountsJSON), &item.Accounts)
	}
	if item.Details == nil {
		item.Details = make([]*salemodels.OrderDetail, 0)
	}
	if item.Accounts == nil {
		item.Accounts = make([]*salemodels.OrderAccount, 0)
	}
}

func convertDetailSets(items []*salemodels.OrderDetailSet) []*salemodels.OrderDetail {
	result := make([]*salemodels.OrderDetail, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		result = append(result, &salemodels.OrderDetail{
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
		})
	}
	return result
}

func convertAccountSets(items []*salemodels.OrderAccountSet) []*salemodels.OrderAccount {
	result := make([]*salemodels.OrderAccount, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		result = append(result, &salemodels.OrderAccount{
			FinanceCode: item.FinanceCode,
			Total:       item.Total,
		})
	}
	return result
}

package saleserviceimpl

import (
	"errors"
	"strings"

	"nova-factory-server/app/business/erp/sale/saledao"
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/business/erp/sale/saleservice"

	"github.com/gin-gonic/gin"
)

// OrderAuditServiceImpl ERP订单审核服务实现。
type OrderAuditServiceImpl struct {
	dao      saledao.IOrderAuditDao
	orderDao saledao.IOrderDao
}

// NewOrderAuditService 创建 ERP 订单审核服务。
func NewOrderAuditService(dao saledao.IOrderAuditDao, orderDao saledao.IOrderDao) saleservice.IOrderAuditService {
	return &OrderAuditServiceImpl{dao: dao, orderDao: orderDao}
}

func (o *OrderAuditServiceImpl) Set(c *gin.Context, req *salemodels.OrderAuditSet) (*salemodels.OrderAudit, error) {
	if err := o.validateSet(req); err != nil {
		return nil, err
	}
	return o.dao.Set(c, req)
}

func (o *OrderAuditServiceImpl) ImportData(c *gin.Context, req *salemodels.OrderAuditSet) (*salemodels.OrderAudit, error) {

	return o.dao.Set(c, req)
}
func (o *OrderAuditServiceImpl) Import(c *gin.Context, req *salemodels.OrderAuditImportReq) (*salemodels.OrderAuditImportResult, error) {
	if req == nil || len(req.Records) == 0 {
		return nil, errors.New("导入记录不能为空")
	}
	result := &salemodels.OrderAuditImportResult{
		Total: len(req.Records),
		Items: make([]*salemodels.OrderAuditImportItem, 0, len(req.Records)),
	}
	for index, record := range req.Records {
		item := &salemodels.OrderAuditImportItem{
			Index: index,
		}
		if record != nil {
			item.Tid = strings.TrimSpace(record.Tid)
		}
		data, err := o.ImportData(c, record)
		if err != nil {
			item.Success = false
			item.Message = err.Error()
			result.FailureCount++
			result.Items = append(result.Items, item)
			continue
		}
		item.Success = true
		item.Message = "导入成功"
		if data != nil {
			item.ID = data.ID
			if item.Tid == "" {
				item.Tid = data.Tid
			}
		}
		result.SuccessCount++
		result.Items = append(result.Items, item)
	}
	return result, nil
}

func (o *OrderAuditServiceImpl) GetByID(c *gin.Context, id uint64) (*salemodels.OrderAudit, error) {
	if id == 0 {
		return nil, errors.New("id不能为空")
	}
	return o.dao.GetByID(c, id)
}

func (o *OrderAuditServiceImpl) List(c *gin.Context, req *salemodels.OrderAuditQuery) (*salemodels.OrderAuditListData, error) {
	if req == nil {
		req = new(salemodels.OrderAuditQuery)
	}
	req.Tid = strings.TrimSpace(req.Tid)
	req.ReceiverName = strings.TrimSpace(req.ReceiverName)
	return o.dao.List(c, req)
}

func (o *OrderAuditServiceImpl) DeleteByIDs(c *gin.Context, ids []uint64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的订单审核记录")
	}
	return o.dao.DeleteByIDs(c, ids)
}

func (o *OrderAuditServiceImpl) Approve(c *gin.Context, req *salemodels.OrderAuditAction) (*salemodels.OrderAudit, error) {
	if req == nil || req.ID == 0 {
		return nil, errors.New("id不能为空")
	}
	item, err := o.dao.GetByID(c, req.ID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("订单审核记录不存在")
	}
	if item.AuditStatus == 1 && item.TransferStatus == 1 && item.ERPOrderID > 0 {
		return item, nil
	}
	order, err := o.orderDao.Set(c, o.toOrderSet(item))
	if err != nil {
		return nil, err
	}
	if err := o.dao.Approve(c, req.ID, strings.TrimSpace(req.AuditRemark), order.ID); err != nil {
		return nil, err
	}
	return o.dao.GetByID(c, req.ID)
}

func (o *OrderAuditServiceImpl) Reject(c *gin.Context, req *salemodels.OrderAuditAction) (*salemodels.OrderAudit, error) {
	if req == nil || req.ID == 0 {
		return nil, errors.New("id不能为空")
	}
	item, err := o.dao.GetByID(c, req.ID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("订单审核记录不存在")
	}
	if err := o.dao.Reject(c, req.ID, strings.TrimSpace(req.AuditRemark)); err != nil {
		return nil, err
	}
	return o.dao.GetByID(c, req.ID)
}

func (o *OrderAuditServiceImpl) validateSet(req *salemodels.OrderAuditSet) error {
	if req == nil {
		return errors.New("参数不能为空")
	}
	req.Tid = strings.TrimSpace(req.Tid)
	req.ReceiverName = strings.TrimSpace(req.ReceiverName)
	req.ReceiverProvince = strings.TrimSpace(req.ReceiverProvince)
	req.ReceiverCity = strings.TrimSpace(req.ReceiverCity)
	req.ReceiverDistrict = strings.TrimSpace(req.ReceiverDistrict)
	req.ReceiverAddress = strings.TrimSpace(req.ReceiverAddress)
	req.ReceiverMobile = strings.TrimSpace(req.ReceiverMobile)
	req.Status = strings.TrimSpace(req.Status)
	req.OrderType = strings.TrimSpace(req.OrderType)
	if req.ID == 0 && req.Tid == "" {
		return errors.New("tid不能为空")
	}
	if req.ReceiverName == "" {
		return errors.New("收货人名称不能为空")
	}
	if req.ReceiverProvince == "" {
		return errors.New("收货省不能为空")
	}
	if req.ReceiverCity == "" {
		return errors.New("收货市不能为空")
	}
	if req.ReceiverDistrict == "" {
		return errors.New("收货区不能为空")
	}
	if req.ReceiverAddress == "" {
		return errors.New("收货地址不能为空")
	}
	if req.ReceiverMobile == "" {
		return errors.New("收货人手机号不能为空")
	}
	if req.Status == "" {
		return errors.New("status不能为空")
	}
	if req.OrderType == "" {
		return errors.New("订单type不能为空")
	}
	if len(req.Details) == 0 {
		return errors.New("details不能为空")
	}
	return nil
}

func (o *OrderAuditServiceImpl) toOrderSet(item *salemodels.OrderAudit) *salemodels.OrderSet {
	details := make([]*salemodels.OrderDetailSet, 0, len(item.Details))
	for _, detail := range item.Details {
		if detail == nil {
			continue
		}
		details = append(details, &salemodels.OrderDetailSet{
			OID:            detail.OID,
			Barcode:        detail.Barcode,
			EShopGoodsID:   detail.EShopGoodsID,
			OuterIID:       detail.OuterIID,
			EShopGoodsName: detail.EShopGoodsName,
			EShopSkuID:     detail.EShopSkuID,
			EShopSkuName:   detail.EShopSkuName,
			NumIID:         detail.NumIID,
			SkuID:          detail.SkuID,
			Num:            detail.Num,
			Payment:        detail.Payment,
			PicPath:        detail.PicPath,
			Weight:         detail.Weight,
			Size:           detail.Size,
			UnitID:         detail.UnitID,
			UnitQty:        detail.UnitQty,
		})
	}
	accounts := make([]*salemodels.OrderAccountSet, 0, len(item.Accounts))
	for _, account := range item.Accounts {
		if account == nil {
			continue
		}
		accounts = append(accounts, &salemodels.OrderAccountSet{
			FinanceCode: account.FinanceCode,
			Total:       account.Total,
		})
	}
	payTime := ""
	if item.PayTime != nil {
		payTime = item.PayTime.Format("2006-01-02 15:04:05")
	}
	return &salemodels.OrderSet{
		Tid:                  item.Tid,
		Weight:               item.Weight,
		Size:                 item.Size,
		BuyerNick:            item.BuyerNick,
		BuyerMessage:         item.BuyerMessage,
		SellerMemo:           item.SellerMemo,
		Total:                item.Total,
		Privilege:            item.Privilege,
		PostFee:              item.PostFee,
		ReceiverName:         item.ReceiverName,
		ReceiverProvince:     item.ReceiverProvince,
		ReceiverProvinceName: item.ReceiverProvinceName,
		ReceiverCity:         item.ReceiverCity,
		ReceiverCityName:     item.ReceiverCityName,
		ReceiverDistrict:     item.ReceiverDistrict,
		ReceiverDistrictName: item.ReceiverDistrictName,
		ReceiverStreet:       item.ReceiverStreet,
		ReceiverStreetName:   item.ReceiverStreetName,
		ReceiverAddress:      item.ReceiverAddress,
		ReceiverPhone:        item.ReceiverPhone,
		ReceiverMobile:       item.ReceiverMobile,
		ReceiverZip:          item.ReceiverZip,
		Status:               item.Status,
		OrderType:            item.Type,
		InvoiceName:          item.InvoiceName,
		SellerFlag:           item.SellerFlag,
		PayTime:              payTime,
		LogistBTypeCode:      item.LogistBTypeCode,
		LogistBillCode:       item.LogistBillCode,
		BTypeCode:            item.BTypeCode,
		Details:              details,
		Accounts:             accounts,
	}
}

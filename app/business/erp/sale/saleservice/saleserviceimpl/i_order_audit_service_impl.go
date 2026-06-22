package saleserviceimpl

import (
	"errors"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/admin/system/systemdao"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	mastermodels "nova-factory-server/app/business/erp/master/mastermodels"
	masterservice "nova-factory-server/app/business/erp/master/masterservice"
	"nova-factory-server/app/business/erp/setting/settingdao"
	"nova-factory-server/app/datasource/cache"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"nova-factory-server/app/business/erp/sale/saledao"
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/business/erp/sale/saleservice"
	searchutil "nova-factory-server/app/utils/vectorsearch"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OrderAuditServiceImpl ERP订单审核服务实现。
type OrderAuditServiceImpl struct {
	dao                  saledao.IOrderAuditDao
	orderDao             saledao.IOrderDao
	productService       masterservice.IProductService
	integrationConfigDao settingdao.IIntegrationConfigDao
	modelService         aidatasetservice.IAiModelProviderService
	dictDataDao          systemdao.IDictDataDao
	db                   *gorm.DB
	cache                cache.Cache
}

// NewOrderAuditService 创建 ERP 订单审核服务。
func NewOrderAuditService(dao saledao.IOrderAuditDao, orderDao saledao.IOrderDao,
	productService masterservice.IProductService, db *gorm.DB,
	integrationConfigDao settingdao.IIntegrationConfigDao, cache cache.Cache,
	modelService aidatasetservice.IAiModelProviderService,
	dictDataDao systemdao.IDictDataDao) saleservice.IOrderAuditService {
	return &OrderAuditServiceImpl{
		dao:                  dao,
		orderDao:             orderDao,
		productService:       productService,
		db:                   db,
		integrationConfigDao: integrationConfigDao,
		cache:                cache,
		modelService:         modelService,
		dictDataDao:          dictDataDao,
	}
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
	info, err := o.dao.GetByID(c, id)
	if err != nil || info == nil {
		return info, err
	}
	if info.AuditStatus == 0 {
		if err = o.fillOrderAuditDetails(c, info); err != nil {
			zap.L().Warn("fill order audit details failed", zap.Uint64("id", id), zap.Error(err))
		}
	}
	info.Status = o.normalizeOrderStatus(c, info.Status)

	// 付款时间
	if info.PayTime == nil {
		kt := baize.NewTime().ToString()
		// 按当前服务时区解析时间字符串，避免默认按 UTC 解析导致时区偏差。
		if v, err := time.ParseInLocation("2006-01-02 15:04:05", kt, time.Local); err == nil {
			info.PayTime = &v
		} else {
			now := time.Now().In(time.Local)
			info.PayTime = &now
		}
	}
	return info, nil
}

func (o *OrderAuditServiceImpl) normalizeOrderStatus(c *gin.Context, status string) string {
	if o.dictDataDao == nil {
		return strings.TrimSpace(status)
	}
	rows := o.dictDataDao.SelectDictDataByType(c, "erp_order_status")
	if len(rows) == 0 {
		return strings.TrimSpace(status)
	}
	status = strings.TrimSpace(status)
	firstValue := ""
	for _, row := range rows {
		if row == nil {
			continue
		}
		value := strings.TrimSpace(row.DictValue)
		if value == "" {
			continue
		}
		if firstValue == "" {
			firstValue = value
		}
		if value == status {
			return status
		}
	}
	if firstValue != "" {
		return firstValue
	}
	return status
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
	var result *salemodels.OrderAudit
	db := o.db.WithContext(c)
	err = db.Transaction(func(tx *gorm.DB) error {
		order, setErr := o.orderDao.SetWithTx(c, tx, o.toOrderSet(item))
		if setErr != nil {
			return setErr
		}
		if order == nil {
			return errors.New("正式订单保存成功但事务内未返回订单")
		}
		if approveErr := o.dao.ApproveWithTx(c, tx, req.ID, strings.TrimSpace(req.AuditRemark), order.ID); approveErr != nil {
			return approveErr
		}
		result, err = o.dao.GetByIDWithTx(c, tx, req.ID)
		if err != nil {
			zap.L().Error("get id error", zap.Error(err))
			return err
		}

		// CheckLoginState 检查当前 ERP 集成客户端的登录状态。
		cfg, err := o.integrationConfigDao.GetEnabled(c)
		if err != nil {
			return err
		}
		if cfg == nil {
			return nil
		}
		service, err := cfg.Service()
		if err != nil {
			return err
		}
		if service == nil {
			return errors.New("没有配置集成商")
		}
		data := o.toOrderSyncRequest(item)
		data.WithConfig(cfg)
		data.WithCache(o.cache)
		_, err = service.OrderSyncer().SyncOrders(c, data)
		return err
	})
	if err != nil {
		zap.L().Error("sync order error", zap.Error(err))
	}
	return result, err
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
	var result *salemodels.OrderAudit
	err = o.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if rejectErr := o.dao.RejectWithTx(c, tx, req.ID, strings.TrimSpace(req.AuditRemark)); rejectErr != nil {
			return rejectErr
		}
		result, err = o.dao.GetByIDWithTx(c, tx, req.ID)
		return err
	})
	return result, err
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

func (o *OrderAuditServiceImpl) toOrderSyncRequest(item *salemodels.OrderAudit) *salemodels.OrderSyncRequest {
	if item == nil {
		return nil
	}
	details := make([]*salemodels.OrderSyncDetail, 0, len(item.Details))
	for _, detail := range item.Details {
		if detail == nil {
			continue
		}
		details = append(details, &salemodels.OrderSyncDetail{
			OID:            detail.OID,
			Barcode:        detail.Barcode,
			EShopGoodsID:   detail.EShopGoodsID,
			OuterIID:       detail.OuterIID,
			EShopGoodsName: detail.EShopGoodsName,
			EShopSKUId:     detail.EShopSkuID,
			EShopSKUName:   detail.EShopSkuName,
			NumIID:         detail.NumIID,
			SKUId:          detail.SkuID,
			Num:            detail.Num,
			Payment:        detail.Payment,
			PicPath:        detail.PicPath,
			Weight:         detail.Weight,
			Size:           detail.Size,
			UnitID:         detail.UnitID,
			UnitQty:        detail.UnitQty,
		})
	}
	accounts := make([]*salemodels.OrderSyncAccount, 0, len(item.Accounts))
	for _, account := range item.Accounts {
		if account == nil {
			continue
		}
		accounts = append(accounts, &salemodels.OrderSyncAccount{
			FinanceCode: account.FinanceCode,
			Total:       account.Total,
		})
	}
	var created string
	if item.CreateTime != nil {
		created = item.CreateTime.Format("2006-01-02 15:04:05")
	}
	var payTime string
	if item.PayTime != nil {
		formatted := item.PayTime.Format("2006-01-02 15:04:05")
		payTime = formatted
	}
	return &salemodels.OrderSyncRequest{
		Orders: []*salemodels.OrderSyncOrder{
			{
				Tid:              item.Tid,
				Weight:           item.Weight,
				Size:             item.Size,
				BuyerNick:        item.BuyerNick,
				BuyerMessage:     item.BuyerMessage,
				SellerMemo:       item.SellerMemo,
				Total:            item.Total,
				Privilege:        item.Privilege,
				PostFee:          item.PostFee,
				ReceiverName:     item.ReceiverName,
				ReceiverState:    item.ReceiverProvince,
				ReceiverCity:     item.ReceiverCity,
				ReceiverDistrict: item.ReceiverDistrict,
				ReceiverAddress:  item.ReceiverAddress,
				ReceiverPhone:    item.ReceiverPhone,
				ReceiverMobile:   item.ReceiverMobile,
				Created:          created,
				Status:           item.Status,
				Type:             item.Type,
				InvoiceName:      item.InvoiceName,
				SellerFlag:       item.SellerFlag,
				PayTime:          payTime,
				LogistBTypeCode:  item.LogistBTypeCode,
				LogistBillCode:   item.LogistBillCode,
				BTypeCode:        item.BTypeCode,
				Details:          details,
				Accounts:         accounts,
			},
		},
	}
}

// fillOrderAuditDetails 填充商品详情，
func (o *OrderAuditServiceImpl) fillOrderAuditDetails(c *gin.Context, info *salemodels.OrderAudit) error {
	if info == nil || len(info.Details) == 0 || o.productService == nil {
		return nil
	}
	queries := make([]string, 0, len(info.Details))
	indexes := make([]int, 0, len(info.Details))
	for index, detail := range info.Details {
		query := buildOrderAuditDetailSearchQuery(detail)
		if query == "" {
			continue
		}
		queries = append(queries, query)
		indexes = append(indexes, index)
	}
	if len(queries) == 0 {
		return nil
	}
	embeddingInfo, err := o.modelService.EmbeddingWithLLM(c)
	if err != nil {
		return err
	}
	data, err := o.productService.BatchSearchVector(c, &mastermodels.ProductVectorBatchSearchReq{
		Queries: queries,
		Limit:   3,
		Embedding: &mastermodels.ProductEmbeddingConfig{
			ProviderType: embeddingInfo.GetAPIType(),
			ProviderID:   embeddingInfo.GetLLMFactory(),
			APIEndpoint:  embeddingInfo.GetAPIBase(),
			ModelID:      embeddingInfo.GetLLMName(),
			ApiKey:       embeddingInfo.GetAPIKey(),
		},
	})
	if err != nil {
		return err
	}
	if data == nil {
		return nil
	}
	for i, row := range data.Rows {
		if i >= len(indexes) {
			break
		}
		fillOrderAuditDetailMatchedProduct(info.Details[indexes[i]], row)
	}
	return nil
}

func buildOrderAuditDetailSearchQuery(detail *salemodels.OrderDetail) string {
	return searchutil.BuildLabeledContentFromProvider(detail, 0)
}

func fillOrderAuditDetailMatchedProduct(detail *salemodels.OrderDetail, item *mastermodels.ProductVectorBatchSearchItem) {
	if detail == nil || item == nil || len(item.Rows) == 0 || item.Rows[0] == nil {
		return
	}
	match := item.Rows[0]
	if match.BarCode != "" {
		detail.OldBarcode = detail.Barcode
		detail.Barcode = match.BarCode
	}

	if match.Name != "" {
		detail.OldEShopGoodsName = detail.EShopGoodsName
		detail.EShopGoodsName = match.Name
	}
	if match.Standard != "" {
		detail.OldEShopSkuName = detail.EShopSkuName
		detail.EShopSkuName = match.Name
	}
	if match.UnitId > 0 {
		detail.UnitID = match.UnitId
	}
	if match.ProductID != 0 {
		detail.OldEShopSkuID = detail.EShopSkuID
		detail.EShopSkuID = strconv.FormatInt(match.ProductID, 10)
	}

	if match.SalePrice != 0.0 {
		detail.OldPayment = detail.Payment
		detail.Payment = match.SalePrice
	}

}

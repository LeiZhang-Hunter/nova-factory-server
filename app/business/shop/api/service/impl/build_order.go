package impl

import (
	"fmt"
	"nova-factory-server/app/business/shop/api/models"
	shopordermodels "nova-factory-server/app/business/shop/order/models"
	shopusermodels "nova-factory-server/app/business/shop/user/models"
	orderConstant "nova-factory-server/app/constant/order"
	shopConstant "nova-factory-server/app/constant/shop"
	"strings"
)

// buildERPOrderSet 将商城订单请求转换为 ERP 订单保存参数。
func (s *IApiShopOrderServiceImpl) buildERPOrderSet(
	orderNo string,
	shopUser *shopusermodels.User,
	address *models.ShopUserAddressApp,
	cacheData *models.OrderCacheData,
	req *models.OrderCreateReq,
) *shopordermodels.OrderSet {
	details := make([]*shopordermodels.OrderDetailSet, 0, len(cacheData.Items))
	for index, item := range cacheData.Items {
		if item == nil {
			continue
		}
		details = append(details, &shopordermodels.OrderDetailSet{
			OID:            fmt.Sprintf("%s-%d", orderNo, index+1),
			EShopGoodsID:   fmt.Sprintf("%d", item.GoodsID),
			EShopGoodsName: item.GoodsName,
			EShopSkuID:     fmt.Sprintf("%d", item.SkuID),
			EShopSkuName:   item.SkuName,
			NumIID:         item.GoodsID,
			SkuID:          item.SkuID,
			Num:            float64(item.Quantity),
			Payment:        item.TotalAmount,
			PicPath:        item.ImageURL,
		})
	}

	return &shopordermodels.OrderSet{
		UserID: shopUser.ID,
		//PayTime:              time.Now().Format("2006-01-02 15:04:05"),
		Tid:                  orderNo,
		BuyerNick:            s.buildOrderBuyerNick(shopUser),
		BuyerMessage:         strings.TrimSpace(req.Remark),
		SellerMemo:           strings.TrimSpace(req.Remark),
		Total:                cacheData.GoodsAmount,
		Privilege:            cacheData.DiscountAmount,
		PostFee:              cacheData.FreightAmount,
		ReceiverName:         address.ReceiverName,
		ReceiverProvince:     address.ProvinceCode,
		ReceiverProvinceName: address.ProvinceName,
		ReceiverCity:         address.CityCode,
		ReceiverCityName:     address.CityName,
		ReceiverDistrict:     address.DistrictCode,
		ReceiverDistrictName: address.DistrictName,
		ReceiverAddress:      address.DetailAddress,
		ReceiverPhone:        address.ReceiverMobile,
		ReceiverMobile:       address.ReceiverMobile,
		Status:               orderConstant.ShopStatusToERPStatus(orderConstant.OrderStatusPending),
		OrderType:            shopConstant.NoCod,
		Details:              details,
		Accounts: []*shopordermodels.OrderAccountSet{
			{
				FinanceCode: "0168",
				Total:       cacheData.PayAmount,
			},
		},
	}
}

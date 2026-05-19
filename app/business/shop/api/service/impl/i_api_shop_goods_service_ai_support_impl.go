//go:build ai
// +build ai

package impl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	discountservice "nova-factory-server/app/business/shop/discount/service"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/utils/baizeContext"
)

// IApiShopGoodsServiceImpl 商品服务实现
type IApiShopGoodsServiceImpl struct {
	dao              dao.IApiShopGoodsDao
	shopGoodsService shopservice.IShopGoodsService
	discountService  discountservice.IDiscountCalculateService
	// service 读取模型
	service aidatasetservice.IAiModelProviderService
}

// NewIApiShopGoodsServiceImpl  创建商品服务
func NewIApiShopGoodsServiceImpl(dao dao.IApiShopGoodsDao,
	shopGoodsService shopservice.IShopGoodsService,
	discountService discountservice.IDiscountCalculateService, service aidatasetservice.IAiModelProviderService) service.IApiShopGoodsService {
	return &IApiShopGoodsServiceImpl{
		dao:              dao,
		shopGoodsService: shopGoodsService,
		discountService:  discountService,
		service:          service,
	}
}

// Search 按多个商品名称检索相似商品，并回填数据库中的最新商品数据
func (s *IApiShopGoodsServiceImpl) Search(c *gin.Context, req *models.GoodsSearchReq) (*models.GoodsSearchData, error) {
	if req == nil {
		return nil, errors.New("检索参数不能为空")
	}
	if len(req.GoodsNames) == 0 {
		return nil, errors.New("商品名称不能为空")
	}

	embeddingInfo, err := s.service.EmbeddingWithLLM(c)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return nil, err
	}

	if embeddingInfo == nil {
		return nil, errors.New("嵌入式模型没设置")
	}

	embeddingModel := shopmodels.EmbeddingConfig{
		ProviderType: embeddingInfo.APIType,
		ProviderID:   embeddingInfo.LLMFactory,
		APIEndpoint:  embeddingInfo.APIBase,
		ModelID:      embeddingInfo.LLMName,
		ApiKey:       embeddingInfo.APIKey,
	}

	limit := normalizeGoodsSearchLimit(req.Limit)
	vectorData, err := s.shopGoodsService.BatchSearchVector(c, &shopmodels.GoodsVectorBatchSearchReq{
		Queries:   req.GoodsNames,
		Limit:     buildGoodsVectorSearchLimit(limit),
		Embedding: &embeddingModel,
	})
	if err != nil {
		return nil, err
	}
	if vectorData == nil || len(vectorData.Rows) == 0 {
		return &models.GoodsSearchData{
			Rows:  make([]*models.GoodsSearchItem, 0),
			Total: 0,
		}, nil
	}

	goodsMap, err := s.loadSearchGoodsMap(c, vectorData.Rows)
	if err != nil {
		return nil, err
	}

	items := make([]*models.GoodsSearchItem, 0, len(vectorData.Rows))
	for _, row := range vectorData.Rows {
		if row == nil {
			continue
		}
		matches := make([]*models.GoodsSearchMatch, 0, limit)
		seen := make(map[int64]struct{}, len(row.Rows))
		for _, hit := range row.Rows {
			if hit == nil || hit.GoodsDBID == 0 {
				continue
			}
			if _, ok := seen[hit.GoodsDBID]; ok {
				continue
			}
			goods, ok := goodsMap[hit.GoodsDBID]
			if !ok || goods == nil {
				continue
			}
			seen[hit.GoodsDBID] = struct{}{}
			matches = append(matches, &models.GoodsSearchMatch{
				Score: hit.Score,
				Goods: goods,
			})
			if len(matches) >= limit {
				break
			}
		}
		items = append(items, &models.GoodsSearchItem{
			Query: row.Query,
			Rows:  matches,
			Total: int64(len(matches)),
		})
	}

	return &models.GoodsSearchData{
		Rows:  items,
		Total: int64(len(items)),
	}, nil
}

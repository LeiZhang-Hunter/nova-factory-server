package shopserviceimpl

import (
	"fmt"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/utils/fileUtils"
	embeddingutil "nova-factory-server/app/utils/llm/embedding"
	"strings"

	"github.com/gin-gonic/gin"
)

type ShopSkuServiceImpl struct {
	dao                  shopdao.IShopSkuDao
	goodsDao             shopdao.IShopGoodsDao
	categoryDao          shopdao.IShopCategoryDao
	vectorDao            shopdao.IShopGoodsVectorDao
	modelProviderService aidatasetservice.IAiModelProviderService
}

// NewShopSkuService 创建商品规格服务。
// 除了常规 SKU 读写 DAO 外，还额外注入了商品、分类、向量 DAO 和模型配置服务，
// 用于在规格发生变化后同步维护商品向量数据。
func NewShopSkuService(dao shopdao.IShopSkuDao, goodsDao shopdao.IShopGoodsDao,
	categoryDao shopdao.IShopCategoryDao, vectorDao shopdao.IShopGoodsVectorDao,
	modelProviderService aidatasetservice.IAiModelProviderService) shopservice.IShopSkuService {
	return &ShopSkuServiceImpl{
		dao:                  dao,
		goodsDao:             goodsDao,
		categoryDao:          categoryDao,
		vectorDao:            vectorDao,
		modelProviderService: modelProviderService,
	}
}

// Create 创建单条商品规格。
// 这里会把“保存 SKU”与“重建并写回商品向量”放到同一个数据库事务流程里：
// 1. 先写入关系库；
// 2. 再同步 Milvus 向量；
// 3. 如果向量同步失败，则返回错误并回滚本次数据库写入。
func (s *ShopSkuServiceImpl) Create(c *gin.Context, req *shopmodels.GoodsSkuUpsert) (*shopmodels.GoodsSku, error) {
	var (
		data *shopmodels.GoodsSku
		err  error
	)
	err = s.dao.Transaction(c, func(txDao shopdao.IShopSkuDao) error {
		data, err = txDao.Create(c, req)
		if err != nil {
			return err
		}
		return s.syncGoodsVectorAfterSkuChange(c, txDao, data, req)
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Update 更新单条商品规格。
// 与 Create 一样，更新成功后会立即重建该商品的向量数据；
// 若向量同步失败，则当前事务整体回滚，避免数据库与向量库状态不一致。
func (s *ShopSkuServiceImpl) Update(c *gin.Context, req *shopmodels.GoodsSkuUpsert) (*shopmodels.GoodsSku, error) {
	var (
		data *shopmodels.GoodsSku
		err  error
	)
	err = s.dao.Transaction(c, func(txDao shopdao.IShopSkuDao) error {
		data, err = txDao.Update(c, req)
		if err != nil {
			return err
		}
		return s.syncGoodsVectorAfterSkuChange(c, txDao, data, req)
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteByIDs 批量删除商品规格。
// 删除流程同样走事务控制：
// 1. 先删除关系库中的 SKU 记录；
// 2. 再按同一批 SKU 主键删除 Milvus 向量行；
// 3. 若向量删除失败，则返回错误并回滚本次数据库删除。
func (s *ShopSkuServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return s.dao.Transaction(c, func(txDao shopdao.IShopSkuDao) error {
		if err := txDao.DeleteByIDs(c, ids); err != nil {
			return err
		}
		if s.vectorDao == nil {
			return nil
		}
		if err := s.vectorDao.DeleteBySkuIDs(c, ids); err != nil {
			return fmt.Errorf("同步删除商品向量失败: %w", err)
		}
		return nil
	})
}

// GetByID 根据主键查询单条商品规格。
func (s *ShopSkuServiceImpl) GetByID(c *gin.Context, id int64) (*shopmodels.GoodsSku, error) {
	return s.dao.GetByID(c, id)
}

// List 查询商品规格列表，并补全图片绝对地址，方便接口层直接回传给前端。
func (s *ShopSkuServiceImpl) List(c *gin.Context, req *shopmodels.GoodsSkuQuery) (*shopmodels.GoodsSkuListData, error) {
	ret, err := s.dao.List(c, req)
	if err != nil {
		return nil, err
	}

	if ret == nil || len(ret.Rows) == 0 {
		return ret, nil
	}

	for k, v := range ret.Rows {
		ret.Rows[k].ImageURL = fileUtils.BuildAbsoluteURL(c, v.ImageURL)
	}
	return ret, nil
}

// syncGoodsVectorAfterSkuChange 在 SKU 新增/修改后重建对应商品的向量。
// 当前实现不是只更新单条 SKU 向量，而是按商品整体重建：
// 1. 先回查商品；
// 2. 再挂载该商品下全部 SKU 和分类；
// 3. 使用 embedding 配置生成最新向量；
// 4. 最终通过 vectorDao.Upsert 覆盖写回 Milvus。
// 这样可以避免只更新局部规格后，商品整体检索文本、metadata 与上下架状态出现不一致。
func (s *ShopSkuServiceImpl) syncGoodsVectorAfterSkuChange(c *gin.Context, skuDao shopdao.IShopSkuDao,
	sku *shopmodels.GoodsSku, req *shopmodels.GoodsSkuUpsert) error {
	if s.goodsDao == nil || s.categoryDao == nil || s.vectorDao == nil {
		return nil
	}

	goodsID := ""
	if sku != nil {
		goodsID = strings.TrimSpace(sku.GoodsID)
	}
	if goodsID == "" && req != nil {
		goodsID = strings.TrimSpace(req.GoodsID)
	}
	if goodsID == "" {
		return nil
	}

	goods, err := s.goodsDao.GetByGoodsID(c, goodsID)
	if err != nil {
		return fmt.Errorf("查询关联商品失败: %w", err)
	}
	if goods == nil {
		return fmt.Errorf("关联商品不存在，goodsId=%s", goodsID)
	}

	helper := &ShopGoodsServiceImpl{
		skuDao:      skuDao,
		categoryDao: s.categoryDao,
		vectorDao:   s.vectorDao,
	}
	if err = helper.attachSkus(c, []*shopmodels.Goods{goods}); err != nil {
		return fmt.Errorf("加载商品规格失败: %w", err)
	}
	if err = helper.attachCategoryNames(c, []*shopmodels.Goods{goods}); err != nil {
		return fmt.Errorf("加载商品分类失败: %w", err)
	}

	cfg, err := loadEmbeddingProviderConfig(s.loadGoodsVectorEmbeddingConfig(c))
	if err != nil {
		return err
	}
	requestCtx := buildRequestContext(c)
	embedder, err := embeddingutil.NewEmbedder(requestCtx, cfg)
	if err != nil {
		return fmt.Errorf("初始化向量模型失败: %w", err)
	}
	if _, err = helper.generateGoodsVectorWithEmbedder(c, requestCtx, embedder, goods); err != nil {
		return fmt.Errorf("同步商品向量失败: %w", err)
	}
	return nil
}

// loadGoodsVectorEmbeddingConfig 从当前用户可用的 embedding 模型配置中提取向量模型参数。
// 这里优先走模型配置服务返回的 SysUserLLM，避免在 SKU service 里再额外维护一套静态配置读取逻辑。
func (s *ShopSkuServiceImpl) loadGoodsVectorEmbeddingConfig(c *gin.Context) *shopmodels.EmbeddingConfig {
	if s == nil || s.modelProviderService == nil {
		return nil
	}
	info, err := s.modelProviderService.EmbeddingWithLLM(c)
	if err != nil || info == nil {
		return nil
	}
	return &shopmodels.EmbeddingConfig{
		ProviderType: strings.TrimSpace(info.APIType),
		ProviderID:   strings.TrimSpace(info.LLMFactory),
		APIEndpoint:  strings.TrimSpace(info.APIBase),
		ModelID:      strings.TrimSpace(info.LLMName),
		ApiKey:       strings.TrimSpace(info.APIKey),
	}
}

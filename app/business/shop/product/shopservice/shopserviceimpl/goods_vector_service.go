package shopserviceimpl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudwego/eino/components/embedding"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/datasource/cache/cacheError"
	"nova-factory-server/app/utils/baizeContext"
	embeddingutil "nova-factory-server/app/utils/llm/embedding"
	"nova-factory-server/app/utils/snowflake"
)

const goodsVectorContentMaxLength = 16384
const (
	defaultGoodsVectorSearchLimit = 10
	maxGoodsVectorSearchLimit     = 50
	defaultGoodsVectorBatchSize   = 20
	maxGoodsVectorBatchSize       = 100
	goodsVectorTaskTTL            = 0 * time.Hour
	goodsVectorTaskStatusPending  = "pending"
	goodsVectorTaskStatusRunning  = "running"
	goodsVectorTaskStatusDone     = "completed"
	goodsVectorTaskStatusFailed   = "failed"
)

type goodsVectorConfig struct {
	Embedding embeddingutil.ProviderConfig `mapstructure:"embedding"`
}

type goodsEmbeddingPayload struct {
	sku     *shopmodels.GoodsSku
	content string
}

func (s *ShopGoodsServiceImpl) GenerateVector(c *gin.Context, req *shopmodels.GenVectorReq) (*shopmodels.GoodsVectorResult, error) {
	goods, err := s.dao.GetByID(c, req.ID)
	if err != nil {
		return nil, fmt.Errorf("查询商品失败: %w", err)
	}
	if goods == nil {
		return nil, errors.New("商品不存在")
	}
	if goods.IsOnSale <= 0 {
		return nil, errors.New("商品未上架，不能生成向量")
	}

	if err = s.attachSkus(c, []*shopmodels.Goods{goods}); err != nil {
		return nil, err
	}

	if err = s.attachCategoryNames(c, []*shopmodels.Goods{goods}); err != nil {
		return nil, err
	}

	cfg, err := s.loadGoodsVectorConfig(req)
	if err != nil {
		return nil, err
	}

	requestCtx := buildRequestContext(c)
	embedder, err := embeddingutil.NewEmbedder(requestCtx, &cfg.Embedding)
	if err != nil {
		return nil, fmt.Errorf("初始化向量模型失败: %w", err)
	}

	result, err := s.generateGoodsVectorWithEmbedder(c, requestCtx, embedder, goods)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("未生成有效的商品向量数据")
	}
	return result, nil
}

func (s *ShopGoodsServiceImpl) GenerateAllOnSaleVectors(c *gin.Context,
	req *shopmodels.GenAllVectorReq) (*shopmodels.GoodsVectorTaskData, error) {
	if req == nil {
		return nil, errors.New("全量生成参数不能为空")
	}
	if _, err := loadEmbeddingProviderConfig(req.Embedding); err != nil {
		return nil, err
	}

	operatorID := baizeContext.GetUserIdSafe(c)
	operator := baizeContext.GetUserName(c)
	batchSize := resolveGoodsVectorTaskBatchSize(req.BatchSize, 0)
	taskID := strings.TrimSpace(req.TaskID)
	if s.cache == nil {
		return nil, errors.New("缓存未初始化")
	}

	var (
		progress *shopmodels.GoodsVectorTaskProgress
		err      error
	)
	if taskID != "" {
		progress, err = s.loadGoodsVectorTaskProgress(taskID)
		if err != nil {
			if errors.Is(err, cacheError.Nil) {
				return nil, errors.New("续跑任务不存在或已过期")
			}
			return nil, err
		}
		if progress != nil && progress.OperatorID > 0 && operatorID > 0 && progress.OperatorID != operatorID {
			return nil, errors.New("该任务不属于当前操作人")
		}
		if progress != nil && !isGoodsVectorTaskUnfinished(progress) {
			return nil, errors.New("该任务已执行完成，请重新创建新任务")
		}
		if progress != nil && progress.Processed > 0 {
			batchSize = resolveGoodsVectorTaskBatchSize(0, progress.BatchSize)
		}
	} else {
		progress, err = s.loadLatestUnfinishedGoodsVectorTask(operatorID)
		if err != nil {
			return nil, err
		}
		if progress != nil {
			taskID = progress.TaskID
			if progress.Processed > 0 {
				batchSize = resolveGoodsVectorTaskBatchSize(0, progress.BatchSize)
			} else {
				batchSize = resolveGoodsVectorTaskBatchSize(req.BatchSize, progress.BatchSize)
			}
		}
	}

	if progress == nil {
		taskID = strconv.FormatInt(snowflake.GenID(), 10)
		now := time.Now().Format(time.DateTime)
		progress = &shopmodels.GoodsVectorTaskProgress{
			TaskID:     taskID,
			Status:     goodsVectorTaskStatusPending,
			BatchSize:  batchSize,
			Message:    "任务已创建，等待执行",
			OperatorID: operatorID,
			Operator:   operator,
			StartedAt:  now,
			UpdatedAt:  now,
		}
	} else {
		progress.Status = goodsVectorTaskStatusPending
		progress.BatchSize = batchSize
		progress.OperatorID = operatorID
		progress.Operator = operator
		progress.UpdatedAt = time.Now().Format(time.DateTime)
		progress.FinishedAt = ""
		progress.Progress = calcGoodsVectorTaskProgress(progress.Processed, progress.Total)
		if progress.StartedAt == "" {
			progress.StartedAt = progress.UpdatedAt
		}
		if progress.Processed > 0 {
			progress.Message = fmt.Sprintf("检测到未完成任务，准备从第%d条继续执行", progress.Processed+1)
		} else {
			progress.Message = "检测到未完成任务，准备继续执行"
		}
	}

	lockKey := goodsVectorTaskLockKey(taskID)
	if !s.cache.SetNX(context.Background(), lockKey, strconv.FormatInt(operatorID, 10), goodsVectorTaskTTL) {
		return &shopmodels.GoodsVectorTaskData{
			TaskID:      taskID,
			ProgressKey: goodsVectorTaskCacheKey(taskID),
			Status:      goodsVectorTaskStatusRunning,
		}, nil
	}
	if err = s.saveGoodsVectorTaskProgress(taskID, progress); err != nil {
		s.cache.Del(context.Background(), lockKey)
		return nil, err
	}
	s.saveLatestGoodsVectorTask(operatorID, taskID)

	copiedReq := *req
	copiedReq.TaskID = taskID
	copiedReq.BatchSize = batchSize
	backgroundCtx := cloneGinContext(c)
	go s.runGenerateAllOnSaleVectors(backgroundCtx, taskID, &copiedReq, progress.OperatorID, progress.Operator)

	return &shopmodels.GoodsVectorTaskData{
		TaskID:      taskID,
		ProgressKey: goodsVectorTaskCacheKey(taskID),
		Status:      progress.Status,
	}, nil
}

func (s *ShopGoodsServiceImpl) GetGenerateAllOnSaleVectorsProgress(c *gin.Context,
	taskID string) (*shopmodels.GoodsVectorTaskProgress, error) {
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, errors.New("任务ID不能为空")
	}
	progress, err := s.loadGoodsVectorTaskProgress(taskID)
	if err != nil {
		if errors.Is(err, cacheError.Nil) {
			return nil, errors.New("任务不存在或已过期")
		}
		return nil, err
	}
	return progress, nil
}

func (s *ShopGoodsServiceImpl) ListGenerateAllOnSaleVectorTasks(c *gin.Context) (*shopmodels.GoodsVectorTaskListData, error) {
	if s.cache == nil {
		return nil, errors.New("缓存未初始化")
	}

	operatorID := baizeContext.GetUserIdSafe(c)
	keys, err := s.scanGoodsVectorTaskKeys()
	if err != nil {
		return nil, err
	}

	rows := make([]*shopmodels.GoodsVectorTaskProgress, 0, len(keys))
	for _, key := range keys {
		taskID := strings.TrimPrefix(key, goodsVectorTaskCacheKey(""))
		taskID = strings.TrimSpace(taskID)
		if taskID == "" {
			continue
		}

		progress, loadErr := s.loadGoodsVectorTaskProgress(taskID)
		if loadErr != nil {
			if errors.Is(loadErr, cacheError.Nil) {
				continue
			}
			return nil, loadErr
		}
		if progress == nil || !isGoodsVectorTaskUnfinished(progress) {
			continue
		}
		if operatorID > 0 && progress.OperatorID != operatorID {
			continue
		}
		rows = append(rows, progress)
	}

	sort.Slice(rows, func(i, j int) bool {
		return goodsVectorTaskSortTime(rows[i]).After(goodsVectorTaskSortTime(rows[j]))
	})

	return &shopmodels.GoodsVectorTaskListData{
		Rows:  rows,
		Total: int64(len(rows)),
	}, nil
}

func (s *ShopGoodsServiceImpl) SearchVector(c *gin.Context, req *shopmodels.GoodsVectorSearchReq) (*shopmodels.GoodsVectorSearchData, error) {
	if req == nil {
		return nil, errors.New("搜索参数不能为空")
	}
	queries, err := normalizeGoodsVectorSearchQueries([]string{req.Query})
	if err != nil {
		return nil, err
	}
	items, err := s.batchSearchVector(c, queries, req.Limit, req.Embedding)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return &shopmodels.GoodsVectorSearchData{
			Rows:  make([]*shopmodels.GoodsVectorSearchItem, 0),
			Total: 0,
		}, nil
	}
	return &shopmodels.GoodsVectorSearchData{
		Rows:  items[0].Rows,
		Total: items[0].Total,
	}, nil
}

func (s *ShopGoodsServiceImpl) BatchSearchVector(c *gin.Context,
	req *shopmodels.GoodsVectorBatchSearchReq) (*shopmodels.GoodsVectorBatchSearchData, error) {
	if req == nil {
		return nil, errors.New("批量搜索参数不能为空")
	}
	queries, err := normalizeGoodsVectorSearchQueries(req.Queries)
	if err != nil {
		return nil, err
	}
	rows, err := s.batchSearchVector(c, queries, req.Limit, req.Embedding)
	if err != nil {
		return nil, err
	}
	return &shopmodels.GoodsVectorBatchSearchData{
		Rows:  rows,
		Total: int64(len(rows)),
	}, nil
}

func (s *ShopGoodsServiceImpl) batchSearchVector(c *gin.Context, queries []string, limit int,
	embedding *shopmodels.EmbeddingConfig) ([]*shopmodels.GoodsVectorBatchSearchItem, error) {
	cfg, err := loadEmbeddingProviderConfig(embedding)
	if err != nil {
		return nil, err
	}

	requestCtx := buildRequestContext(c)
	embedder, err := embeddingutil.NewEmbedder(requestCtx, cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化向量模型失败: %w", err)
	}

	vectors, err := embedder.EmbedStrings(requestCtx, queries)
	if err != nil {
		return nil, fmt.Errorf("生成搜索向量失败: %w", err)
	}
	if len(vectors) != len(queries) {
		return nil, fmt.Errorf("向量模型返回结果数量不匹配，expected=%d actual=%d", len(queries), len(vectors))
	}

	normalizedVectors := make([][]float32, 0, len(vectors))
	for idx := range queries {
		if len(vectors[idx]) == 0 {
			return nil, fmt.Errorf("第%d条搜索向量为空", idx+1)
		}
		normalizedVectors = append(normalizedVectors, float64ToFloat32(vectors[idx]))
	}
	data, err := s.vectorDao.BatchSearch(c, &shopmodels.GoodsVectorBatchSearchReq{
		Queries: queries,
		Limit:   normalizeGoodsVectorSearchLimit(limit),
	}, normalizedVectors)
	if err != nil {
		return nil, err
	}
	if data == nil || data.Rows == nil {
		return make([]*shopmodels.GoodsVectorBatchSearchItem, 0), nil
	}
	return data.Rows, nil
}

func (s *ShopGoodsServiceImpl) loadGoodsVectorConfig(req *shopmodels.GenVectorReq) (*goodsVectorConfig, error) {
	cfg := &goodsVectorConfig{}
	embeddingCfg, err := loadEmbeddingProviderConfig(req.Embedding)
	if err != nil {
		return nil, err
	}
	cfg.Embedding = *embeddingCfg
	return cfg, nil
}

func (s *ShopGoodsServiceImpl) generateGoodsVectorWithEmbedder(c *gin.Context, requestCtx context.Context,
	embedder embedding.Embedder, goods *shopmodels.Goods) (*shopmodels.GoodsVectorResult, error) {
	payloads := buildGoodsEmbeddingPayloads(goods)
	if len(payloads) == 0 {
		return nil, errors.New("商品内容为空，无法生成向量")
	}

	texts := make([]string, 0, len(payloads))
	for _, payload := range payloads {
		texts = append(texts, payload.content)
	}

	vectors, err := embedder.EmbedStrings(requestCtx, texts)
	if err != nil {
		return nil, fmt.Errorf("生成商品向量失败: %w", err)
	}
	if len(vectors) != len(payloads) {
		return nil, fmt.Errorf("向量模型返回结果数量不匹配，expected=%d actual=%d", len(payloads), len(vectors))
	}
	if len(vectors) == 0 || len(vectors[0]) == 0 {
		return nil, errors.New("向量模型未返回有效结果")
	}

	items := make([]*shopmodels.GoodsVectorUpsertItem, 0, len(payloads))
	for idx, payload := range payloads {
		if len(vectors[idx]) == 0 {
			return nil, fmt.Errorf("第%d条SKU向量为空", idx+1)
		}
		item := &shopmodels.GoodsVectorUpsertItem{
			Content: payload.content,
			Vector:  float64ToFloat32(vectors[idx]),
		}
		if payload.sku != nil {
			item.SkuID = int64(payload.sku.ID)
			item.SkuName = payload.sku.SkuName
			item.SkuDescription = payload.sku.Description
			item.RetailPrice = payload.sku.RetailPrice
			item.Weight = payload.sku.Weight
			item.Quantity = payload.sku.Quantity
		}
		items = append(items, item)
	}
	result, err := s.vectorDao.Upsert(c, goods, items)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("未生成有效的商品向量数据")
	}
	return result, nil
}

func (s *ShopGoodsServiceImpl) runGenerateAllOnSaleVectors(c *gin.Context, taskID string,
	req *shopmodels.GenAllVectorReq, operatorID int64, operator string) {
	progress, err := s.loadGoodsVectorTaskProgress(taskID)
	if err != nil {
		if !errors.Is(err, cacheError.Nil) {
			zap.L().Error("load goods vector task progress fail", zap.String("taskId", taskID), zap.Error(err))
		}
		progress = nil
	}
	now := time.Now().Format(time.DateTime)
	if progress == nil {
		progress = &shopmodels.GoodsVectorTaskProgress{
			TaskID:     taskID,
			BatchSize:  resolveGoodsVectorTaskBatchSize(req.BatchSize, 0),
			OperatorID: operatorID,
			Operator:   operator,
			StartedAt:  now,
		}
	} else {
		progress.BatchSize = resolveGoodsVectorTaskBatchSize(req.BatchSize, progress.BatchSize)
		progress.OperatorID = operatorID
		progress.Operator = operator
		if progress.StartedAt == "" {
			progress.StartedAt = now
		}
	}
	if progress.Processed > 0 {
		progress.Message = fmt.Sprintf("正在继续生成商品向量，已处理%d条", progress.Processed)
	} else {
		progress.Message = "正在初始化向量模型"
	}
	progress.Status = goodsVectorTaskStatusRunning
	progress.UpdatedAt = now
	if err := s.saveGoodsVectorTaskProgress(taskID, progress); err != nil {
		zap.L().Error("save goods vector task progress fail", zap.String("taskId", taskID), zap.Error(err))
		s.cache.Del(context.Background(), goodsVectorTaskLockKey(taskID))
		return
	}

	defer func() {
		s.cache.Del(context.Background(), goodsVectorTaskLockKey(taskID))
		if r := recover(); r != nil {
			progress.Status = goodsVectorTaskStatusFailed
			progress.Message = fmt.Sprintf("任务执行异常: %v", r)
			progress.UpdatedAt = time.Now().Format(time.DateTime)
			progress.FinishedAt = progress.UpdatedAt
			if err := s.saveGoodsVectorTaskProgress(taskID, progress); err != nil {
				zap.L().Error("save goods vector task progress after panic fail", zap.String("taskId", taskID), zap.Error(err))
			}
		}
	}()

	cfg, err := loadEmbeddingProviderConfig(req.Embedding)
	if err != nil {
		s.failGoodsVectorTask(taskID, progress, err)
		return
	}
	requestCtx := buildRequestContext(c)
	embedder, err := embeddingutil.NewEmbedder(requestCtx, cfg)
	if err != nil {
		s.failGoodsVectorTask(taskID, progress, fmt.Errorf("初始化向量模型失败: %w", err))
		return
	}

	size := int64(resolveGoodsVectorTaskBatchSize(req.BatchSize, progress.BatchSize))
	page := int64(1)
	resumeSkip := 0
	resumeCurrentID := progress.CurrentID
	needResumeTrim := progress.Processed > 0 && size > 0
	if needResumeTrim {
		page = progress.Processed/size + 1
		resumeSkip = int(progress.Processed % size)
	}
	isOnSale := true

	for {
		listReq := &shopmodels.GoodsQuery{
			IsOnSale: &isOnSale,
			Page:     page,
			Size:     size,
		}
		data, listErr := s.List(c, listReq)
		if listErr != nil {
			s.failGoodsVectorTask(taskID, progress, fmt.Errorf("查询上架商品失败: %w", listErr))
			return
		}
		if data == nil {
			data = &shopmodels.GoodsListData{}
		}
		if progress.Total == 0 {
			progress.Total = data.Total
			if data.Total == 0 {
				progress.Status = goodsVectorTaskStatusDone
				progress.Progress = 100
				progress.Message = "没有可生成向量的上架商品"
				progress.UpdatedAt = time.Now().Format(time.DateTime)
				progress.FinishedAt = progress.UpdatedAt
				_ = s.saveGoodsVectorTaskProgress(taskID, progress)
				return
			}
		}
		rows := data.Rows
		if needResumeTrim {
			rows = trimGoodsVectorResumeRows(rows, resumeCurrentID, resumeSkip)
			needResumeTrim = false
			if len(rows) == 0 && len(data.Rows) > 0 {
				page++
				continue
			}
		}

		for _, goods := range rows {
			if goods == nil {
				continue
			}
			progress.CurrentID = goods.ID
			progress.CurrentName = goods.GoodsName
			progress.Message = "正在生成商品向量"
			progress.UpdatedAt = time.Now().Format(time.DateTime)
			progress.Progress = calcGoodsVectorTaskProgress(progress.Processed, progress.Total)
			_ = s.saveGoodsVectorTaskProgress(taskID, progress)

			_, generateErr := s.generateGoodsVectorWithEmbedder(c, requestCtx, embedder, goods)
			progress.Processed++
			if generateErr != nil {
				progress.Failed++
				progress.Message = fmt.Sprintf("商品[%d-%s]生成失败: %v", goods.ID, goods.GoodsName, generateErr)
				zap.L().Error("generate goods vector in batch fail",
					zap.String("taskId", taskID),
					zap.Int64("goodsId", goods.ID),
					zap.String("goodsName", goods.GoodsName),
					zap.Error(generateErr))
			} else {
				progress.Success++
				progress.Message = fmt.Sprintf("商品[%d-%s]生成成功", goods.ID, goods.GoodsName)
			}
			progress.Progress = calcGoodsVectorTaskProgress(progress.Processed, progress.Total)
			progress.UpdatedAt = time.Now().Format(time.DateTime)
			_ = s.saveGoodsVectorTaskProgress(taskID, progress)
		}

		if len(data.Rows) == 0 || progress.Processed >= progress.Total {
			break
		}
		page++
	}

	progress.Status = goodsVectorTaskStatusDone
	progress.Progress = 100
	progress.Message = fmt.Sprintf("任务完成，成功%d，失败%d", progress.Success, progress.Failed)
	progress.UpdatedAt = time.Now().Format(time.DateTime)
	progress.FinishedAt = progress.UpdatedAt
	if err = s.saveGoodsVectorTaskProgress(taskID, progress); err != nil {
		zap.L().Error("save completed goods vector task progress fail", zap.String("taskId", taskID), zap.Error(err))
	}
}

func (s *ShopGoodsServiceImpl) failGoodsVectorTask(taskID string, progress *shopmodels.GoodsVectorTaskProgress, err error) {
	if progress == nil {
		return
	}
	progress.Status = goodsVectorTaskStatusFailed
	progress.Message = err.Error()
	progress.UpdatedAt = time.Now().Format(time.DateTime)
	progress.FinishedAt = progress.UpdatedAt
	if saveErr := s.saveGoodsVectorTaskProgress(taskID, progress); saveErr != nil {
		zap.L().Error("save failed goods vector task progress fail",
			zap.String("taskId", taskID),
			zap.Error(saveErr))
	}
}

func (s *ShopGoodsServiceImpl) saveGoodsVectorTaskProgress(taskID string,
	progress *shopmodels.GoodsVectorTaskProgress) error {
	if s.cache == nil {
		return errors.New("缓存未初始化")
	}
	if progress == nil {
		return errors.New("任务进度不能为空")
	}
	if progress.TaskID == "" {
		progress.TaskID = taskID
	}
	if progress.UpdatedAt == "" {
		progress.UpdatedAt = time.Now().Format(time.DateTime)
	}
	body, err := json.Marshal(progress)
	if err != nil {
		return err
	}
	s.cache.Set(context.Background(), goodsVectorTaskCacheKey(taskID), string(body), goodsVectorTaskTTL)
	return nil
}

func (s *ShopGoodsServiceImpl) loadGoodsVectorTaskProgress(taskID string) (*shopmodels.GoodsVectorTaskProgress, error) {
	if s.cache == nil {
		return nil, errors.New("缓存未初始化")
	}
	val, err := s.cache.Get(context.Background(), goodsVectorTaskCacheKey(taskID))
	if err != nil {
		return nil, err
	}
	progress := &shopmodels.GoodsVectorTaskProgress{}
	if err = json.Unmarshal([]byte(val), progress); err != nil {
		return nil, err
	}
	return progress, nil
}

func loadEmbeddingProviderConfig(req *shopmodels.EmbeddingConfig) (*embeddingutil.ProviderConfig, error) {
	if req == nil {
		return nil, errors.New("未配置 embedding.model_id 或 shop.goods_vector.embedding.model_id")
	}
	cfg := &embeddingutil.ProviderConfig{
		ProviderType: req.ProviderType,
		ProviderID:   req.ProviderID,
		APIEndpoint:  req.APIEndpoint,
		ModelID:      req.ModelID,
		APIKey:       req.ApiKey,
	}
	if cfg.ProviderType == "" {
		cfg.ProviderType = "openai"
	}
	if cfg.ProviderID == "" {
		cfg.ProviderID = cfg.ProviderType
	}
	if cfg.ModelID == "" {
		return nil, errors.New("未配置 embedding.model_id 或 shop.goods_vector.embedding.model_id")
	}
	return cfg, nil
}

func buildRequestContext(c *gin.Context) context.Context {
	if c != nil && c.Request != nil {
		return c.Request.Context()
	}
	return context.Background()
}

func float64ToFloat32(vector []float64) []float32 {
	result := make([]float32, 0, len(vector))
	for _, value := range vector {
		result = append(result, float32(value))
	}
	return result
}

func trimRunes(value string, max int) string {
	value = strings.TrimSpace(value)
	if max <= 0 {
		return value
	}
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return string(runes[:max])
}

func normalizeGoodsVectorSearchQueries(queries []string) ([]string, error) {
	result := make([]string, 0, len(queries))
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		result = append(result, query)
	}
	if len(result) == 0 {
		return nil, errors.New("搜索内容不能为空")
	}
	return result, nil
}

func cloneGinContext(c *gin.Context) *gin.Context {
	if c == nil {
		return nil
	}
	clone := c.Copy()
	if c.Request != nil {
		clone.Request = c.Request.Clone(context.Background())
	}
	return clone
}

func goodsVectorTaskCacheKey(taskID string) string {
	return "shop:goods:vector:task:" + strings.TrimSpace(taskID)
}

func goodsVectorTaskLockKey(taskID string) string {
	return goodsVectorTaskCacheKey(taskID) + ":lock"
}

func goodsVectorLatestTaskKey(operatorID int64) string {
	return "shop:goods:vector:latest:" + strconv.FormatInt(operatorID, 10)
}

func calcGoodsVectorTaskProgress(processed, total int64) int {
	if total <= 0 {
		return 0
	}
	if processed <= 0 {
		return 0
	}
	if processed >= total {
		return 100
	}
	return int(processed * 100 / total)
}

func normalizeGoodsVectorBatchSize(batchSize int) int {
	if batchSize <= 0 {
		return defaultGoodsVectorBatchSize
	}
	if batchSize > maxGoodsVectorBatchSize {
		return maxGoodsVectorBatchSize
	}
	return batchSize
}

func resolveGoodsVectorTaskBatchSize(reqBatchSize, progressBatchSize int) int {
	if reqBatchSize > 0 {
		return normalizeGoodsVectorBatchSize(reqBatchSize)
	}
	if progressBatchSize > 0 {
		return normalizeGoodsVectorBatchSize(progressBatchSize)
	}
	return defaultGoodsVectorBatchSize
}

func isGoodsVectorTaskUnfinished(progress *shopmodels.GoodsVectorTaskProgress) bool {
	if progress == nil {
		return false
	}
	return progress.Status != goodsVectorTaskStatusDone
}

func trimGoodsVectorResumeRows(rows []*shopmodels.Goods, currentID int64, skipCount int) []*shopmodels.Goods {
	if len(rows) == 0 {
		return rows
	}
	if currentID > 0 {
		for index, goods := range rows {
			if goods != nil && goods.ID == currentID {
				if index+1 >= len(rows) {
					return rows[:0]
				}
				return rows[index+1:]
			}
		}
	}
	if skipCount <= 0 {
		return rows
	}
	if skipCount >= len(rows) {
		return rows[:0]
	}
	return rows[skipCount:]
}

func (s *ShopGoodsServiceImpl) loadLatestUnfinishedGoodsVectorTask(operatorID int64) (*shopmodels.GoodsVectorTaskProgress, error) {
	if operatorID <= 0 || s.cache == nil {
		return nil, nil
	}
	taskID, err := s.cache.Get(context.Background(), goodsVectorLatestTaskKey(operatorID))
	if err != nil {
		if errors.Is(err, cacheError.Nil) {
			return nil, nil
		}
		return nil, err
	}
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, nil
	}
	progress, err := s.loadGoodsVectorTaskProgress(taskID)
	if err != nil {
		if errors.Is(err, cacheError.Nil) {
			s.cache.Del(context.Background(), goodsVectorLatestTaskKey(operatorID))
			return nil, nil
		}
		return nil, err
	}
	if !isGoodsVectorTaskUnfinished(progress) {
		return nil, nil
	}
	return progress, nil
}

func (s *ShopGoodsServiceImpl) saveLatestGoodsVectorTask(operatorID int64, taskID string) {
	if operatorID <= 0 || s.cache == nil {
		return
	}
	s.cache.Set(context.Background(), goodsVectorLatestTaskKey(operatorID), strings.TrimSpace(taskID), goodsVectorTaskTTL)
}

func (s *ShopGoodsServiceImpl) scanGoodsVectorTaskKeys() ([]string, error) {
	if s.cache == nil {
		return nil, errors.New("缓存未初始化")
	}

	match := goodsVectorTaskCacheKey("") + "*"
	keys := make([]string, 0)
	seen := make(map[string]struct{})
	var cursor uint64

	for {
		scannedKeys, nextCursor := s.cache.Scan(context.Background(), cursor, match, 100)
		for _, key := range scannedKeys {
			if strings.HasSuffix(key, ":lock") {
				continue
			}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			keys = append(keys, key)
		}
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	return keys, nil
}

func goodsVectorTaskSortTime(progress *shopmodels.GoodsVectorTaskProgress) time.Time {
	if progress == nil {
		return time.Time{}
	}
	if ts := parseGoodsVectorTaskTime(progress.UpdatedAt); !ts.IsZero() {
		return ts
	}
	if ts := parseGoodsVectorTaskTime(progress.StartedAt); !ts.IsZero() {
		return ts
	}
	return time.Time{}
}

func parseGoodsVectorTaskTime(value string) time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}
	}
	ts, err := time.ParseInLocation(time.DateTime, value, time.Local)
	if err != nil {
		return time.Time{}
	}
	return ts
}

func buildGoodsEmbeddingPayloads(goods *shopmodels.Goods) []goodsEmbeddingPayload {
	if goods == nil {
		return nil
	}

	baseLines := make([]string, 0, 9)
	appendLine := func(label, value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		baseLines = append(baseLines, label+": "+value)
	}

	appendLine("商品名称", goods.GoodsName)
	appendLine("商品分类", goods.ShopCategoryName)
	appendLine("商品描述", goods.Description)
	//if goods.RetailPrice > 0 {
	//	baseLines = append(baseLines, fmt.Sprintf("零售价: %.2f", goods.RetailPrice))
	//}
	//if goods.Weight > 0 {
	//	baseLines = append(baseLines, fmt.Sprintf("重量: %.2f%s", goods.Weight, strings.TrimSpace(goods.WeightUnit)))
	//}
	if goods.Unit != "" {
		baseLines = append(baseLines, "销售单位: "+strings.TrimSpace(goods.Unit))
	}
	//baseLines = append(baseLines, fmt.Sprintf("库存: %d", goods.Quantity))
	//if goods.IsOnSale > 0 {
	//	baseLines = append(baseLines, "上架状态: 上架")
	//} else {
	//	baseLines = append(baseLines, "上架状态: 下架")
	//}

	payloads := make([]goodsEmbeddingPayload, 0, maxInt(len(goods.Skus), 1))
	for _, sku := range goods.Skus {
		if sku == nil {
			continue
		}
		lines := append([]string{}, baseLines...)
		appendSkuLines(&lines, sku)
		content := trimRunes(strings.Join(lines, "\n"), goodsVectorContentMaxLength)
		if strings.TrimSpace(content) == "" {
			continue
		}
		payloads = append(payloads, goodsEmbeddingPayload{
			sku:     sku,
			content: content,
		})
	}
	if len(payloads) > 0 {
		return payloads
	}
	content := trimRunes(strings.Join(baseLines, "\n"), goodsVectorContentMaxLength)
	if strings.TrimSpace(content) == "" {
		return nil
	}
	return []goodsEmbeddingPayload{{content: content}}
}

func appendSkuLines(lines *[]string, sku *shopmodels.GoodsSku) {
	if sku == nil {
		return
	}
	appendLine := func(label, value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		*lines = append(*lines, label+": "+value)
	}
	appendLine("SKU名称", sku.SkuName)
	appendLine("SKU描述", sku.Description)
	if sku.Unit != "" {
		*lines = append(*lines, "SKU单位: "+strings.TrimSpace(sku.Unit))
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func normalizeGoodsVectorSearchLimit(limit int) int {
	if limit <= 0 {
		return defaultGoodsVectorSearchLimit
	}
	if limit > maxGoodsVectorSearchLimit {
		return maxGoodsVectorSearchLimit
	}
	return limit
}

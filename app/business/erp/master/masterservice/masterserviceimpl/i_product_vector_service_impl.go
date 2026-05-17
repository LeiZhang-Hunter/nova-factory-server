package masterserviceimpl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/datasource/cache/cacheError"
	"nova-factory-server/app/utils/baizeContext"
	embeddingutil "nova-factory-server/app/utils/llm/embedding"
	"nova-factory-server/app/utils/snowflake"
)

const (
	productVectorContentMaxLengthInService = 16384

	defaultProductVectorSearchLimit = 10
	maxProductVectorSearchLimit     = 50
	defaultProductVectorBatchSize   = 20
	maxProductVectorBatchSize       = 100
	productVectorTaskTTL            = 0 * time.Hour

	productVectorTaskStatusPending = "pending"
	productVectorTaskStatusRunning = "running"
	productVectorTaskStatusDone    = "completed"
	productVectorTaskStatusFailed  = "failed"
)

type productVectorServiceConfig struct {
	Embedding embeddingutil.ProviderConfig
}

func (s *ProductServiceImpl) GenerateVector(c *gin.Context, req *mastermodels.ProductGenVectorReq) (*mastermodels.ProductVectorResult, error) {
	if req == nil {
		return nil, errors.New("生成参数不能为空")
	}
	product, err := s.dao.GetByID(c, req.ID)
	if err != nil {
		return nil, fmt.Errorf("查询产品失败: %w", err)
	}
	if product == nil {
		return nil, errors.New("产品不存在")
	}
	if product.Status != 1 {
		return nil, errors.New("产品未启用，不能生成向量")
	}
	if err = s.fillProductRelations(c, []*mastermodels.Product{product}); err != nil {
		return nil, err
	}

	cfg, err := s.loadProductVectorConfig(req)
	if err != nil {
		return nil, err
	}

	requestCtx := buildProductRequestContext(c)
	embedder, err := embeddingutil.NewEmbedder(requestCtx, &cfg.Embedding)
	if err != nil {
		return nil, fmt.Errorf("初始化向量模型失败: %w", err)
	}
	return s.generateProductVectorWithEmbedder(c, requestCtx, embedder, product)
}

func (s *ProductServiceImpl) GenerateAllVectors(c *gin.Context, req *mastermodels.ProductGenAllVectorReq) (*mastermodels.ProductVectorTaskData, error) {
	if req == nil {
		return nil, errors.New("全量生成参数不能为空")
	}
	if _, err := loadProductEmbeddingProviderConfig(req.Embedding); err != nil {
		return nil, err
	}

	operatorID := baizeContext.GetUserIdSafe(c)
	operator := baizeContext.GetUserName(c)
	batchSize := resolveProductVectorTaskBatchSize(req.BatchSize, 0)
	taskID := strings.TrimSpace(req.TaskID)
	if s.cache == nil {
		return nil, errors.New("缓存未初始化")
	}

	var (
		progress *mastermodels.ProductVectorTaskProgress
		err      error
	)
	if taskID != "" {
		progress, err = s.loadProductVectorTaskProgress(taskID)
		if err != nil {
			if errors.Is(err, cacheError.Nil) {
				return nil, errors.New("续跑任务不存在或已过期")
			}
			return nil, err
		}
		if progress != nil && progress.OperatorID > 0 && operatorID > 0 && progress.OperatorID != operatorID {
			return nil, errors.New("该任务不属于当前操作人")
		}
		if progress != nil && !isProductVectorTaskUnfinished(progress) {
			return nil, errors.New("该任务已执行完成，请重新创建新任务")
		}
		if progress != nil && progress.Processed > 0 {
			batchSize = resolveProductVectorTaskBatchSize(0, progress.BatchSize)
		}
	} else {
		progress, err = s.loadLatestUnfinishedProductVectorTask(operatorID)
		if err != nil {
			return nil, err
		}
		if progress != nil {
			taskID = progress.TaskID
			if progress.Processed > 0 {
				batchSize = resolveProductVectorTaskBatchSize(0, progress.BatchSize)
			} else {
				batchSize = resolveProductVectorTaskBatchSize(req.BatchSize, progress.BatchSize)
			}
		}
	}

	if progress == nil {
		taskID = strconv.FormatInt(snowflake.GenID(), 10)
		now := time.Now().Format(time.DateTime)
		progress = &mastermodels.ProductVectorTaskProgress{
			TaskID:     taskID,
			Status:     productVectorTaskStatusPending,
			BatchSize:  batchSize,
			Message:    "任务已创建，等待执行",
			OperatorID: operatorID,
			Operator:   operator,
			StartedAt:  now,
			UpdatedAt:  now,
		}
	} else {
		progress.Status = productVectorTaskStatusPending
		progress.BatchSize = batchSize
		progress.OperatorID = operatorID
		progress.Operator = operator
		progress.UpdatedAt = time.Now().Format(time.DateTime)
		progress.FinishedAt = ""
		progress.Progress = calcProductVectorTaskProgress(progress.Processed, progress.Total)
		if progress.StartedAt == "" {
			progress.StartedAt = progress.UpdatedAt
		}
		if progress.Processed > 0 {
			progress.Message = fmt.Sprintf("检测到未完成任务，准备从第%d条继续执行", progress.Processed+1)
		} else {
			progress.Message = "检测到未完成任务，准备继续执行"
		}
	}

	lockKey := productVectorTaskLockKey(taskID)
	if !s.cache.SetNX(context.Background(), lockKey, strconv.FormatInt(operatorID, 10), productVectorTaskTTL) {
		return &mastermodels.ProductVectorTaskData{
			TaskID:      taskID,
			ProgressKey: productVectorTaskCacheKey(taskID),
			Status:      productVectorTaskStatusRunning,
		}, nil
	}
	if err = s.saveProductVectorTaskProgress(taskID, progress); err != nil {
		s.cache.Del(context.Background(), lockKey)
		return nil, err
	}
	s.saveLatestProductVectorTask(operatorID, taskID)

	copiedReq := *req
	copiedReq.TaskID = taskID
	copiedReq.BatchSize = batchSize
	backgroundCtx := cloneProductGinContext(c)
	go s.runGenerateAllVectors(backgroundCtx, taskID, &copiedReq, progress.OperatorID, progress.Operator)

	return &mastermodels.ProductVectorTaskData{
		TaskID:      taskID,
		ProgressKey: productVectorTaskCacheKey(taskID),
		Status:      progress.Status,
	}, nil
}

func (s *ProductServiceImpl) GetGenerateAllVectorsProgress(c *gin.Context, taskID string) (*mastermodels.ProductVectorTaskProgress, error) {
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, errors.New("任务ID不能为空")
	}
	progress, err := s.loadProductVectorTaskProgress(taskID)
	if err != nil {
		if errors.Is(err, cacheError.Nil) {
			return nil, errors.New("任务不存在或已过期")
		}
		return nil, err
	}
	return progress, nil
}

func (s *ProductServiceImpl) ListGenerateAllVectorTasks(c *gin.Context) (*mastermodels.ProductVectorTaskListData, error) {
	if s.cache == nil {
		return nil, errors.New("缓存未初始化")
	}

	operatorID := baizeContext.GetUserIdSafe(c)
	keys, err := s.scanProductVectorTaskKeys()
	if err != nil {
		return nil, err
	}

	rows := make([]*mastermodels.ProductVectorTaskProgress, 0, len(keys))
	for _, key := range keys {
		taskID := strings.TrimPrefix(key, productVectorTaskCacheKey(""))
		taskID = strings.TrimSpace(taskID)
		if taskID == "" {
			continue
		}

		progress, loadErr := s.loadProductVectorTaskProgress(taskID)
		if loadErr != nil {
			if errors.Is(loadErr, cacheError.Nil) {
				continue
			}
			return nil, loadErr
		}
		if progress == nil || !isProductVectorTaskUnfinished(progress) {
			continue
		}
		if operatorID > 0 && progress.OperatorID != operatorID {
			continue
		}
		rows = append(rows, progress)
	}

	sort.Slice(rows, func(i, j int) bool {
		return productVectorTaskSortTime(rows[i]).After(productVectorTaskSortTime(rows[j]))
	})

	return &mastermodels.ProductVectorTaskListData{
		Rows:  rows,
		Total: int64(len(rows)),
	}, nil
}

func (s *ProductServiceImpl) SearchVector(c *gin.Context, req *mastermodels.ProductVectorSearchReq) (*mastermodels.ProductVectorSearchData, error) {
	if req == nil {
		return nil, errors.New("搜索参数不能为空")
	}
	queries, err := normalizeProductVectorSearchQueries([]string{req.Query})
	if err != nil {
		return nil, err
	}
	items, err := s.batchSearchProductVector(c, queries, req.Limit, req.Embedding)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return &mastermodels.ProductVectorSearchData{
			Rows:  make([]*mastermodels.ProductVectorSearchItem, 0),
			Total: 0,
		}, nil
	}
	return &mastermodels.ProductVectorSearchData{
		Rows:  items[0].Rows,
		Total: items[0].Total,
	}, nil
}

func (s *ProductServiceImpl) BatchSearchVector(c *gin.Context, req *mastermodels.ProductVectorBatchSearchReq) (*mastermodels.ProductVectorBatchSearchData, error) {
	if req == nil {
		return nil, errors.New("批量搜索参数不能为空")
	}
	queries, err := normalizeProductVectorSearchQueries(req.Queries)
	if err != nil {
		return nil, err
	}
	rows, err := s.batchSearchProductVector(c, queries, req.Limit, req.Embedding)
	if err != nil {
		return nil, err
	}
	return &mastermodels.ProductVectorBatchSearchData{
		Rows:  rows,
		Total: int64(len(rows)),
	}, nil
}

func (s *ProductServiceImpl) batchSearchProductVector(c *gin.Context, queries []string, limit int, embeddingCfg *mastermodels.ProductEmbeddingConfig) ([]*mastermodels.ProductVectorBatchSearchItem, error) {
	cfg, err := loadProductEmbeddingProviderConfig(embeddingCfg)
	if err != nil {
		return nil, err
	}

	requestCtx := buildProductRequestContext(c)
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
		normalizedVectors = append(normalizedVectors, float64SliceToFloat32(vectors[idx]))
	}

	data, err := s.vectorDao.BatchSearch(c, &mastermodels.ProductVectorBatchSearchReq{
		Queries: queries,
		Limit:   normalizeProductVectorSearchLimit(limit),
	}, normalizedVectors)
	if err != nil {
		return nil, err
	}
	if data == nil || data.Rows == nil {
		return make([]*mastermodels.ProductVectorBatchSearchItem, 0), nil
	}
	return data.Rows, nil
}

func (s *ProductServiceImpl) loadProductVectorConfig(req *mastermodels.ProductGenVectorReq) (*productVectorServiceConfig, error) {
	cfg := &productVectorServiceConfig{}
	embeddingCfg, err := loadProductEmbeddingProviderConfig(req.Embedding)
	if err != nil {
		return nil, err
	}
	cfg.Embedding = *embeddingCfg
	return cfg, nil
}

func (s *ProductServiceImpl) generateProductVectorWithEmbedder(c *gin.Context, requestCtx context.Context, embedder embedding.Embedder, product *mastermodels.Product) (*mastermodels.ProductVectorResult, error) {
	content := buildProductEmbeddingContent(product)
	if strings.TrimSpace(content) == "" {
		return nil, errors.New("产品内容为空，无法生成向量")
	}

	vectors, err := embedder.EmbedStrings(requestCtx, []string{content})
	if err != nil {
		return nil, fmt.Errorf("生成产品向量失败: %w", err)
	}
	if len(vectors) != 1 {
		return nil, fmt.Errorf("向量模型返回结果数量不匹配，expected=1 actual=%d", len(vectors))
	}
	if len(vectors[0]) == 0 {
		return nil, errors.New("向量模型未返回有效结果")
	}

	return s.vectorDao.Upsert(c, product, &mastermodels.ProductVectorUpsertItem{
		Content: content,
		Vector:  float64SliceToFloat32(vectors[0]),
	})
}

func (s *ProductServiceImpl) runGenerateAllVectors(c *gin.Context, taskID string, req *mastermodels.ProductGenAllVectorReq, operatorID int64, operator string) {
	progress, err := s.loadProductVectorTaskProgress(taskID)
	if err != nil {
		if !errors.Is(err, cacheError.Nil) {
			zap.L().Error("load product vector task progress fail", zap.String("taskId", taskID), zap.Error(err))
		}
		progress = nil
	}
	now := time.Now().Format(time.DateTime)
	if progress == nil {
		progress = &mastermodels.ProductVectorTaskProgress{
			TaskID:     taskID,
			BatchSize:  resolveProductVectorTaskBatchSize(req.BatchSize, 0),
			OperatorID: operatorID,
			Operator:   operator,
			StartedAt:  now,
		}
	} else {
		progress.BatchSize = resolveProductVectorTaskBatchSize(req.BatchSize, progress.BatchSize)
		progress.OperatorID = operatorID
		progress.Operator = operator
		if progress.StartedAt == "" {
			progress.StartedAt = now
		}
	}
	if progress.Processed > 0 {
		progress.Message = fmt.Sprintf("正在继续生成产品向量，已处理%d条", progress.Processed)
	} else {
		progress.Message = "正在初始化向量模型"
	}
	progress.Status = productVectorTaskStatusRunning
	progress.UpdatedAt = now
	if err := s.saveProductVectorTaskProgress(taskID, progress); err != nil {
		zap.L().Error("save product vector task progress fail", zap.String("taskId", taskID), zap.Error(err))
		s.cache.Del(context.Background(), productVectorTaskLockKey(taskID))
		return
	}

	defer func() {
		s.cache.Del(context.Background(), productVectorTaskLockKey(taskID))
		if r := recover(); r != nil {
			progress.Status = productVectorTaskStatusFailed
			progress.Message = fmt.Sprintf("任务执行异常: %v", r)
			progress.UpdatedAt = time.Now().Format(time.DateTime)
			progress.FinishedAt = progress.UpdatedAt
			if err := s.saveProductVectorTaskProgress(taskID, progress); err != nil {
				zap.L().Error("save product vector task progress after panic fail", zap.String("taskId", taskID), zap.Error(err))
			}
		}
	}()

	cfg, err := loadProductEmbeddingProviderConfig(req.Embedding)
	if err != nil {
		s.failProductVectorTask(taskID, progress, err)
		return
	}

	requestCtx := buildProductRequestContext(c)
	embedder, err := embeddingutil.NewEmbedder(requestCtx, cfg)
	if err != nil {
		s.failProductVectorTask(taskID, progress, fmt.Errorf("初始化向量模型失败: %w", err))
		return
	}

	size := int64(resolveProductVectorTaskBatchSize(req.BatchSize, progress.BatchSize))
	page := int64(1)
	resumeSkip := 0
	resumeCurrentID := progress.CurrentID
	needResumeTrim := progress.Processed > 0 && size > 0
	if needResumeTrim {
		page = progress.Processed/size + 1
		resumeSkip = int(progress.Processed % size)
	}
	status := int32(1)

	for {
		data, listErr := s.List(c, &mastermodels.ProductQuery{
			Status: &status,
			Page:   page,
			Size:   size,
		})
		if listErr != nil {
			s.failProductVectorTask(taskID, progress, fmt.Errorf("查询启用产品失败: %w", listErr))
			return
		}
		if data == nil {
			data = &mastermodels.ProductListData{}
		}
		if progress.Total == 0 {
			progress.Total = data.Total
			if data.Total == 0 {
				progress.Status = productVectorTaskStatusDone
				progress.Progress = 100
				progress.Message = "没有可生成向量的启用产品"
				progress.UpdatedAt = time.Now().Format(time.DateTime)
				progress.FinishedAt = progress.UpdatedAt
				_ = s.saveProductVectorTaskProgress(taskID, progress)
				return
			}
		}
		rows := data.Rows
		if needResumeTrim {
			rows = trimProductVectorResumeRows(rows, resumeCurrentID, resumeSkip)
			needResumeTrim = false
			if len(rows) == 0 && len(data.Rows) > 0 {
				page++
				continue
			}
		}

		for _, product := range rows {
			if product == nil {
				continue
			}
			progress.CurrentID = product.ID
			progress.CurrentName = product.Name
			progress.Message = "正在生成产品向量"
			progress.UpdatedAt = time.Now().Format(time.DateTime)
			progress.Progress = calcProductVectorTaskProgress(progress.Processed, progress.Total)
			_ = s.saveProductVectorTaskProgress(taskID, progress)

			_, generateErr := s.generateProductVectorWithEmbedder(c, requestCtx, embedder, product)
			progress.Processed++
			if generateErr != nil {
				progress.Failed++
				progress.Message = fmt.Sprintf("产品[%d-%s]生成失败: %v", product.ID, product.Name, generateErr)
				zap.L().Error("generate product vector in batch fail",
					zap.String("taskId", taskID),
					zap.Int64("productId", product.ID),
					zap.String("productName", product.Name),
					zap.Error(generateErr))
			} else {
				progress.Success++
				progress.Message = fmt.Sprintf("产品[%d-%s]生成成功", product.ID, product.Name)
			}
			progress.Progress = calcProductVectorTaskProgress(progress.Processed, progress.Total)
			progress.UpdatedAt = time.Now().Format(time.DateTime)
			_ = s.saveProductVectorTaskProgress(taskID, progress)
		}

		if len(data.Rows) == 0 || progress.Processed >= progress.Total {
			break
		}
		page++
	}

	progress.Status = productVectorTaskStatusDone
	progress.Progress = 100
	progress.Message = fmt.Sprintf("任务完成，成功%d，失败%d", progress.Success, progress.Failed)
	progress.UpdatedAt = time.Now().Format(time.DateTime)
	progress.FinishedAt = progress.UpdatedAt
	if err = s.saveProductVectorTaskProgress(taskID, progress); err != nil {
		zap.L().Error("save completed product vector task progress fail", zap.String("taskId", taskID), zap.Error(err))
	}
}

func (s *ProductServiceImpl) failProductVectorTask(taskID string, progress *mastermodels.ProductVectorTaskProgress, err error) {
	if progress == nil {
		return
	}
	progress.Status = productVectorTaskStatusFailed
	progress.Message = err.Error()
	progress.UpdatedAt = time.Now().Format(time.DateTime)
	progress.FinishedAt = progress.UpdatedAt
	if saveErr := s.saveProductVectorTaskProgress(taskID, progress); saveErr != nil {
		zap.L().Error("save failed product vector task progress fail", zap.String("taskId", taskID), zap.Error(saveErr))
	}
}

func (s *ProductServiceImpl) saveProductVectorTaskProgress(taskID string, progress *mastermodels.ProductVectorTaskProgress) error {
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
	s.cache.Set(context.Background(), productVectorTaskCacheKey(taskID), string(body), productVectorTaskTTL)
	return nil
}

func (s *ProductServiceImpl) loadProductVectorTaskProgress(taskID string) (*mastermodels.ProductVectorTaskProgress, error) {
	if s.cache == nil {
		return nil, errors.New("缓存未初始化")
	}
	val, err := s.cache.Get(context.Background(), productVectorTaskCacheKey(taskID))
	if err != nil {
		return nil, err
	}
	progress := &mastermodels.ProductVectorTaskProgress{}
	if err = json.Unmarshal([]byte(val), progress); err != nil {
		return nil, err
	}
	return progress, nil
}

func loadProductEmbeddingProviderConfig(req *mastermodels.ProductEmbeddingConfig) (*embeddingutil.ProviderConfig, error) {
	if req == nil {
		return nil, errors.New("未配置 embedding.model_id 或 erp.master.product_vector.embedding.model_id")
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
		return nil, errors.New("未配置 embedding.model_id 或 erp.master.product_vector.embedding.model_id")
	}
	return cfg, nil
}

func buildProductRequestContext(c *gin.Context) context.Context {
	if c != nil && c.Request != nil {
		return c.Request.Context()
	}
	return context.Background()
}

func cloneProductGinContext(c *gin.Context) *gin.Context {
	if c == nil {
		return nil
	}
	clone := c.Copy()
	if c.Request != nil {
		clone.Request = c.Request.Clone(context.Background())
	}
	return clone
}

func buildProductEmbeddingContent(product *mastermodels.Product) string {
	if product == nil {
		return ""
	}
	lines := make([]string, 0, 10)
	appendLine := func(label, value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		lines = append(lines, label+": "+value)
	}

	appendLine("产品名称", product.Name)
	appendLine("产品分类", product.CategoryName)
	appendLine("单位", product.UnitName)
	appendLine("条码", product.BarCode)
	appendLine("规格", product.Standard)
	appendLine("备注", product.Remark)
	if product.ExpiryDay > 0 {
		lines = append(lines, fmt.Sprintf("保质期: %d天", product.ExpiryDay))
	}
	if product.Weight > 0 {
		lines = append(lines, fmt.Sprintf("重量: %.3fkg", product.Weight))
	}
	if product.PurchasePrice > 0 {
		lines = append(lines, fmt.Sprintf("采购价: %.2f", product.PurchasePrice))
	}
	if product.SalePrice > 0 {
		lines = append(lines, fmt.Sprintf("销售价: %.2f", product.SalePrice))
	}
	if product.MinPrice > 0 {
		lines = append(lines, fmt.Sprintf("最低价: %.2f", product.MinPrice))
	}
	return trimProductEmbeddingContent(strings.Join(lines, "\n"), productVectorContentMaxLengthInService)
}

func trimProductEmbeddingContent(value string, max int) string {
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

func normalizeProductVectorSearchQueries(queries []string) ([]string, error) {
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

func (s *ProductServiceImpl) scanProductVectorTaskKeys() ([]string, error) {
	if s.cache == nil {
		return nil, errors.New("缓存未初始化")
	}

	match := productVectorTaskCacheKey("") + "*"
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

func productVectorTaskCacheKey(taskID string) string {
	return "erp:master:product:vector:task:" + strings.TrimSpace(taskID)
}

func productVectorTaskLockKey(taskID string) string {
	return productVectorTaskCacheKey(taskID) + ":lock"
}

func productVectorLatestTaskKey(operatorID int64) string {
	return "erp:master:product:vector:latest:" + strconv.FormatInt(operatorID, 10)
}

func calcProductVectorTaskProgress(processed, total int64) int {
	if total <= 0 || processed <= 0 {
		return 0
	}
	if processed >= total {
		return 100
	}
	return int(processed * 100 / total)
}

func normalizeProductVectorBatchSize(batchSize int) int {
	if batchSize <= 0 {
		return defaultProductVectorBatchSize
	}
	if batchSize > maxProductVectorBatchSize {
		return maxProductVectorBatchSize
	}
	return batchSize
}

func resolveProductVectorTaskBatchSize(reqBatchSize, progressBatchSize int) int {
	if reqBatchSize > 0 {
		return normalizeProductVectorBatchSize(reqBatchSize)
	}
	if progressBatchSize > 0 {
		return normalizeProductVectorBatchSize(progressBatchSize)
	}
	return defaultProductVectorBatchSize
}

func isProductVectorTaskUnfinished(progress *mastermodels.ProductVectorTaskProgress) bool {
	if progress == nil {
		return false
	}
	return progress.Status != productVectorTaskStatusDone
}

func trimProductVectorResumeRows(rows []*mastermodels.Product, currentID int64, skipCount int) []*mastermodels.Product {
	if len(rows) == 0 {
		return rows
	}
	if currentID > 0 {
		for index, product := range rows {
			if product != nil && product.ID == currentID {
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

func (s *ProductServiceImpl) loadLatestUnfinishedProductVectorTask(operatorID int64) (*mastermodels.ProductVectorTaskProgress, error) {
	if operatorID <= 0 || s.cache == nil {
		return nil, nil
	}
	taskID, err := s.cache.Get(context.Background(), productVectorLatestTaskKey(operatorID))
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
	progress, err := s.loadProductVectorTaskProgress(taskID)
	if err != nil {
		if errors.Is(err, cacheError.Nil) {
			s.cache.Del(context.Background(), productVectorLatestTaskKey(operatorID))
			return nil, nil
		}
		return nil, err
	}
	if !isProductVectorTaskUnfinished(progress) {
		return nil, nil
	}
	return progress, nil
}

func (s *ProductServiceImpl) saveLatestProductVectorTask(operatorID int64, taskID string) {
	if operatorID <= 0 || s.cache == nil {
		return
	}
	s.cache.Set(context.Background(), productVectorLatestTaskKey(operatorID), strings.TrimSpace(taskID), productVectorTaskTTL)
}

func productVectorTaskSortTime(progress *mastermodels.ProductVectorTaskProgress) time.Time {
	if progress == nil {
		return time.Time{}
	}
	if ts := parseProductVectorTaskTime(progress.UpdatedAt); !ts.IsZero() {
		return ts
	}
	if ts := parseProductVectorTaskTime(progress.StartedAt); !ts.IsZero() {
		return ts
	}
	return time.Time{}
}

func parseProductVectorTaskTime(value string) time.Time {
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

func normalizeProductVectorSearchLimit(limit int) int {
	if limit <= 0 {
		return defaultProductVectorSearchLimit
	}
	if limit > maxProductVectorSearchLimit {
		return maxProductVectorSearchLimit
	}
	return limit
}

func float64SliceToFloat32(vector []float64) []float32 {
	result := make([]float32, 0, len(vector))
	for _, value := range vector {
		result = append(result, float32(value))
	}
	return result
}

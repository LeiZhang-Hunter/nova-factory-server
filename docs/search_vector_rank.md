# 商品向量检索与重排架构说明

## 1. 文档目标

本文档用于说明当前商品向量检索与重排方案的整体架构、模块职责、检索链路、重排机制以及后续扩展方向。

当前实现目标不是只做单一路径的向量检索，而是构建一套可演进的商品搜索架构，具备以下能力：

- Query 预处理与语义增强
- 多意图查询识别
- 向量检索 + BM25 混合召回
- 应用层字段加权重排
- 模块化分层，便于后续继续扩展品牌、分类、规格等策略

---

## 2. 总体架构

当前架构可以分为 5 层：

1. 接口与服务编排层
2. Query 理解层
3. 检索策略层
4. Milvus 执行层
5. 应用层重排层

整体数据流如下：

```text
用户输入 Query
    |
    v
Service 层
    - NormalizeQueries
    - ProcessQueries
    - 生成 EmbeddingText / HybridText
    - 调用向量模型生成向量
    |
    v
DAO 主流程
    - BatchSearch
    - 构建 SearchPlan
    |
    v
策略层
    - code / spec / category / semantic
    - 按意图增强 SearchText
    |
    v
执行层
    - HybridSearch: Dense Vector + BM25 Sparse
    - 或旧 collection 降级为纯向量检索
    |
    v
结果解析层
    - 解析 Milvus ResultSet
    - 回填业务字段
    |
    v
应用层重排
    - 标题/编码/分类/规格/正文字段加权
    - 动态阈值裁剪
    |
    v
最终商品结果
```

---

## 3. 模块分层

### 3.1 Service 编排层

核心文件：

- `app/business/erp/master/masterservice/masterserviceimpl/i_product_vector_service_impl.go`

主要职责：

- 接收单条或批量搜索请求
- 调用 `vectorsearch` 对 query 做清洗和结构化处理
- 生成两套检索文本：
  - `EmbeddingText`：给向量模型使用
  - `HybridText`：给 BM25 / 稀疏检索使用
- 调用 embedding 模型生成稠密向量
- 调用 DAO 层进行批量搜索

关键流程：

```text
SearchVector / BatchSearchVector
    -> batchSearchProductVector
        -> ProcessQueries
        -> EmbedStrings(EmbeddingText)
        -> vectorDao.BatchSearch(Queries + SearchTexts + vectors)
```

设计意图：

- Service 层不直接参与 Milvus 细节
- Service 层负责把用户查询转换成“适合检索的输入”
- DAO 层只关注“怎么查”和“怎么排”

---

### 3.2 Query 理解层

核心目录：

- `app/utils/vectorsearch`

当前已拆分模块如下：

- `types.go`
- `query.go`
- `dictionary.go`
- `parser.go`
- `text.go`
- `spec.go`
- `category.go`
- `rank.go`

#### 3.2.1 `types.go`

职责：

- 定义查询处理中使用的核心结构

核心结构：

- `LabeledValue`
  - 用于把结构化字段拼接成 embedding 文本
- `ProcessedQuery`
  - 表示 query 处理后的结构化结果

`ProcessedQuery` 当前包含：

- `Original`：原始 query
- `Normalized`：归一化后 query
- `Tokens`：基础分词结果
- `ExpandedTokens`：同义词扩展结果
- `Keywords`：去重后的关键词集合
- `CategoryTerms`：分类词
- `SpecTerms`：规格词
- `CodeTerms`：编码词
- `EmbeddingText`：向量检索文本
- `HybridText`：混合检索文本
- `IsCodeLike`：是否像编码类查询
- `IsShortQuery`：是否短 query

#### 3.2.2 `query.go`

职责：

- 提供 Query 处理入口

核心函数：

- `NormalizeQueries`
  - 清洗 query，过滤空字符串
- `ProcessQueries`
  - 批量处理 query
- `ProcessQuery`
  - 处理单条 query，生成 `ProcessedQuery`

`ProcessQuery` 当前做的事情：

1. 统一空白与大小写
2. 分词
3. 同义词扩展
4. 分类词提取
5. 规格词提取
6. 编码词提取
7. 生成 `EmbeddingText`
8. 生成 `HybridText`
9. 标记 `IsCodeLike` / `IsShortQuery`

#### 3.2.3 `dictionary.go`

职责：

- 提供轻量内置词典

当前词典包括：

- `defaultStopWords`
  - 停用词
- `defaultSynonyms`
  - 同义词词典
- `defaultCategoryTerms`
  - 分类词词典

设计意图：

- 先用代码内置词典快速建立基础能力
- 后续可升级为配置中心、数据库或 Redis 管理

#### 3.2.4 `parser.go`

职责：

- 提供 query 基础切词与词语加工能力

核心函数：

- `tokenize`
  - 提取中文、字母、数字
  - 对中文短词补充 2-gram / 3-gram
- `expandTokens`
  - 基于同义词词典扩展 query
- `dedupeKeywords`
  - 统一大小写、过滤停用词、去重
- `buildQueryText`
  - 构造最终的 embedding / hybrid 检索文本

设计意图：

- 把 query 从“用户自然输入”转为“适合召回的检索表达”

#### 3.2.5 `text.go`

职责：

- 提供通用文本清洗与拼装能力

核心函数：

- `BuildLabeledContent`
  - 把商品结构化字段拼接成带标签文本
- `TrimRunes`
  - 按 rune 裁剪文本，避免中文截断问题
- `NormalizeWhitespace`
  - 统一空白字符
- `isCodeLike`
  - 判断 query 是否更像编码、条码、SKU

设计意图：

- 统一文本规范，避免各层自己处理字符串导致逻辑分散

#### 3.2.6 `spec.go`

职责：

- 识别和标准化商品规格

核心能力：

- 识别 `550ml`、`0.55L`、`5kg`、`1斤` 等规格表达
- 将不同单位统一为标准表达
- 提取编码类词
- 为重排提供规格匹配分

核心函数：

- `extractSpecTerms`
- `extractCodeTerms`
- `findNormalizedSpecs`
- `normalizeSpec`
- `scoreSpecTerms`

设计意图：

- 解决商品搜索里“同名不同规格”的误召回问题

#### 3.2.7 `category.go`

职责：

- 从 query 和 token 中提取分类词

核心函数：

- `extractCategoryTerms`
- `matchCategoryTerms`

设计意图：

- 支持分类类检索意图识别
- 为后续分类优先召回和分类加权重排做准备

#### 3.2.8 `rank.go`

职责：

- 对第一阶段召回候选做应用层精排

核心结构：

- `RankCandidate`
- `RankedCandidate`

核心函数：

- `RerankCandidates`
- `combineCandidateScore`
- `calcDynamicThreshold`
- `normalizeBaseScores`

设计意图：

- 保留 Milvus 初排价值
- 增强业务字段可解释性
- 对编码类、规格类、短 query 做差异化重排

---

## 4. DAO 分层设计

当前 DAO 层已经拆成多个文件，避免所有逻辑堆积在一个文件里。

### 4.1 `product_vector_constants.go`

职责：

- 定义 collection 名、字段名、索引名、长度限制、搜索参数
- 提供 `NewProductVectorDao`

这里统一管理了：

- Milvus 字段定义
- varchar 最大长度
- 候选召回倍数
- 最大/最小候选数
- 默认返回数量和最大返回数量

设计意图：

- 避免 magic number 分散在各个文件中

### 4.2 `product_vector_collection.go`

职责：

- 处理 Milvus collection 相关逻辑

核心函数：

- `loadProductVectorConfig`
- `supportsProductVectorHybridSearch`
- `ensureProductVectorCollection`

能力说明：

- 自动读取 collection 配置
- 自动创建 collection
- 自动创建 dense vector index
- 自动创建 sparse vector index
- 检查旧 collection 是否支持 hybrid search
- 校验已有 collection 字段和向量维度

设计意图：

- 把表结构管理从业务查询逻辑中剥离

### 4.3 `i_product_vector_dao_impl.go`

职责：

- DAO 主流程编排

核心函数：

- `Upsert`
- `Search`
- `BatchSearch`

其中：

- `Upsert`
  - 负责写入商品结构化字段、全文文本和向量
- `Search`
  - 单条查询复用批量查询逻辑
- `BatchSearch`
  - 批量搜索总入口
  - 调用策略层生成检索计划
  - 调用执行层查询 Milvus
  - 调用结果解析与重排逻辑

设计意图：

- 保持“主流程清晰、细节下沉”

### 4.4 `product_vector_search_strategy.go`

职责：

- 识别 query 意图
- 根据不同意图构建不同的检索计划

核心结构：

- `productVectorSearchPlan`
  - 单条 query 的检索计划
- `productVectorSearchStrategy`
  - 检索策略接口

当前已实现策略：

- `productVectorCodeSearchStrategy`
  - 编码、条码、SKU 类查询
- `productVectorSpecSearchStrategy`
  - 规格类查询
- `productVectorCategorySearchStrategy`
  - 分类类查询
- `productVectorSemanticSearchStrategy`
  - 通用语义查询兜底

策略优先级：

```text
code -> spec -> category -> semantic
```

设计意图：

- 不同 query 走不同策略
- 减少一个函数里堆积大量 if/else
- 后续可继续新增 `brand`、`price`、`supplier` 等策略

### 4.5 `product_vector_search_executor.go`

职责：

- 统一执行 Milvus 查询

核心函数：

- `executeProductVectorSearch`
- `buildProductVectorOutputFields`

执行逻辑：

1. 判断 collection 是否支持 hybrid search
2. 如果支持：
   - Dense Vector 检索
   - Sparse BM25 检索
   - RRF 融合
3. 如果不支持：
   - 降级为纯向量检索

设计意图：

- 将“查询计划”和“查询执行”完全分离

### 4.6 `product_vector_result.go`

职责：

- 解析 Milvus 返回结果
- 处理空结果
- 做应用层重排

核心函数：

- `buildEmptyProductVectorBatchSearchData`
- `parseProductVectorSearchResultSet`
- `normalizeProductVectorSearchLimit`
- `resolveProductVectorSearchCandidateLimit`
- `rerankProductVectorSearchRows`

设计意图：

- 统一结果回填逻辑
- 在 DAO 层集中处理二次重排

---

## 5. 检索链路详解

### 5.1 写入链路

写入链路如下：

```text
Product
    -> 组装结构化文本
    -> 生成向量
    -> Upsert 到 Milvus
        - product_id
        - name
        - bar_code
        - category_id/category_name
        - unit_id/unit_name
        - standard
        - remark
        - content
        - vector
        - text_sparse_vector(由 BM25 function 生成)
```

说明：

- `content` 是完整商品语义文本
- `vector` 是 embedding 模型生成的稠密向量
- `text_sparse_vector` 是由 Milvus BM25 function 自动生成的稀疏向量

### 5.2 查询链路

查询链路如下：

```text
原始 query
    -> NormalizeQueries
    -> ProcessQuery
        - tokenize
        - expandTokens
        - extractCategoryTerms
        - extractSpecTerms
        - extractCodeTerms
    -> 生成 EmbeddingText / HybridText
    -> 调用 embedding 模型生成向量
    -> buildProductVectorSearchPlans
    -> executeProductVectorSearch
    -> parseProductVectorSearchResultSet
    -> rerankProductVectorSearchRows
    -> 输出最终结果
```

### 5.3 Query 表达

当前 query 最终会形成两套文本：

- `EmbeddingText`
  - 用于生成向量
  - 更强调扩展词和语义补全
- `HybridText`
  - 用于 sparse / BM25
  - 更强调关键词匹配

这样可以让：

- 向量检索负责语义召回
- BM25 负责关键词召回
- 两者优势互补

---

## 6. 查询策略设计

### 6.1 Code Strategy

适用场景：

- 条码
- SKU
- 货号
- 编码类查询

策略特点：

- 判断 `IsCodeLike`
- 提取 `CodeTerms`
- 将编码词增强到 `SearchText`
- 重排时大幅提高编码命中权重

目标：

- 把“精确目标商品”顶到前面

### 6.2 Spec Strategy

适用场景：

- `矿泉水 550ml`
- `牛奶 250ml`
- `大米 5kg`

策略特点：

- 提取规格词
- 做规格标准化
- 将规格词增强到检索文本
- 重排时对规格命中加分
- 对 query 带规格但候选不匹配规格的结果进行惩罚

目标：

- 拉开“同名不同规格”商品差异

### 6.3 Category Strategy

适用场景：

- `饮料`
- `纸品`
- `日化`

策略特点：

- 提取分类词
- 将分类词增强到检索文本
- 为后续分类优先召回打基础

目标：

- 让分类类 query 更容易返回相应类目商品

### 6.4 Semantic Strategy

适用场景：

- 普通商品名称或语义 query

策略特点：

- 不额外做强意图增强
- 作为兜底策略存在

目标：

- 保证所有 query 都有可执行路径

---

## 7. 混合召回设计

当前查询执行使用 Milvus 的 HybridSearch。

召回组成：

- Dense Vector
  - 基于 embedding 向量相似度
  - 更适合语义召回
- Sparse BM25
  - 基于 `content` 的稀疏向量
  - 更适合关键词命中
- RRF
  - 对多路结果进行融合

整体逻辑：

```text
Embedding 向量召回
      +
BM25 稀疏召回
      |
      v
RRF 融合
      |
      v
候选集
```

为什么要这样设计：

- 纯向量检索可能语义对了但关键词不够准
- 纯 BM25 检索可能关键词对了但语义泛化差
- HybridSearch 可以兼顾两者

兼容逻辑：

- 如果旧 collection 不支持 sparse 字段，则自动降级为纯向量检索

---

## 8. 应用层重排设计

Milvus 返回候选后，会进入应用层精排。

核心入口：

- `vectorsearch.RerankCandidates`

### 8.1 为什么要应用层重排

底层召回分数只代表相似度，不一定完全符合商品搜索场景。

例如：

- 商品名命中比备注命中更重要
- 编码匹配比正文语义更重要
- 规格匹配对同名商品区分非常关键

因此需要在业务层额外融合字段信号。

### 8.2 当前重排信号

当前综合分由以下几类信号组成：

- `baseScore`
  - 底层召回原始分数
- `fieldScore`
  - 标题、编码、分类、单位、规格、备注、正文的字段加权得分
- `tagScore`
  - 分类、单位、规格中的标签命中
- `rankBonus`
  - 初始排序位置奖励

### 8.3 字段权重思路

当前权重倾向如下：

- `Title` 权重最高
- `Code` 次高
- `Standard`、`Spec` 较高
- `Category` 中等
- `Remark`、`Content` 较低

设计原则：

- 商品名是最核心的业务字段
- 编码类查询要更偏精确匹配
- 长文本只能作为补充信号，不能压过高价值结构化字段

### 8.4 特殊 query 的差异化重排

#### 编码类 query

- 提高 `codeScore`
- 降低正文噪声影响
- 对精确码命中追加额外奖励

#### 规格类 query

- 规格命中额外加分
- 候选完全不命中规格时给予惩罚

#### 短 query

- 额外提高标题命中权重

### 8.5 动态阈值裁剪

精排后并不会简单返回前 N 条，而是会根据本次 query 的结果质量动态计算阈值。

阈值依据：

- top1 的绝对得分
- 是否编码类查询
- 是否短 query
- top1 / top2 的分差

设计目的：

- 防止阈值过低导致大量弱相关结果进入结果集
- 防止阈值过高导致本来可用的结果被过滤掉

---

## 9. Milvus 数据结构设计

当前 collection 主要字段如下：

- `product_id`
- `name`
- `bar_code`
- `category_id`
- `category_name`
- `unit_id`
- `unit_name`
- `standard`
- `remark`
- `expiry_day`
- `weight`
- `purchase_price`
- `sale_price`
- `min_price`
- `content`
- `vector`
- `text_sparse_vector`

索引设计：

- `vector`
  - AutoIndex
  - Cosine
- `text_sparse_vector`
  - SparseInvertedIndex
  - BM25

说明：

- `content` 字段开启 analyzer / match
- `text_sparse_vector` 由 BM25 function 自动派生
- 当前 sparse 路径主要依赖 `content`

---

## 10. 当前设计的优点

### 10.1 分层清晰

- Query 理解、策略调度、执行、解析、重排分开
- DAO 主流程不再承担所有细节

### 10.2 可扩展性强

- 可以继续扩展新的 query 策略
- 可以继续增强词典体系
- 可以继续接入模型 reranker

### 10.3 兼容存量数据

- 旧 collection 不支持 hybrid 时可自动降级

### 10.4 业务可解释性更强

- 重排规则可解释
- 字段权重可控
- 便于根据业务反馈调整

---

## 11. 当前限制

虽然当前架构已经比最初版本清晰很多，但仍有一些限制：

### 11.1 词典仍是内置词典

当前：

- 同义词
- 分类词

仍然写在代码里。

影响：

- 更新词典需要发版
- 业务运营无法自主调整

### 11.2 字段化 sparse 召回尚未完全展开

当前 sparse 召回主要基于：

- `content`

后续更适合继续拆成：

- `name_sparse_vector`
- `category_sparse_vector`
- `standard_sparse_vector`
- `code_sparse_vector`

### 11.3 品牌策略尚未接入

当前已有：

- code
- spec
- category
- semantic

后续可以继续加：

- brand

### 11.4 尚未接入模型级 reranker

当前精排是规则型重排。

后续可继续接入：

- BGE Reranker
- CrossEncoder
- 外部 rerank 服务

---

## 12. 推荐扩展方向

建议后续按以下顺序继续演进：

### 第一阶段：品牌与词典体系

- 增加 `BrandTerms`
- 增加品牌词典
- 增加 `brand` 策略
- 在重排中加入品牌强命中信号

### 第二阶段：字段化召回

- 将 `name`、`category`、`standard`、`code` 分拆为独立 sparse 召回通道
- 支持多路字段化 hybrid search

### 第三阶段：模型级精排

- 在应用层规则重排之后，对 TopK 候选接入 CrossEncoder / Reranker 精排

### 第四阶段：反馈闭环

- 记录 query、召回结果、点击商品
- 结合点击数据调优字段权重和词典

---

## 13. 当前文件结构建议理解方式

建议从以下顺序理解整套代码：

1. `i_product_vector_service_impl.go`
   - 看请求如何进入搜索链路
2. `query.go`
   - 看 query 如何变成结构化表达
3. `product_vector_search_strategy.go`
   - 看 query 如何被识别为不同策略
4. `product_vector_search_executor.go`
   - 看策略结果如何执行到 Milvus
5. `product_vector_result.go`
   - 看结果如何解析和重排
6. `rank.go`
   - 看精排规则如何融合业务信号

---

## 14. 架构总结

当前商品向量搜索架构已经从“单文件堆逻辑”演进为“多层分治”的可维护结构：

- Service 层负责输入编排与向量生成
- `vectorsearch` 负责 query 理解与重排能力
- Strategy 层负责按意图生成检索计划
- Executor 层负责统一执行 Milvus 查询
- Result 层负责结果解析与应用层重排

这套架构的核心价值在于：

- 检索效果更准
- 代码职责更清楚
- 后续更容易扩展品牌、分类、规格、模型重排等能力

如果后续继续演进，建议优先做：

1. 品牌策略
2. 词典外置化
3. 字段化 sparse 召回
4. 模型级 reranker

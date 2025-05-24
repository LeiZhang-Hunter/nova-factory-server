package aiDataSetModels

import (
	"nova-factory-server/app/baize"
	"time"
)

// DataSetRequest 创建知识库请求
type DataSetRequest struct {
	Name           string        `json:"name" binding:"required"`
	Avatar         string        `json:"avatar"`
	Description    string        `json:"description"`
	EmbeddingModel string        `json:"embedding_model"`
	Permission     string        `json:"permission"`
	ChunkMethod    string        `json:"chunk_method"`
	Pagerank       string        `json:"pagerank"`
	ParserConfig   *ParserConfig `json:"parser_config,omitempty"`
}

type Graphrag struct {
	UseGraphrag bool     `json:"use_graphrag"`
	EntityTypes []string `json:"entity_types"`
	Method      string   `json:"method"`
	Resolution  bool     `json:"resolution"`
	Community   bool     `json:"community"`
}

type ParserConfig struct {
	ChunkTokenNum   int    `json:"chunk_token_num"`
	AutoKeywords    int    `json:"auto_keywords"`
	AutoQuestions   int    `json:"auto_questions"`
	TagKbIds        int    `json:"tag_kb_ids,omitempty"`
	Delimiter       string `json:"delimiter"`
	Html4Excel      bool   `json:"html4excel"`
	LayoutRecognize string `json:"layout_recognize"`
	Raptor          struct {
		UseRaptor  bool    `json:"use_raptor"`
		Prompt     string  `json:"prompt"`
		MaxToken   int     `json:"max_token"`
		Threshold  float64 `json:"threshold"`
		MaxCluster int     `json:"max_cluster"`
		RandomSeed int     `json:"random_seed"`
	} `json:"raptor"`
	TaskPageSize int       `json:"task_page_size"`
	Graphrag     *Graphrag `json:"graphrag,omitempty"`
}

type DataSetData struct {
	Avatar                 string        `json:"avatar"`
	ChunkCount             int           `json:"chunk_count"`
	ChunkMethod            string        `json:"chunk_method"`
	CreateDate             string        `json:"create_date"`
	CreateTime             int64         `json:"create_time"`
	CreatedBy              string        `json:"created_by"`
	Description            string        `json:"description"`
	DocumentCount          int           `json:"document_count"`
	EmbeddingModel         string        `json:"embedding_model"`
	Id                     string        `json:"id"`
	Language               string        `json:"language"`
	Name                   string        `json:"name"`
	Pagerank               int           `json:"pagerank"`
	ParserConfig           *ParserConfig `json:"parser_config,omitempty"`
	Permission             string        `json:"permission"`
	SimilarityThreshold    float64       `json:"similarity_threshold"`
	Status                 string        `json:"status"`
	TenantId               string        `json:"tenant_id"`
	TokenNum               int           `json:"token_num"`
	UpdateDate             string        `json:"update_date"`
	UpdateTime             int64         `json:"update_time"`
	VectorSimilarityWeight float64       `json:"vector_similarity_weight"`
}

// DataSetCreateResponse 创建知识库响应
type DataSetCreateResponse struct {
	Code    int         `json:"code"`
	Data    DataSetData `json:"data"`
	Message string      `json:"message"`
}

type SysDataset struct {
	DatasetID                     int64         `gorm:"column:dataset_id;primaryKey;comment:出库id" json:"datasetId,string"`                                                // 出库id
	DatasetAvatar                 string        `gorm:"column:dataset_avatar;comment:base64 编码的头像。" json:"dataset_avatar"`                                                // base64 编码的头像。
	DatasetChunkMethod            string        `gorm:"column:dataset_chunk_method;comment:要创建的数据集的分块方法" json:"dataset_chunk_method"`                                     // 要创建的数据集的分块方法
	DatasetCreateDate             time.Time     `gorm:"column:dataset_create_date;comment:创建时间" json:"dataset_create_date"`                                               // 创建时间
	DatasetCreateTime             int64         `gorm:"column:dataset_create_time;not null;comment:创建时间" json:"dataset_create_time"`                                      // 创建时间
	DatasetCreatedBy              string        `gorm:"column:dataset_created_by;comment:创建人" json:"dataset_created_by"`                                                  // 创建人
	DatasetDescription            string        `gorm:"column:dataset_description;comment:描述" json:"dataset_description"`                                                 // 描述
	DatasetDocumentCount          int64         `gorm:"column:dataset_document_count;not null;comment:文档数" json:"dataset_document_count"`                                 // 文档数
	DatasetEmbeddingModel         string        `gorm:"column:dataset_embedding_model;comment:嵌入模型" json:"dataset_embedding_model"`                                       // 嵌入模型
	DatasetUUID                   string        `gorm:"column:dataset_uuid;comment:知识库id" json:"dataset_uuid"`                                                            // 知识库id
	DatasetLanguage               string        `gorm:"column:dataset_language;comment:知识库语言" json:"dataset_language"`                                                    // 知识库语言
	DatasetName                   string        `gorm:"column:dataset_name;comment:要创建的数据集的唯一名称" json:"dataset_name"`                                                     // 要创建的数据集的唯一名称
	DatasetPagerank               int           `gorm:"column:dataset_pagerank;comment:页面排名" json:"dataset_pagerank"`                                                     // 页面排名
	DatasetParserConfig           string        `gorm:"column:dataset_parser_config;comment:数据集解析器的配置设置。" json:"dataset_parser_config"`                                   // 数据集解析器的配置设置。
	DatasetPermission             string        `gorm:"column:dataset_permission;comment:指定谁可以访问要创建的数据集。" json:"dataset_permission"`                                      // 指定谁可以访问要创建的数据集。
	DatasetSimilarityThreshold    float64       `gorm:"column:dataset_similarity_threshold;default:0.00;comment:最小相似度分数" json:"dataset_similarity_threshold"`             // 最小相似度分数
	DatasetTokenNum               int64         `gorm:"column:dataset_token_num;not null;comment:token num" json:"dataset_token_num"`                                     // token num
	DatasetUpdateDate             time.Time     `gorm:"column:dataset_update_date;comment:更新时间" json:"dataset_update_date"`                                               // 更新时间
	DatasetUpdateTime             int64         `gorm:"column:dataset_update_time;not null;comment:更新时间" json:"dataset_update_time"`                                      // 更新时间
	DatasetVectorSimilarityWeight float64       `gorm:"column:dataset_vector_similarity_weight;default:0.00;comment:矢量余弦相似性的权重。" json:"dataset_vector_similarity_weight"` // 矢量余弦相似性的权重。
	ParserConfig                  *ParserConfig `gorm:"-" json:"parser_config,omitempty"`                                                                                 // 数据集解析器的配置设置。
	State                         bool          `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                                                                 // 操作状态（0正常 -1删除）
	DeptID                        int64         `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                                                       // 部门ID
	CreateUserName                string        `json:"createUserName" gorm:"-"`
	UpdateUserName                string        `json:"updateUserName" gorm:"-"`
	baize.BaseEntity
}

// UpdateDataSetRequest 更新知识库请求
type UpdateDataSetRequest struct {
	Name           string        `json:"name" binding:"required"`
	EmbeddingModel string        `json:"embedding_model,omitempty"`
	ChunkMethod    string        `json:"chunk_method"`
	ParserConfig   *ParserConfig `json:"parser_config,omitempty"`
	Description    string        `json:"description,omitempty"`
	Pagerank       int           `json:"pagerank"`
}

// DatasetListReq 知识库列表集
type DatasetListReq struct {
	Name string `json:"name,string" db:"name" form:"name"`
	baize.BaseEntityDQL
}

// SysDatasetListData 知识库列表返回结果
type SysDatasetListData struct {
	Rows  []*SysDataset `json:"rows"`
	Total int64         `json:"total"`
}

// SysDatasetDeleteReq 知识库列表返回结果
type SysDatasetDeleteReq struct {
	Ids []string `json:"ids,string" db:"ids" form:"ids"`
}

type DatasetListResponse struct {
	Code    int           `json:"code"`
	Data    []DataSetData `json:"data"`
	Message string        `json:"message"`
}

// GetDatasetInfoReq 知识库列表集
type GetDatasetInfoReq struct {
	Id int64 `json:"dataset_id,string" db:"dataset_id" form:"dataset_id"`
}

type DataSetConfig struct {
	LayoutRecognize []string `json:"layout_recognize"`
	Delimiter       string   `json:"delimiter"`
	ChunkMethod     []string `json:"chunk_method"`
	GraphragMethod  []string `json:"graphrag_method"`
}

func NewDataSetConfig() *DataSetConfig {
	return &DataSetConfig{
		LayoutRecognize: []string{
			"DeepDOC",
			"Naive",
		},
		Delimiter: "\\n!?;。；！？",
		ChunkMethod: []string{
			"naive",
			"manual",
			"qa",
			"table",
			"paper",
			"book",
			"laws",
			"presentation",
			"picture",
			"one",
			"email",
		},
		GraphragMethod: []string{
			"Light",
			"General",
		},
	}
}

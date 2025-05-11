package aiDataSetModels

import "nova-factory-server/app/baize"

type UploadDocumentData struct {
	ChunkMethod  string       `json:"chunk_method"`
	CreatedBy    string       `json:"created_by"`
	DatasetId    string       `json:"dataset_id"`
	Id           string       `json:"id"`
	Location     string       `json:"location"`
	Name         string       `json:"name"`
	ParserConfig ParserConfig `json:"parser_config"`
	Run          string       `json:"run"`
	Size         int          `json:"size"`
	Thumbnail    string       `json:"thumbnail"`
	Type         string       `json:"type"`
}

type UploadDocumentResponse struct {
	Code    int                   `json:"code"`
	Data    []*UploadDocumentData `json:"data"`
	Message string                `json:"message"`
}

// SysDatasetDocument 物料出库管理
type SysDatasetDocument struct {
	DocumentID          int64  `gorm:"column:document_id;not null;comment:文档id" json:"document_id"`                    // 文档id
	DatasetID           int64  `gorm:"column:dataset_id;primaryKey;comment:数据集id" json:"dataset_id"`                   // 数据集id
	DatasetChunkMethod  string `gorm:"column:dataset_chunk_method;comment:要创建的数据集的分块方法" json:"dataset_chunk_method"`   // 要创建的数据集的分块方法
	DatasetCreatedBy    string `gorm:"column:dataset_created_by;comment:创建人" json:"dataset_created_by"`                // 创建人
	DatasetDocumentUUID string `gorm:"column:dataset_document_uuid;comment:文档id" json:"dataset_document_uuid"`         // 文档id
	DatasetDatasetUUID  string `gorm:"column:dataset_dataset_uuid;comment:知识库id" json:"dataset_dataset_uuid"`          // 知识库id
	DatasetLanguage     string `gorm:"column:dataset_language;comment:知识库语言" json:"dataset_language"`                  // 知识库语言
	DatasetLocation     string `gorm:"column:dataset_location;comment:知识库定位" json:"dataset_location"`                  // 知识库定位
	DatasetName         string `gorm:"column:dataset_name;comment:要创建的数据集的唯一名称" json:"dataset_name"`                   // 要创建的数据集的唯一名称
	DatasetParserConfig string `gorm:"column:dataset_parser_config;comment:数据集解析器的配置设置。" json:"dataset_parser_config"` // 数据集解析器的配置设置。
	DatasetRun          string `gorm:"column:dataset_run;comment:知识库运行状态" json:"dataset_run"`                          // 知识库运行状态
	DatasetSize         int64  `gorm:"column:dataset_size;not null;comment:知识库尺寸大小" json:"dataset_size"`               // 知识库尺寸大小
	DatasetThumbnail    string `gorm:"column:dataset_thumbnail;comment:知识库锁略图" json:"dataset_thumbnail"`               // 知识库锁略图
	DatasetType         string `gorm:"column:dataset_type;comment:知识库类型" json:"dataset_type"`                          // 知识库类型
	DeptID              int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                     // 部门ID
	State               bool   `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`                               // 操作状态（0正常 -1删除）
	baize.BaseEntity
}

// PutDocumentRequest 更新文档请求
type PutDocumentRequest struct {
	Name         string                 `json:"name"`
	ChunkMethod  string                 `json:"chunk_method"`
	ParserConfig *ParserConfig          `json:"parser_config,omitempty"`
	MetaFields   map[string]interface{} `json:"meta_fields,omitempty"`
}

type PutDocumentResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ListDocumentRequest 文档列表请求
type ListDocumentRequest struct {
	Keywords     string `form:"keywords"`
	DatasetId    int64  `form:"dataset_id,string"`
	DocumentId   string `form:"document_id"`
	DocumentName string `form:"document_name"`
	baize.BaseEntityDQL
}

// DeleteDocumentRequest 删除文档
type DeleteDocumentRequest struct {
	DatasetId   int64    `json:"dataset_id,string"`
	DocumentIds []string `json:"document_ids,string"`
}

// DeleteDocumentApiRequest 删除文档
type DeleteDocumentApiRequest struct {
	Ids []string `json:"ids"`
}

type DeleteDocumentApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ParseDocumentApiRequest 删除文档
type ParseDocumentApiRequest struct {
	DatasetId   int64    `json:"dataset_id,string"`
	DocumentIds []string `json:"document_ids"`
}

// ParseDocumentApiResponse 解析文档详情
type ParseDocumentApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type DocumentData struct {
	ChunkCount      int64        `json:"chunk_count"`
	CreateDate      string       `json:"create_date"`
	CreateTime      int64        `json:"create_time"`
	CreatedBy       string       `json:"created_by"`
	Id              string       `json:"id"`
	KnowledgebaseId string       `json:"knowledgebase_id"`
	Location        string       `json:"location"`
	Name            string       `json:"name"`
	ParserConfig    ParserConfig `json:"parser_config"`
	ChunkMethod     string       `json:"chunk_method"`
	ProcessBeginAt  interface{}  `json:"process_begin_at"`
	ProcessDuation  float64      `json:"process_duation"`
	Progress        float64      `json:"progress"`
	ProgressMsg     string       `json:"progress_msg"`
	Run             string       `json:"run"`
	Size            int          `json:"size"`
	SourceType      string       `json:"source_type"`
	Status          string       `json:"status"`
	Thumbnail       string       `json:"thumbnail"`
	TokenCount      int          `json:"token_count"`
	Type            string       `json:"type"`
	UpdateDate      string       `json:"update_date"`
	UpdateTime      int64        `json:"update_time"`
}

// ListDocumentResponseData 文档列表请求
type ListDocumentResponseData struct {
	Code int `json:"code"`
	Data struct {
		Docs  []DocumentData `json:"docs"`
		Total int            `json:"total"`
	} `json:"data"`
	Message string `json:"message"`
}

type ListDocumentData struct {
	Rows  []*SysDatasetDocumentData `json:"rows"`
	Total int64                     `json:"total"`
}

// SysDatasetDocumentData 物料出库管理
type SysDatasetDocumentData struct {
	DocumentID          int64        `gorm:"column:document_id;not null;comment:文档id" json:"document_id,string"`             // 文档id
	DatasetID           int64        `gorm:"column:dataset_id;primaryKey;comment:数据集id" json:"dataset_id,string"`            // 数据集id
	DatasetChunkMethod  string       `gorm:"column:dataset_chunk_method;comment:要创建的数据集的分块方法" json:"dataset_chunk_method"`   // 要创建的数据集的分块方法
	DatasetCreatedBy    string       `gorm:"column:dataset_created_by;comment:创建人" json:"dataset_created_by"`                // 创建人
	DatasetDocumentUUID string       `gorm:"column:dataset_document_uuid;comment:文档id" json:"dataset_document_uuid"`         // 文档id
	DatasetDatasetUUID  string       `gorm:"column:dataset_dataset_uuid;comment:知识库id" json:"dataset_dataset_uuid"`          // 知识库id
	DatasetLanguage     string       `gorm:"column:dataset_language;comment:知识库语言" json:"dataset_language"`                  // 知识库语言
	DatasetLocation     string       `gorm:"column:dataset_location;comment:知识库定位" json:"dataset_location"`                  // 知识库定位
	DatasetName         string       `gorm:"column:dataset_name;comment:要创建的数据集的唯一名称" json:"dataset_name"`                   // 要创建的数据集的唯一名称
	DatasetParserConfig ParserConfig `gorm:"column:dataset_parser_config;comment:数据集解析器的配置设置。" json:"dataset_parser_config"` // 数据集解析器的配置设置。
	DatasetRun          string       `gorm:"column:dataset_run;comment:知识库运行状态" json:"dataset_run"`                          // 知识库运行状态
	DatasetSize         int64        `gorm:"column:dataset_size;not null;comment:知识库尺寸大小" json:"dataset_size"`               // 知识库尺寸大小
	DatasetThumbnail    string       `gorm:"column:dataset_thumbnail;comment:知识库锁略图" json:"dataset_thumbnail"`               // 知识库锁略图
	DatasetType         string       `gorm:"column:dataset_type;comment:知识库类型" json:"dataset_type"`                          // 知识库类型
	ChunkCount          int64        `gorm:"-" json:"chunk_count"`                                                           // 知识库类型
	Progress            float64      `gorm:"-" json:"progress"`
	DeptID              int64        `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`       // 部门ID
	State               bool         `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
	Run                 string       `json:"run"`
	baize.BaseEntity
}

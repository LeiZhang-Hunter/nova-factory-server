package aiDataSetModels

import "nova-factory-server/app/baize"

// ChunkListReq chunk列表集
type ChunkListReq struct {
	DocumentUuid string `json:"document_uuid"  form:"document_uuid"`
	DatasetUuid  string `json:"dataset_uuid"  form:"dataset_uuid"`
	Keywords     string `json:"keywords,string" form:"keywords"`
	Id           string `json:"id,string" form:"id"`
	baize.BaseEntityDQL
}

type ChunkData struct {
	Available         bool     `json:"available"`
	Content           string   `json:"content"`
	DocnmKwd          string   `json:"docnm_kwd"`
	DocumentId        string   `json:"document_id"`
	Id                string   `json:"id"`
	ImageId           string   `json:"image_id"`
	ImportantKeywords []string `json:"important_keywords"`
	Positions         []string `json:"positions"`
}

type ChunkDoc struct {
	ChunkCount     int          `json:"chunk_count"`
	ChunkMethod    string       `json:"chunk_method"`
	CreateDate     string       `json:"create_date"`
	CreateTime     int64        `json:"create_time"`
	CreatedBy      string       `json:"created_by"`
	DatasetId      string       `json:"dataset_id"`
	Id             string       `json:"id"`
	Location       string       `json:"location"`
	Name           string       `json:"name"`
	ParserConfig   ParserConfig `json:"parser_config"`
	ProcessBeginAt string       `json:"process_begin_at"`
	ProcessDuation float64      `json:"process_duation"`
	Progress       float64      `json:"progress"`
	ProgressMsg    string       `json:"progress_msg"`
	Run            string       `json:"run"`
	Size           int          `json:"size"`
	SourceType     string       `json:"source_type"`
	Status         string       `json:"status"`
	Thumbnail      string       `json:"thumbnail"`
	TokenCount     int          `json:"token_count"`
	Type           string       `json:"type"`
	UpdateDate     string       `json:"update_date"`
	UpdateTime     int64        `json:"update_time"`
}
type ChunkListResponse struct {
	Code int `json:"code"`
	Data struct {
		Chunks []*ChunkData `json:"chunks"`
		Doc    *ChunkDoc    `json:"doc"`
		Total  int          `json:"total"`
	} `json:"data"`
	Message string `json:"message"`
}

// AddChunkReq 添加chunk请求
type AddChunkReq struct {
	DocumentUuid      string   `json:"document_uuid"  form:"document_uuid"`
	DatasetUuid       string   `json:"dataset_uuid"  form:"dataset_uuid"`
	Content           string   `json:"content"  form:"content"`
	ImportantKeywords []string `json:"important_keywords"  form:"important_keywords"`
	Questions         []string `json:"questions,string" form:"questions"`
}

type AddChunkData struct {
	Content           string   `json:"content"`
	CreateTime        string   `json:"create_time"`
	CreateTimestamp   float64  `json:"create_timestamp"`
	DatasetId         string   `json:"dataset_id"`
	DocumentId        string   `json:"document_id"`
	Id                string   `json:"id"`
	ImportantKeywords []string `json:"important_keywords,omitempty"`
	Questions         []string `json:"questions,omitempty"`
}

type AddChunkResponse struct {
	Code int `json:"code"`
	Data struct {
		Chunk *AddChunkData `json:"chunk"`
	} `json:"data"`
	Message string `json:"message"`
}

// RemoveChunkReq 移除chunk
type RemoveChunkReq struct {
	DatasetUuid  string   `json:"dataset_uuid"  form:"dataset_uuid"`
	DocumentUuid string   `json:"document_uuid"  form:"document_uuid"`
	ChunkIds     []string `json:"chunk_ids"  form:"chunk_ids"`
}

// UpdateChunkReq 更新chunk请求
type UpdateChunkReq struct {
	DocumentUuid      string   `json:"document_uuid"  form:"document_uuid"`
	DatasetUuid       string   `json:"dataset_uuid"  form:"dataset_uuid"`
	ChunkUuid         string   `json:"chunk_id"  form:"chunk_id"`
	Content           string   `json:"content"  form:"content"`
	ImportantKeywords []string `json:"important_keywords"  form:"important_keywords"`
	Questions         []string `json:"questions,string" form:"questions"`
}

type UpdateChunkResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RetrievalListReq 检索req
type RetrievalListReq struct {
	Question               string   `json:"question"  form:"question"`
	DatasetIds             []string `json:"dataset_uuids"  form:"dataset_uuids"`
	DocumentIds            []string `json:"document_uuids"  form:"document_uuids"`
	Page                   int      `json:"page"  form:"page"`
	PageSize               int      `json:"page_size"  form:"page_size"`
	SimilarityThreshold    float64  `json:"similarity_threshold"  form:"similarity_threshold"`
	VectorSimilarityWeight float64  `json:"vector_similarity_weight"  form:"vector_similarity_weight"`
	TopK                   int      `json:"top_k"  form:"top_k"`
	RerankId               string   `json:"rerank_id"  form:"rerank_id"`
	Keyword                bool     `json:"keyword"  form:"keyword"`
	Highlight              bool     `json:"highlight"  form:"highlight"`
}

type RetrievalApiListReq struct {
	Question               string   `json:"question,omitempty"  form:"question"`
	DatasetIds             []string `json:"dataset_ids"  form:"dataset_uuids"`
	DocumentIds            []string `json:"document_ids"  form:"document_uuids"`
	Page                   int      `json:"page"  form:"page"`
	PageSize               int      `json:"page_size"  form:"page_size"`
	SimilarityThreshold    float64  `json:"similarity_threshold,omitempty"  form:"similarity_threshold"`
	VectorSimilarityWeight float64  `json:"vector_similarity_weight,omitempty"  form:"vector_similarity_weight"`
	TopK                   int      `json:"top_k,omitempty"  form:"top_k"`
	RerankId               string   `json:"rerank_id,omitempty"  form:"rerank_id"`
	Keyword                bool     `json:"keyword,omitempty"  form:"keyword"`
	Highlight              bool     `json:"highlight,omitempty"  form:"highlight"`
}

func (api *RetrievalApiListReq) Of(req *RetrievalListReq) {
	api.Question = req.Question
	api.DatasetIds = req.DatasetIds
	api.DocumentIds = req.DocumentIds
	api.Page = req.Page
	api.PageSize = req.PageSize
	api.SimilarityThreshold = req.SimilarityThreshold
	api.VectorSimilarityWeight = req.VectorSimilarityWeight
	api.TopK = req.TopK
	api.RerankId = req.RerankId
	api.Keyword = req.Keyword
	api.Highlight = req.Highlight
}

type RetrievalApiData struct {
	Content           string   `json:"content"`
	ContentLtks       string   `json:"content_ltks"`
	DocumentId        string   `json:"document_id"`
	DocumentKeyword   string   `json:"document_keyword"`
	Highlight         string   `json:"highlight"`
	Id                string   `json:"id"`
	ImageId           string   `json:"image_id"`
	ImportantKeywords []string `json:"important_keywords"`
	KbId              string   `json:"kb_id"`
	Positions         []string `json:"positions"`
	Similarity        float64  `json:"similarity"`
	TermSimilarity    float64  `json:"term_similarity"`
	VectorSimilarity  float64  `json:"vector_similarity"`
}

type DocAggs struct {
	Count   int    `json:"count"`
	DocId   string `json:"doc_id"`
	DocName string `json:"doc_name"`
}

type RetrievalApiListResponse struct {
	Code int `json:"code"`
	Data struct {
		Chunks  []RetrievalApiData `json:"chunks"`
		DocAggs interface{}        `json:"doc_aggs"`
		Total   int                `json:"total"`
	} `json:"data"`
	Message string `json:"message"`
}

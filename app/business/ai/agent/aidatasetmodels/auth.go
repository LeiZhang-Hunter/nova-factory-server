package aidatasetmodels

// RagflowAuthOutput defines the output of the Ragflow auth tool.
type RagflowAuthOutput struct {
	Configured  bool     `json:"configured"`
	DatasetIDs  []string `json:"dataset_ids,omitempty"`
	DocumentIDs []string `json:"document_ids,omitempty"`
	Message     string   `json:"message,omitempty"`
}

// DatasetAccessData 聚合后的知识库与文档访问数据。
type DatasetAccessData struct {
	DatasetIDs    []string `json:"-"`
	DatasetUuIDs  []string `json:"dataset_ids"`
	DocumentIDs   []string `json:"-"`
	DocumentUuIDs []string `json:"document_ids"`
}

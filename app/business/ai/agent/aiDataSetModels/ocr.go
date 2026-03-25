package aiDataSetModels

type OCRPageResult struct {
	Page int    `json:"page"`
	Text string `json:"text"`
}

type OCRExtractResult struct {
	FileName string           `json:"file_name"`
	Text     string           `json:"text"`
	Pages    []*OCRPageResult `json:"pages"`
}

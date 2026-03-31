package aiDataSetModels

import "nova-factory-server/app/baize"

type SysAiLLMSetting struct {
	ID        int64  `gorm:"column:id;primaryKey;comment:主键ID" json:"id,string"`
	Name      string `gorm:"column:name;comment:名称" json:"name"`
	PublicKey string `gorm:"column:public_key;comment:公钥" json:"public_key"`
	LlmID     string `gorm:"column:llm_id;comment:llm_id" json:"llm_id"`
	EmbdID    string `gorm:"column:embd_id;comment:embd_id" json:"embd_id"`
	AsrID     string `gorm:"column:asr_id;comment:asr_id" json:"asr_id"`
	Img2txtID string `gorm:"column:img2txt_id;comment:img2txt_id" json:"img2txt_id"`
	RerankID  string `gorm:"column:rerank_id;comment:rerank_id" json:"rerank_id"`
	TtsID     string `gorm:"column:tts_id;comment:tts_id" json:"tts_id"`
	ParserIDs string `gorm:"column:parser_ids;comment:parser_ids" json:"parser_ids"`
	Credit    int64  `gorm:"column:credit;comment:额度" json:"credit"`
	Status    string `gorm:"column:status;comment:状态" json:"status"`
	DeptID    int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id,string"`
	baize.BaseEntity
	State int32 `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"`
}

type SetSysAiLLMSetting struct {
	ID        int64  `json:"id,string"`
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
	LlmID     string `json:"llm_id"`
	EmbdID    string `json:"embd_id"`
	AsrID     string `json:"asr_id"`
	Img2txtID string `json:"img2txt_id"`
	RerankID  string `json:"rerank_id"`
	TtsID     string `json:"tts_id"`
	ParserIDs string `json:"parser_ids"`
	Credit    int64  `json:"credit"`
	Status    string `json:"status"`
}

type GetSysAiLLMSettingReq struct {
	ID int64 `form:"id" json:"id,string"`
}

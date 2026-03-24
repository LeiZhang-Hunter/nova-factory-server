package ai

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"nova-factory-server/app/utils/snowflake"

	"gorm.io/gorm"
)

type FactoryBootstrap struct {
	db            *gorm.DB
	providerTable string
	llmTable      string
}

type factoryPayload struct {
	FactoryLLMInfos []*factoryItem `json:"factory_llm_infos"`
}

type factoryItem struct {
	Name   string     `json:"name"`
	Logo   string     `json:"logo"`
	Tags   string     `json:"tags"`
	Status string     `json:"status"`
	Rank   string     `json:"rank"`
	LLM    []*llmItem `json:"llm"`
}

type llmItem struct {
	LLMName   string `json:"llm_name"`
	Tags      string `json:"tags"`
	MaxTokens int64  `json:"max_tokens"`
	ModelType string `json:"model_type"`
	IsTools   bool   `json:"is_tools"`
}

func NewFactoryBootstrap(db *gorm.DB) *FactoryBootstrap {
	return &FactoryBootstrap{
		db:            db,
		providerTable: "ai_model_provider",
		llmTable:      "ai_llm",
	}
}

func (f *FactoryBootstrap) Init() error {
	configPath, err := resolveFactoryConfigPath()
	if err != nil {
		return err
	}
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	payload := new(factoryPayload)
	if err = json.Unmarshal(raw, payload); err != nil {
		return err
	}
	return f.db.WithContext(context.Background()).Transaction(func(tx *gorm.DB) error {
		for _, item := range payload.FactoryLLMInfos {
			if item == nil || item.Name == "" {
				continue
			}
			now := time.Now()
			status, _ := strconv.Atoi(item.Status)
			rank, _ := strconv.Atoi(item.Rank)
			providerID, err := f.upsertProvider(tx, item, int32(status), int32(rank), now)
			if err != nil {
				return err
			}
			_ = providerID
			if err = f.upsertLLM(tx, item, item.Status, now); err != nil {
				return err
			}
		}
		return nil
	})
}

func (f *FactoryBootstrap) upsertProvider(tx *gorm.DB, item *factoryItem, status int32, rank int32, now time.Time) (int64, error) {
	type providerRow struct {
		ID    int64 `db:"id"`
		State int32 `db:"state"`
	}
	var row providerRow
	err := tx.Table(f.providerTable).Select("id", "state").Where("name = ?", item.Name).Where("state = 0").Take(&row).Error
	if err == nil {
		update := map[string]interface{}{
			"logo":           item.Logo,
			"tags":           item.Tags,
			"status":         status,
			"rank":           rank,
			"update_time_db": now,
		}
		return row.ID, tx.Table(f.providerTable).Where("id = ?", row.ID).Updates(update).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	err = tx.Table(f.providerTable).Select("id", "state").Where("name = ?", item.Name).Take(&row).Error
	if err == nil {
		update := map[string]interface{}{
			"name":           item.Name,
			"logo":           item.Logo,
			"tags":           item.Tags,
			"status":         status,
			"rank":           rank,
			"state":          0,
			"update_time_db": now,
		}
		return row.ID, tx.Table(f.providerTable).Where("id = ?", row.ID).Updates(update).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	id := snowflake.GenID()
	create := map[string]interface{}{
		"id":             id,
		"name":           item.Name,
		"logo":           item.Logo,
		"tags":           item.Tags,
		"status":         status,
		"rank":           rank,
		"state":          0,
		"create_time_db": now,
		"update_time_db": now,
	}
	return id, tx.Table(f.providerTable).Create(create).Error
}

func (f *FactoryBootstrap) upsertLLM(tx *gorm.DB, item *factoryItem, status string, now time.Time) error {
	type llmRow struct {
		FID     string `db:"fid"`
		LlmName string `db:"llm_name"`
	}
	for _, llm := range item.LLM {
		if llm == nil || llm.LLMName == "" {
			continue
		}
		toolFlag := 0
		if llm.IsTools {
			toolFlag = 1
		}
		var row llmRow
		err := tx.Table(f.llmTable).Select("fid", "llm_name").
			Where("fid = ?", item.Name).Where("llm_name = ?", llm.LLMName).Take(&row).Error
		if err == nil {
			update := map[string]interface{}{
				"model_type":     llm.ModelType,
				"max_tokens":     llm.MaxTokens,
				"tags":           llm.Tags,
				"is_tools":       toolFlag,
				"status":         status,
				"state":          0,
				"update_time_db": now,
			}
			if err = tx.Table(f.llmTable).Where("fid = ?", item.Name).Where("llm_name = ?", llm.LLMName).Updates(update).Error; err != nil {
				return err
			}
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		create := map[string]interface{}{
			"fid":            item.Name,
			"llm_name":       llm.LLMName,
			"model_type":     llm.ModelType,
			"max_tokens":     llm.MaxTokens,
			"tags":           llm.Tags,
			"is_tools":       toolFlag,
			"status":         status,
			"state":          0,
			"create_time_db": now,
			"update_time_db": now,
		}
		if err = tx.Table(f.llmTable).Create(create).Error; err != nil {
			return err
		}
	}
	return nil
}

func resolveFactoryConfigPath() (string, error) {
	candidates := []string{
		"config/llm_factories.json",
		"../config/llm_factories.json",
		"../../config/llm_factories.json",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", os.ErrNotExist
}

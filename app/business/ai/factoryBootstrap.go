package ai

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"

	"gorm.io/gorm"
)

type FactoryBootstrap struct {
	db          *gorm.DB
	providerDao aidatasetdao.IAiModelProviderDao
	llmDao      aidatasetdao.IAiLLMDao
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

func NewFactoryBootstrap(db *gorm.DB, providerDao aidatasetdao.IAiModelProviderDao, llmDao aidatasetdao.IAiLLMDao) *FactoryBootstrap {
	return &FactoryBootstrap{
		db:          db,
		providerDao: providerDao,
		llmDao:      llmDao,
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
			providerID, err := f.providerDao.UpsertFactoryProvider(tx, &aidatasetmodels.FactoryProviderUpsert{
				Name: item.Name,
				Logo: item.Logo,
				Tags: item.Tags,
			}, int32(status), int32(rank), now)
			if err != nil {
				return err
			}
			_ = providerID
			seedLLMs := make([]*aidatasetmodels.FactoryLLMUpsert, 0, len(item.LLM))
			for _, llm := range item.LLM {
				if llm == nil {
					continue
				}
				seedLLMs = append(seedLLMs, &aidatasetmodels.FactoryLLMUpsert{
					LLMName:   llm.LLMName,
					Tags:      llm.Tags,
					MaxTokens: llm.MaxTokens,
					ModelType: llm.ModelType,
					IsTools:   llm.IsTools,
				})
			}
			if err = f.llmDao.UpsertFactoryLLMs(tx, item.Name, seedLLMs, item.Status, now); err != nil {
				return err
			}
		}
		return nil
	})
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

package aiDataSetDaoImpl

import (
	"testing"

	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
)

func TestFilterEmbeddingProviders(t *testing.T) {
	rows := []*aidatasetmodels.SysAiModelProvider{
		{
			Name: "provider-a",
			LLMs: []*aidatasetmodels.SysAiLLM{
				{Fid: "provider-a", LlmName: "chat-model", ModelType: "chat"},
				{Fid: "provider-a", LlmName: "embedding-model", ModelType: "embedding"},
				{Fid: "provider-a", LlmName: "embedding-model-2", ModelType: " EMBEDDING "},
			},
		},
		{
			Name: "provider-b",
			LLMs: []*aidatasetmodels.SysAiLLM{
				{Fid: "provider-b", LlmName: "rerank-model", ModelType: "rerank"},
			},
		},
	}

	filtered := filterEmbeddingProviders(rows)
	if len(filtered) != 1 {
		t.Fatalf("unexpected provider count: %d", len(filtered))
	}
	if filtered[0].Name != "provider-a" {
		t.Fatalf("unexpected provider name: %s", filtered[0].Name)
	}
	if len(filtered[0].LLMs) != 2 {
		t.Fatalf("unexpected embedding llm count: %d", len(filtered[0].LLMs))
	}
	if filtered[0].LLMs[0].LlmName != "embedding-model" {
		t.Fatalf("unexpected first embedding llm: %s", filtered[0].LLMs[0].LlmName)
	}
	if filtered[0].LLMs[1].LlmName != "embedding-model-2" {
		t.Fatalf("unexpected second embedding llm: %s", filtered[0].LLMs[1].LlmName)
	}
}

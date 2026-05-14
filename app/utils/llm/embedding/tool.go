package embedding

import (
	"fmt"
	"strings"
)

// ParseProviderModelID 解析 modelID@providerID 格式的 embedding 模型标识。
func ParseProviderModelID(value string) (providerID string, modelID string, err error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", "", fmt.Errorf("embedding model config is empty")
	}

	modelID, providerID, found := strings.Cut(value, "@")
	if !found {
		return "", "", fmt.Errorf("invalid embedding model config %q: missing @", value)
	}

	modelID = strings.TrimSpace(modelID)
	providerID = strings.TrimSpace(providerID)
	if modelID == "" || providerID == "" {
		return "", "", fmt.Errorf("invalid embedding model config %q: modelID or providerID is empty", value)
	}

	return providerID, modelID, nil
}

package aidatasetcontroller

import (
	"sync"
	"time"
)

type AgentRegistryInfo struct {
	GatewayID  int64
	Username   string
	RegTime    time.Time
	LastActive time.Time
}

type AgentRegistryManager struct {
	mu       sync.RWMutex
	registry map[int64]*AgentRegistryInfo
}

func NewAgentRegistryManager() *AgentRegistryManager {
	return &AgentRegistryManager{
		registry: make(map[int64]*AgentRegistryInfo),
	}
}

func (m *AgentRegistryManager) Register(gatewayID int64, username string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	m.registry[gatewayID] = &AgentRegistryInfo{
		GatewayID:  gatewayID,
		Username:   username,
		RegTime:    now,
		LastActive: now,
	}
}

func (m *AgentRegistryManager) Unregister(gatewayID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.registry, gatewayID)
}

func (m *AgentRegistryManager) IsRegistered(gatewayID int64) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.registry[gatewayID]
	return ok
}

func (m *AgentRegistryManager) GetInfo(gatewayID int64) (*AgentRegistryInfo, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	info, ok := m.registry[gatewayID]
	return info, ok
}

func (m *AgentRegistryManager) UpdateLastActive(gatewayID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if info, ok := m.registry[gatewayID]; ok {
		info.LastActive = time.Now()
	}
}

func (m *AgentRegistryManager) GetAllGatewayIDs() []int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ids := make([]int64, 0, len(m.registry))
	for id := range m.registry {
		ids = append(ids, id)
	}
	return ids
}

func (m *AgentRegistryManager) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.registry)
}

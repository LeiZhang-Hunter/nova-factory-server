package aidatasetcontroller

import (
	"errors"
	client "nova-factory-server/app/utils/grpc/confighotload/v1"
	"sync"
)

type AgentConfigManager struct {
	mu      sync.RWMutex
	clients map[string]client.AgentControllerService_WatchAgentChangesServer
}

func NewAgentConfigManager() *AgentConfigManager {
	return &AgentConfigManager{
		clients: make(map[string]client.AgentControllerService_WatchAgentChangesServer),
	}
}

func (m *AgentConfigManager) AddClient(node string, stream client.AgentControllerService_WatchAgentChangesServer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[node] = stream
}

func (m *AgentConfigManager) DeleteClient(node string, stream client.AgentControllerService_WatchAgentChangesServer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	current, ok := m.clients[node]
	if !ok {
		return
	}
	if current == stream {
		delete(m.clients, node)
	}
}

func (m *AgentConfigManager) GetClient(node string) (client.AgentControllerService_WatchAgentChangesServer, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	stream, ok := m.clients[node]
	return stream, ok
}

func (m *AgentConfigManager) GetNodes() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	nodes := make([]string, 0, len(m.clients))
	for node := range m.clients {
		nodes = append(nodes, node)
	}
	return nodes
}

func (m *AgentConfigManager) PublishToNode(node string, message *client.ConfigChangeEvent) (uint32, error) {
	stream, ok := m.GetClient(node)
	if !ok {
		return 0, errors.New("agent config subscriber not found")
	}
	if err := stream.Send(message); err != nil {
		m.DeleteClient(node, stream)
		return 0, errors.New("agent config subscriber send failed")
	}
	return 1, nil
}

func (m *AgentConfigManager) BroadcastToAllNodes(message *client.ConfigChangeEvent) (uint32, []string) {
	nodes := m.GetNodes()
	deliveredNodes := make([]string, 0, len(nodes))
	var totalDelivered uint32 = 0
	for _, node := range nodes {
		delivered, err := m.PublishToNode(node, message)
		if err != nil {
			continue
		}
		if delivered > 0 {
			deliveredNodes = append(deliveredNodes, node)
			totalDelivered += delivered
		}
	}
	return totalDelivered, deliveredNodes
}

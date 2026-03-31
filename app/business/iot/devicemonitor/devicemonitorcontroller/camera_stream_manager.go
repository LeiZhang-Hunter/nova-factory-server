package devicemonitorcontroller

import (
	"errors"
	"sync"

	client "github.com/novawatcher-io/nova-factory-payload/camera/v1"
)

type CameraSubscribeManager struct {
	mu      sync.RWMutex
	clients map[string]client.CameraService_WebrtcSubscribeServer
}

func NewCameraSubscribeManager() *CameraSubscribeManager {
	return &CameraSubscribeManager{
		clients: make(map[string]client.CameraService_WebrtcSubscribeServer),
	}
}

func (m *CameraSubscribeManager) AddClient(node string, stream client.CameraService_WebrtcSubscribeServer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[node] = stream
}

func (m *CameraSubscribeManager) DeleteClient(node string, stream client.CameraService_WebrtcSubscribeServer) {
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

func (m *CameraSubscribeManager) GetClient(node string) (client.CameraService_WebrtcSubscribeServer, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	stream, ok := m.clients[node]
	return stream, ok
}

func (m *CameraSubscribeManager) GetNodes() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	nodes := make([]string, 0, len(m.clients))
	for node := range m.clients {
		nodes = append(nodes, node)
	}
	return nodes
}

func (m *CameraSubscribeManager) PublishToNode(node string, message *client.SubscribeMessage) (uint32, error) {
	stream, ok := m.GetClient(node)
	if !ok {
		return 0, errors.New("camera subscriber not found")
	}
	if err := stream.Send(message); err != nil {
		m.DeleteClient(node, stream)
		return 0, errors.New("camera subscriber send failed")
	}
	return 1, nil
}

func (m *CameraSubscribeManager) BroadcastToAllNodes(message *client.SubscribeMessage) (uint32, []string) {
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

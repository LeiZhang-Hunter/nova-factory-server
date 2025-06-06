package daemonizeServiceImpl

import (
	v1 "github.com/novawatcher-io/nova-factory-payload/daemonize/grpc/v1"
	"sync"
)

type ManagerServiceImpl struct {
	clients sync.Map
	stop    bool
}

func NewManagerServiceImpl() *ManagerServiceImpl {
	return &ManagerServiceImpl{}
}

// AddClient 添加agent客户端连接
func (s *ManagerServiceImpl) AddClient(clientId uint64, stream v1.AgentControllerService_AgentOperateServer) {
	s.clients.Store(clientId, stream)
}

// DeleteClient 删除agent客户端连接
func (s *ManagerServiceImpl) DeleteClient(clientId uint64) {
	s.clients.Delete(clientId)
}

// Stop 设置停止标识
func (s *ManagerServiceImpl) Stop() {
	s.stop = true
}

// IsStopped 是否停止
func (s *ManagerServiceImpl) IsStopped() bool {
	return s.stop
}

func (s *ManagerServiceImpl) getClient(objectId uint64) (stream v1.AgentControllerService_AgentOperateServer) {
	value, ok := s.clients.Load(objectId)
	if !ok {
		return nil
	}
	return value.(v1.AgentControllerService_AgentOperateServer)
}

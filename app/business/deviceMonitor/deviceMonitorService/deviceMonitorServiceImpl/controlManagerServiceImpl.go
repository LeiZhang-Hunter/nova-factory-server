package deviceMonitorServiceImpl

import (
	v1 "github.com/novawatcher-io/nova-factory-payload/control/v1"
	"sync"
)

type ControlManagerServiceImpl struct {
	clients sync.Map
	stop    bool
}

func NewControlManagerServiceImpl() *ControlManagerServiceImpl {
	return &ControlManagerServiceImpl{}
}

// AddClient 添加agent客户端连接
func (s *ControlManagerServiceImpl) AddClient(clientId uint64, stream v1.ControlService_OperateServer) {
	s.clients.Store(clientId, stream)
}

// DeleteClient 删除agent客户端连接
func (s *ControlManagerServiceImpl) DeleteClient(clientId uint64) {
	s.clients.Delete(clientId)
}

// Stop 设置停止标识
func (s *ControlManagerServiceImpl) Stop() {
	s.stop = true
}

// IsStopped 是否停止
func (s *ControlManagerServiceImpl) IsStopped() bool {
	return s.stop
}

func (s *ControlManagerServiceImpl) getClient(objectId uint64) (stream v1.ControlService_OperateServer) {
	value, ok := s.clients.Load(objectId)
	if !ok {
		return nil
	}
	return value.(v1.ControlService_OperateServer)
}

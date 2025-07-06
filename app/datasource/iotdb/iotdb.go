package iotdb

import (
	"github.com/apache/iotdb-client-go/client"
	"go.uber.org/zap"
)

type IotDb struct {
	pool client.SessionPool
}

func NewIotDb() *IotDb {
	config := &client.PoolConfig{
		Host:     "192.168.2.100",
		Port:     "6667",
		UserName: "root",
		Password: "root",
	}
	pool := client.NewSessionPool(config, 3, 60000, 60000, false)
	//defer sessionPool.Close()

	session, err := pool.GetSession()
	if err != nil {
		zap.L().Error("get session", zap.Error(err))
		return nil
	}

	defer pool.PutBack(session)

	session.ExecuteStatement("create device template nova_device_template (template_id INT64, data_id INT64, value FLOAT)")

	return &IotDb{
		pool: pool,
	}
}

func (i *IotDb) GetSession() (client.Session, error) {
	return i.pool.GetSession()
}

func (i *IotDb) PutSession(session client.Session) {
	i.pool.PutBack(session)
}

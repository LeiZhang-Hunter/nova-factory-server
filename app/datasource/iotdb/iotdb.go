package iotdb

import (
	"errors"
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
	db := &IotDb{}
	pool := client.NewSessionPool(config, 3, 60000, 60000, false)
	//defer sessionPool.Close()
	db.pool = pool
	session, err := pool.GetSession()
	if err != nil {
		zap.L().Error("get session", zap.Error(err))
		return db
	}

	defer pool.PutBack(session)

	session.ExecuteStatement("create device template nova_device_template ALIGNED (value DOUBLE)")

	return db
}

func (i *IotDb) GetSession() (client.Session, error) {
	if i == nil {
		return client.Session{}, errors.New("<UNK>")
	}
	return i.pool.GetSession()
}

func (i *IotDb) PutSession(session client.Session) {
	i.pool.PutBack(session)
}

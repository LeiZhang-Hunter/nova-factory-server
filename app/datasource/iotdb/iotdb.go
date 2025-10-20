package iotdb

import (
	"errors"
	"fmt"
	"github.com/apache/iotdb-client-go/client"
	"go.uber.org/zap"
	iotdb2 "nova-factory-server/app/constant/iotdb"
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

	// 创建设备数据采集模板
	session.ExecuteStatement(fmt.Sprintf("create device template %s ALIGNED (value DOUBLE)", iotdb2.NOVA_DEVICE_TEMPLATE))
	// 创建设备运行时间统计模板
	session.ExecuteStatement(fmt.Sprintf("create device template %s ALIGNED (duration INT64, status INT64)", iotdb2.NOVA_DEVICE_RUN_TEMPLATE))

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

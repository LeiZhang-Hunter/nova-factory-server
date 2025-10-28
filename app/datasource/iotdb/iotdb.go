package iotdb

import (
	"errors"
	"fmt"
	"github.com/apache/iotdb-client-go/client"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	iotdb2 "nova-factory-server/app/constant/iotdb"
)

type IotDb struct {
	pool client.SessionPool
}

func NewIotDb() *IotDb {
	type IotDbConfig struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		UserName string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}
	// 把读取到的配置信息反序列化到 Conf 变量中
	var d IotDbConfig
	if err := viper.UnmarshalKey("iotdb", &d); err != nil {
		panic(err)
	}
	config := &client.PoolConfig{
		Host:     d.Host,
		Port:     d.Port,
		UserName: d.UserName,
		Password: d.Password,
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

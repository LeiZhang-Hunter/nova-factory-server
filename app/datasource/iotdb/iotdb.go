package iotdb

import (
	"errors"
	"sync"
	"time"

	"github.com/apache/iotdb-client-go/client"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type IotDb struct {
	pool *client.SessionPool
	mtx  sync.Mutex
}

func NewIotDb() *IotDb {
	return &IotDb{}
}

// connect 连接数据库
func (i *IotDb) connect() *IotDb {
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

	pool := client.NewSessionPool(config, 3, 60000, 60000, false)
	//defer sessionPool.Close()
	i.pool = &pool
	var session client.Session
	var err error
	for {
		session, err = pool.GetSession()
		if err == nil {
			//defer pool.PutBack(session)
			return i
		}
		if err.Error() == "get session timeout" {
			zap.L().Error("get session error", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}
		if err != nil {
			zap.L().Error("get session error", zap.Error(err))
			i.pool.PutBack(session)
			if err.Error() == "sessionPool has closed" {
				return i
			}
			time.Sleep(5 * time.Second)
			continue
		}
	}

	return i
}

func (i *IotDb) GetSession() (client.Session, error) {
	if i == nil {
		return client.Session{}, errors.New("<UNK>")
	}
	if i.pool == nil {
		i.mtx.Lock()
		if i.pool == nil {
			i.connect()
		}
		i.mtx.Unlock()
	}
	session, err := i.pool.GetSession()
	if err == nil {
		return session, nil
	}

	if err.Error() == "get session timeout" {
		zap.L().Error("get session error", zap.Error(err))
		return client.Session{}, errors.New("sessionPool has closed")
	}
	if err != nil {
		zap.L().Error("get session error", zap.Error(err))
		i.pool.PutBack(session)
		zap.L().Error("get session", zap.Error(err))
		if err.Error() == "sessionPool has closed" {
			return session, nil

		}
		return client.Session{}, err
	}
	return session, nil
}

func (i *IotDb) PutSession(session client.Session) {
	i.pool.PutBack(session)
}

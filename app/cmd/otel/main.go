package main

import (
	"github.com/spf13/viper"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"net"
	"time"
)

func main() {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	// 设置为中国时区
	time.Local = location

	// Automatically set GOMAXPROCS to match Linux container CPU quota
	if _, err := maxprocs.Set(); err != nil {
		zap.L().Fatal("set maxprocs error: %v", zap.Error(err))
	}

	// 创建grpc服务
	s, cleanup, err := wireApp()
	if err != nil {
		panic(err)
	}
	defer cleanup()
	var host string
	if viper.GetString("metric.host") == "" {
		host = "0.0.0.0:6000"
	} else {
		host = viper.GetString("metric.host")
	}
	listen, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}

	zap.L().Info("start grpc server", zap.String("host", host))
	err = s.Serve(listen)
	if err != nil {
		panic(err)
	} // grpc服务启动

	return
}

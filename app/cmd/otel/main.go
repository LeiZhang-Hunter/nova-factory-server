package main

import (
	"context"
	"fmt"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net"
	"time"
)

type MetricServer struct {
	v1.UnimplementedDeviceReportServiceServer
}

func (s MetricServer) ReportContainer(context.Context, *v1.ExportMetricsServiceRequest) (*v1.NodeRes, error) {

	fmt.Println("article 进来了")
	return &v1.NodeRes{
		Code: 0,
	}, nil
}

func main() {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	// 设置为中国时区
	time.Local = location

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

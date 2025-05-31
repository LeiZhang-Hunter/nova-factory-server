package main

import (
	"context"
	"fmt"
	"net"
	v1 "nova-factory-server/app/pkg/metric/grpc/v1"
	"time"
)

type MetricServer struct {
	v1.UnimplementedDeviceReportServiceServer
}

func (s MetricServer) ReportContainer(context.Context, *v1.DeviceData) (*v1.NodeRes, error) {

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
	listen, err := net.Listen("tcp", "127.0.0.1:6002")
	fmt.Println("监听6002端口。。。")
	if err != nil {
		fmt.Println("网络错误")
	}
	s.Serve(listen) // grpc服务启动

	return
}

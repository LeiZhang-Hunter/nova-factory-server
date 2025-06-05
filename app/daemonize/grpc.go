package daemonize

import (
	"context"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"time"
)

func CreateGRpcServer() (grpcServer *grpcx.GrpcServer) {
	grpcServerConfig := grpcx.Server.NewConfig()
	value := viper.GetString("daemonize.address")
	var address string
	if value != "" {
		address = value
	} else {
		address = ":0"
	}
	glog.Infof(context.TODO(), "grpc server start")
	grpcServerConfig.Address = address
	grpcServerConfig.Options = append(grpcServerConfig.Options, []grpc.ServerOption{
		grpcx.Server.ChainUnary(
			grpcx.Server.UnaryValidate,
			CompanyValidate,
		),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    10 * time.Second, // 服务器在发送最后一个keepalive PING后将等待响应的时间
			Timeout: 8 * time.Second,  // 如果服务器在此时间内没有收到PING的响应，将关闭连接
		}),
	}...,
	)
	grpcServer = grpcx.Server.New(grpcServerConfig)
	valueDebug := viper.GetBool("daemonize.debug")
	debug := false
	if valueDebug != false {
		debug = valueDebug
	}
	if debug {
		reflection.Register(grpcServer.Server)
	}
	return
}

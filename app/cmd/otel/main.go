package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/soheilhy/cmux"
	"github.com/spf13/viper"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
	otelCleanup, err := registerMonitorOTLPGRPC(context.Background(), s)
	if err != nil {
		panic(err)
	}
	defer otelCleanup()

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

	mux := cmux.New(listen)
	grpcListener := mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := mux.Match(cmux.HTTP1Fast())
	httpServer := newMonitorHTTPServer()

	errCh := make(chan error, 3)
	go func() { errCh <- s.Serve(grpcListener) }()
	go func() { errCh <- httpServer.Serve(httpListener) }()
	go func() { errCh <- mux.Serve() }()

	zap.L().Info("start grpc/http server", zap.String("host", host), zap.String("http_endpoint", "/otel/monitor"))
	err = <-errCh
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

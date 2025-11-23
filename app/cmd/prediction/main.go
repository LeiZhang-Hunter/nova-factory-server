package main

import (
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
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
}

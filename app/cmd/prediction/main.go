package main

import (
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
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

	app, f, err := wireApp()
	if err != nil {
		return
	}
	defer f()
	app.Run()
	//等待一个INT或TERM信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.stop()
	return
}

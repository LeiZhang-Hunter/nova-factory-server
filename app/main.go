package main

import (
	"fmt"
	"nova-factory-server/app/setting"
	"time"
)

// @title baize
// @version 2.0.x
// @description baize接口文档

// @contact.name danny
// @contact.email zhao_402295440@126.com

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	// 设置为中国时区
	time.Local = location
	app, cleanup, err := wireApp()
	if err != nil {
		panic(err)
	}
	defer cleanup()
	//redisListener.StartRedisListener()
	//monitorServiceImpl.GetJobService().InitJobRun()
	app.Run(fmt.Sprintf(":%d", setting.Conf.Port))

}

package main

import (
	api2 "cloudpan/internal/api"
	conf2 "cloudpan/internal/conf"
	model2 "cloudpan/internal/model"
	"net/http"
)

func main() {
	// 前端服务
	go http.ListenAndServe(":8080", http.FileServer(http.Dir("frontend/")))

	// 读取配置文件
	conf2.Init()
	// 连接数据库
	model2.Init()

	// 装载后端路由
	r := api2.NewRouter()
	r.Run(":8081")
}

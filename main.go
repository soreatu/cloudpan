package main

import (
	api "cloudpan/internal/api"
	conf "cloudpan/internal/conf"
	model "cloudpan/internal/model"
	"net/http"
)

func main() {
	// 前端服务
	go http.ListenAndServe(":8080", http.FileServer(http.Dir("frontend/")))

	// 读取配置文件
	conf.Init()
	// 连接数据库
	model.Init()

	// 装载后端路由
	r := api.NewRouter()
	r.Run(":8081")
}

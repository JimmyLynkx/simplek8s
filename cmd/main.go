package main

import (
	"fmt"
	"go_code/simplek8s/server"
	"net/http"
)

func main() {
	// 初始化日志
	server.InitLogger()
	defer server.Logger.Sync()

	// 初始化路由
	httpHandler, err := server.InitializeRouter()
	if err != nil {
		server.Logger.Fatal(err.Error())
	}

	// 启动 HTTP 服务器
	addr := ":8080"
	server.Logger.Info(fmt.Sprintf("Server is running at %s...", addr))
	if err := http.ListenAndServe(addr, httpHandler); err != nil {
		server.Logger.Fatal(err.Error())
	}
}

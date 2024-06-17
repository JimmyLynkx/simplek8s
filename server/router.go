package server

import (
	"go_code/simplek8s/core/application/handler"
	"go_code/simplek8s/middleware"
	"net/http"
)

func NewRouter(clusterHandler *handler.ClusterHandler) http.Handler {
	mux := http.NewServeMux()
	RegisterRoutes(mux, clusterHandler)

	// 使用中间件的顺序：先恢复 panic，再记录日志，最后处理 JSON 响应
	recoverMiddleware := middleware.RecoverMiddleware(mux, Logger)
	loggingMiddleware := middleware.LoggingMiddleware(recoverMiddleware, Logger)
	finalHandler := middleware.JSONResponseMiddleware(loggingMiddleware)

	return finalHandler
}

func RegisterRoutes(mux *http.ServeMux, clusterHandler *handler.ClusterHandler) {
	// 添加路由，并将请求通过中间件处理
	mux.Handle("/cluster/add", http.HandlerFunc(clusterHandler.AddCluster))
	mux.Handle("/deployment/create", http.HandlerFunc(clusterHandler.CreateDeployment))
	mux.Handle("/deployment/update", http.HandlerFunc(clusterHandler.UpdateDeployment))
	mux.Handle("/deployment/get", http.HandlerFunc(clusterHandler.GetDeployment))
	mux.Handle("/statefulset/create", http.HandlerFunc(clusterHandler.CreateStatefulSet))
	mux.Handle("/statefulset/update", http.HandlerFunc(clusterHandler.UpdateStatefulSet))
	mux.Handle("/statefulset/get", http.HandlerFunc(clusterHandler.GetStatefulSet))

	Logger.Info("Routes registered")
}

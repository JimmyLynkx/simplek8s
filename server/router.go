package server

import (
	"go_code/simplek8s/core/application/handler"
	"go_code/simplek8s/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func InitializeRouter() (http.Handler, error) {
	r := mux.NewRouter()

	// 添加中间件
	r.Use(middleware.JSONResponseMiddleware)
	r.Use(func(next http.Handler) http.Handler {
		return middleware.LoggingMiddleware(next, Logger)
	})
	r.Use(func(next http.Handler) http.Handler {
		return middleware.RecoverMiddleware(next, Logger)
	})

	// 添加路由
	r.HandleFunc("/deployments", handler.CreateDeployment).Methods("POST")
	r.HandleFunc("/deployments/delete", handler.DeleteDeployment).Methods("POST") // 从请求体接收参数
	r.HandleFunc("/statefulsets", handler.CreateStatefulSet).Methods("POST")
	r.HandleFunc("/statefulsets/delete", handler.DeleteStatefulSet).Methods("POST")
	r.HandleFunc("/pods", handler.GetPod).Methods("POST")

	return r, nil
}

package server

import (
	"go_code/simplek8s/handler"
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
	r.HandleFunc("/deployments/{name}", handler.DeleteDeployment).Methods("POST")
	r.HandleFunc("/statefulsets", handler.CreateStatefulSet).Methods("POST")
	r.HandleFunc("/statefulsets/{name}", handler.DeleteStatefulSet).Methods("POST")
	r.HandleFunc("/pod/{namespace}/{name}", handler.GetPod).Methods("GET") // 查询指定 Pod
	r.HandleFunc("/pod/{namespace}", handler.GetPod).Methods("GET")        // 查询指定命名空间下的所有 Pod

	return r, nil
}

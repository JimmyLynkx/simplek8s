package server

import (
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
	r.HandleFunc("/deployments", CreateDeployment).Methods("POST")
	r.HandleFunc("/deployments/{name}", DeleteDeployment).Methods("DELETE")
	r.HandleFunc("/statefulsets", CreateStatefulSet).Methods("POST")
	r.HandleFunc("/statefulsets/{name}", DeleteStatefulSet).Methods("DELETE")

	return r, nil
}

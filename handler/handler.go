package handler

import (
	"context"
	"encoding/json"
	"go_code/simplek8s/internal/utils"
	"net/http"

	"github.com/gorilla/mux"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func CreateDeployment(w http.ResponseWriter, r *http.Request) {
	var deployment v1.Deployment
	if err := json.NewDecoder(r.Body).Decode(&deployment); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// 获取集群配置
	config, err := rest.InClusterConfig()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// 使用获得的配置创建一个Kubernetes客户端clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// 创建deployment
	_, err = clientset.AppsV1().Deployments(deployment.Namespace).Create(context.TODO(), &deployment, metav1.CreateOptions{})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, deployment)
}

func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	config, err := rest.InClusterConfig()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = clientset.AppsV1().Deployments("default").Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func CreateStatefulSet(w http.ResponseWriter, r *http.Request) {
	var statefulSet v1.StatefulSet
	if err := json.NewDecoder(r.Body).Decode(&statefulSet); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// 获取集群配置
	config, err := rest.InClusterConfig()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// 使用获得的配置创建一个Kubernetes客户端clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// 创建statefulset
	_, err = clientset.AppsV1().StatefulSets(statefulSet.Namespace).Create(context.TODO(), &statefulSet, metav1.CreateOptions{})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, statefulSet)
}

func DeleteStatefulSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	config, err := rest.InClusterConfig()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = clientset.AppsV1().StatefulSets("default").Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func GetPod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	name := vars["name"]

	// 获取集群配置
	config, err := rest.InClusterConfig()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 使用获得的配置创建一个 Kubernetes 客户端 clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if name != "" {
		// 查询特定 Pod
		pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, pod)
	} else {
		// 查询所有 Pod
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, pods.Items)
	}
}

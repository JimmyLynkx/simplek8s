package handler

import (
	"context"
	"encoding/json"
	"go_code/simplek8s/internal/utils"
	"net/http"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getClientset() (*kubernetes.Clientset, error) {
	// 获取集群配置
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// 使用获得的配置创建一个 Kubernetes 客户端 clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func CreateDeployment(w http.ResponseWriter, r *http.Request) {
	var deployment v1.Deployment
	if err := json.NewDecoder(r.Body).Decode(&deployment); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// 获取 Kubernetes 客户端 clientset
	clientset, err := getClientset()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 创建 deployment
	_, err = clientset.AppsV1().Deployments(deployment.Namespace).Create(context.TODO(), &deployment, metav1.CreateOptions{})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, deployment)
}

type DeleteRequest struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	var req DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// 获取 Kubernetes 客户端 clientset
	clientset, err := getClientset()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if req.Namespace == "" {
		req.Namespace = "default"
	}

	// 删除 deployment
	err = clientset.AppsV1().Deployments(req.Namespace).Delete(context.TODO(), req.Name, metav1.DeleteOptions{})
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

	// 获取 Kubernetes 客户端 clientset
	clientset, err := getClientset()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 创建 statefulset
	_, err = clientset.AppsV1().StatefulSets(statefulSet.Namespace).Create(context.TODO(), &statefulSet, metav1.CreateOptions{})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, statefulSet)
}

func DeleteStatefulSet(w http.ResponseWriter, r *http.Request) {
	var req DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// 获取 Kubernetes 客户端 clientset
	clientset, err := getClientset()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if req.Namespace == "" {
		req.Namespace = "default"
	}

	// 删除 statefulset
	err = clientset.AppsV1().StatefulSets(req.Namespace).Delete(context.TODO(), req.Name, metav1.DeleteOptions{})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

type GetPodRequest struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

func GetPod(w http.ResponseWriter, r *http.Request) {
	var req GetPodRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// 获取 Kubernetes 客户端 clientset
	clientset, err := getClientset()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if req.Name != "" {
		// 查询特定 Pod
		pod, err := clientset.CoreV1().Pods(req.Namespace).Get(context.TODO(), req.Name, metav1.GetOptions{})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, pod)
	} else {
		// 查询所有 Pod
		pods, err := clientset.CoreV1().Pods(req.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, pods.Items)
	}
}

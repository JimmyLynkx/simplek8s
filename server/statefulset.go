package server

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

func CreateStatefulSet(w http.ResponseWriter, r *http.Request) {
	var statefulSet v1.StatefulSet
	if err := json.NewDecoder(r.Body).Decode(&statefulSet); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

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

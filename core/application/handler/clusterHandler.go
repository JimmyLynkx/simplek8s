package handler

import (
	"encoding/json"
	"net/http"

	"go_code/simplek8s/core/application/service"
	"go_code/simplek8s/core/entity"
	"go_code/simplek8s/internal/utils"
)

type ClusterHandler struct {
	ClusterService service.ClusterService
}

func NewClusterHandler(clusterService service.ClusterService) *ClusterHandler {
	return &ClusterHandler{ClusterService: clusterService}
}

func (h *ClusterHandler) AddCluster(w http.ResponseWriter, r *http.Request) {
	var cluster entity.Cluster
	if err := json.NewDecoder(r.Body).Decode(&cluster); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	err := h.ClusterService.AddCluster(cluster)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, "Cluster added successfully")
}

type CreateDeploymentRequest struct {
	ClusterID      int    `json:"cluster_id"`
	DeploymentYAML string `json:"deploymentYAML"`
}

func (h *ClusterHandler) CreateDeployment(w http.ResponseWriter, r *http.Request) {
	var req CreateDeploymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	err := h.ClusterService.CreateDeployment(req.ClusterID, req.DeploymentYAML)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "Deployment created successfully")
}

type UpdateDeploymentRequest CreateDeploymentRequest

func (h *ClusterHandler) UpdateDeployment(w http.ResponseWriter, r *http.Request) {
	var req UpdateDeploymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	err := h.ClusterService.UpdateDeployment(req.ClusterID, req.DeploymentYAML)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "Deployment updated successfully")
}

type GetDeploymentRequest struct {
	ClusterID      int    `json:"cluster_id"`
	Namespace      string `json:"namespace"`
	DeploymentName string `json:"deploymentName"`
}

func (h *ClusterHandler) GetDeployment(w http.ResponseWriter, r *http.Request) {
	var req GetDeploymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.DeploymentName == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Deployment name is required")
		return
	}

	deployment, err := h.ClusterService.GetDeployment(req.ClusterID, req.Namespace, req.DeploymentName)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, deployment)
}

type DeleteDeploymentRequest GetDeploymentRequest

func (h *ClusterHandler) DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	var req DeleteDeploymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.DeploymentName == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Deployment name is required")
		return
	}

	err := h.ClusterService.DeleteDeployment(req.ClusterID, req.Namespace, req.DeploymentName)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Deployment delete successfully"})
}

type CreateStatefulSetRequest struct {
	ClusterID       int    `json:"cluster_id"`
	StatefulSetYAML string `json:"statefulSetYAML"`
}

// CreateStatefulSet 创建 StatefulSet 的处理函数
func (h *ClusterHandler) CreateStatefulSet(w http.ResponseWriter, r *http.Request) {
	var req CreateStatefulSetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.StatefulSetYAML == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "StatefulSet YAML is required")
		return
	}

	err := h.ClusterService.CreateStatefulSet(req.ClusterID, req.StatefulSetYAML)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "StatefulSet created successfully"})
}

type UpdateStatefulSetRequest CreateStatefulSetRequest

// UpdateStatefulSet 更新 StatefulSet 的处理函数
func (h *ClusterHandler) UpdateStatefulSet(w http.ResponseWriter, r *http.Request) {
	var req UpdateStatefulSetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.StatefulSetYAML == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "StatefulSet YAML is required")
		return
	}

	err := h.ClusterService.UpdateStatefulSet(req.ClusterID, req.StatefulSetYAML)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "StatefulSet updated successfully"})
}

type GetStatefulSetRequest struct {
	ClusterID       int    `json:"cluster_id"`
	Namespace       string `json:"namespace"`
	StatefulSetName string `json:"statefulSetName"`
}

// GetStatefulSet 获取 StatefulSet 的处理函数
func (h *ClusterHandler) GetStatefulSet(w http.ResponseWriter, r *http.Request) {
	var req GetStatefulSetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.StatefulSetName == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "StatefulSet name is required")
		return
	}

	statefulSet, err := h.ClusterService.GetStatefulSet(req.ClusterID, req.Namespace, req.StatefulSetName)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, statefulSet)
}

type DeleteStatefulSetRequest GetStatefulSetRequest

func (h *ClusterHandler) DeleteStatefulSet(w http.ResponseWriter, r *http.Request) {
	var req DeleteStatefulSetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.StatefulSetName == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "StatefulSet name is required")
		return
	}

	err := h.ClusterService.DeleteStatefulSet(req.ClusterID, req.Namespace, req.StatefulSetName)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "StatefulSet delete successfully"})
}

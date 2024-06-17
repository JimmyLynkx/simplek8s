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
	ClusterID  int    `json:"cluster_id"`
	Deployment string `json:"deployment"`
}

func (h *ClusterHandler) CreateDeployment(w http.ResponseWriter, r *http.Request) {
	var req CreateDeploymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	err := h.ClusterService.CreateDeployment(req.ClusterID, req.Deployment)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "Deployment created successfully")
}

type UpdateDeploymentRequest struct {
	ClusterID          int    `json:"cluster_id"`
	DeploymentYAMLPath string `json:"deploymentYAMLPath"`
}

func (h *ClusterHandler) UpdateDeployment(w http.ResponseWriter, r *http.Request) {
	var req UpdateDeploymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	err := h.ClusterService.UpdateDeployment(req.ClusterID, req.DeploymentYAMLPath)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "Deployment updated successfully")
}

// 请求体结构体
type GetDeploymentRequest struct {
	ClusterID      int    `json:"cluster_id"`
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

	deployment, err := h.ClusterService.GetDeployment(req.ClusterID, req.DeploymentName)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, deployment)
}

// 请求体结构体
type CreateStatefulSetRequest struct {
	ClusterID       int    `json:"cluster_id"`
	StatefulSetYAML string `json:"statefulSetYAML"`
}

type UpdateStatefulSetRequest struct {
	ClusterID       int    `json:"cluster_id"`
	StatefulSetYAML string `json:"statefulSetYAML"`
}

type GetStatefulSetRequest struct {
	ClusterID       int    `json:"cluster_id"`
	StatefulSetName string `json:"statefulSetName"`
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

	statefulSet, err := h.ClusterService.GetStatefulSet(req.ClusterID, req.StatefulSetName)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, statefulSet)
}

package service

import (
	"context"
	"fmt"
	"go_code/simplek8s/core/application/repository"
	"go_code/simplek8s/core/entity"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type ClusterService struct {
	ClusterRepo repository.ClusterRepo
}

func NewClusterService(clusterRepo repository.ClusterRepo) ClusterService {
	return ClusterService{ClusterRepo: clusterRepo}
}

// AddCluster 添加新的集群信息
func (s *ClusterService) AddCluster(cluster entity.Cluster) error {
	_, err := s.ClusterRepo.Create(cluster)
	return err
}

// CreateDeployment 在指定集群上创建 Deployment
func (s *ClusterService) CreateDeployment(clusterID int, deploymentYAML string) error {
	// 从存储库中获取集群信息
	cluster, err := s.ClusterRepo.GetByID(clusterID)
	if err != nil {
		return fmt.Errorf("failed to get cluster: %v", err)
	}

	// 从字符串创建 REST 配置
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(cluster.Config))
	if err != nil {
		return fmt.Errorf("failed to create rest config: %v", err)
	}

	// 创建动态客户端
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %v", err)
	}

	// 解析部署 YAML 字符串
	deployment := &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(deploymentYAML), deployment); err != nil {
		return fmt.Errorf("failed to unmarshal deployment YAML: %v", err)
	}

	// 获取 GVK 和 GVR
	gvk := deployment.GroupVersionKind()
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: "deployments",
	}

	// 获取命名空间
	namespace := deployment.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}

	// 创建 Deployment
	_, err = dynamicClient.Resource(gvr).Namespace(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create deployment: %v", err)
	}

	return nil
}

// UpdateDeployment 在指定集群上更新 Deployment
func (s *ClusterService) UpdateDeployment(clusterID int, deploymentYAML string) error {
	cluster, err := s.ClusterRepo.GetByID(clusterID)
	if err != nil {
		return fmt.Errorf("failed to get cluster: %v", err)
	}

	configBytes := []byte(cluster.Config)
	config, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		return fmt.Errorf("failed to create rest config: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %v", err)
	}

	// 将字符串形式的 YAML 解析为 Unstructured 对象
	deployment := &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(deploymentYAML), deployment); err != nil {
		return fmt.Errorf("failed to unmarshal deployment YAML: %v", err)
	}

	// 从 YAML 中解析出 namespace 和 deploymentName
	namespace := deployment.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}
	deploymentName := deployment.GetName()
	if deploymentName == "" {
		return fmt.Errorf("deployment name is required in the YAML")
	}

	gvk := deployment.GroupVersionKind()
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: "deployments",
	}

	// 获取现有的 Deployment
	existingDeployment, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get existing deployment: %v", err)
	}

	// 更新现有 Deployment 的 spec
	existingDeployment.Object["spec"] = deployment.Object["spec"]

	// 更新 Deployment
	_, err = dynamicClient.Resource(gvr).Namespace(namespace).Update(context.Background(), existingDeployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update deployment: %v", err)
	}

	return nil
}

// GetDeployment 获取指定集群的 Deployment
func (s *ClusterService) GetDeployment(clusterID int, namespace, deploymentName string) (*appsv1.Deployment, error) {
	// 从存储库中获取集群信息
	cluster, err := s.ClusterRepo.GetByID(clusterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster: %v", err)
	}

	// 创建 REST 配置
	configBytes := []byte(cluster.Config)
	config, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create rest config: %v", err)
	}

	// 创建 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	if namespace == "" {
		namespace = "default"
	}

	// 获取 Deployment
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %v", err)
	}

	return deployment, nil
}

// DeleteDeployment 删除指定集群的 Deployment
func (s *ClusterService) DeleteDeployment(clusterID int, namespace, deploymentName string) error {
	// 从存储库中获取集群信息
	cluster, err := s.ClusterRepo.GetByID(clusterID)
	if err != nil {
		return fmt.Errorf("failed to get cluster: %v", err)
	}

	// 创建 REST 配置
	configBytes := []byte(cluster.Config)
	config, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		return fmt.Errorf("failed to create rest config: %v", err)
	}

	// 创建 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	if namespace == "" {
		namespace = "default"
	}

	// 删除 Deployment
	err = clientset.AppsV1().Deployments(namespace).Delete(context.Background(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete deployment: %v", err)
	}

	return nil
}

// CreateStatefulSet 在指定集群上创建 StatefulSet
func (s *ClusterService) CreateStatefulSet(clusterID int, statefulSetYAML string) error {
	cluster, err := s.ClusterRepo.GetByID(clusterID)
	if err != nil {
		return fmt.Errorf("failed to get cluster: %v", err)
	}

	configBytes := []byte(cluster.Config)

	config, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		return fmt.Errorf("failed to create rest config: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %v", err)
	}

	statefulSet := &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(statefulSetYAML), statefulSet); err != nil {
		return fmt.Errorf("failed to unmarshal statefulSet YAML: %v", err)
	}

	gvk := statefulSet.GroupVersionKind()
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: "statefulsets",
	}

	namespace := statefulSet.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}

	_, err = dynamicClient.Resource(gvr).Namespace(namespace).Create(context.Background(), statefulSet, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create statefulSet: %v", err)
	}

	return nil
}

// UpdateStatefulSet 在指定集群上更新 StatefulSet
func (s *ClusterService) UpdateStatefulSet(clusterID int, statefulSetYAML string) error {
	cluster, err := s.ClusterRepo.GetByID(clusterID)
	if err != nil {
		return fmt.Errorf("failed to get cluster: %v", err)
	}

	configBytes := []byte(cluster.Config)

	config, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		return fmt.Errorf("failed to create rest config: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %v", err)
	}

	statefulSet := &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(statefulSetYAML), statefulSet); err != nil {
		return fmt.Errorf("failed to unmarshal statefulSet YAML: %v", err)
	}

	// 从 YAML 中解析出 namespace 和 statefulSetName
	namespace := statefulSet.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}
	statefulSetName := statefulSet.GetName()
	if statefulSetName == "" {
		return fmt.Errorf("statefulSet name is required in the YAML")
	}

	gvk := statefulSet.GroupVersionKind()
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: "statefulsets",
	}

	existingStatefulSet, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(context.Background(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get existing statefulSet: %v", err)
	}

	// 更新现有 StatefulSet 的 spec
	existingStatefulSet.Object["spec"] = statefulSet.Object["spec"]

	_, err = dynamicClient.Resource(gvr).Namespace(namespace).Update(context.Background(), existingStatefulSet, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update statefulSet: %v", err)
	}

	return nil
}

// GetStatefulSet 获取指定集群的 StatefulSet
func (s *ClusterService) GetStatefulSet(clusterID int, namespace, statefulSetName string) (*appsv1.StatefulSet, error) {
	// 从存储库中获取集群信息
	cluster, err := s.ClusterRepo.GetByID(clusterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster: %v", err)
	}

	configBytes := []byte(cluster.Config)

	// 从配置文件创建 REST 配置
	config, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create rest config: %v", err)
	}

	// 创建 Kubernetes 客户端集
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %v", err)
	}

	if namespace == "" {
		namespace = "default"
	}

	// 获取 StatefulSet
	statefulSet, err := clientset.AppsV1().StatefulSets(namespace).Get(context.Background(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get statefulSet: %v", err)
	}

	return statefulSet, nil
}

// DeleteStatefulSet 删除指定集群的 StatefulSet
func (s *ClusterService) DeleteStatefulSet(clusterID int, namespace, statefulSetName string) error {
	// 从存储库中获取集群信息
	cluster, err := s.ClusterRepo.GetByID(clusterID)
	if err != nil {
		return fmt.Errorf("failed to get cluster: %v", err)
	}

	// 创建 REST 配置
	configBytes := []byte(cluster.Config)
	config, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		return fmt.Errorf("failed to create rest config: %v", err)
	}

	// 创建 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %v", err)
	}

	if namespace == "" {
		namespace = "default"
	}

	// 删除 StatefulSet
	err = clientset.AppsV1().StatefulSets(namespace).Delete(context.Background(), statefulSetName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete statefulSet: %v", err)
	}

	return nil
}

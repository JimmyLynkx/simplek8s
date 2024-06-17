package service

import (
	"context"
	"fmt"
	"go_code/simplek8s/core/application/repository"
	"go_code/simplek8s/core/entity"
	"os"

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
func (s *ClusterService) CreateDeployment(clusterID int, deploymentYAMLPath string) error {
	cluster, err := s.ClusterRepo.GetByID(clusterID)
	if err != nil {
		return fmt.Errorf("failed to get cluster: %v", err)
	}

	configBytes, err := os.ReadFile(cluster.Config)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		return fmt.Errorf("failed to create rest config: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %v", err)
	}

	// 读取并解析 YAML 文件
	yamlFile, err := os.ReadFile(deploymentYAMLPath)
	if err != nil {
		return fmt.Errorf("failed to read deployment YAML file: %v", err)
	}

	deployment := &unstructured.Unstructured{}
	if err := yaml.Unmarshal(yamlFile, deployment); err != nil {
		return fmt.Errorf("failed to unmarshal deployment YAML: %v", err)
	}

	gvk := deployment.GroupVersionKind()
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: "deployments",
	}

	namespace := deployment.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}

	_, err = dynamicClient.Resource(gvr).Namespace(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create deployment: %v", err)
	}

	return nil
}

// UpdateDeployment 在指定集群上更新 Deployment
func (s *ClusterService) UpdateDeployment(clusterID int, deploymentYAMLPath string) error {
	cluster, err := s.ClusterRepo.GetByID(clusterID)
	if err != nil {
		return fmt.Errorf("failed to get cluster: %v", err)
	}

	configBytes, err := os.ReadFile(cluster.Config)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		return fmt.Errorf("failed to create rest config: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %v", err)
	}

	// 读取并解析 YAML 文件
	yamlFile, err := os.ReadFile(deploymentYAMLPath)
	if err != nil {
		return fmt.Errorf("failed to read deployment YAML file: %v", err)
	}

	deployment := &unstructured.Unstructured{}
	if err := yaml.Unmarshal(yamlFile, deployment); err != nil {
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
func (s *ClusterService) GetDeployment(clusterID int, deploymentName string) (*appsv1.Deployment, error) {
	// 从存储库中获取集群信息
	cluster, err := s.ClusterRepo.GetByID(clusterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster: %v", err)
	}

	// 读取集群配置文件
	configBytes, err := os.ReadFile(cluster.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// 从 kubeconfig 文件创建 REST 配置
	config, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create REST config from kubeconfig: %v", err)
	}

	// 创建 Kubernetes 客户端集
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	// 从 kubeconfig 文件中获取命名空间
	apiConfig, err := clientcmd.Load(configBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to load API config: %v", err)
	}
	namespace := apiConfig.Contexts[apiConfig.CurrentContext].Namespace
	if namespace == "" {
		namespace = "default" // 使用默认命名空间
	}

	// 获取指定命名空间中的 Deployment
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %v", err)
	}

	return deployment, nil
}

// CreateStatefulSet 在指定集群上创建 StatefulSet
func (s *ClusterService) CreateStatefulSet(clusterID int, statefulSetYAMLPath string) error {
	cluster, err := s.ClusterRepo.GetByID(clusterID)
	if err != nil {
		return fmt.Errorf("failed to get cluster: %v", err)
	}

	configBytes, err := os.ReadFile(cluster.Config)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		return fmt.Errorf("failed to create rest config: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %v", err)
	}

	// 读取并解析 YAML 文件
	yamlFile, err := os.ReadFile(statefulSetYAMLPath)
	if err != nil {
		return fmt.Errorf("failed to read statefulSet YAML file: %v", err)
	}

	statefulSet := &unstructured.Unstructured{}
	if err := yaml.Unmarshal(yamlFile, statefulSet); err != nil {
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
func (s *ClusterService) UpdateStatefulSet(clusterID int, statefulSetYAMLPath string) error {
	cluster, err := s.ClusterRepo.GetByID(clusterID)
	if err != nil {
		return fmt.Errorf("failed to get cluster: %v", err)
	}

	configBytes, err := os.ReadFile(cluster.Config)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig(configBytes)
	if err != nil {
		return fmt.Errorf("failed to create rest config: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %v", err)
	}

	// 读取并解析 YAML 文件
	yamlFile, err := os.ReadFile(statefulSetYAMLPath)
	if err != nil {
		return fmt.Errorf("failed to read statefulSet YAML file: %v", err)
	}

	statefulSet := &unstructured.Unstructured{}
	if err := yaml.Unmarshal(yamlFile, statefulSet); err != nil {
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
func (s *ClusterService) GetStatefulSet(clusterID int, statefulSetName string) (*appsv1.StatefulSet, error) {
	// 从存储库中获取集群信息
	cluster, err := s.ClusterRepo.GetByID(clusterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster: %v", err)
	}

	// 从配置文件读取 Kubernetes 配置
	configBytes, err := os.ReadFile(cluster.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

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

	// 从 kubeconfig 中提取默认命名空间
	apiConfig, err := clientcmd.Load(configBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to load api config: %v", err)
	}
	namespace := "default"
	if ctx, exists := apiConfig.Contexts[apiConfig.CurrentContext]; exists {
		if ctx.Namespace != "" {
			namespace = ctx.Namespace
		}
	}

	// 获取 StatefulSet
	statefulSet, err := clientset.AppsV1().StatefulSets(namespace).Get(context.Background(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get statefulSet: %v", err)
	}

	return statefulSet, nil
}

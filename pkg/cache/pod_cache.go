package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// PodInfo 存储 Pod 的关键信息
type PodInfo struct {
	UID         string            `json:"uid"`
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	NodeName    string            `json:"nodeName"`
	Phase       string            `json:"phase"`
	PodIP       string            `json:"podIP"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
}

// PodCache 管理 Pod 信息的缓存
type PodCache struct {
	clientset *kubernetes.Clientset
	informer  cache.SharedIndexInformer
	mu        sync.RWMutex
	pods      map[types.UID]*PodInfo
	synced    bool
}

// NewPodCache 创建新的 Pod 缓存实例
func NewPodCache() (*PodCache, error) {
	config, err := getKubernetesConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)
	informer := factory.Core().V1().Pods().Informer()

	pc := &PodCache{
		clientset: clientset,
		informer:  informer,
		pods:      make(map[types.UID]*PodInfo),
		synced:    false,
	}

	// 注册事件处理器
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    pc.onAdd,
		UpdateFunc: pc.onUpdate,
		DeleteFunc: pc.onDelete,
	})

	return pc, nil
}

// getKubernetesConfig 获取 Kubernetes 配置
func getKubernetesConfig() (*rest.Config, error) {
	// 优先使用集群内配置
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// 回退到 kubeconfig
	kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Start 启动 informer
func (pc *PodCache) Start(ctx context.Context) error {
	go pc.informer.Run(ctx.Done())
	return nil
}

// WaitForCacheSync 等待缓存同步
func (pc *PodCache) WaitForCacheSync(ctx context.Context) bool {
	synced := cache.WaitForCacheSync(ctx.Done(), pc.informer.HasSynced)
	if synced {
		pc.mu.Lock()
		pc.synced = true
		pc.mu.Unlock()
	}
	return synced
}

// IsSynced 检查缓存是否已同步
func (pc *PodCache) IsSynced() bool {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return pc.synced
}

// GetPodByUID 根据 UID 获取 Pod 信息
func (pc *PodCache) GetPodByUID(uid types.UID) (*PodInfo, error) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	pod, exists := pc.pods[uid]
	if !exists {
		return nil, fmt.Errorf("pod with UID %s not found", uid)
	}

	return pod, nil
}

// GetPodCount 获取缓存中的 Pod 数量
func (pc *PodCache) GetPodCount() int {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return len(pc.pods)
}

// onAdd 处理 Pod 添加事件
func (pc *PodCache) onAdd(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return
	}

	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.pods[pod.UID] = convertPodToPodInfo(pod)
}

// onUpdate 处理 Pod 更新事件
func (pc *PodCache) onUpdate(oldObj, newObj interface{}) {
	pod, ok := newObj.(*corev1.Pod)
	if !ok {
		return
	}

	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.pods[pod.UID] = convertPodToPodInfo(pod)
}

// onDelete 处理 Pod 删除事件
func (pc *PodCache) onDelete(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		// 处理 DeletedFinalStateUnknown 情况
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			return
		}
		pod, ok = tombstone.Obj.(*corev1.Pod)
		if !ok {
			return
		}
	}

	pc.mu.Lock()
	defer pc.mu.Unlock()
	delete(pc.pods, pod.UID)
}

// convertPodToPodInfo 将 Pod 对象转换为 PodInfo
func convertPodToPodInfo(pod *corev1.Pod) *PodInfo {
	return &PodInfo{
		UID:         string(pod.UID),
		Name:        pod.Name,
		Namespace:   pod.Namespace,
		NodeName:    pod.Spec.NodeName,
		Phase:       string(pod.Status.Phase),
		PodIP:       pod.Status.PodIP,
		Labels:      pod.Labels,
		Annotations: pod.Annotations,
		CreatedAt:   pod.CreationTimestamp.Time,
	}
}

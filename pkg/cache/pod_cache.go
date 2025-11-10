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

// PodInfo stores pod information
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

// PodCache manages pod information cache
type PodCache struct {
	clientset *kubernetes.Clientset
	informer  cache.SharedIndexInformer
	mu        sync.RWMutex
	pods      map[types.UID]*PodInfo
	synced    bool
}

// NewPodCache creates a new pod cache instance
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

	// Register event handlers
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    pc.onAdd,
		UpdateFunc: pc.onUpdate,
		DeleteFunc: pc.onDelete,
	})

	return pc, nil
}

// getKubernetesConfig retrieves Kubernetes configuration
func getKubernetesConfig() (*rest.Config, error) {
	// Try in-cluster config first
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}

	// Fallback to kubeconfig
	kubeconfig := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Start starts the informer
func (pc *PodCache) Start(ctx context.Context) error {
	go pc.informer.Run(ctx.Done())
	return nil
}

// WaitForCacheSync waits for cache synchronization
func (pc *PodCache) WaitForCacheSync(ctx context.Context) bool {
	synced := cache.WaitForCacheSync(ctx.Done(), pc.informer.HasSynced)
	if synced {
		pc.mu.Lock()
		pc.synced = true
		pc.mu.Unlock()
	}
	return synced
}

// IsSynced checks if cache is synchronized
func (pc *PodCache) IsSynced() bool {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return pc.synced
}

// GetPodByUID retrieves pod information by UID
func (pc *PodCache) GetPodByUID(uid types.UID) (*PodInfo, error) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	pod, exists := pc.pods[uid]
	if !exists {
		return nil, fmt.Errorf("pod with UID %s not found", uid)
	}

	return pod, nil
}

// GetPodCount returns the number of pods in cache
func (pc *PodCache) GetPodCount() int {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return len(pc.pods)
}

// onAdd handles pod add events
func (pc *PodCache) onAdd(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return
	}

	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.pods[pod.UID] = convertPodToPodInfo(pod)
}

// onUpdate handles pod update events
func (pc *PodCache) onUpdate(oldObj, newObj interface{}) {
	pod, ok := newObj.(*corev1.Pod)
	if !ok {
		return
	}

	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.pods[pod.UID] = convertPodToPodInfo(pod)
}

// onDelete handles pod delete events
func (pc *PodCache) onDelete(obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		// Handle DeletedFinalStateUnknown case
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

// convertPodToPodInfo converts Pod object to PodInfo
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

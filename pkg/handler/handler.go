package handler

import (
	"encoding/json"
	"net/http"

	"k8s.io/apimachinery/pkg/types"

	"github.com/weipengyu/pod-index/pkg/cache"
)

// Handler HTTP 请求处理器
type Handler struct {
	podCache *cache.PodCache
}

// NewHandler 创建新的处理器实例
func NewHandler(podCache *cache.PodCache) *Handler {
	return &Handler{
		podCache: podCache,
	}
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Error string `json:"error"`
}

// GetPodByUID 根据 UID 获取 Pod 信息
// Query parameter: uid
func (h *Handler) GetPodByUID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uid := r.URL.Query().Get("uid")
	if uid == "" {
		h.respondError(w, "uid parameter is required", http.StatusBadRequest)
		return
	}

	podInfo, err := h.podCache.GetPodByUID(types.UID(uid))
	if err != nil {
		h.respondError(w, err.Error(), http.StatusNotFound)
		return
	}

	h.respondJSON(w, podInfo, http.StatusOK)
}

// Health 健康检查端点
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.respondJSON(w, map[string]string{"status": "healthy"}, http.StatusOK)
}

// Ready 就绪检查端点
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !h.podCache.IsSynced() {
		h.respondError(w, "cache not synced", http.StatusServiceUnavailable)
		return
	}

	h.respondJSON(w, map[string]interface{}{
		"status":   "ready",
		"podCount": h.podCache.GetPodCount(),
	}, http.StatusOK)
}

// respondJSON 发送 JSON 响应
func (h *Handler) respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// respondError 发送错误响应
func (h *Handler) respondError(w http.ResponseWriter, message string, statusCode int) {
	h.respondJSON(w, ErrorResponse{Error: message}, statusCode)
}

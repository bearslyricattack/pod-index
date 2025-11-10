package handler

import (
	"encoding/json"
	"net/http"

	"k8s.io/apimachinery/pkg/types"

	"github.com/weipengyu/pod-index/pkg/cache"
)

// Handler handles HTTP requests
type Handler struct {
	podCache *cache.PodCache
}

// NewHandler creates a new handler instance
func NewHandler(podCache *cache.PodCache) *Handler {
	return &Handler{
		podCache: podCache,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// GetPodByUID retrieves pod information by UID
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

// Health returns health status
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.respondJSON(w, map[string]string{"status": "healthy"}, http.StatusOK)
}

// Ready returns readiness status
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

// respondJSON sends a JSON response
func (h *Handler) respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// respondError sends an error response
func (h *Handler) respondError(w http.ResponseWriter, message string, statusCode int) {
	h.respondJSON(w, ErrorResponse{Error: message}, statusCode)
}

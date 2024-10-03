package http

import (
	"encoding/json"
	"github.com/pedrokunz/distributed_cache_go/internal/cache_manager"
	"net/http"
)

type Handler struct {
	cacheManager *cache_manager.CacheManager
}

func New(cacheManager *cache_manager.CacheManager) *Handler {
	return &Handler{
		cacheManager: cacheManager,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodDelete:
		h.handleDelete(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key not provided", http.StatusBadRequest)
		return
	}

	node, err := h.cacheManager.GetNodeForKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	value, _ := node.Get(key)

	w.WriteHeader(http.StatusOK)

	response := map[string]string{"key": key, "value": value, "node": node.ID()}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	if key == "" || value == "" {
		http.Error(w, "key or value not provided", http.StatusBadRequest)
		return
	}

	err := h.cacheManager.SetNodeForKey(key, value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	response := map[string]string{"key": key, "value": value}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key not provided", http.StatusBadRequest)
		return
	}

	err := h.cacheManager.DeleteNodeForKey(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	response := map[string]string{"key": key}
	_ = json.NewEncoder(w).Encode(response)
}

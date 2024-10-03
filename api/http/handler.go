package http

import (
	"encoding/json"
	"github.com/pedrokunz/distributed_cache_go/internal/cache_manager"
	"net/http"
)

// Handler is an HTTP handler for the cache manager
type Handler struct {
	cacheManager *cache_manager.CacheManager
}

// New creates a new Handler instance
func New(cacheManager *cache_manager.CacheManager) *Handler {
	return &Handler{
		cacheManager: cacheManager,
	}
}

// ServeHTTP implements the http.Handler interface
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

// handleGet handles GET requests to retrieve a key-value pair from the cache
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

// handlePost handles POST requests to set a key-value pair in the cache
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

// handleDelete handles DELETE requests to delete a key-value pair from the cache
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

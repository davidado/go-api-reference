// Package product : Product service
package product

import (
	"net/http"

	"github.com/davidado/go-api-reference/netjson"
	"github.com/davidado/go-api-reference/types"
	"github.com/gorilla/mux"
)

// Handler : Product handler
type Handler struct {
	store types.ProductStore
}

// NewHandler creates a new product handler
func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

// RegisterRoutes registers product routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
}

// handleGetProducts gets products
func (h *Handler) handleGetProducts(w http.ResponseWriter, _ *http.Request) {
	ps, err := h.store.GetProducts()
	if err != nil {
		netjson.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	netjson.Write(w, http.StatusOK, ps)
}

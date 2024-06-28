// Package cart : Cart service
package cart

import (
	"fmt"
	"net/http"

	"github.com/davidado/go-api-reference/netjson"
	"github.com/davidado/go-api-reference/service/auth"
	"github.com/davidado/go-api-reference/types"
	vd "github.com/davidado/go-api-reference/validator"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// Handler : Cart handler
type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

// NewHandler creates a new cart handler
func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:        store,
		productStore: productStore,
		userStore:    userStore,
	}
}

// RegisterRoutes registers cart routes
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

// handleCheckout handles the checkout of a cart
func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	var cart types.CartCheckoutPayload
	if err := netjson.Parse(r, &cart); err != nil {
		netjson.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := vd.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		netjson.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// get products
	productIDs, err := getCartItemsIDs(cart.Items)
	if err != nil {
		netjson.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ps, err := h.productStore.GetProductsByID(productIDs)
	if err != nil {
		netjson.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	orderID, totalPrice, err := h.createOrder(ps, cart.Items, userID)
	if err != nil {
		netjson.WriteError(w, http.StatusBadRequest, err)
		return
	}

	netjson.Write(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id":    orderID,
	})
}

package cart

import (
	"fmt"
	"net/http"
	"scratch/service/auth"
	"scratch/types"
	"scratch/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	orderStore   types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(orderStore types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{orderStore, productStore, userStore}
}

// Register my routes with a router or subrouter
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.HandleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.USER_ID).(int)
	var cart types.CartCheckoutPayload
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		err = err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %w", err))
		return
	}
	productIds, err := getCartItemsIDs(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %w", err))
		return
	}
	ps, err := h.productStore.GetProductsByIDs(productIds)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	orderID, totalPrice, err := h.CreateOrder(ps, cart.Items, userID)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON[map[string]any](w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id":    orderID,
	})
}

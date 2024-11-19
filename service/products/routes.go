package products

import (
	"net/http"
	"scratch/types"
	"scratch/utils"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store}
}

// Register my routes with a router or subrouter
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.HandleProducts).Methods(http.MethodGet)
	// router.HandleFunc("/products", h.HandleCreateProduct).Methods(http.MethodPost)
}

func (h *Handler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	// fetch products
	ps, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// return products
	utils.WriteJSON[[]types.Product](w, http.StatusOK, ps)
}

func (h *Handler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	// TODO: create
	// parse payload
	// validate payload
	// store product
	// return created
}

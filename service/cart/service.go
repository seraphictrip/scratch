package cart

import (
	"errors"
	"fmt"
	"scratch/types"
)

var (
	ErrInvalidQuantity = errors.New("invalid quantity")
	ErrOutOfStock      = errors.New("out of stock")
	ErrEmptyCart       = errors.New("empty cart")
)

func getCartItemsIDs(items []types.CartItem) ([]int, error) {
	productIds := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("%w: for product %d", ErrInvalidQuantity, item.ProductID)
		}
		productIds[i] = item.ProductID
	}
	return productIds, nil
}

func (h *Handler) CreateOrder(products []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}
	// TODO: this should be a transaction
	// check if products are in stock
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}

	// calculate total price
	totalPrice, err := calculateTotalPrice(items, productMap)
	if err != nil {
		return 0, 0, err
	}
	// reduce quantity of product
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity
		h.productStore.UpdateProduct(product)
	}

	// create the order
	orderID, err := h.orderStore.CreateOrder(types.Order{
		UserID: userID,
		Total:  totalPrice,
		Status: "pending",
		// TODO: provide or get from user
		Address: "some address",
	})
	if err != nil {
		return 0, 0, err
	}
	// create order item
	for _, item := range items {
		h.orderStore.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil
}

func checkIfCartIsInStock(items []types.CartItem, productMap map[int]types.Product) error {
	if len(items) == 0 {
		return ErrEmptyCart
	}

	for _, item := range items {
		product, ok := productMap[item.ProductID]
		if !ok {
			return fmt.Errorf("%w: product %v", ErrOutOfStock, item.ProductID)
		}
		if product.Quantity < item.Quantity {
			return fmt.Errorf("%w: low stock on %v", ErrInvalidQuantity, product.Name)
		}
	}
	return nil
}

func calculateTotalPrice(items []types.CartItem, productMap map[int]types.Product) (float64, error) {
	var totalPrice float64
	for _, item := range items {
		product, ok := productMap[item.ProductID]
		if !ok {
			return 0, fmt.Errorf("%w: product %v", ErrOutOfStock, item.ProductID)
		}
		totalPrice += float64(item.Quantity) * product.Price
	}
	return totalPrice, nil
}

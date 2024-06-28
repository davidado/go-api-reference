package cart

import (
	"fmt"

	"github.com/davidado/go-api-reference/types"
)

func getCartItemsIDs(items []types.CartItem) ([]int, error) {
	productIDs := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %d", item.ProductID)
		}
		productIDs[i] = item.ProductID
	}

	return productIDs, nil
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}

	// Check if all products are actually in stock.
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}

	// Calculate the total price.
	totalPrice := calculateTotalPrice(items, productMap)

	// Reduce the quantity of the products in the db.
	// Warning: This is a naive implementation if there are multiple requests.
	// Use a separate join table like Orders_Items to store order quantities.
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		h.productStore.UpdateProduct(product)
	}

	// Create the order.
	orderID, err := h.store.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "123 Main St", // Create an address table
	})
	if err != nil {
		return 0, 0, err
	}

	// Create the order items.
	for _, item := range items {
		h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}

	// TODO: Wrap all the above statements in a transaction.

	return orderID, totalPrice, nil
}

func checkIfCartIsInStock(cartItems []types.CartItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %d is out of stock", item.ProductID)
		}
	}

	return nil
}

func calculateTotalPrice(cartItems []types.CartItem, products map[int]types.Product) float64 {
	totalPrice := 0.0
	for _, item := range cartItems {
		product := products[item.ProductID]
		totalPrice += product.Price * float64(item.Quantity)
	}

	return totalPrice
}

// Package types contains the types used in the application.
package types

import "time"

// UserStore : User store interface
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(u User) error
}

// ProductStore : Product store interface
type ProductStore interface {
	GetProducts() ([]Product, error)
	GetProductsByID(ids []int) ([]Product, error)
	UpdateProduct(Product) error
}

// OrderStore : Order store interface
type OrderStore interface {
	CreateOrder(o Order) (int, error)
	CreateOrderItem(oi OrderItem) error
}

// Order : Order type
type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

// OrderItem : Order item type
type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"orderId"`
	ProductID int       `json:"productId"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

// Product : Product type
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"` // Good enough for now but in a real-world scenario, this should be atomic.
	CreatedAt   time.Time `json:"createdAt"`
}

// User : User type
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

// RegisterUserPayload : Register user payload
type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
}

// LoginUserPayload : Login user payload
type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// CartItem : Cart item type
type CartItem struct {
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

// CartCheckoutPayload : Cart checkout payload
type CartCheckoutPayload struct {
	Items []CartItem `json:"items" validate:"required"`
}

package models

import (
	"time"
)

// User represents a user in the system (PostgreSQL)
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Product represents a product (MongoDB)
type Product struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Price       float64   `json:"price" bson:"price"`
	Category    string    `json:"category" bson:"category"`
	Brand       string    `json:"brand" bson:"brand"`
	ImageURL    string    `json:"image_url" bson:"image_url"`
	Rating      float64   `json:"rating" bson:"rating"`
	Tags        []string  `json:"tags" bson:"tags"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// Order represents an order (PostgreSQL)
type Order struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	TotalAmount     float64   `json:"total_amount"`
	Status          string    `json:"status"`
	PaymentMethod   string    `json:"payment_method"`
	ShippingAddress string    `json:"shipping_address"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Items           []OrderItem `json:"items,omitempty"`
}

// OrderItem represents an item in an order (PostgreSQL)
type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// Inventory represents inventory data (MySQL)
type Inventory struct {
	ID                 int       `json:"id"`
	ProductID          string    `json:"product_id"`
	Quantity           int       `json:"quantity"`
	WarehouseLocation  string    `json:"warehouse_location"`
	LastRestocked      time.Time `json:"last_restocked"`
	LowStockThreshold  int       `json:"low_stock_threshold"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// Review represents a product review (MongoDB)
type Review struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	ProductID string    `json:"product_id" bson:"product_id"`
	UserID    int       `json:"user_id" bson:"user_id"`
	Rating    int       `json:"rating" bson:"rating"`
	Comment   string    `json:"comment" bson:"comment"`
	Helpful   int       `json:"helpful" bson:"helpful"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// Category represents a product category (MongoDB)
type Category struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	ParentID    string    `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	ImageURL    string    `json:"image_url" bson:"image_url"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// CartItem represents an item in shopping cart (PostgreSQL)
type CartItem struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProductID string    `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SalesAnalytics represents sales data (MySQL)
type SalesAnalytics struct {
	ID           int       `json:"id"`
	ProductID    string    `json:"product_id"`
	QuantitySold int       `json:"quantity_sold"`
	Revenue      float64   `json:"revenue"`
	SaleDate     time.Time `json:"sale_date"`
	CreatedAt    time.Time `json:"created_at"`
}

// Wishlist represents user wishlist items (MongoDB)
type Wishlist struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    int       `json:"user_id" bson:"user_id"`
	ProductID string    `json:"product_id" bson:"product_id"`
	AddedAt   time.Time `json:"added_at" bson:"added_at"`
}

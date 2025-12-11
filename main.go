package main

import (
	"log"
	"net/http"
	"os"

	"sample-application/config"
	"sample-application/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database connections
	config.InitDatabases()
	defer config.CloseDatabases()

	// Create router
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// User routes (PostgreSQL)
	router.HandleFunc("/api/users", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/api/users", handlers.GetAllUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", handlers.GetUserByID).Methods("GET")
	router.HandleFunc("/api/users/{id}", handlers.UpdateUser).Methods("PUT")
	router.HandleFunc("/api/users/{id}", handlers.DeleteUser).Methods("DELETE")
	router.HandleFunc("/api/users/{id}/orders", handlers.GetUserOrders).Methods("GET")

	// Product routes (MongoDB)
	// Note: Specific routes must come before parameterized routes
	router.HandleFunc("/api/products/search", handlers.SearchProducts).Methods("GET")
	router.HandleFunc("/api/products/category/{category}", handlers.GetProductsByCategory).Methods("GET")
	router.HandleFunc("/api/products", handlers.CreateProduct).Methods("POST")
	router.HandleFunc("/api/products", handlers.GetAllProducts).Methods("GET")
	router.HandleFunc("/api/products/{id}", handlers.GetProductByID).Methods("GET")
	router.HandleFunc("/api/products/{id}", handlers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/api/products/{id}", handlers.DeleteProduct).Methods("DELETE")

	// Order routes (PostgreSQL)
	router.HandleFunc("/api/orders", handlers.CreateOrder).Methods("POST")
	router.HandleFunc("/api/orders", handlers.GetAllOrders).Methods("GET")
	router.HandleFunc("/api/orders/{id}", handlers.GetOrderByID).Methods("GET")
	router.HandleFunc("/api/orders/{id}/status", handlers.UpdateOrderStatus).Methods("PATCH")
	router.HandleFunc("/api/orders/{id}/cancel", handlers.CancelOrder).Methods("POST")

	// Inventory routes (MySQL)
	router.HandleFunc("/api/inventory", handlers.GetAllInventory).Methods("GET")
	router.HandleFunc("/api/inventory/{product_id}", handlers.GetInventoryByProduct).Methods("GET")
	router.HandleFunc("/api/inventory/{product_id}", handlers.UpdateInventory).Methods("PUT")
	router.HandleFunc("/api/inventory/{product_id}/restock", handlers.RestockInventory).Methods("POST")
	router.HandleFunc("/api/inventory/low-stock", handlers.GetLowStockItems).Methods("GET")

	// Review routes (MongoDB)
	router.HandleFunc("/api/reviews", handlers.CreateReview).Methods("POST")
	router.HandleFunc("/api/reviews/product/{product_id}", handlers.GetProductReviews).Methods("GET")
	router.HandleFunc("/api/reviews/{id}", handlers.DeleteReview).Methods("DELETE")
	router.HandleFunc("/api/reviews/{id}/helpful", handlers.MarkReviewHelpful).Methods("POST")

	// Category routes (MongoDB)
	router.HandleFunc("/api/categories", handlers.CreateCategory).Methods("POST")
	router.HandleFunc("/api/categories", handlers.GetAllCategories).Methods("GET")
	router.HandleFunc("/api/categories/{id}", handlers.GetCategoryByID).Methods("GET")
	router.HandleFunc("/api/categories/{id}", handlers.UpdateCategory).Methods("PUT")
	router.HandleFunc("/api/categories/{id}", handlers.DeleteCategory).Methods("DELETE")

	// Cart routes (PostgreSQL)
	router.HandleFunc("/api/cart/{user_id}", handlers.GetCart).Methods("GET")
	router.HandleFunc("/api/cart/{user_id}/items", handlers.AddToCart).Methods("POST")
	router.HandleFunc("/api/cart/{user_id}/items/{item_id}", handlers.RemoveFromCart).Methods("DELETE")
	router.HandleFunc("/api/cart/{user_id}/clear", handlers.ClearCart).Methods("DELETE")

	// Analytics routes (MySQL)
	router.HandleFunc("/api/analytics/sales", handlers.GetSalesAnalytics).Methods("GET")
	router.HandleFunc("/api/analytics/popular-products", handlers.GetPopularProducts).Methods("GET")
	router.HandleFunc("/api/analytics/revenue", handlers.GetRevenueStats).Methods("GET")

	// Wishlist routes (MongoDB)
	router.HandleFunc("/api/wishlist/{user_id}", handlers.GetWishlist).Methods("GET")
	router.HandleFunc("/api/wishlist/{user_id}/items", handlers.AddToWishlist).Methods("POST")
	router.HandleFunc("/api/wishlist/{user_id}/items/{product_id}", handlers.RemoveFromWishlist).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

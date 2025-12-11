package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"sample-application/config"
	"sample-application/models"
)

// Order Handlers (PostgreSQL)
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := config.PostgresDB.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	query := `INSERT INTO orders (user_id, total_amount, status, payment_method, shipping_address) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	err = tx.QueryRow(query, order.UserID, order.TotalAmount, order.Status, order.PaymentMethod, order.ShippingAddress).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert order items
	if len(order.Items) > 0 {
		for _, item := range order.Items {
			itemQuery := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)`
			_, err = tx.Exec(itemQuery, order.ID, item.ProductID, item.Quantity, item.Price)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func GetAllOrders(w http.ResponseWriter, r *http.Request) {
	rows, err := config.PostgresDB.Query(`SELECT id, user_id, total_amount, status, payment_method, shipping_address, created_at, updated_at FROM orders LIMIT 100`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	orders := []models.Order{}
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.PaymentMethod, &order.ShippingAddress, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			continue
		}
		orders = append(orders, order)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func GetOrderByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var order models.Order
	query := `SELECT id, user_id, total_amount, status, payment_method, shipping_address, created_at, updated_at FROM orders WHERE id = $1`
	err := config.PostgresDB.QueryRow(query, id).Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.PaymentMethod, &order.ShippingAddress, &order.CreatedAt, &order.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get order items
	itemRows, err := config.PostgresDB.Query(`SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1`, order.ID)
	if err == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var item models.OrderItem
			if err := itemRows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err == nil {
				order.Items = append(order.Items, item)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := data["status"]
	query := `UPDATE orders SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	result, err := config.PostgresDB.Exec(query, status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Order status updated successfully"})
}

func CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	query := `UPDATE orders SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	result, err := config.PostgresDB.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Order cancelled successfully"})
}

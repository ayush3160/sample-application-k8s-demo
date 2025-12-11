package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"sample-application/config"
	"sample-application/models"

	"github.com/gorilla/mux"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// User Handlers (PostgreSQL)
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO users (name, email, password, address, phone) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	err := config.PostgresDB.QueryRow(query, user.Name, user.Email, user.Password, user.Address, user.Phone).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Password = "" // Don't return password
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := config.PostgresDB.Query(`SELECT id, name, email, address, phone, created_at, updated_at FROM users LIMIT 100`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.Phone, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			continue
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user models.User
	query := `SELECT id, name, email, address, phone, created_at, updated_at FROM users WHERE id = $1`
	err := config.PostgresDB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.Phone, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE users SET name = $1, email = $2, address = $3, phone = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`
	result, err := config.PostgresDB.Exec(query, user.Name, user.Email, user.Address, user.Phone, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := config.PostgresDB.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

func GetUserOrders(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	rows, err := config.PostgresDB.Query(`SELECT id, user_id, total_amount, status, payment_method, shipping_address, created_at, updated_at FROM orders WHERE user_id = $1`, userID)
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

// Cart Handlers (PostgreSQL)
func GetCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	rows, err := config.PostgresDB.Query(`SELECT id, user_id, product_id, quantity, created_at, updated_at FROM cart WHERE user_id = $1`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	cartItems := []models.CartItem{}
	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(&item.ID, &item.UserID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			continue
		}
		cartItems = append(cartItems, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartItems)
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	var item models.CartItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userIDInt, _ := strconv.Atoi(userID)
	item.UserID = userIDInt

	query := `INSERT INTO cart (user_id, product_id, quantity) VALUES ($1, $2, $3) 
			  ON CONFLICT (user_id, product_id) DO UPDATE SET quantity = cart.quantity + $3, updated_at = CURRENT_TIMESTAMP
			  RETURNING id, created_at, updated_at`

	// Note: This requires a unique constraint on (user_id, product_id)
	err := config.PostgresDB.QueryRow(query, item.UserID, item.ProductID, item.Quantity).
		Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		// Fallback to simple insert if constraint doesn't exist
		query = `INSERT INTO cart (user_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
		err = config.PostgresDB.QueryRow(query, item.UserID, item.ProductID, item.Quantity).
			Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	itemID := vars["item_id"]

	result, err := config.PostgresDB.Exec(`DELETE FROM cart WHERE id = $1 AND user_id = $2`, itemID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Cart item not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Item removed from cart"})
}

func ClearCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	_, err := config.PostgresDB.Exec(`DELETE FROM cart WHERE user_id = $1`, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Cart cleared successfully"})
}

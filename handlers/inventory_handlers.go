package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"sample-application/config"
	"sample-application/models"
)

// Inventory Handlers (MySQL)
func GetAllInventory(w http.ResponseWriter, r *http.Request) {
	rows, err := config.MySQLDB.Query(`SELECT id, product_id, quantity, warehouse_location, last_restocked, low_stock_threshold, created_at, updated_at FROM inventory LIMIT 100`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	inventory := []models.Inventory{}
	for rows.Next() {
		var item models.Inventory
		err := rows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.WarehouseLocation, &item.LastRestocked, &item.LowStockThreshold, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			continue
		}
		inventory = append(inventory, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventory)
}

func GetInventoryByProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["product_id"]

	var item models.Inventory
	query := `SELECT id, product_id, quantity, warehouse_location, last_restocked, low_stock_threshold, created_at, updated_at FROM inventory WHERE product_id = ?`
	err := config.MySQLDB.QueryRow(query, productID).Scan(&item.ID, &item.ProductID, &item.Quantity, &item.WarehouseLocation, &item.LastRestocked, &item.LowStockThreshold, &item.CreatedAt, &item.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Inventory not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func UpdateInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["product_id"]

	var item models.Inventory
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE inventory SET quantity = ?, warehouse_location = ?, updated_at = NOW() WHERE product_id = ?`
	result, err := config.MySQLDB.Exec(query, item.Quantity, item.WarehouseLocation, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		// Insert if not exists
		insertQuery := `INSERT INTO inventory (product_id, quantity, warehouse_location, last_restocked) VALUES (?, ?, ?, NOW())`
		_, err := config.MySQLDB.Exec(insertQuery, productID, item.Quantity, item.WarehouseLocation)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Inventory updated successfully"})
}

func RestockInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["product_id"]

	var data map[string]int
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	quantity := data["quantity"]
	query := `UPDATE inventory SET quantity = quantity + ?, last_restocked = NOW(), updated_at = NOW() WHERE product_id = ?`
	result, err := config.MySQLDB.Exec(query, quantity, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Inventory not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Inventory restocked successfully"})
}

func GetLowStockItems(w http.ResponseWriter, r *http.Request) {
	rows, err := config.MySQLDB.Query(`SELECT id, product_id, quantity, warehouse_location, last_restocked, low_stock_threshold, created_at, updated_at FROM inventory WHERE quantity <= low_stock_threshold`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	inventory := []models.Inventory{}
	for rows.Next() {
		var item models.Inventory
		err := rows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.WarehouseLocation, &item.LastRestocked, &item.LowStockThreshold, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			continue
		}
		inventory = append(inventory, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventory)
}

// Analytics Handlers (MySQL)
func GetSalesAnalytics(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	query := `SELECT id, product_id, quantity_sold, revenue, sale_date, created_at FROM sales_analytics WHERE sale_date BETWEEN ? AND ? LIMIT 1000`
	
	if startDate == "" {
		startDate = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	rows, err := config.MySQLDB.Query(query, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	analytics := []models.SalesAnalytics{}
	for rows.Next() {
		var item models.SalesAnalytics
		err := rows.Scan(&item.ID, &item.ProductID, &item.QuantitySold, &item.Revenue, &item.SaleDate, &item.CreatedAt)
		if err != nil {
			continue
		}
		analytics = append(analytics, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

func GetPopularProducts(w http.ResponseWriter, r *http.Request) {
	query := `SELECT product_id, SUM(quantity_sold) as total_sold, SUM(revenue) as total_revenue 
			  FROM sales_analytics 
			  GROUP BY product_id 
			  ORDER BY total_sold DESC 
			  LIMIT 20`

	rows, err := config.MySQLDB.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type PopularProduct struct {
		ProductID   string  `json:"product_id"`
		TotalSold   int     `json:"total_sold"`
		TotalRevenue float64 `json:"total_revenue"`
	}

	products := []PopularProduct{}
	for rows.Next() {
		var item PopularProduct
		err := rows.Scan(&item.ProductID, &item.TotalSold, &item.TotalRevenue)
		if err != nil {
			continue
		}
		products = append(products, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func GetRevenueStats(w http.ResponseWriter, r *http.Request) {
	query := `SELECT 
				DATE(sale_date) as date, 
				SUM(revenue) as daily_revenue, 
				COUNT(DISTINCT product_id) as products_sold 
			  FROM sales_analytics 
			  WHERE sale_date >= DATE_SUB(CURDATE(), INTERVAL 30 DAY)
			  GROUP BY DATE(sale_date) 
			  ORDER BY date DESC`

	rows, err := config.MySQLDB.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type RevenueStats struct {
		Date         string  `json:"date"`
		DailyRevenue float64 `json:"daily_revenue"`
		ProductsSold int     `json:"products_sold"`
	}

	stats := []RevenueStats{}
	for rows.Next() {
		var item RevenueStats
		err := rows.Scan(&item.Date, &item.DailyRevenue, &item.ProductsSold)
		if err != nil {
			continue
		}
		stats = append(stats, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

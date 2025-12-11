package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	baseURL              = "http://localhost:8080"
	totalRequests        = 1000
	concurrentRequests   = 50
	delayBetweenRequests = 50 * time.Millisecond // Delay between each request
)

var (
	successCount uint64
	failureCount uint64
	totalLatency uint64
)

type Stats struct {
	TotalRequests      int
	SuccessfulRequests int
	FailedRequests     int
	AverageLatency     time.Duration
	TotalDuration      time.Duration
	RequestsPerSecond  float64
}

// Sample data structures
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
}

type Product struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Category    string   `json:"category"`
	Brand       string   `json:"brand"`
	ImageURL    string   `json:"image_url"`
	Rating      float64  `json:"rating"`
	Tags        []string `json:"tags"`
}

type Order struct {
	UserID          int     `json:"user_id"`
	TotalAmount     float64 `json:"total_amount"`
	Status          string  `json:"status"`
	PaymentMethod   string  `json:"payment_method"`
	ShippingAddress string  `json:"shipping_address"`
}

type Review struct {
	ProductID string `json:"product_id"`
	UserID    int    `json:"user_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}

func main() {
	log.Println("Starting load test...")
	log.Printf("Target: %s", baseURL)
	log.Printf("Total Requests: %d", totalRequests)
	log.Printf("Concurrent Requests: %d", concurrentRequests)
	log.Printf("Delay Between Requests: %v", delayBetweenRequests)

	// Check if server is up
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		log.Fatalf("Server is not reachable: %v", err)
	}
	resp.Body.Close()
	log.Println("Server is healthy, starting load test...")

	startTime := time.Now()

	// Create a semaphore to limit concurrent requests
	semaphore := make(chan struct{}, concurrentRequests)
	var wg sync.WaitGroup

	// Generate and execute requests
	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire semaphore

		go func(requestNum int) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			makeRandomRequest(requestNum)
		}(i)

		// Add delay between each request
		time.Sleep(delayBetweenRequests)
	}

	wg.Wait()
	totalDuration := time.Since(startTime)

	// Print statistics
	printStats(Stats{
		TotalRequests:      totalRequests,
		SuccessfulRequests: int(atomic.LoadUint64(&successCount)),
		FailedRequests:     int(atomic.LoadUint64(&failureCount)),
		AverageLatency:     time.Duration(atomic.LoadUint64(&totalLatency) / uint64(totalRequests)),
		TotalDuration:      totalDuration,
		RequestsPerSecond:  float64(totalRequests) / totalDuration.Seconds(),
	})
}

func makeRandomRequest(requestNum int) {
	// Randomly select an endpoint to test
	endpoints := []func() (string, string, []byte){
		generateHealthCheck,
		generateGetUsers,
		generateGetProducts,
		generateGetOrders,
		generateGetInventory,
		generateGetCategories,
		generateSearchProducts,
		generateGetAnalytics,
		generateGetReviews,
		generateCreateUser,
		generateCreateProduct,
		generateCreateOrder,
		generateCreateReview,
	}

	requestStart := time.Now()

	endpoint := endpoints[rand.Intn(len(endpoints))]
	method, url, body := endpoint()

	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, baseURL+url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, baseURL+url, nil)
	}

	if err != nil {
		atomic.AddUint64(&failureCount, 1)
		log.Printf("Request #%d: Failed to create request: %v", requestNum, err)
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	latency := time.Since(requestStart)
	atomic.AddUint64(&totalLatency, uint64(latency))

	if err != nil {
		atomic.AddUint64(&failureCount, 1)
		log.Printf("Request #%d: Failed (%s %s): %v", requestNum, method, url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		atomic.AddUint64(&successCount, 1)
		if requestNum%100 == 0 {
			log.Printf("Request #%d: Success (%s %s) - Status: %d, Latency: %v",
				requestNum, method, url, resp.StatusCode, latency)
		}
	} else {
		atomic.AddUint64(&failureCount, 1)
		log.Printf("Request #%d: Failed (%s %s) - Status: %d, Latency: %v",
			requestNum, method, url, resp.StatusCode, latency)
	}
}

// Request generators
func generateHealthCheck() (string, string, []byte) {
	return "GET", "/health", nil
}

func generateGetUsers() (string, string, []byte) {
	return "GET", "/api/users", nil
}

func generateGetProducts() (string, string, []byte) {
	return "GET", "/api/products", nil
}

func generateGetOrders() (string, string, []byte) {
	return "GET", "/api/orders", nil
}

func generateGetInventory() (string, string, []byte) {
	return "GET", "/api/inventory", nil
}

func generateGetCategories() (string, string, []byte) {
	return "GET", "/api/categories", nil
}

func generateSearchProducts() (string, string, []byte) {
	queries := []string{"laptop", "phone", "book", "shoes", "watch", "camera"}
	query := queries[rand.Intn(len(queries))]
	return "GET", "/api/products/search?q=" + query, nil
}

func generateGetAnalytics() (string, string, []byte) {
	endpoints := []string{
		"/api/analytics/sales?start_date=2024-01-01&end_date=2024-12-31",
		"/api/analytics/popular-products",
		"/api/analytics/revenue",
	}
	return "GET", endpoints[rand.Intn(len(endpoints))], nil
}

func generateGetReviews() (string, string, []byte) {
	productID := fmt.Sprintf("prod%d", rand.Intn(100)+1)
	return "GET", "/api/reviews/product/" + productID, nil
}

func generateCreateUser() (string, string, []byte) {
	user := User{
		Name:     fmt.Sprintf("User_%d", rand.Intn(10000)),
		Email:    fmt.Sprintf("user%d@example.com", rand.Intn(10000)),
		Password: "password123",
		Address:  fmt.Sprintf("%d Main St, City, State", rand.Intn(1000)),
		Phone:    fmt.Sprintf("+1-555-%04d", rand.Intn(10000)),
	}
	body, _ := json.Marshal(user)
	return "POST", "/api/users", body
}

func generateCreateProduct() (string, string, []byte) {
	categories := []string{"Electronics", "Clothing", "Books", "Home", "Sports"}
	brands := []string{"BrandA", "BrandB", "BrandC", "BrandD", "BrandE"}
	tags := [][]string{
		{"new", "sale", "popular"},
		{"featured", "bestseller"},
		{"limited", "exclusive"},
	}

	product := Product{
		Name:        fmt.Sprintf("Product_%d", rand.Intn(10000)),
		Description: fmt.Sprintf("Description for product %d", rand.Intn(10000)),
		Price:       float64(rand.Intn(1000)) + 0.99,
		Category:    categories[rand.Intn(len(categories))],
		Brand:       brands[rand.Intn(len(brands))],
		ImageURL:    fmt.Sprintf("https://example.com/image%d.jpg", rand.Intn(100)),
		Rating:      float64(rand.Intn(5)) + 1.0,
		Tags:        tags[rand.Intn(len(tags))],
	}
	body, _ := json.Marshal(product)
	return "POST", "/api/products", body
}

func generateCreateOrder() (string, string, []byte) {
	statuses := []string{"pending", "processing", "shipped", "delivered"}
	paymentMethods := []string{"credit_card", "debit_card", "paypal", "cash"}

	order := Order{
		UserID:          rand.Intn(100) + 1,
		TotalAmount:     float64(rand.Intn(500)) + 0.99,
		Status:          statuses[rand.Intn(len(statuses))],
		PaymentMethod:   paymentMethods[rand.Intn(len(paymentMethods))],
		ShippingAddress: fmt.Sprintf("%d Shipping St, City, State", rand.Intn(1000)),
	}
	body, _ := json.Marshal(order)
	return "POST", "/api/orders", body
}

func generateCreateReview() (string, string, []byte) {
	comments := []string{
		"Great product!",
		"Excellent quality",
		"Very satisfied",
		"Could be better",
		"Amazing purchase",
	}

	review := Review{
		ProductID: fmt.Sprintf("prod%d", rand.Intn(100)+1),
		UserID:    rand.Intn(100) + 1,
		Rating:    rand.Intn(5) + 1,
		Comment:   comments[rand.Intn(len(comments))],
	}
	body, _ := json.Marshal(review)
	return "POST", "/api/reviews", body
}

func printStats(stats Stats) {
	separator := "============================================================"
	fmt.Println("\n" + separator)
	fmt.Println("LOAD TEST RESULTS")
	fmt.Println(separator)
	fmt.Printf("Total Requests:       %d\n", stats.TotalRequests)
	fmt.Printf("Successful Requests:  %d (%.2f%%)\n",
		stats.SuccessfulRequests,
		float64(stats.SuccessfulRequests)/float64(stats.TotalRequests)*100)
	fmt.Printf("Failed Requests:      %d (%.2f%%)\n",
		stats.FailedRequests,
		float64(stats.FailedRequests)/float64(stats.TotalRequests)*100)
	fmt.Printf("Average Latency:      %v\n", stats.AverageLatency)
	fmt.Printf("Total Duration:       %v\n", stats.TotalDuration)
	fmt.Printf("Requests/Second:      %.2f\n", stats.RequestsPerSecond)
	fmt.Println(separator)
}

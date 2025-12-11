# E-commerce API - Go Application

A comprehensive e-commerce REST API built with Go, featuring multiple databases (PostgreSQL, MySQL, MongoDB) and multiple routes for a realistic production-like application.

## 🚀 Features

- **Multiple Database Integration**:
  - PostgreSQL: User management, Orders, Shopping Cart
  - MySQL: Inventory management, Sales Analytics
  - MongoDB: Products, Categories, Reviews, Wishlist

- **50+ API Endpoints** covering:
  - User Management (CRUD)
  - Product Management (CRUD + Search)
  - Order Processing
  - Shopping Cart
  - Inventory Tracking
  - Product Reviews
  - Categories
  - Wishlist
  - Sales Analytics

- **Load Testing Tool**: Generate thousands of concurrent requests
- **Containerized**: Docker support with multi-stage builds
- **Kubernetes Ready**: Full K8s manifests with auto-scaling

## 📋 Prerequisites

- Go 1.21+
- Docker & Docker Compose (optional)
- Kubernetes cluster (for K8s deployment)
- PostgreSQL, MySQL, MongoDB (if running locally)

## 🏗️ Project Structure

```
sample-application/
├── main.go                 # Application entry point
├── config/
│   └── database.go        # Database connection management
├── models/
│   └── models.go          # Data models
├── handlers/
│   ├── user_handlers.go   # User & cart endpoints
│   ├── product_handlers.go # Product & category endpoints
│   ├── order_handlers.go  # Order management endpoints
│   ├── inventory_handlers.go # Inventory & analytics endpoints
│   └── review_handlers.go # Review & wishlist endpoints
├── load_test.go           # Load testing program
├── Dockerfile             # Multi-stage Docker build
├── k8s/                   # Kubernetes manifests
│   ├── 00-namespace-config.yaml
│   ├── 01-postgres.yaml
│   ├── 02-mysql.yaml
│   ├── 03-mongodb.yaml
│   └── 04-api-deployment.yaml
├── go.mod
├── go.sum
└── README.md
```

## 🔧 Setup & Installation

### Local Development

1. **Clone the repository**
   ```bash
   cd sample-application
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. **Start databases** (using Docker)
   ```bash
   # PostgreSQL
   docker run -d --name postgres -p 5432:5432 \
     -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=ecommerce postgres:15-alpine

   # MySQL
   docker run -d --name mysql -p 3306:3306 \
     -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=ecommerce mysql:8.0

   # MongoDB
   docker run -d --name mongodb -p 27017:27017 mongo:7.0
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

   The API will be available at `http://localhost:8080`

### Using Docker

1. **Build the Docker image**
   ```bash
   docker build -t ecommerce-api:latest .
   ```

2. **Run with Docker Compose** (create docker-compose.yml first)
   ```bash
   docker-compose up -d
   ```

## 🐳 Docker Deployment

### Build Image
```bash
docker build -t ecommerce-api:latest .
```

### Run Container
```bash
docker run -d -p 8080:8080 \
  -e POSTGRES_HOST=host.docker.internal \
  -e MYSQL_HOST=host.docker.internal \
  -e MONGO_HOST=host.docker.internal \
  ecommerce-api:latest
```

## ☸️ Kubernetes Deployment

### Prerequisites
- Kubernetes cluster (minikube, kind, or cloud provider)
- kubectl configured

### Deploy to Kubernetes

1. **Apply all manifests**
   ```bash
   kubectl apply -f k8s/
   ```

2. **Check deployment status**
   ```bash
   kubectl get pods -n ecommerce
   kubectl get services -n ecommerce
   ```

3. **Access the API**
   ```bash
   # Get the LoadBalancer IP
   kubectl get service ecommerce-api-service -n ecommerce
   ```

### Kubernetes Features

- **3 API replicas** (scales 3-10 based on load)
- **Horizontal Pod Autoscaling** based on CPU/Memory
- **Persistent volumes** for all databases
- **Health checks** (liveness & readiness probes)
- **Resource limits** for all pods
- **Separate namespace** for isolation

### Scale manually
```bash
kubectl scale deployment ecommerce-api -n ecommerce --replicas=5
```

## 📊 Load Testing

Run the included load testing program to generate thousands of requests:

```bash
# Build the load tester
cd loadtest
go build -o load_test main.go

# Run load test (1000 requests, 50 concurrent)
./load_test
```

Or use Make from the root directory:
```bash
make load-test
```

The load tester will:
- Make 1000 total requests
- Use 50 concurrent workers
- Test random endpoints (GET, POST, PUT, DELETE)
- Generate realistic data
- Report statistics (success rate, latency, throughput)

### Customize load test
Edit `loadtest/main.go` and modify:
```go
const (
    totalRequests      = 1000  // Change total requests
    concurrentRequests = 50    // Change concurrency
)
```

## 🔌 API Endpoints

### Health Check
- `GET /health` - Service health status

### Users
- `POST /api/users` - Create user
- `GET /api/users` - List all users
- `GET /api/users/{id}` - Get user by ID
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user
- `GET /api/users/{id}/orders` - Get user's orders

### Products
- `POST /api/products` - Create product
- `GET /api/products` - List all products
- `GET /api/products/{id}` - Get product by ID
- `PUT /api/products/{id}` - Update product
- `DELETE /api/products/{id}` - Delete product
- `GET /api/products/search?q={query}` - Search products
- `GET /api/products/category/{category}` - Get products by category

### Orders
- `POST /api/orders` - Create order
- `GET /api/orders` - List all orders
- `GET /api/orders/{id}` - Get order by ID
- `PATCH /api/orders/{id}/status` - Update order status
- `POST /api/orders/{id}/cancel` - Cancel order

### Inventory
- `GET /api/inventory` - List all inventory
- `GET /api/inventory/{product_id}` - Get inventory for product
- `PUT /api/inventory/{product_id}` - Update inventory
- `POST /api/inventory/{product_id}/restock` - Restock item
- `GET /api/inventory/low-stock` - Get low stock items

### Reviews
- `POST /api/reviews` - Create review
- `GET /api/reviews/product/{product_id}` - Get product reviews
- `DELETE /api/reviews/{id}` - Delete review
- `POST /api/reviews/{id}/helpful` - Mark review as helpful

### Categories
- `POST /api/categories` - Create category
- `GET /api/categories` - List all categories
- `GET /api/categories/{id}` - Get category by ID
- `PUT /api/categories/{id}` - Update category
- `DELETE /api/categories/{id}` - Delete category

### Shopping Cart
- `GET /api/cart/{user_id}` - Get user's cart
- `POST /api/cart/{user_id}/items` - Add item to cart
- `DELETE /api/cart/{user_id}/items/{item_id}` - Remove item from cart
- `DELETE /api/cart/{user_id}/clear` - Clear cart

### Analytics
- `GET /api/analytics/sales` - Get sales analytics
- `GET /api/analytics/popular-products` - Get popular products
- `GET /api/analytics/revenue` - Get revenue statistics

### Wishlist
- `GET /api/wishlist/{user_id}` - Get user's wishlist
- `POST /api/wishlist/{user_id}/items` - Add to wishlist
- `DELETE /api/wishlist/{user_id}/items/{product_id}` - Remove from wishlist

## 📝 Example Requests

### Create a User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "address": "123 Main St",
    "phone": "+1-555-0100"
  }'
```

### Create a Product
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99,
    "category": "Electronics",
    "brand": "TechBrand",
    "rating": 4.5,
    "tags": ["new", "sale"]
  }'
```

### Search Products
```bash
curl "http://localhost:8080/api/products/search?q=laptop"
```

## 🛠️ Development

### Run tests
```bash
go test ./...
```

### Format code
```bash
go fmt ./...
```

### Build for production
```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
```

## 📦 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | API server port | `8080` |
| `POSTGRES_HOST` | PostgreSQL host | `localhost` |
| `POSTGRES_PORT` | PostgreSQL port | `5432` |
| `POSTGRES_USER` | PostgreSQL user | `postgres` |
| `POSTGRES_PASSWORD` | PostgreSQL password | `postgres` |
| `POSTGRES_DB` | PostgreSQL database | `ecommerce` |
| `MYSQL_HOST` | MySQL host | `localhost` |
| `MYSQL_PORT` | MySQL port | `3306` |
| `MYSQL_USER` | MySQL user | `root` |
| `MYSQL_PASSWORD` | MySQL password | `root` |
| `MYSQL_DB` | MySQL database | `ecommerce` |
| `MONGO_HOST` | MongoDB host | `localhost` |
| `MONGO_PORT` | MongoDB port | `27017` |
| `MONGO_USER` | MongoDB user | `` |
| `MONGO_PASSWORD` | MongoDB password | `` |
| `MONGO_DB` | MongoDB database | `ecommerce` |

## 🎯 Performance

- Supports thousands of concurrent requests
- Connection pooling for all databases
- Horizontal scaling with Kubernetes
- Auto-scaling based on CPU/Memory metrics

## 📄 License

MIT License

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

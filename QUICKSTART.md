# Quick Start Guide

## Prerequisites
- Go 1.21+ installed
- Docker & Docker Compose installed
- kubectl and Kubernetes cluster (for K8s deployment)

## 🚀 Quick Start with Docker Compose (Recommended)

This is the fastest way to get everything running:

```bash
# 1. Navigate to the project directory
cd sample-application

# 2. Start all services (databases + API)
docker-compose up -d

# 3. Wait for services to be healthy (~30 seconds)
docker-compose ps

# 4. Test the API
curl http://localhost:8080/health

# 5. View logs
docker-compose logs -f api
```

The API will be available at `http://localhost:8080`

### Stop Services
```bash
docker-compose down
```

### With Volumes Cleanup
```bash
docker-compose down -v
```

## 🔧 Local Development (Without Docker)

### 1. Start Databases

**PostgreSQL:**
```bash
docker run -d --name postgres -p 5432:5432 \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=ecommerce \
  postgres:15-alpine
```

**MySQL:**
```bash
docker run -d --name mysql -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=ecommerce \
  mysql:8.0
```

**MongoDB:**
```bash
docker run -d --name mongodb -p 27017:27017 \
  mongo:7.0
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Run the Application
```bash
go run main.go
```

Or using Make:
```bash
make run
```

## 📊 Load Testing

### Build and Run Load Test
```bash
# Using Make
make load-test

# Or manually
cd loadtest
go build -o load_test main.go
./load_test
```

This will:
- Generate 1000 requests
- Use 50 concurrent workers
- Test all endpoints randomly
- Display performance statistics

## ☸️ Kubernetes Deployment

### Using the Deploy Script (Easy)
```bash
./deploy-k8s.sh
```

### Manual Deployment
```bash
# 1. Build Docker image
docker build -t ecommerce-api:latest .

# 2. If using minikube
minikube image load ecommerce-api:latest

# 3. Apply manifests
kubectl apply -f k8s/

# 4. Check status
kubectl get pods -n ecommerce
kubectl get services -n ecommerce

# 5. Access the service (minikube)
minikube service ecommerce-api-service -n ecommerce
```

### Delete Deployment
```bash
kubectl delete namespace ecommerce
```

## 🧪 Testing the API

### Health Check
```bash
curl http://localhost:8080/health
```

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

### Get All Users
```bash
curl http://localhost:8080/api/users
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

### Get All Products
```bash
curl http://localhost:8080/api/products
```

### Create an Order
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "total_amount": 999.99,
    "status": "pending",
    "payment_method": "credit_card",
    "shipping_address": "123 Main St"
  }'
```

### Get Inventory
```bash
curl http://localhost:8080/api/inventory
```

### Get Analytics
```bash
curl http://localhost:8080/api/analytics/sales
curl http://localhost:8080/api/analytics/popular-products
curl http://localhost:8080/api/analytics/revenue
```

## 🛠️ Using Make Commands

The project includes a Makefile for common tasks:

```bash
make help              # Show all available commands
make build             # Build the application
make run               # Run the application
make test              # Run tests
make clean             # Clean build artifacts
make load-test         # Run load tests
make docker-build      # Build Docker image
make docker-run        # Run with Docker Compose
make docker-stop       # Stop Docker Compose
make docker-logs       # View Docker logs
make k8s-deploy        # Deploy to Kubernetes
make k8s-delete        # Delete K8s deployment
make k8s-status        # Check K8s status
make k8s-logs          # View K8s logs
make fmt               # Format code
```

## 🐛 Troubleshooting

### Database Connection Issues
1. Check if databases are running:
   ```bash
   docker ps
   ```

2. Check database logs:
   ```bash
   docker logs postgres
   docker logs mysql
   docker logs mongodb
   ```

3. Verify environment variables in `.env` file

### Port Already in Use
If port 8080 is already in use, change it in `.env`:
```bash
PORT=3000
```

### Kubernetes Pod Not Starting
1. Check pod status:
   ```bash
   kubectl describe pod <pod-name> -n ecommerce
   ```

2. Check logs:
   ```bash
   kubectl logs <pod-name> -n ecommerce
   ```

3. Ensure databases are ready before API starts

## 📚 Next Steps

1. **Explore API Endpoints**: Check `README.md` for complete API documentation
2. **Run Load Tests**: Test performance with the load testing tool
3. **Deploy to Production**: Use Kubernetes manifests for production deployment
4. **Monitor**: Set up monitoring and logging for production use
5. **Scale**: Adjust replicas and resource limits in K8s manifests

## 🔗 Resources

- [Go Documentation](https://golang.org/doc/)
- [Docker Documentation](https://docs.docker.com/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [MySQL Documentation](https://dev.mysql.com/doc/)
- [MongoDB Documentation](https://docs.mongodb.com/)

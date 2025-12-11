#!/bin/bash

# Script to deploy the application to Kubernetes

set -e

echo "======================================"
echo "E-commerce API Kubernetes Deployment"
echo "======================================"

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "Error: kubectl is not installed"
    exit 1
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Error: docker is not installed"
    exit 1
fi

# Build Docker image
echo ""
echo "Building Docker image..."
docker build -t ecommerce-api:latest .

# If using minikube, load image into minikube
if command -v minikube &> /dev/null && minikube status &> /dev/null; then
    echo ""
    echo "Loading image into minikube..."
    minikube image load ecommerce-api:latest
fi

# Apply Kubernetes manifests
echo ""
echo "Applying Kubernetes manifests..."
kubectl apply -f k8s/00-namespace-config.yaml
echo "Waiting for namespace to be ready..."
sleep 2

kubectl apply -f k8s/01-postgres.yaml
echo "Waiting for PostgreSQL to be ready..."
sleep 5

kubectl apply -f k8s/02-mysql.yaml
echo "Waiting for MySQL to be ready..."
sleep 5

kubectl apply -f k8s/03-mongodb.yaml
echo "Waiting for MongoDB to be ready..."
sleep 5

kubectl apply -f k8s/04-api-deployment.yaml

echo ""
echo "======================================"
echo "Deployment completed!"
echo "======================================"

echo ""
echo "Checking deployment status..."
kubectl get pods -n ecommerce

echo ""
echo "Waiting for pods to be ready (this may take a few minutes)..."
kubectl wait --for=condition=ready pod -l app=postgres -n ecommerce --timeout=120s || true
kubectl wait --for=condition=ready pod -l app=mysql -n ecommerce --timeout=120s || true
kubectl wait --for=condition=ready pod -l app=mongodb -n ecommerce --timeout=120s || true
kubectl wait --for=condition=ready pod -l app=ecommerce-api -n ecommerce --timeout=120s || true

echo ""
echo "======================================"
echo "Final Status"
echo "======================================"
kubectl get all -n ecommerce

echo ""
echo "======================================"
echo "Access Information"
echo "======================================"
echo ""
echo "To access the API, get the service URL:"
echo "  kubectl get service ecommerce-api-service -n ecommerce"
echo ""
echo "If using minikube, run:"
echo "  minikube service ecommerce-api-service -n ecommerce"
echo ""
echo "To view logs:"
echo "  kubectl logs -f deployment/ecommerce-api -n ecommerce"
echo ""
echo "To delete the deployment:"
echo "  kubectl delete namespace ecommerce"
echo ""

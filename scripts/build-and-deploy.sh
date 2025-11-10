#!/bin/bash
set -e

echo "======================================"
echo "Pod Index - Build and Deploy Script"
echo "======================================"

# Check if necessary tools are installed
command -v docker >/dev/null 2>&1 || { echo "Error: Docker is required"; exit 1; }
command -v kubectl >/dev/null 2>&1 || { echo "Error: kubectl is required"; exit 1; }

# Build Docker image
echo ""
echo "[1/3] Building Docker image..."
docker build -t pod-index:latest .

# Load image to cluster if using Minikube or Kind
if command -v minikube >/dev/null 2>&1 && minikube status >/dev/null 2>&1; then
    echo ""
    echo "[2/3] Loading image to Minikube..."
    minikube image load pod-index:latest
elif command -v kind >/dev/null 2>&1; then
    CLUSTER_NAME=${KIND_CLUSTER_NAME:-kind}
    if kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
        echo ""
        echo "[2/3] Loading image to Kind..."
        kind load docker-image pod-index:latest --name ${CLUSTER_NAME}
    fi
else
    echo ""
    echo "[2/3] Skipping image load (not a local cluster)"
fi

# Deploy to Kubernetes
echo ""
echo "[3/3] Deploying to Kubernetes..."
kubectl apply -k deploy/

echo ""
echo "======================================"
echo "Deployment completed!"
echo "======================================"
echo ""
echo "Check deployment status:"
echo "  kubectl get pods -l app=pod-index"
echo ""
echo "View logs:"
echo "  kubectl logs -l app=pod-index -f"
echo ""
echo "Port forward for testing:"
echo "  kubectl port-forward svc/pod-index 8080:80"
echo ""

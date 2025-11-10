# Quick Start Guide

## Deploy to Kubernetes in Three Steps

### 1. Build Image
```bash
docker build -t pod-index:latest .
```

If using Minikube:
```bash
minikube image load pod-index:latest
```

If using Kind:
```bash
kind load docker-image pod-index:latest
```

### 2. Deploy to Cluster
```bash
kubectl apply -k deploy/
```

### 3. Verify Deployment
```bash
# Check Pod status
kubectl get pods -l app=pod-index

# Port forward
kubectl port-forward svc/pod-index 8080:80
```

## Test API

```bash
# Health check
curl http://localhost:8080/health

# Get any Pod UID and query
POD_UID=$(kubectl get pod -A -o jsonpath='{.items[0].metadata.uid}')
curl "http://localhost:8080/api/v1/pod?uid=${POD_UID}"
```

## One-Click Deploy (Recommended)

```bash
# Automatically build, load image and deploy
./scripts/build-and-deploy.sh

# Automatically test API
./scripts/test-api.sh http://localhost:8080
```

## Local Development

```bash
# Run service (requires valid kubeconfig)
go run main.go

# Or use Make
make run
```

## Uninstall

```bash
kubectl delete -k deploy/
```

## Common Commands

```bash
make help           # View all available commands
make build          # Compile
make docker-build   # Build Docker image
make deploy         # Deploy to K8s
make undeploy       # Uninstall
```

For complete documentation, see [README_CN.md](README_CN.md)

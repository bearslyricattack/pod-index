# Pod Index

A high-performance Pod information caching service based on Kubernetes Informer.

## Project Overview

Pod Index uses the Kubernetes Informer mechanism to cache information for all Pods in the cluster in real-time and provides a simple HTTP API for querying. The service features:

- **Efficient Caching**: Leverages Kubernetes Informer mechanism to automatically sync and cache all Pod information
- **Fast Queries**: Quickly retrieve Pod details by Pod UID without accessing the API Server each time
- **Resource Friendly**: Lightweight design with minimal resource footprint
- **Production Ready**: Complete health checks, graceful shutdown, and error handling mechanisms
- **Security Compliant**: Follows the principle of least privilege with container security hardening

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Kubernetes cluster (local or remote)
- kubectl configured and accessible to cluster

### Local Execution

1. **Clone Project**
```bash
git clone <your-repo-url>
cd pod-index
```

2. **Download Dependencies**
```bash
go mod download
```

3. **Run Service**
```bash
go run main.go
```

The service will start at `http://localhost:8080`.

### Test API

```bash
# Health check
curl http://localhost:8080/health

# Readiness check
curl http://localhost:8080/ready

# Query Pod information
POD_UID=$(kubectl get pod -n default -o jsonpath='{.items[0].metadata.uid}')
curl "http://localhost:8080/api/v1/pod?uid=${POD_UID}"
```

## API Documentation

### 1. Query Pod Information

**Request**
```
GET /api/v1/pod?uid={pod-uid}
```

**Parameters**
- `uid`: Pod UID (required)

**Success Response** (200 OK)
```json
{
  "uid": "1234abcd-5678-90ef-ghij-klmnopqrstuv",
  "name": "example-pod",
  "namespace": "default",
  "nodeName": "node-1",
  "phase": "Running",
  "podIP": "10.244.1.5",
  "labels": {
    "app": "example"
  },
  "annotations": {
    "kubernetes.io/created-by": "..."
  },
  "createdAt": "2024-01-01T00:00:00Z"
}
```

**Error Response** (404 Not Found)
```json
{
  "error": "pod with UID xxx not found"
}
```

### 2. Health Check

**Request**
```
GET /health
```

**Response** (200 OK)
```json
{
  "status": "healthy"
}
```

### 3. Readiness Check

**Request**
```
GET /ready
```

**Response** (200 OK)
```json
{
  "status": "ready",
  "podCount": 150
}
```

## Docker Deployment

### Build Image

```bash
docker build -t pod-index:latest .
```

### Local Run (using kubeconfig)

```bash
docker run -d \
  -p 8080:8080 \
  -v ~/.kube/config:/root/.kube/config:ro \
  --name pod-index \
  pod-index:latest
```

### Test Container

```bash
# View logs
docker logs -f pod-index

# Test API
curl http://localhost:8080/health
```

## Kubernetes Deployment

### Method 1: Using Kustomize (Recommended)

```bash
# Deploy all resources
kubectl apply -k deploy/

# View deployment status
kubectl get pods -l app=pod-index
kubectl logs -l app=pod-index -f
```

### Method 2: Using Automation Script

```bash
# Build and deploy
./scripts/build-and-deploy.sh
```

### Method 3: Manual Deployment

```bash
# 1. Create RBAC permissions
kubectl apply -f deploy/rbac.yaml

# 2. Deploy application
kubectl apply -f deploy/deployment.yaml

# 3. Create service
kubectl apply -f deploy/service.yaml
```

### Verify Deployment

```bash
# Check Pod status
kubectl get pods -l app=pod-index

# View logs
kubectl logs -l app=pod-index -f

# Port forward for testing
kubectl port-forward svc/pod-index 8080:80

# Test in another terminal
./scripts/test-api.sh http://localhost:8080
```

### Uninstall

```bash
kubectl delete -k deploy/
# or
make undeploy
```

## Makefile Commands

The project provides convenient Makefile commands:

```bash
make help           # Display all available commands
make build          # Compile application
make run            # Run locally
make test           # Run tests
make docker-build   # Build Docker image
make docker-run     # Run Docker container
make deploy         # Deploy to Kubernetes
make undeploy       # Uninstall from Kubernetes
make clean          # Clean build files
make deps           # Download and tidy dependencies
```

## Project Structure

```
pod-index/
├── main.go                      # Application entry point
├── pkg/
│   ├── cache/                  # Informer cache layer
│   │   └── pod_cache.go       # Pod cache implementation
│   └── handler/               # HTTP handler layer
│       └── handler.go         # API handler
├── deploy/                    # Kubernetes deployment files
│   ├── rbac.yaml             # RBAC permission configuration
│   ├── deployment.yaml       # Deployment configuration
│   ├── service.yaml          # Service configuration
│   └── kustomization.yaml    # Kustomize configuration
├── scripts/                  # Helper scripts
│   ├── build-and-deploy.sh  # Automated build and deploy
│   └── test-api.sh          # API test script
├── Dockerfile               # Docker build file
├── .dockerignore           # Docker ignore file
├── Makefile                # Make commands
├── go.mod                  # Go module definition
├── go.sum                  # Go dependency lock
├── LICENSE                 # Open source license
└── README.md              # Project documentation
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | HTTP service listening port | 8080 |

### Kubernetes Configuration

The application looks for Kubernetes configuration in the following order:

1. In-Cluster Config - for deployment in Kubernetes
2. Kubeconfig file (`~/.kube/config`) - for local development

### RBAC Permissions

The service requires the following minimum permissions:

```yaml
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "watch"]
```

## Security Best Practices

This project implements the following security measures:

1. **Principle of Least Privilege**: Only requests necessary RBAC permissions
2. **Non-root Execution**: Container runs as UID 65534
3. **Read-only Root Filesystem**: Prevents runtime file tampering
4. **Disable Privilege Escalation**: Prevents privilege escalation attacks
5. **Drop All Capabilities**: Minimizes container capabilities
6. **Resource Limits**: Sets CPU and memory limits to prevent resource exhaustion

## FAQ

### Q: Service cannot connect to Kubernetes API

**A:** Check the following:
- Ensure kubeconfig is configured correctly (local run)
- Ensure ServiceAccount and RBAC are configured correctly (cluster run)
- Check network connection and firewall settings

### Q: Pod query returns 404

**A:** Possible reasons:
- Cache not yet synced, retry after a few seconds
- Pod UID is incorrect, use `kubectl get pod -o yaml` to view the correct UID
- Pod has been deleted

### Q: How to deploy in production?

**A:** Recommendations:
1. Change image pull policy to `Always`
2. Adjust resource limits based on cluster size
3. Configure HPA (Horizontal Pod Autoscaler) if needed
4. Enable monitoring and log collection
5. Use Ingress or LoadBalancer to expose service

### Q: Does it support multi-cluster?

**A:** Current version only supports single cluster. For multi-cluster support, you can:
- Deploy one instance in each cluster
- Use a unified frontend to aggregate queries

## Performance Considerations

- **Memory Usage**: Approximately 100-200MB + (number of Pods × 2KB)
- **Startup Time**: Usually 5-10 seconds to complete cache sync
- **Query Latency**: < 1ms (in-memory query)
- **Concurrency**: Supports high concurrent reads

## Contributing

Contributions are welcome! Please:

1. Fork this project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see [LICENSE](LICENSE) file for details.

## Contact

For questions or suggestions, please submit an Issue.

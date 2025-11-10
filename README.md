# Pod Index

Pod Index is a Pod information caching service based on Kubernetes Informer, providing an efficient HTTP API for Pod queries.

## Features

- Uses Kubernetes Informer to cache all Pod information from the cluster in real-time
- Provides HTTP API to query Pod details by Pod UID
- Health check and readiness check endpoints
- Lightweight with low resource usage
- Supports both in-cluster and out-of-cluster operation

## API Endpoints

### Query Pod Information

**GET** `/api/v1/pod?uid={pod-uid}`

Query Parameters:
- `uid`: Pod UID (required)

Response Example:
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
  "annotations": {},
  "createdAt": "2024-01-01T00:00:00Z"
}
```

### Health Check

**GET** `/health`

Response Example:
```json
{
  "status": "healthy"
}
```

### Readiness Check

**GET** `/ready`

Response Example:
```json
{
  "status": "ready",
  "podCount": 150
}
```

## Local Development

### Prerequisites

- Go 1.21+
- kubectl configured to access a Kubernetes cluster
- kubeconfig file located at `~/.kube/config`

### Run

```bash
# Download dependencies
go mod download

# Run service
go run main.go
```

The service will start at `http://localhost:8080`.

### Test API

```bash
# Get a Pod UID
POD_UID=$(kubectl get pod -n default -o jsonpath='{.items[0].metadata.uid}')

# Query Pod information
curl "http://localhost:8080/api/v1/pod?uid=${POD_UID}"

# Health check
curl http://localhost:8080/health

# Readiness check
curl http://localhost:8080/ready
```

## Docker Build

```bash
# Build image
docker build -t pod-index:latest .

# Run locally (mount kubeconfig)
docker run -d \
  -p 8080:8080 \
  -v ~/.kube/config:/root/.kube/config:ro \
  pod-index:latest
```

## Kubernetes Deployment

### Deploy to Cluster

```bash
# Apply all deployment files
kubectl apply -k deploy/

# Or apply individually
kubectl apply -f deploy/rbac.yaml
kubectl apply -f deploy/deployment.yaml
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
```

### Test Service

```bash
# Get a Pod UID
POD_UID=$(kubectl get pod -n default -o jsonpath='{.items[0].metadata.uid}')

# Test through port forward
curl "http://localhost:8080/api/v1/pod?uid=${POD_UID}"
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | HTTP service listening port | 8080 |

## Security

- Follows the principle of least privilege, only requires `get`, `list`, `watch` pods permissions
- Container runs as non-root user
- Enables read-only root filesystem
- Disables privilege escalation

## Project Structure

```
.
├── main.go                 # Main program entry point
├── pkg/
│   ├── cache/             # Informer cache implementation
│   │   └── pod_cache.go
│   └── handler/           # HTTP handlers
│       └── handler.go
├── deploy/                # Kubernetes deployment files
│   ├── rbac.yaml         # RBAC permission configuration
│   ├── deployment.yaml   # Deployment configuration
│   ├── service.yaml      # Service configuration
│   └── kustomization.yaml
├── Dockerfile            # Docker build file
├── go.mod               # Go module definition
└── README.md            # Project documentation
```

## License

MIT

# Pod Index

Pod Index 是一个基于 Kubernetes Informer 的 Pod 信息缓存服务，提供高效的 Pod 查询 HTTP API。

## 功能特性

- 使用 Kubernetes Informer 实时缓存集群中所有 Pod 信息
- 提供 HTTP API 通过 Pod UID 查询 Pod 详细信息
- 健康检查和就绪检查端点
- 轻量级，资源占用低
- 支持集群内和集群外运行

## API 接口

### 查询 Pod 信息

**GET** `/api/v1/pod?uid={pod-uid}`

查询参数：
- `uid`: Pod 的 UID（必填）

响应示例：
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

### 健康检查

**GET** `/health`

响应示例：
```json
{
  "status": "healthy"
}
```

### 就绪检查

**GET** `/ready`

响应示例：
```json
{
  "status": "ready",
  "podCount": 150
}
```

## 本地开发

### 前置要求

- Go 1.21+
- kubectl 配置可访问的 Kubernetes 集群
- kubeconfig 文件位于 `~/.kube/config`

### 运行

```bash
# 下载依赖
go mod download

# 运行服务
go run main.go
```

服务将在 `http://localhost:8080` 启动。

### 测试 API

```bash
# 获取一个 Pod 的 UID
POD_UID=$(kubectl get pod -n default -o jsonpath='{.items[0].metadata.uid}')

# 查询 Pod 信息
curl "http://localhost:8080/api/v1/pod?uid=${POD_UID}"

# 健康检查
curl http://localhost:8080/health

# 就绪检查
curl http://localhost:8080/ready
```

## Docker 构建

```bash
# 构建镜像
docker build -t pod-index:latest .

# 本地运行（需要挂载 kubeconfig）
docker run -d \
  -p 8080:8080 \
  -v ~/.kube/config:/root/.kube/config:ro \
  pod-index:latest
```

## Kubernetes 部署

### 部署到集群

```bash
# 应用所有部署文件
kubectl apply -k deploy/

# 或者分别应用
kubectl apply -f deploy/rbac.yaml
kubectl apply -f deploy/deployment.yaml
kubectl apply -f deploy/service.yaml
```

### 验证部署

```bash
# 检查 Pod 状态
kubectl get pods -l app=pod-index

# 查看日志
kubectl logs -l app=pod-index -f

# 端口转发测试
kubectl port-forward svc/pod-index 8080:80
```

### 测试服务

```bash
# 获取一个 Pod 的 UID
POD_UID=$(kubectl get pod -n default -o jsonpath='{.items[0].metadata.uid}')

# 通过端口转发测试
curl "http://localhost:8080/api/v1/pod?uid=${POD_UID}"
```

## 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| PORT | HTTP 服务监听端口 | 8080 |

## 安全性

- 遵循最小权限原则，仅需要 `get`, `list`, `watch` pods 权限
- 容器以非 root 用户运行
- 启用只读根文件系统
- 禁用特权提升

## 项目结构

```
.
├── main.go                 # 主程序入口
├── pkg/
│   ├── cache/             # Informer 缓存实现
│   │   └── pod_cache.go
│   └── handler/           # HTTP 处理器
│       └── handler.go
├── deploy/                # Kubernetes 部署文件
│   ├── rbac.yaml         # RBAC 权限配置
│   ├── deployment.yaml   # Deployment 配置
│   ├── service.yaml      # Service 配置
│   └── kustomization.yaml
├── Dockerfile            # Docker 构建文件
├── go.mod               # Go 模块定义
└── README.md            # 项目文档
```

## 许可证

MIT

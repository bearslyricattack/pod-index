# Pod Index

一个基于 Kubernetes Informer 的高性能 Pod 信息缓存服务。

## 项目简介

Pod Index 使用 Kubernetes Informer 机制实时缓存集群中所有 Pod 的信息，并提供简洁的 HTTP API 供查询。该服务具有以下特点：

- **高效缓存**：利用 Kubernetes Informer 机制，自动同步并缓存所有 Pod 信息
- **快速查询**：通过 Pod UID 快速检索 Pod 详细信息，无需每次访问 API Server
- **资源友好**：轻量级设计，最小资源占用
- **生产就绪**：完整的健康检查、优雅关闭、错误处理机制
- **安全合规**：遵循最小权限原则，容器安全加固

## 快速开始

### 前置条件

- Go 1.21 或更高版本
- Kubernetes 集群（本地或远程）
- kubectl 已配置并可访问集群

### 本地运行

1. **克隆项目**
```bash
git clone <your-repo-url>
cd pod-index
```

2. **下载依赖**
```bash
go mod download
```

3. **运行服务**
```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

### 测试 API

```bash
# 健康检查
curl http://localhost:8080/health

# 就绪检查
curl http://localhost:8080/ready

# 查询 Pod 信息
POD_UID=$(kubectl get pod -n default -o jsonpath='{.items[0].metadata.uid}')
curl "http://localhost:8080/api/v1/pod?uid=${POD_UID}"
```

## API 文档

### 1. 查询 Pod 信息

**请求**
```
GET /api/v1/pod?uid={pod-uid}
```

**参数**
- `uid`: Pod 的 UID（必填）

**成功响应** (200 OK)
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

**错误响应** (404 Not Found)
```json
{
  "error": "pod with UID xxx not found"
}
```

### 2. 健康检查

**请求**
```
GET /health
```

**响应** (200 OK)
```json
{
  "status": "healthy"
}
```

### 3. 就绪检查

**请求**
```
GET /ready
```

**响应** (200 OK)
```json
{
  "status": "ready",
  "podCount": 150
}
```

## Docker 部署

### 构建镜像

```bash
docker build -t pod-index:latest .
```

### 本地运行（使用 kubeconfig）

```bash
docker run -d \
  -p 8080:8080 \
  -v ~/.kube/config:/root/.kube/config:ro \
  --name pod-index \
  pod-index:latest
```

### 测试容器

```bash
# 查看日志
docker logs -f pod-index

# 测试 API
curl http://localhost:8080/health
```

## Kubernetes 部署

### 方式一：使用 Kustomize（推荐）

```bash
# 部署所有资源
kubectl apply -k deploy/

# 查看部署状态
kubectl get pods -l app=pod-index
kubectl logs -l app=pod-index -f
```

### 方式二：使用自动化脚本

```bash
# 构建并部署
./scripts/build-and-deploy.sh
```

### 方式三：手动部署

```bash
# 1. 创建 RBAC 权限
kubectl apply -f deploy/rbac.yaml

# 2. 部署应用
kubectl apply -f deploy/deployment.yaml

# 3. 创建服务
kubectl apply -f deploy/service.yaml
```

### 验证部署

```bash
# 检查 Pod 状态
kubectl get pods -l app=pod-index

# 查看日志
kubectl logs -l app=pod-index -f

# 端口转发进行测试
kubectl port-forward svc/pod-index 8080:80

# 在另一个终端测试
./scripts/test-api.sh http://localhost:8080
```

### 卸载

```bash
kubectl delete -k deploy/
# 或
make undeploy
```

## Makefile 命令

项目提供了便捷的 Makefile 命令：

```bash
make help           # 显示所有可用命令
make build          # 编译应用
make run            # 本地运行
make test           # 运行测试
make docker-build   # 构建 Docker 镜像
make docker-run     # 运行 Docker 容器
make deploy         # 部署到 Kubernetes
make undeploy       # 从 Kubernetes 卸载
make clean          # 清理构建文件
make deps           # 下载并整理依赖
```

## 项目结构

```
pod-index/
├── main.go                      # 应用入口
├── pkg/
│   ├── cache/                  # Informer 缓存层
│   │   └── pod_cache.go       # Pod 缓存实现
│   └── handler/               # HTTP 处理层
│       └── handler.go         # API 处理器
├── deploy/                    # Kubernetes 部署文件
│   ├── rbac.yaml             # RBAC 权限配置
│   ├── deployment.yaml       # Deployment 配置
│   ├── service.yaml          # Service 配置
│   └── kustomization.yaml    # Kustomize 配置
├── scripts/                  # 辅助脚本
│   ├── build-and-deploy.sh  # 自动构建部署
│   └── test-api.sh          # API 测试脚本
├── Dockerfile               # Docker 构建文件
├── .dockerignore           # Docker 忽略文件
├── Makefile                # Make 命令
├── go.mod                  # Go 模块定义
├── go.sum                  # Go 依赖锁定
├── LICENSE                 # 开源协议
└── README.md              # 项目文档
```

## 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| PORT | HTTP 服务监听端口 | 8080 |

### Kubernetes 配置

应用会按以下顺序查找 Kubernetes 配置：

1. 集群内配置（In-Cluster Config）- 适用于部署在 Kubernetes 中
2. Kubeconfig 文件（`~/.kube/config`）- 适用于本地开发

### RBAC 权限

服务需要以下最小权限：

```yaml
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "watch"]
```

## 安全最佳实践

本项目实施了以下安全措施：

1. **最小权限原则**：仅请求必要的 RBAC 权限
2. **非 root 运行**：容器以 UID 65534 运行
3. **只读根文件系统**：防止运行时文件篡改
4. **禁用特权提升**：防止权限提升攻击
5. **删除所有 Capabilities**：最小化容器能力
6. **资源限制**：设置 CPU 和内存限制，防止资源耗尽

## 常见问题

### Q: 服务无法连接到 Kubernetes API

**A:** 检查以下几点：
- 确保 kubeconfig 配置正确（本地运行）
- 确保 ServiceAccount 和 RBAC 配置正确（集群运行）
- 检查网络连接和防火墙设置

### Q: Pod 查询返回 404

**A:** 可能原因：
- 缓存尚未同步完成，等待几秒后重试
- Pod UID 不正确，使用 `kubectl get pod -o yaml` 查看正确的 UID
- Pod 已被删除

### Q: 如何在生产环境部署？

**A:** 建议：
1. 修改镜像拉取策略为 `Always`
2. 根据集群规模调整资源限制
3. 配置 HPA（Horizontal Pod Autoscaler）如需要
4. 启用监控和日志收集
5. 使用 Ingress 或 LoadBalancer 暴露服务

### Q: 支持多集群吗？

**A:** 当前版本仅支持单集群。如需多集群支持，可以：
- 在每个集群部署一个实例
- 使用统一的前端聚合查询

## 性能考虑

- **内存占用**：约 100-200MB + (Pod 数量 × 2KB)
- **启动时间**：通常 5-10 秒完成缓存同步
- **查询延迟**：< 1ms（内存查询）
- **并发能力**：支持高并发读取

## 贡献指南

欢迎贡献！请：

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

如有问题或建议，请提交 Issue。

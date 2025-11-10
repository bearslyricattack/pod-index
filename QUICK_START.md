# 快速开始指南

## 三步部署到 Kubernetes

### 1. 构建镜像
```bash
docker build -t pod-index:latest .
```

如果使用 Minikube：
```bash
minikube image load pod-index:latest
```

如果使用 Kind：
```bash
kind load docker-image pod-index:latest
```

### 2. 部署到集群
```bash
kubectl apply -k deploy/
```

### 3. 验证部署
```bash
# 检查 Pod 状态
kubectl get pods -l app=pod-index

# 端口转发
kubectl port-forward svc/pod-index 8080:80
```

## 测试 API

```bash
# 健康检查
curl http://localhost:8080/health

# 获取任意 Pod 的 UID 并查询
POD_UID=$(kubectl get pod -A -o jsonpath='{.items[0].metadata.uid}')
curl "http://localhost:8080/api/v1/pod?uid=${POD_UID}"
```

## 一键部署（推荐）

```bash
# 自动构建、加载镜像并部署
./scripts/build-and-deploy.sh

# 自动测试 API
./scripts/test-api.sh http://localhost:8080
```

## 本地开发

```bash
# 运行服务（需要有效的 kubeconfig）
go run main.go

# 或使用 Make
make run
```

## 卸载

```bash
kubectl delete -k deploy/
```

## 常用命令

```bash
make help           # 查看所有可用命令
make build          # 编译
make docker-build   # 构建 Docker 镜像
make deploy         # 部署到 K8s
make undeploy       # 卸载
```

完整文档请参考 [README_CN.md](README_CN.md)

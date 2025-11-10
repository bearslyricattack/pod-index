#!/bin/bash
set -e

echo "======================================"
echo "Pod Index - 构建和部署脚本"
echo "======================================"

# 检查是否安装了必要工具
command -v docker >/dev/null 2>&1 || { echo "错误: 需要安装 Docker"; exit 1; }
command -v kubectl >/dev/null 2>&1 || { echo "错误: 需要安装 kubectl"; exit 1; }

# 构建 Docker 镜像
echo ""
echo "[1/3] 构建 Docker 镜像..."
docker build -t pod-index:latest .

# 如果使用 Minikube 或 Kind，加载镜像到集群
if command -v minikube >/dev/null 2>&1 && minikube status >/dev/null 2>&1; then
    echo ""
    echo "[2/3] 加载镜像到 Minikube..."
    minikube image load pod-index:latest
elif command -v kind >/dev/null 2>&1; then
    CLUSTER_NAME=${KIND_CLUSTER_NAME:-kind}
    if kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
        echo ""
        echo "[2/3] 加载镜像到 Kind..."
        kind load docker-image pod-index:latest --name ${CLUSTER_NAME}
    fi
else
    echo ""
    echo "[2/3] 跳过镜像加载（非本地集群）"
fi

# 部署到 Kubernetes
echo ""
echo "[3/3] 部署到 Kubernetes..."
kubectl apply -k deploy/

echo ""
echo "======================================"
echo "部署完成！"
echo "======================================"
echo ""
echo "检查部署状态："
echo "  kubectl get pods -l app=pod-index"
echo ""
echo "查看日志："
echo "  kubectl logs -l app=pod-index -f"
echo ""
echo "端口转发测试："
echo "  kubectl port-forward svc/pod-index 8080:80"
echo ""

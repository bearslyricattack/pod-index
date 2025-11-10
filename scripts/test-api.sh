#!/bin/bash
set -e

BASE_URL=${1:-"http://localhost:8080"}

echo "======================================"
echo "Pod Index - API 测试脚本"
echo "======================================"
echo "测试地址: $BASE_URL"
echo ""

# 测试健康检查
echo "[1/3] 测试健康检查..."
curl -s "${BASE_URL}/health" | jq .
echo ""

# 测试就绪检查
echo "[2/3] 测试就绪检查..."
curl -s "${BASE_URL}/ready" | jq .
echo ""

# 测试 Pod 查询
echo "[3/3] 测试 Pod 查询..."
POD_UID=$(kubectl get pod -A -o jsonpath='{.items[0].metadata.uid}' 2>/dev/null || echo "")

if [ -z "$POD_UID" ]; then
    echo "警告: 无法获取 Pod UID，跳过查询测试"
else
    POD_NAME=$(kubectl get pod -A -o jsonpath='{.items[0].metadata.name}')
    POD_NS=$(kubectl get pod -A -o jsonpath='{.items[0].metadata.namespace}')
    echo "查询 Pod: $POD_NS/$POD_NAME (UID: $POD_UID)"
    curl -s "${BASE_URL}/api/v1/pod?uid=${POD_UID}" | jq .
fi

echo ""
echo "======================================"
echo "测试完成！"
echo "======================================"

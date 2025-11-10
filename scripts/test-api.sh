#!/bin/bash
set -e

BASE_URL=${1:-"http://localhost:8080"}

echo "======================================"
echo "Pod Index - API Test Script"
echo "======================================"
echo "Test URL: $BASE_URL"
echo ""

# Test health check
echo "[1/3] Testing health check..."
curl -s "${BASE_URL}/health" | jq .
echo ""

# Test readiness check
echo "[2/3] Testing readiness check..."
curl -s "${BASE_URL}/ready" | jq .
echo ""

# Test Pod query
echo "[3/3] Testing Pod query..."
POD_UID=$(kubectl get pod -A -o jsonpath='{.items[0].metadata.uid}' 2>/dev/null || echo "")

if [ -z "$POD_UID" ]; then
    echo "Warning: Unable to get Pod UID, skipping query test"
else
    POD_NAME=$(kubectl get pod -A -o jsonpath='{.items[0].metadata.name}')
    POD_NS=$(kubectl get pod -A -o jsonpath='{.items[0].metadata.namespace}')
    echo "Querying Pod: $POD_NS/$POD_NAME (UID: $POD_UID)"
    curl -s "${BASE_URL}/api/v1/pod?uid=${POD_UID}" | jq .
fi

echo ""
echo "======================================"
echo "Test completed!"
echo "======================================"

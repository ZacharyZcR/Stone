#!/bin/bash

echo "=== 测试Stone WAF速率限制功能（快速请求） ==="
echo "规则：每分钟最多3个请求"
echo "测试：快速发送10个请求，几乎同时"
echo

# 快速发送多个请求，不间隔等待
for i in {1..10}; do
  echo -n "请求 $i: "
  RESPONSE=$(curl -s -w "HTTP:%{http_code}" -m 5 http://localhost:8082/test-$i 2>/dev/null)
  echo "$RESPONSE"
done

echo
echo "=== 查看速率限制统计信息 ==="
sleep 1
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTY5OTk1NTksInJvbGUiOiJ1c2VyIiwidXNlcm5hbWUiOiJ0ZXN0dXNlciJ9.4be1lZwJuYvmGlE28souyt12hOSVZf21E9mU7RZA0LQ"
curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8083/rate-limit/stats
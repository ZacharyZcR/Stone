#!/bin/bash

echo "=== 测试Stone WAF速率限制功能 ==="
echo "规则：每分钟最多3个请求"
echo "目标端口：8082 (WAF流量捕获端口)"
echo

# 测试正常请求（应该通过）
echo "--- 发送正常请求（前3个应该成功） ---"
for i in {1..3}; do
  echo -n "请求 $i: "
  RESPONSE=$(curl -s -w "HTTP:%{http_code}" -m 5 http://localhost:8082/test 2>/dev/null)
  echo "$RESPONSE"
  sleep 1
done

echo
echo "--- 发送超限请求（应该被阻止） ---"
for i in {4..6}; do
  echo -n "请求 $i: "
  RESPONSE=$(curl -s -w "HTTP:%{http_code}" -m 5 http://localhost:8082/test 2>/dev/null)
  echo "$RESPONSE"
  sleep 1
done

echo
echo "=== 等待5秒后测试令牌桶恢复 ==="
sleep 5

for i in {7..8}; do
  echo -n "请求 $i: "
  RESPONSE=$(curl -s -w "HTTP:%{http_code}" -m 5 http://localhost:8082/test 2>/dev/null)
  echo "$RESPONSE"
  sleep 1
done
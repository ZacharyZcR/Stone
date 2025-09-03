# scripts/setup_iptables.sh

#!/bin/bash

# 配置iptables规则，将80端口的流量重定向到8080端口
iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 8080

echo "iptables规则已设置: 将80端口流量重定向到8080端口"
package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// 定义Prometheus指标
var (
	WebsiteRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "website_requests_total",
			Help: "防火墙所保护的网站的访问总次数",
		},
		[]string{"status"},
	)
	BlockedByBlacklistTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "blocked_by_blacklist_total",
			Help: "被黑名单拦截的请求总次数",
		},
	)
	BlockedByRulesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "blocked_by_rules_total",
			Help: "被规则拦截的请求总次数",
		},
	)
)

// InitMonitoring 初始化Prometheus监控
func InitMonitoring() {
	// 注册指标
	prometheus.MustRegister(WebsiteRequestsTotal)
	prometheus.MustRegister(BlockedByBlacklistTotal)
	prometheus.MustRegister(BlockedByRulesTotal)

	// 启动HTTP服务器以暴露指标
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)
}

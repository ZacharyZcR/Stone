// pkg/monitoring/monitoring.go

package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// 定义Prometheus指标
var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "stone_requests_total",
			Help: "处理的总请求数",
		},
		[]string{"status"},
	)
)

// InitMonitoring 初始化Prometheus监控
func InitMonitoring() {
	// 注册指标
	prometheus.MustRegister(RequestsTotal)

	// 启动HTTP服务器以暴露指标
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)
}

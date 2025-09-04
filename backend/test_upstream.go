// test_upstream.go - 简单的上游服务器用于测试WAF速率限制
package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s %s from %s\n", 
			time.Now().Format("15:04:05"), 
			r.Method, 
			r.URL.Path, 
			r.RemoteAddr)
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, `{"message": "Hello from upstream server", "path": "%s", "time": "%s"}`, 
			r.URL.Path, time.Now().Format("15:04:05"))
	})

	fmt.Println("上游测试服务器启动在端口 9001")
	if err := http.ListenAndServe(":9001", nil); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
	}
}
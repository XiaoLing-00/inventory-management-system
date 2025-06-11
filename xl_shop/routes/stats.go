package routes

import (
	"encoding/json"
	"net/http"
	"xl_shop/utils"
)

// 注册统计数据路由
func RegisterStatsRoute(warehouse *utils.Warehouse) {
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "方法不支持", http.StatusMethodNotAllowed)
			return
		}
		handleGetStats(w, r, warehouse)
	})
}

// 处理获取统计数据的请求
func handleGetStats(w http.ResponseWriter, _ *http.Request, warehouse *utils.Warehouse) {
	stats := warehouse.CalculateStats()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

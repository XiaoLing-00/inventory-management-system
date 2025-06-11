// routes/records.go

package routes

import (
	"encoding/json"
	"net/http"
	"xl_shop/utils"
)

// 注册操作记录查询路由
func RegisterRecordsRoute(warehouse *utils.Warehouse) {
	http.HandleFunc("/records", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "方法不支持", http.StatusMethodNotAllowed)
			return
		}
		handleGetRecords(w, r, warehouse)
	})
}

// 处理获取所有记录的请求
func handleGetRecords(w http.ResponseWriter, _ *http.Request, warehouse *utils.Warehouse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(warehouse.Records)
}

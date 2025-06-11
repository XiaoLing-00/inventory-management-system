package routes

import (
	"encoding/json"
	"net/http"
	"xl_shop/utils"
)

// 注册库存查询路由
func RegisterInventoryRoute(warehouse *utils.Warehouse) {
	http.HandleFunc("/inventory", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "方法不支持", http.StatusMethodNotAllowed)
			return
		}
		handleInventory(w, r, warehouse)
	})
}

func handleInventory(w http.ResponseWriter, r *http.Request, warehouse *utils.Warehouse) {
	keyword := r.URL.Query().Get("keyword")
	results := warehouse.QueryInventory(keyword)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"xl_shop/utils"
)

// 注册出库路由
func RegisterStockOutRoute(warehouse *utils.Warehouse) {
	http.HandleFunc("/stock/out", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "方法不支持", http.StatusMethodNotAllowed)
			return
		}
		handleStockOut(w, r, warehouse)
	})
}

func handleStockOut(w http.ResponseWriter, r *http.Request, warehouse *utils.Warehouse) {
	productID := r.URL.Query().Get("product_id")
	opt := r.URL.Query().Get("operator")
	quantityStr := r.URL.Query().Get("quantity")

	if productID == "" || opt == "" || quantityStr == "" {
		http.Error(w, "缺少必要参数", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		http.Error(w, "数量必须为整数", http.StatusBadRequest)
		return
	}

	if err := warehouse.StockOut(productID, opt, quantity); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "出库成功"})
}

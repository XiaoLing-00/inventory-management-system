package routes

import (
    "encoding/json"
    "net/http"
    "xl_shop/utils"
)

// 注册商品相关路由
func RegisterProductRoutes(warehouse *utils.Warehouse) {
    http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "POST":
            handleProductPost(w, r, warehouse)
        case "GET":
            handleProductGet(w, r, warehouse)
        default:
            http.Error(w, "方法不支持", http.StatusMethodNotAllowed)
        }
    })
}

// 处理商品创建（POST）
func handleProductPost(w http.ResponseWriter, r *http.Request, warehouse *utils.Warehouse) {
    var product utils.Product
    if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
        http.Error(w, "无效的请求数据", http.StatusBadRequest)
        return
    }

    warehouse.AddProduct(&product)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(product)
}

// 处理商品查询（GET）
func handleProductGet(w http.ResponseWriter, _ *http.Request, warehouse *utils.Warehouse) {
    products := make([]utils.Product, 0, len(warehouse.Products))
    for _, product := range warehouse.Products {
        products = append(products, *product)
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(products)
}
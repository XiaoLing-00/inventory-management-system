// main.go

package main

import (
	"fmt"
	"net/http"
	"xl_shop/routes"
	"xl_shop/utils"
)

func main() {
	// 初始化仓库系统
	inventorySystem := utils.NewWarehouse()

	// 注册所有路由
	routes.RegisterProductRoutes(inventorySystem)
	routes.RegisterStockInRoute(inventorySystem)
	routes.RegisterStockOutRoute(inventorySystem)
	routes.RegisterInventoryRoute(inventorySystem)
	routes.RegisterRecordsRoute(inventorySystem) 
	routes.RegisterStatsRoute(inventorySystem)   // 注册统计路由

	// 静态文件服务
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// 启动服务器
	fmt.Println("服务器正在运行在端口 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
	}
}

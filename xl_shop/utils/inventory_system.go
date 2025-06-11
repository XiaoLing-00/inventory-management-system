package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// 商品结构体
type Product struct {
	ID          string `json:"ID"`
	Name        string `json:"Name"`
	SafetyStock int    `json:"SafetyStock"`
	Stock       int    `json:"Stock"`
}

// 操作记录结构体
type Record struct {
	ProductID string `json:"ProductID"`
	Operation string `json:"Operation"` 
	Quantity  int    `json:"Quantity"`
	Operator  string `json:"Operator"`
	Timestamp int64  `json:"Timestamp"`
}

// 仓库结构体
type Warehouse struct {
	Products map[string]*Product
	Records  []Record
	NextID   int
}

// 统计数据结构体
type Stats struct {
	TotalProducts int `json:"totalProducts"`
	TodayIn       int `json:"todayIn"`
	TodayOut      int `json:"todayOut"`
	LowStock      int `json:"lowStock"`
}

// 创建新的仓库实例
func NewWarehouse() *Warehouse {
	return &Warehouse{
		Products: make(map[string]*Product),
		Records:  []Record{},
		NextID:   1,
	}
}

// 添加商品(自动生成ID)
func (w *Warehouse) AddProduct(p *Product) {
	p.ID = strconv.Itoa(w.NextID)
	w.NextID++
	p.Stock = 0
	w.Products[p.ID] = p
}


func (w *Warehouse) StockIn(id, operator string, qty int) error {
	product, exists := w.Products[id]
	if !exists {
		return errors.New("商品不存在")
	}
	product.Stock += qty
	w.Records = append(w.Records, Record{
		ProductID: id,
		Operation: "in",
		Quantity:  qty,
		Operator:  operator,

		Timestamp: time.Now().Unix(),
	})
	return nil
}

func (w *Warehouse) StockOut(id, operator string, qty int) error {
	product, exists := w.Products[id]
	if !exists {
		return errors.New("商品不存在")
	}
	if product.Stock < qty {
		return errors.New("库存不足")
	}
	product.Stock -= qty
	w.Records = append(w.Records, Record{
		ProductID: id,
		Operation: "out",
		Quantity:  qty,
		Operator:  operator,
	
		Timestamp: time.Now().Unix(),
	})
	return nil
}

// 查询库存(根据商品名模糊匹配)
func (w *Warehouse) QueryInventory(keyword string) []*Product {
	var result []*Product
	// 如果关键字为空，返回所有商品
	if keyword == "" {
		for _, p := range w.Products {
			result = append(result, p)
		}
		return result
	}
	// 如果有关键字，则进行模糊搜索
	for _, p := range w.Products {
		if strings.Contains(p.Name, keyword) {
			result = append(result, p)
		}
	}
	return result
}

// 计算统计数据的方法
func (w *Warehouse) CalculateStats() Stats {
	var stats Stats
	stats.TotalProducts = len(w.Products)

	// 计算低库存商品
	for _, product := range w.Products {
		if product.Stock <= product.SafetyStock {
			stats.LowStock++
		}
	}

	// 计算今日出入库
	now := time.Now()
	year, month, day := now.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, now.Location()).Unix()

	for _, record := range w.Records {
		if record.Timestamp >= startOfDay {
			if record.Operation == "in" {
				stats.TodayIn += record.Quantity
			} else if record.Operation == "out" {
				stats.TodayOut += record.Quantity
			}
		}
	}
	return stats
}

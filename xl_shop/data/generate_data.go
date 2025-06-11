package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Product represents the product structure for creating test data.
type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

// createProduct sends a POST request to /products endpoint to create a new product.
func createProduct(product Product) error {
	url := "http://localhost:8080/products"
	body, _ := json.Marshal(product)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create product, status code: %d", resp.StatusCode)
	}

	return nil
}

// stockIn sends a POST request to /stock/in endpoint to increase stock of a product.
func stockIn(productID, operator string, quantity int) error {
	url := "http://localhost:8080/stock/in?product_id=" + productID + "&operator=" + operator + "&quantity=" + strconv.Itoa(quantity)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to stock in, status code: %d", resp.StatusCode)
	}

	return nil
}

// stockOut sends a POST request to /stock/out endpoint to decrease stock of a product.
func stockOut(productID, operator string, quantity int) error {
	url := "http://localhost:8080/stock/out?product_id=" + productID + "&operator=" + operator + "&quantity=" + strconv.Itoa(quantity)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to stock out, status code: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	// 创建商品
	products := []Product{
		{ID: "1", Name: "Test Product 1", Price: 99.99, Category: "TestCategory"},
		{ID: "2", Name: "Test Product 2", Price: 199.99, Category: "TestCategory"},
		{ID: "3", Name: "Test Product 3", Price: 299.99, Category: "TestCategory"},
		{ID: "4", Name: "Test Product 4", Price: 399.99, Category: "TestCategory"},
		{ID: "5", Name: "Test Product 5", Price: 499.99, Category: "TestCategory"},
	}

	for _, product := range products {
		if err := createProduct(product); err != nil {
			fmt.Printf("Error creating product %s: %v\n", product.ID, err)
		} else {
			fmt.Printf("Product %s created successfully\n", product.ID)
		}
	}

	// 等待一段时间以确保服务器有时间处理上一个请求
	time.Sleep(2 * time.Second)

	// 商品入库/出库操作
	operations := []struct {
		productID string
		operator  string
		quantity  int
		isStockIn bool
	}{
		{"1", "admin", 100, true},
		{"2", "admin", 200, true},
		{"3", "admin", 300, true},
		{"4", "admin", 400, true},
		{"5", "admin", 500, true},
		{"1", "admin", 50, false},
		{"2", "admin", 100, false},
		{"3", "admin", 150, false},
		{"4", "admin", 200, false},
		{"5", "admin", 250, false},
		{"1", "admin", 30, true},
		{"2", "admin", 60, true},
		{"3", "admin", 90, true},
		{"4", "admin", 120, true},
		{"5", "admin", 150, true},
	}

	for i, op := range operations {
		var err error
		if op.isStockIn {
			err = stockIn(op.productID, op.operator, op.quantity)
		} else {
			err = stockOut(op.productID, op.operator, op.quantity)
		}

		if err != nil {
			fmt.Printf("Error operation %d: %v\n", i+1, err)
		} else {
			operationType := "in"
			if !op.isStockIn {
				operationType = "out"
			}
			fmt.Printf("Operation %d successful: Product %s stock %s by %d units\n", i+1, op.productID, operationType, op.quantity)
		}
		// 短暂等待，避免请求过于密集
		time.Sleep(500 * time.Millisecond)
	}
}

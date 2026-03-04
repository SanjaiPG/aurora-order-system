package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OrderRequest struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	Price     int `json:"price"`
}

func main() {

	req := OrderRequest{
		UserID:    1,
		ProductID: 1,
		Price:     1000,
	}

	data, _ := json.Marshal(req)

	resp, err := http.Post(
		"http://ORDER_SERVICE_IP:8080/place-order",
		"application/json",
		bytes.NewBuffer(data),
	)

	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println("Order Response:", resp.Status)
}
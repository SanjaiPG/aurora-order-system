package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

type InventoryRequest struct {
	ProductID int `json:"product_id"`
}

type StockResponse struct {
	Stock int `json:"stock"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	fmt.Println("Inventory Service connected to DB")

	http.HandleFunc("/check-stock", checkStock)
	http.HandleFunc("/update-stock", updateStock)

	fmt.Println("Inventory Service running on port 8081")

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func checkStock(w http.ResponseWriter, r *http.Request) {

	var req InventoryRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", 400)
		return
	}

	var stock int

	err = db.QueryRow(
		"SELECT stock FROM inventory WHERE product_id=$1",
		req.ProductID,
	).Scan(&stock)

	if err != nil {
		http.Error(w, "Product not found", 404)
		return
	}

	resp := StockResponse{Stock: stock}

	json.NewEncoder(w).Encode(resp)
}

func updateStock(w http.ResponseWriter, r *http.Request) {

	var req InventoryRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", 400)
		return
	}

	_, err = db.Exec(
		"UPDATE inventory SET stock=stock-1 WHERE product_id=$1",
		req.ProductID,
	)

	if err != nil {
		http.Error(w, "Stock update failed", 500)
		return
	}

	w.Write([]byte("Stock updated"))
}
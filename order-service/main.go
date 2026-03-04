package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var writerDB *sql.DB
var readerDB *sql.DB

type OrderRequest struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	Price     int `json:"price"`
}

type StockResponse struct {
	Stock int `json:"stock"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	inventoryService := os.Getenv("INVENTORY_SERVICE_URL")

	// Primary DB (writer)
	writerConn := fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname,
	)

	writerDB, err = sql.Open("postgres", writerConn)
	if err != nil {
		log.Fatal(err)
	}

	// Replica DB (reader)
	readerConn := fmt.Sprintf(
		"host=%s port=5433 user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname,
	)

	readerDB, err = sql.Open("postgres", readerConn)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/place-order", enableCORS(placeOrder(inventoryService)))

	fmt.Println("Order Service running on port 8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		if r.Method == "OPTIONS" {
			return
		}

		next(w, r)
	}
}

func placeOrder(inventoryService string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var req OrderRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON", 400)
			return
		}

		// Check stock
		body, _ := json.Marshal(map[string]int{
			"product_id": req.ProductID,
		})

		resp, err := http.Post(
			inventoryService+"/check-stock",
			"application/json",
			bytes.NewBuffer(body),
		)

		if err != nil {
			http.Error(w, "Inventory service error", 500)
			return
		}
		defer resp.Body.Close()

		var stockResp StockResponse
		json.NewDecoder(resp.Body).Decode(&stockResp)

		if stockResp.Stock <= 0 {
			http.Error(w, "Out of stock", 400)
			return
		}

		// Read balance from replica
		var balance int

		err = readerDB.QueryRow(
			"SELECT balance FROM users WHERE user_id=$1",
			req.UserID,
		).Scan(&balance)

		if err != nil {
			http.Error(w, "User not found", 400)
			return
		}

		if balance < req.Price {
			http.Error(w, "Insufficient balance", 400)
			return
		}

		// Write transaction to primary
		tx, err := writerDB.Begin()
		if err != nil {
			http.Error(w, "Transaction failed", 500)
			return
		}

		_, err = tx.Exec(
			"UPDATE users SET balance = balance - $1 WHERE user_id = $2",
			req.Price,
			req.UserID,
		)

		if err != nil {
			tx.Rollback()
			http.Error(w, "Payment failed", 500)
			return
		}

		_, err = tx.Exec(
			"INSERT INTO orders (user_id,product_id,status) VALUES ($1,$2,'CONFIRMED')",
			req.UserID,
			req.ProductID,
		)

		if err != nil {
			tx.Rollback()
			http.Error(w, "Order failed", 500)
			return
		}

		tx.Commit()

		// Update stock
		http.Post(
			inventoryService+"/update-stock",
			"application/json",
			bytes.NewBuffer(body),
		)

		w.Write([]byte("Order placed successfully"))
	}
}
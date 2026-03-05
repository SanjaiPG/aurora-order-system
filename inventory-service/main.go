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

var pendingProduct int

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname,
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

	http.HandleFunc("/prepare-stock", prepareStock)
	http.HandleFunc("/commit-stock", commitStock)
	http.HandleFunc("/abort-stock", abortStock)

	fmt.Println("Inventory Service running on port 8081")

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func prepareStock(w http.ResponseWriter, r *http.Request) {

	var req InventoryRequest
	json.NewDecoder(r.Body).Decode(&req)

	var stock int

	err := db.QueryRow(
		"SELECT stock FROM inventory WHERE product_id=$1",
		req.ProductID,
	).Scan(&stock)

	if err != nil || stock <= 0 {
		http.Error(w, "ABORT", 400)
		return
	}

	pendingProduct = req.ProductID

	w.Write([]byte("READY"))
}

func commitStock(w http.ResponseWriter, r *http.Request) {

	_, err := db.Exec(
		"UPDATE inventory SET stock = stock - 1 WHERE product_id=$1",
		pendingProduct,
	)

	if err != nil {
		http.Error(w, "Commit failed", 500)
		return
	}

	w.Write([]byte("COMMIT OK"))
}

func abortStock(w http.ResponseWriter, r *http.Request) {

	pendingProduct = 0

	w.Write([]byte("ABORTED"))
}
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

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
		log.Fatal("Database connection failed:", err)
	}

	fmt.Println("Connected to PostgreSQL successfully")

	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/place-order", placeOrder)

	fmt.Println("Order Service running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Order Service is running"))
}

func placeOrder(w http.ResponseWriter, r *http.Request) {

	userID := 1
	productID := 1
	price := 1000

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", 500)
		return
	}

	var stock int
	err = tx.QueryRow("SELECT stock FROM inventory WHERE product_id=$1 FOR UPDATE", productID).Scan(&stock)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Product not found", 400)
		return
	}

	if stock <= 0 {
		tx.Rollback()
		http.Error(w, "Out of stock", 400)
		return
	}

	var balance int
	err = tx.QueryRow("SELECT balance FROM users WHERE user_id=$1 FOR UPDATE", userID).Scan(&balance)
	if err != nil {
		tx.Rollback()
		http.Error(w, "User not found", 400)
		return
	}

	if balance < price {
		tx.Rollback()
		http.Error(w, "Insufficient balance", 400)
		return
	}

	_, err = tx.Exec("UPDATE inventory SET stock=stock-1 WHERE product_id=$1", productID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Inventory update failed", 500)
		return
	}

	_, err = tx.Exec("UPDATE users SET balance=balance-$1 WHERE user_id=$2", price, userID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Payment failed", 500)
		return
	}

	_, err = tx.Exec("INSERT INTO orders (user_id, product_id, status) VALUES ($1,$2,'CONFIRMED')", userID, productID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Order insert failed", 500)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Commit failed", 500)
		return
	}

	w.Write([]byte("Order placed successfully"))
}
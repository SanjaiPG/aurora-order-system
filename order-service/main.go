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

	fmt.Println("Order Service running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Order Service is running"))
}
package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Kudzeri/go-postgres-api/models"
	"log"
	"net/http"
	"os"
	"strconv"
)
import "github.com/joho/godotenv"

type Response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateConnection() *sql.DB {
	funcName := "handlers.CreateConnection():"

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(funcName + "Error loading .env file")
	}
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(funcName + "Error connecting to database")
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(funcName + "Error pinging database")
	}

	fmt.Println("Successfully connected to database")
	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	funcName := "handlers.CreateStock():"
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		fmt.Println(funcName + "Unable to decode the request body")
	}
	insertID := insertStock(stock)

	res := Response{
		ID:      insertID,
		Message: "stock created successfully",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	funcName := "handlers.GetStock():"
	params := r.URL.Query()

	id, err := strconv.Atoi(params.Get("id"))
	if err != nil {
		fmt.Println(funcName + "Unable to convert the string into int")
	}

	stock, err := getStock(int64(id))
	if err != nil {
		fmt.Println(funcName+"Unable to get stock.", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stock)
}

func GetAllStocks(w http.ResponseWriter, r *http.Request) {
	funcName := "handlers.GetAllStocks():"
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	funcName := "handlers.UpdateStock():"
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	funcName := "handlers.DeleteStock():"
}

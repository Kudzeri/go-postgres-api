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
		log.Fatal(funcName + "Error loading .env file.")
	}
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal(funcName + "Error connecting to database.")
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(funcName + "Error pinging database.")
	}

	fmt.Println("Successfully connected to database.")
	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	funcName := "handlers.CreateStock():"
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		fmt.Println(funcName + "Unable to decode the request body.")
	}
	insertID := insertStock(stock)

	res := Response{
		ID:      insertID,
		Message: "Stock created successfully.",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func insertStock(stock models.Stock) int64 {
	funcName := "handlers.insertStock():"
	db := CreateConnection()
	defer db.Close()

	sqlQuery := `INSERT INTO stocks(name,price,company) VALUES ($1,$2,$3) RETURNING stockid`
	var id int64

	err := db.QueryRow(sqlQuery, stock.Name, stock.Price, stock.Company).Scan(&id)
	if err != nil {
		fmt.Println(funcName + "Unable to execute the query.")
	}

	fmt.Printf("Inserted a single record %v", id)
	return id
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	funcName := "handlers.GetStock():"
	params := r.URL.Query()

	id, err := strconv.Atoi(params.Get("id"))
	if err != nil {
		fmt.Println(funcName + "Unable to convert the string into int.")
	}

	stock, err := getStock(int64(id))
	if err != nil {
		fmt.Println(funcName+"Unable to get stock.", err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stock)
}

func getStock(id int64) (models.Stock, error) {
	funcName := "handlers.getStock():"
	db := CreateConnection()
	defer db.Close()

	sqlQuery := `SELECT * FROM stocks WHERE stockid=$1`
	var stock models.Stock

	row := db.QueryRow(sqlQuery, id)
	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println(funcName + "Now rows were returned!")
		return stock, nil
	case nil:
		return stock, nil
	default:
		fmt.Println(funcName + "Unable to scan the row")
	}

	return stock, err
}

func GetAllStocks(w http.ResponseWriter, r *http.Request) {
	funcName := "handlers.GetAllStocks():"

	stocks, err := getAllStocks()
	if err != nil {
		fmt.Println(funcName + "Unable to get stocks.")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stocks)
}

func getAllStocks() ([]models.Stock, error) {
	funcName := "handlers.getAllStocks():"
	db := CreateConnection()
	defer db.Close()

	sqlQuery := `SELECT * FROM stocks`
	var stocks []models.Stock

	rows, err := db.Query(sqlQuery)
	if err != nil {
		fmt.Println(funcName + "Unable to execute the query.")
	}

	defer rows.Close()
	for rows.Next() {
		var stock models.Stock
		err := rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			fmt.Println(funcName + "Unable to scan the row.")
			stocks = append(stocks, stock)
		}
	}

	return stocks, err
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	funcName := "handlers.UpdateStock():"
	params := r.URL.Query()

	id, err := strconv.Atoi(params.Get("id"))
	if err != nil {
		fmt.Println(funcName + "Unable to convert the string into int.")
	}

	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		fmt.Println(funcName + "Unable to decode the request body.")
	}

	updatedRows := updateStock(int64(id), stock)

	msg := fmt.Sprintf("Stock updated successfully. Total rows/records affected %v.", updatedRows)
	res := Response{
		ID:      int64(id),
		Message: msg,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func updateStock(id int64, stock models.Stock) int64 {
	funcName := "handlers.updateStock():"
	db := CreateConnection()
	defer db.Close()

	sqlQuery := `UPDATE stocks SET name = $2, price = $3, company = $4 WHERE stockid = $1`
	res, err := db.Exec(sqlQuery, id, stock.Name, stock.Price, stock.Company)
	if err != nil {
		fmt.Println(funcName + "Unable to execute the query.")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Println(funcName + "Error while checking the affected rows.")
	}

	fmt.Printf("Total rows/records affected %v", rowsAffected)
	return rowsAffected
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	funcName := "handlers.DeleteStock():"
	params := r.URL.Query()

	id, err := strconv.Atoi(params.Get("id"))
	if err != nil {
		fmt.Println(funcName + "Unable to convert the string into int.")
	}

	deletedRows := deleteStock(int64(id))
	msg := fmt.Sprintf("Stock deleted successfully.Total rows/records %v", deletedRows)
	res := Response{
		ID:      int64(id),
		Message: msg,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func deleteStock(id int64) int64 {
	funcName := "handlers.deleteStock():"
	db := CreateConnection()
	defer db.Close()

	sqlQuery := `DELETE FROM stocks WHERE stockid = $1`
	res, err := db.Exec(sqlQuery, id)
	if err != nil {
		fmt.Println(funcName + "Unable to execute the query.")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Println(funcName + "Error while checking the affected rows.")
	}

	fmt.Printf("Total rows/records affected %v", rowsAffected)
	return rowsAffected
}

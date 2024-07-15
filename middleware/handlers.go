package middleware

import (
	"database/sql"
	"log"
)
import "github.com/joho/godotenv"

type Response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func CreateConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("handlers.CreateConnection():Can't upload the env file")
	}
}

package router

import "github.com/gorilla/mux"

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/stocks/{id}", middleware.GetStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/stocks", middleware.GetAllStocks).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/stocks/new", middleware.CreateStock).Methods("POST", "OPTIONS")
	router.HandleFunc("/", middleware.UpdateStock).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/stocks/{id}", middleware.DeleteStock).Methods("DELETE", "OPTIONS")

}

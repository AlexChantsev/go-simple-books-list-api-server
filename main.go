package main

import (
	"database/sql"
	"log"
	"net/http"

	"go-simple-books-list-api-server/controllers"
	"go-simple-books-list-api-server/driver"
	"go-simple-books-list-api-server/models"

	"github.com/subosito/gotenv"

	"github.com/gorilla/mux"
)

var books []models.Book
var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db = driver.ConnectDB()

	router := mux.NewRouter()

	controller := controllers.Controller{}

	router.HandleFunc("/books", controller.GetBooks(db)).Methods("GET")
	router.HandleFunc("/books/{id}", controller.GetBook(db)).Methods("GET")
	router.HandleFunc("/books", controller.AddBook(db)).Methods("POST")
	router.HandleFunc("/books", controller.UpdateBook(db)).Methods("PUT")
	router.HandleFunc("/books/{id}", controller.RemoveBook(db)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
